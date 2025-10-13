package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"

	"aiagentapi/auth"
	"aiagentapi/storage"
)

type chatMessage struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id"`
}

var chatHistory = make([]chatMessage, 0, 1024)

// Chat returns the POST /chat handler using DB + OpenAI and a simple RAG-lite search.
// It will NOT crash if tables are missing; it just answers without context.
func Chat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := auth.GetCurrentUser(c, db)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}
		userID := user.ID

		var req struct{ Message string ` + "`json:\"message\"`" + ` }
		if err := c.BindJSON(&req); err != nil || strings.TrimSpace(req.Message) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "message required"})
			return
		}

		now := time.Now()
		chatHistory = append(chatHistory, chatMessage{
			Role:      "user",
			Content:   req.Message,
			CreatedAt: now,
			UserID:    userID,
		})

		// Demo: enqueue a "tick" task for the worker (non-blocking)
		_, _ = storage.Enqueue(c.Request.Context(), db, userID, "demo_chat_tick", map[string]any{"at": now.UTC().Format(time.RFC3339)}, nil, nil)

		// RAG-lite: try to pull a few snippets from email / notes / contacts
		ctx := c.Request.Context()
		snips := findSnippets(ctx, db, userID, req.Message, 6)

		sys := ` + "`" + `You are an assistant for a financial advisor. 
Use the provided context snippets when relevant. 
If a user asks to act (email, schedule, log note), describe the next steps you will perform.` + "`" + `

		userPrompt := buildUserPrompt(req.Message, snips)

		reply := callLLM(ctx, sys, userPrompt)

		chatHistory = append(chatHistory, chatMessage{
			Role:      "assistant",
			Content:   reply,
			CreatedAt: time.Now(),
			UserID:    userID,
		})

		c.JSON(http.StatusOK, gin.H{"reply": reply, "snippets": snips})
	}
}

// buildUserPrompt formats the user question plus a compact list of snippets.
func buildUserPrompt(question string, snips []string) string {
	var b strings.Builder
	if len(snips) > 0 {
		b.WriteString("Context snippets:\n")
		for i, s := range snips {
			if i >= 6 { break }
			b.WriteString("- ")
			if len(s) > 400 { s = s[:400] + "…" }
			b.WriteString(s)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("User question: ")
	b.WriteString(question)
	return b.String()
}

// callLLM sends to OpenAI if OPENAI_API_KEY is set; otherwise returns a fallback.
func callLLM(ctx context.Context, system, user string) string {
	key := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")
	if model == "" { model = "gpt-4o-mini" }

	if strings.TrimSpace(key) == "" {
		return "I received your message. To enable AI answers, set OPENAI_API_KEY in the environment."
	}

	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
		Temperature: 0.2,
	})
	if err != nil || len(resp.Choices) == 0 {
		if err != nil {
			return fmt.Sprintf("LLM error: %v", err)
		}
		return "No response from model."
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content)
}

// findSnippets performs best-effort text search across likely tables.
// It won't error if tables/columns don't exist; it just returns fewer results.
func findSnippets(ctx context.Context, db *sql.DB, userID, q string, limit int) []string {
	q = strings.TrimSpace(q)
	if q == "" { return nil }

	type row struct{ S string }
	snips := make([]string, 0, limit)

	// helper to run a query safely
	run := func(sqlStr string, args ...any) {
		if len(snips) >= limit { return }
		ctx2, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		rows, err := db.QueryContext(ctx2, sqlStr, args...)
		if err != nil { return }
		defer rows.Close()
		for rows.Next() {
			if len(snips) >= limit { break }
			var r row
			if err := rows.Scan(&r.S); err == nil && strings.TrimSpace(r.S) != "" {
				snips = append(snips, r.S)
			}
		}
	}

	like := "%" + q + "%"

	// Try emails table variants
	run(` + "`" + `SELECT subject || ' — ' || left(coalesce(body_text,snippet,''), 300)
	     FROM email WHERE user_id=$1 AND (subject ILIKE $2 OR snippet ILIKE $2 OR coalesce(body_text,'') ILIKE $2)
	     ORDER BY sent_at DESC LIMIT 5` + "`" + `, userID, like)

	// Try notes table
	run(` + "`" + `SELECT left(body, 300) FROM note WHERE user_id=$1 AND body ILIKE $2 ORDER BY created_at DESC LIMIT 5` + "`" + `, userID, like)

	// Try contacts table
	run(` + "`" + `SELECT coalesce(first_name,'') || ' ' || coalesce(last_name,'') || ' — ' || coalesce(email,'')
	     FROM contact WHERE user_id=$1 AND (email ILIKE $2 OR first_name ILIKE $2 OR last_name ILIKE $2) LIMIT 3` + "`" + `, userID, like)

	if len(snips) > limit {
		return snips[:limit]
	}
	return snips
}
