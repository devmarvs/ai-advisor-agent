package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	ThreadID string `json:"thread_id"`
	Message  string `json:"message" binding:"required"`
}

func Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 1) RAG search across email/note/instruction embeddings
	// 2) LLM call with toolcapabilities
	// 3) Stream tokens via SSE (or return once for now)

	c.JSON(http.StatusOK, gin.H{"reply": "Scaffold is live. Connect Gmail + HubSpot next."})
}
