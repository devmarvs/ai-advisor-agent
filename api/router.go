package main

import (
	"ai-agent/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// OAuth
	r.GET("/oauth/google/start", handlers.GoogleStart)
	r.GET("/oauth/google/callback", handlers.GoogleCallback)
	r.GET("/oauth/hubspot/start", handlers.HubSpotStart)
	r.GET("/oauth/hubspot/callback", handlers.HubSpotCallback)

	// Chat & ingestion
	r.POST("/chat", handlers.Chat)
	r.POST("/ingest/gmail", handlers.ManualGmailIngest)

	// Webhooks
	r.POST("/webhooks/hubspot", handlers.HubSpotWebhook)
	r.POST("/webhooks/calendar", handlers.CalendarWebhook)

	return r
}
