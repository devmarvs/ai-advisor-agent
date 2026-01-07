package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"aiagentapi/auth"
	"aiagentapi/handlers"
	"aiagentapi/storage"
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

	if err := storage.ApplyMigrations(db); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	worker.Start(db)

	if err := storage.EnsureSchema(db); err != nil {
		log.Fatalf("failed to ensure schema: %v", err)
	}

	r := gin.Default()
	webRoot := resolveWebRoot()
	chatTemplate := ""
	if webRoot == "" {
		log.Println("web assets not found; chat UI will not render")
	} else {
		chatTemplate = filepath.Join(webRoot, "templates", "chat.html")
		r.Static("/static", filepath.Join(webRoot, "static"))
	}

	// Public
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.GET("/connect", handlers.ConnectPage) // simple page with a Google OAuth button

	// OAuth routes (to add below)
	r.GET("/oauth/google/start", handlers.GoogleStart())
	r.GET("/oauth/google/callback", handlers.GoogleCallback(db))
	r.GET("/logout", func(c *gin.Context) { auth.Logout(c) })

	// Authed
	authed := r.Group("/")
	authed.Use(auth.RequireAuth())
	authed.GET("/", handlers.Home(chatTemplate))
	authed.POST("/chat", handlers.Chat(db))
	authed.GET("/messages", handlers.Messages(db))
	authed.POST("/internal/cron/tick", handlers.CronTick(db))

	return r
}

func resolveWebRoot() string {
	candidates := []string{
		"web",
		filepath.Join("..", "web"),
	}
	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate
		}
	}
	return ""
}
