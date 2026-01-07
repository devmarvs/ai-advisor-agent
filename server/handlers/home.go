package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home serves the chat UI with History and New Thread actions.
func Home(templatePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if templatePath == "" {
			c.String(http.StatusInternalServerError, "chat UI template not configured")
			return
		}
		c.File(templatePath)
	}
}
