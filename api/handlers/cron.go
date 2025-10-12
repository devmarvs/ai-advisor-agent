package handlers

import (
  "database/sql"
  "net/http"
  "os"
  "github.com/gin-gonic/gin"
)

func CronTick(db *sql.DB) gin.HandlerFunc { return func(c *gin.Context){ if c.GetHeader("Authorization")!="Bearer "+os.Getenv("CRON_TOKEN"){ c.Status(http.StatusUnauthorized); return } c.JSON(http.StatusOK, gin.H{"ok":true}) } }
