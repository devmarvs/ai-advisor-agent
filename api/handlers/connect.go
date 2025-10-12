package handlers

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// Gin-compatible handler
func ConnectPage(c *gin.Context) {
	clientID := os.Getenv("HUBSPOT_CLIENT_ID")
	redirectURI := os.Getenv("HUBSPOT_REDIRECT_URI")
	scopes := os.Getenv("HUBSPOT_SCOPES")
	authBase := os.Getenv("HUBSPOT_AUTH_BASE") // e.g., https://app-eu1.hubspot.com/oauth/authorize

	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", scopes)
	q.Set("state", "connect-"+c.ClientIP())

	authURL := fmt.Sprintf("%s?%s", authBase, q.Encode())

	html := fmt.Sprintf(`
		<h1>Connect your accounts</h1>
		<p><a href="%s">Connect HubSpot</a></p>
		<p><a href="/oauth/google/start">Connect Google (Gmail + Calendar)</a></p>
	`, authURL)

	c.Data(200, "text/html; charset=utf-8", []byte(html))
}
