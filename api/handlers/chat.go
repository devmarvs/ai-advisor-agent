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
	openai "github.com/sashabaranov/go-openai"

	"aiagentapi/auth"
	"aiagentapi/storage"
)

type chatMessage struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id"`
}

// Chat handles POST /chat: saves user message, produces assistant reply, and saves it.
func Chat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := auth.GetCurrentUser(c, db)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}
		userID := user.ID

		var req struct {
			Message string `json:"message"`
		}
		if err := c.BindJSON(&req); err != nil || strings.TrimSpace(req.Message) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "message required"})
			return
		}

		ctx := c.Request.Context()

		// Save the user's message
		if _, err := storage.SaveMessage(ctx, db, "user", req.Message); err != nil {
			// Return the detailed cause to Render logs (and UI JSON for debugging)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "failed to save message",
				"detail": err.Error(),
			})
			return
		}

		// RAG-lite: try to pull snippets
		snips := findSnippets(ctx, db, userID, req.Message, 6)

		sys := "You are an assistant for a financial advisor.\n" +
			"Use the provided context snippets when relevant.\n" +
			"If a user asks to act (email, schedule, log note), describe the next steps you will perform."

		userPrompt := buildUserPrompt(req.Message, snips)
		reply := callLLM(ctx, sys, userPrompt)

		// Save assistant message
		if _, err := storage.SaveMessage(ctx, db, "assistant", reply); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"reply":    reply,
				"snippets": snips,
				"warning":  "failed to save assistant message",
				"detail":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"reply": reply, "snippets": snips})
	}
}

// Messages (grouped) handles GET /messages and returns groups by day for History tab.
func Messages(db *sql.DB) gin.HandlerFunc {
	type item struct {
		Role      string    `json:"role"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}
	type group struct {
		Date  string `json:"date"` // YYYY-MM-DD (user local time not applied here; UTC date)
		Items []item `json:"items"`
	}

	return func(c *gin.Context) {
		user, err := auth.GetCurrentUser(c, db)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		ctx := c.Request.Context()
		msgs, err := storage.ListRecentMessages(ctx, db, 20)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to load history"})
			return
		}

		// Group by UTC date (you can adapt to user tz if needed)
		groups := make([]group, 0, 8)
		var cur group
		var lastDate string

		for _, m := range msgs {
			d := m.CreatedAt.UTC().Format("2006-01-02")
			if d != lastDate {
				if lastDate != "" {
					groups = append(groups, cur)
				}
				cur = group{Date: d, Items: []item{}}
				lastDate = d
			}
			cur.Items = append(cur.Items, item{
				Role:      m.Role,
				Content:   m.Content,
				CreatedAt: m.CreatedAt,
			})
		}
		if lastDate != "" {
			groups = append(groups, cur)
		}

		c.JSON(200, gin.H{"groups": groups})
	}
}

func buildUserPrompt(question string, snips []string) string {
	var b strings.Builder
	if len(snips) > 0 {
		b.WriteString("Context snippets:\n")
		for i, s := range snips {
			if i >= 6 {
				break
			}
			b.WriteString("- ")
			if len(s) > 400 {
				s = s[:400] + "…"
			}
			b.WriteString(s)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("User question: ")
	b.WriteString(question)
	return b.String()
}

func callLLM(ctx context.Context, system, user string) string {
	key := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

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

func findSnippets(ctx context.Context, db *sql.DB, userID, q string, limit int) []string {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil
	}

	type row struct{ S string }
	snips := make([]string, 0, limit)

	run := func(sqlStr string, args ...any) {
		if len(snips) >= limit {
			return
		}
		ctx2, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		rows, err := db.QueryContext(ctx2, sqlStr, args...)
		if err != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			if len(snips) >= limit {
				break
			}
			var r row
			if err := rows.Scan(&r.S); err == nil && strings.TrimSpace(r.S) != "" {
				snips = append(snips, r.S)
			}
		}
	}

	like := "%" + q + "%"

	run(`SELECT subject || ' — ' || left(coalesce(body_text,snippet,''), 300)
	     FROM email WHERE user_id=$1 AND (subject ILIKE $2 OR snippet ILIKE $2 OR coalesce(body_text,'') ILIKE $2)
	     ORDER BY sent_at DESC LIMIT 5`, userID, like)

	run(`SELECT left(body, 300) FROM note WHERE user_id=$1 AND body ILIKE $2 ORDER BY created_at DESC LIMIT 5`, userID, like)

	run(`SELECT coalesce(first_name,'') || ' ' || coalesce(last_name,'') || ' — ' || coalesce(email,'')
	     FROM contact WHERE user_id=$1 AND (email ILIKE $2 OR first_name ILIKE $2 OR last_name ILIKE $2) LIMIT 3`, userID, like)

	if len(snips) > limit {
		return snips[:limit]
	}
	return snips
}
