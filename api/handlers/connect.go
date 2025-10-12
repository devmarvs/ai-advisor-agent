package handlers

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// ConnectPage renders a super simple HTML page with links to begin OAuth flows.
// The HubSpot link is sourced from HUBSPOT_AUTH_URL if present, otherwise we
// assemble a correct authorize URL from the individual env vars.
func ConnectPage(c *gin.Context) {
	hubspotURL := os.Getenv("HUBSPOT_AUTH_URL")
	if hubspotURL == "" {
		// Build a safe default using the standard authorize endpoint.
		clientID := os.Getenv("HUBSPOT_CLIENT_ID")
		redirectURI := os.Getenv("HUBSPOT_REDIRECT_URI")
		scopes := os.Getenv("HUBSPOT_SCOPES") // space-separated list
		portalID := os.Getenv("HUBSPOT_PORTAL_ID")

		if clientID != "" && redirectURI != "" && scopes != "" {
			u, _ := url.Parse("https://app.hubspot.com/oauth/authorize")
			q := u.Query()
			q.Set("client_id", clientID)
			q.Set("redirect_uri", redirectURI)
			q.Set("response_type", "code")
			// HubSpot accepts scopes as a space-separated string in the query.
			q.Set("scope", scopes)
			// Optional but recommended: pin to a portal to skip the chooser and stay in OAuth
			if portalID != "" {
				q.Set("portalId", portalID)
			}
			// Force consent screen if already installed previously
			q.Set("prompt", "consent")
			// Optional state
			q.Set("state", "hubspot_oauth")
			u.RawQuery = q.Encode()
			hubspotURL = u.String()
		}
	}

	googleURL := "/oauth/google/start" // whatever you already have wired for Google

	c.Header("Content-Type", "text/html; charset=utf-8")

	if hubspotURL == "" {
		// If we could not build a URL, show a helpful message.
		c.String(200, `
			<h1>Connect your accounts</h1>
			<p><strong>HubSpot</strong>: Missing env. Set HUBSPOT_AUTH_URL (or HUBSPOT_CLIENT_ID, HUBSPOT_REDIRECT_URI, HUBSPOT_SCOPES) and refresh.</p>
			<p><a href="%s">Connect Google (Gmail + Calendar)</a></p>
		`, googleURL)
		return
	}

	c.String(200, fmt.Sprintf(`
		<h1>Connect your accounts</h1>
		<p><a href="%s">Connect HubSpot</a></p>
		<p><a href="%s">Connect Google (Gmail + Calendar)</a></p>
	`, hubspotURL, googleURL))
}
