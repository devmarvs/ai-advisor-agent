package handlers

import (
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// ConnectPage renders a centered, minimal page with links to begin OAuth flows.
func ConnectPage(c *gin.Context) {
	hubspotURL := os.Getenv("HUBSPOT_AUTH_URL")
	if hubspotURL == "" {
		// Build a safe default if env not provided
		clientID := os.Getenv("HUBSPOT_CLIENT_ID")
		redirectURI := os.Getenv("HUBSPOT_REDIRECT_URI")
		scopes := os.Getenv("HUBSPOT_SCOPES")
		portalID := os.Getenv("HUBSPOT_PORTAL_ID")
		if clientID != "" && redirectURI != "" && scopes != "" {
			u, _ := url.Parse("https://app.hubspot.com/oauth/authorize")
			q := u.Query()
			q.Set("client_id", clientID)
			q.Set("redirect_uri", redirectURI)
			q.Set("response_type", "code")
			q.Set("scope", scopes)
			if portalID != "" {
				q.Set("portalId", portalID)
			}
			q.Set("prompt", "consent")
			q.Set("state", "hubspot_oauth")
			u.RawQuery = q.Encode()
			hubspotURL = u.String()
		}
	}

	googleURL := "/oauth/google/start"

	c.Header("Content-Type", "text/html; charset=utf-8")
	if hubspotURL == "" {
		c.String(200, `<!doctype html>
<html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Connect accounts</title>
<style>
body{font-family:Inter,system-ui,-apple-system,Segoe UI,Roboto,sans-serif;background:#0b0c10;color:#e5e7eb;margin:0}
.center{min-height:100vh;display:grid;place-items:center}
.card{padding:28px; border:1px solid #1f2937; background:#111827; border-radius:14px; text-align:center; max-width:560px}
a.btn{display:inline-block;margin:8px 6px;padding:12px 16px;border-radius:10px;background:#0ea5e9;color:#fff;text-decoration:none}
h1{margin:0 0 12px 0}
p{color:#9ca3af}
</style></head>
<body><div class="center">
  <div class="card">
    <h1>Connect your accounts</h1>
    <p><strong>HubSpot</strong> not configured. Set HUBSPOT_AUTH_URL (or HUBSPOT_CLIENT_ID, HUBSPOT_REDIRECT_URI, HUBSPOT_SCOPES) and refresh.</p>
    <p><a class="btn" href="`+googleURL+`">Connect Google (Gmail + Calendar)</a></p>
  </div>
</div></body></html>`)
		return
	}

	c.String(200, `<!doctype html>
<html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Connect accounts</title>
<style>
body{font-family:Inter,system-ui,-apple-system,Segoe UI,Roboto,sans-serif;background:#0b0c10;color:#e5e7eb;margin:0}
.center{min-height:100vh;display:grid;place-items:center}
.card{padding:28px; border:1px solid #1f2937; background:#111827; border-radius:14px; text-align:center; max-width:560px}
a.btn{display:inline-block;margin:8px 6px;padding:12px 16px;border-radius:10px;background:#0ea5e9;color:#fff;text-decoration:none}
h1{margin:0 0 12px 0}
</style></head>
<body><div class="center">
  <div class="card">
    <h1>Connect your accounts</h1>
    <p><a class="btn" href="`+hubspotURL+`">Connect HubSpot</a></p>
    <p><a class="btn" href="`+googleURL+`">Connect Google (Gmail + Calendar)</a></p>
  </div>
</div></body></html>`)
}
