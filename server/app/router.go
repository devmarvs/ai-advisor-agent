package app

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"aiagentapi/auth"
	"aiagentapi/handlers"
	"aiagentapi/storage"
	"aiagentapi/worker"
)

func SetupRouter() *gin.Engine {
	dsn, err := resolveDatabaseURL()
	if err != nil {
		log.Fatal(err)
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

func resolveDatabaseURL() (string, error) {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn != "" {
		return dsn, nil
	}

	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	name := strings.TrimSpace(os.Getenv("DB_NAME"))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	password := os.Getenv("DB_PASSWORD")
	if host == "" || name == "" || user == "" || password == "" {
		return "", fmt.Errorf("DATABASE_URL or DB_HOST/DB_NAME/DB_USER/DB_PASSWORD must be set")
	}

	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	if port == "" {
		port = "5432"
	}

	sslmode := strings.TrimSpace(os.Getenv("DB_SSLMODE"))
	if sslmode == "" {
		sslmode = "require"
	}

	q := url.Values{}
	if sslmode != "" {
		q.Set("sslmode", sslmode)
	}
	if channelBinding := strings.TrimSpace(os.Getenv("DB_CHANNEL_BINDING")); channelBinding != "" {
		q.Set("channel_binding", channelBinding)
	}

	u := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(user, password),
		Host:     net.JoinHostPort(host, port),
		Path:     "/" + name,
		RawQuery: q.Encode(),
	}
	return u.String(), nil
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
