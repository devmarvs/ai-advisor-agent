package main

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"aiagentapi/auth"
	"aiagentapi/handlers"
	"aiagentapi/worker"
)

func SetupRouter() *gin.Engine {
	dsn := os.Getenv("DATABASE_URL")
	db, _ := sql.Open("pgx", dsn)

	worker.Start(db)

	r := gin.Default()

	// Public
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.GET("/connect", handlers.ConnectPage()) // simple page with “Connect Google / Connect HubSpot” buttons

	// OAuth routes (to add below)
	r.GET("/oauth/google/start", handlers.GoogleStart())
	r.GET("/oauth/google/callback", handlers.GoogleCallback(db))
	r.GET("/oauth/hubspot/start", handlers.HubSpotStart())
	r.GET("/oauth/hubspot/callback", handlers.HubSpotCallback(db))
	r.GET("/logout", func(c *gin.Context) { auth.Logout(c) })

	// Authed
	authed := r.Group("/")
	authed.Use(auth.RequireAuth())
	authed.GET("/", handlers.Home())
	authed.POST("/chat", handlers.Chat(db))
	authed.POST("/internal/cron/tick", handlers.CronTick(db))

	return r
}
