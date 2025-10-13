package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"aiagentapi/auth"
	"aiagentapi/storage"
)

type chatMessage struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id"`
}

type chatRequest struct {
	Message string `json:"message" binding:"required"`
	Thread  string `json:"thread_id"`
}

var chatHistory = make([]chatMessage, 0, 512)

// Chat handles POST /chat requests and returns a basic acknowledgement response.
// It also keeps an in-memory history per user for short-lived demo sessions.
func Chat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req chatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if db == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database not configured"})
			return
		}

		user, err := auth.GetCurrentUser(c, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		userID := user.ID
		now := time.Now()

		// Persist a demo task so the worker pipeline can be exercised locally.
		if db != nil {
			_, _ = storage.Enqueue(
				c.Request.Context(),
				db,
				userID,
				"send_email",
				map[string]any{
					"To":      "test@example.com",
					"Subject": "Hello from chat handler",
					"Body":    req.Message,
					"Thread":  req.Thread,
				},
				nil,
				nil,
			)
		}

		chatHistory = append(chatHistory, chatMessage{
			Role:      "user",
			Content:   req.Message,
			CreatedAt: now,
			UserID:    userID,
		})

		reply := fmt.Sprintf("You said: %s", req.Message)
		chatHistory = append(chatHistory, chatMessage{
			Role:      "assistant",
			Content:   reply,
			CreatedAt: now,
			UserID:    userID,
		})

		c.JSON(http.StatusOK, gin.H{"reply": reply})
	}
}
