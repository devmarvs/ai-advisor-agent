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

	"github.com/gin-gonic/gin"
)

// HubHandlers carries shared deps (e.g., DB) for HubSpot handlers
type HubHandlers struct {
	DB *sql.DB
}

// NewHubHandlers sets up a handler set with DB
func NewHubHandlers(db *sql.DB) *HubHandlers {
	return &HubHandlers{DB: db}
}

// HubSpotStart begins the HubSpot OAuth flow.
// We build the authorization URL from env so we never hard-code scopes.
func (h *HubHandlers) HubSpotStart(c *gin.Context) {
	clientID := os.Getenv("HUBSPOT_CLIENT_ID")
	if clientID == "" {
		c.String(http.StatusInternalServerError, "missing HUBSPOT_CLIENT_ID")
		return
	}

	// HubSpot region base (default global). Set to app-eu1.hubspot.com if you’re in EU.
	authHost := os.Getenv("HUBSPOT_AUTH_HOST")
	if authHost == "" {
		authHost = "app.hubspot.com"
	}

	redirect := strings.TrimRight(os.Getenv("OAUTH_REDIRECT_BASE_URL"), "/") + "/oauth/hubspot/callback"
	if redirect == "" || !strings.HasPrefix(redirect, "http") {
		c.String(http.StatusInternalServerError, "missing or invalid OAUTH_REDIRECT_BASE_URL")
		return
	}

	// Scopes come from env (space- or comma-separated are both accepted).
	scopeEnv := os.Getenv("HUBSPOT_SCOPES")
	if scopeEnv == "" {
		scopeEnv = "oauth crm.objects.contacts.read crm.objects.contacts.write crm.objects.owners.read crm.schemas.contacts.read"
	}
	scopes := strings.Fields(strings.ReplaceAll(scopeEnv, ",", " "))

	v := url.Values{}
	v.Set("client_id", clientID)
	v.Set("redirect_uri", redirect)
	v.Set("response_type", "code")
	v.Set("scope", strings.Join(scopes, " "))
	v.Set("state", "hubspot_oauth") // replace with CSRF token if you have sessions

	authURL := fmt.Sprintf("https://%s/oauth/authorize?%s", authHost, v.Encode())
	c.Redirect(http.StatusFound, authURL)
}

// HubSpotCallback handles the OAuth redirect, exchanges code for tokens, and stores the refresh token in DB.
func (h *HubHandlers) HubSpotCallback(c *gin.Context) {
	if h.DB == nil {
		c.String(http.StatusInternalServerError, "DB not initialized for HubSpot handlers")
		return
	}

	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "missing code")
		return
	}

	clientID := os.Getenv("HUBSPOT_CLIENT_ID")
	clientSecret := os.Getenv("HUBSPOT_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		c.String(http.StatusInternalServerError, "missing HUBSPOT_CLIENT_ID or HUBSPOT_CLIENT_SECRET")
		return
	}

	redirect := strings.TrimRight(os.Getenv("OAUTH_REDIRECT_BASE_URL"), "/") + "/oauth/hubspot/callback"
	tokenURL := "https://api.hubapi.com/oauth/v1/token"

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("redirect_uri", redirect)
	form.Set("code", code)

	req, _ := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		c.String(http.StatusBadGateway, "token exchange error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body)
		c.String(http.StatusBadGateway, "token exchange failed: %s", string(b))
		return
	}

	var payload struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		c.String(http.StatusBadGateway, "decode token response failed: %v", err)
		return
	}

	// Use signed-in user email if you have auth; otherwise demo fallback.
	userEmail := c.GetString("user_email")
	if userEmail == "" {
		userEmail = os.Getenv("DEMO_OWNER_EMAIL")
	}
	if userEmail == "" {
		c.String(http.StatusBadRequest, "no user identity to store HubSpot token (set DEMO_OWNER_EMAIL or wire auth)")
		return
	}

	_, err = h.DB.Exec(`
		INSERT INTO app_user (email, hubspot_refresh_token, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (email)
		DO UPDATE SET hubspot_refresh_token = EXCLUDED.hubspot_refresh_token
	`, userEmail, payload.RefreshToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to save token: %v", err)
		return
	}

	base := strings.TrimRight(os.Getenv("APP_BASE_URL"), "/")
	if base == "" {
		base = "/"
	}
	redirectTo := os.Getenv("POST_CONNECT_REDIRECT")
	if redirectTo == "" {
		redirectTo = "/" // default to your chat UI
	}

	c.Redirect(http.StatusFound, redirectTo)
}
