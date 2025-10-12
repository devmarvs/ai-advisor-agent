package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"aiagentapi/auth"

	"github.com/gin-gonic/gin"
)

// -----------------------------------------------------------------------------
// HubSpot OAuth Start
// -----------------------------------------------------------------------------
func HubSpotStart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Required environment variables
		clientID := os.Getenv("HUBSPOT_CLIENT_ID")
		base := strings.TrimRight(os.Getenv("OAUTH_REDIRECT_BASE_URL"), "/")
		if clientID == "" || base == "" {
			c.String(http.StatusInternalServerError, "Missing HUBSPOT_CLIENT_ID or OAUTH_REDIRECT_BASE_URL")
			return
		}

		// Region-safe auth base (defaults to US; use https://app-eu1.hubspot.com for EU)
		authBase := os.Getenv("HUBSPOT_AUTH_BASE")
		if authBase == "" {
			authBase = "https://app.hubspot.com"
		}

		// Safe scopes — these are valid and supported by HubSpot
		scopes := os.Getenv("HUBSPOT_SCOPES")
		if strings.TrimSpace(scopes) == "" {
			scopes = "crm.objects.contacts.read crm.objects.contacts.write crm.objects.owners.read crm.schemas.contacts.read"
		}
		scopeParam := strings.Join(strings.Fields(scopes), " ")

		redirect := base + "/oauth/hubspot/callback"

		params := url.Values{}
		params.Set("client_id", clientID)
		params.Set("redirect_uri", redirect)
		params.Set("scope", scopeParam)
		params.Set("response_type", "code")

		// Redirect to HubSpot OAuth
		authURL := fmt.Sprintf("%s/oauth/authorize?%s", authBase, params.Encode())
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

// -----------------------------------------------------------------------------
// HubSpot OAuth Callback
// -----------------------------------------------------------------------------
func HubSpotCallback(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.String(http.StatusBadRequest, "Missing code")
			return
		}

		base := strings.TrimRight(os.Getenv("OAUTH_REDIRECT_BASE_URL"), "/")
		clientID := os.Getenv("HUBSPOT_CLIENT_ID")
		clientSecret := os.Getenv("HUBSPOT_CLIENT_SECRET")
		if clientID == "" || clientSecret == "" || base == "" {
			c.String(http.StatusInternalServerError, "Missing HubSpot OAuth env vars")
			return
		}

		redirect := base + "/oauth/hubspot/callback"

		form := url.Values{}
		form.Set("grant_type", "authorization_code")
		form.Set("client_id", clientID)
		form.Set("client_secret", clientSecret)
		form.Set("redirect_uri", redirect)
		form.Set("code", code)

		req, _ := http.NewRequest("POST", "https://api.hubapi.com/oauth/v1/token", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Request error: %v", err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.String(http.StatusBadGateway, fmt.Sprintf("Token exchange failed: %s", string(body)))
			return
		}

		var tok struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int64  `json:"expires_in"`
			TokenType    string `json:"token_type"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Decode error: %v", err))
			return
		}

		if tok.RefreshToken == "" {
			c.String(http.StatusInternalServerError, "No refresh token returned by HubSpot")
			return
		}

		// Example of saving the token for the current user
		user, err := auth.GetCurrentUser(c, db)
		if err != nil {
			c.String(http.StatusUnauthorized, "Not authenticated")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = db.ExecContext(ctx,
			`UPDATE app_user SET hubspot_refresh_token=$1 WHERE id=$2`,
			tok.RefreshToken, user.ID,
		)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("DB save error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `
<!doctype html>
<html>
  <head><meta charset="utf-8"><title>HubSpot connected</title></head>
  <body style="font-family:-apple-system,Segoe UI,Roboto,sans-serif;padding:24px">
    <h1>✅ HubSpot connected!</h1>
    <p>Your refresh token was stored successfully.</p>
  </body>
</html>`)
	}
}
