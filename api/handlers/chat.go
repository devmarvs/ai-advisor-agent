package handlers

import (
  "database/sql"
  "net/http"

  "github.com/gin-gonic/gin"
  "aiagentapi/storage"
)

type ChatRequest struct { ThreadID string `json:"thread_id"`; Message string `json:"message" binding:"required"` }

func Chat(db *sql.DB) gin.HandlerFunc {
  return func(c *gin.Context){
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    _, _ = storage.Enqueue(c, db, "00000000-0000-0000-0000-000000000000", "send_email", map[string]any{"To":"test@example.com","Subject":"Hello","Body":"From scaffold"}, nil, nil)
    c.JSON(200, gin.H{"reply":"Scaffold live. Enqueued a demo task."})
  }
}
