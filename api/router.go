
package main

import (
  "database/sql"
  "os"

  "github.com/gin-gonic/gin"
  _ "github.com/jackc/pgx/v5/stdlib"

  "aiagentapi/handlers"
)

func SetupRouter() *gin.Engine {
  dsn := os.Getenv("DATABASE_URL")
  db, _ := sql.Open("pgx", dsn)

  r := gin.Default()
  r.GET("/healthz", func(c *gin.Context){ c.JSON(200, gin.H{"ok":true}) })
  r.POST("/chat", handlers.Chat(db))
  r.POST("/internal/cron/tick", handlers.CronTick(db))
  return r
}
