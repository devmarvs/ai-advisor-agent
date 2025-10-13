package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"aiagentapi/auth"
	"aiagentapi/handlers"
	"aiagentapi/worker"
)

func SetupRouter() *gin.Engine {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(15)
	db.SetConnMaxIdleTime(5 * time.Minute)

	worker.Start(db)

	r := gin.Default()

	// Public
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.GET("/connect", handlers.ConnectPage) // simple page with “Connect Google / Connect HubSpot” buttons

	hub := handlers.NewHubHandlers(db)

	// OAuth routes (to add below)
	r.GET("/oauth/google/start", handlers.GoogleStart())
	r.GET("/oauth/google/callback", handlers.GoogleCallback(db))
	r.GET("/oauth/hubspot/start", hub.HubSpotStart)
	r.GET("/oauth/hubspot/callback", hub.HubSpotCallback)
	r.GET("/logout", func(c *gin.Context) { auth.Logout(c) })

	// Authed
	authed := r.Group("/")
	authed.Use(auth.RequireAuth())
	authed.GET("/", handlers.Home())
	authed.POST("/chat", handlers.Chat(db))
	authed.GET("/messages", handlers.Messages())
	authed.POST("/internal/cron/tick", handlers.CronTick(db))

	return r
}
