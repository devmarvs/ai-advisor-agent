package handlers

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CronTick(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple token check to protect the endpoint
		if c.GetHeader("Authorization") != "Bearer "+os.Getenv("CRON_TOKEN") {
			c.Status(http.StatusUnauthorized)
			return
		}

		// TODO: enqueue polling tasks, etc.
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
