package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// If your project exposes the DB differently, tweak this.
var db *sql.DB

// Inject a db from main or wherever you initialize it.
func InitHubSpotHandlersDatabase(d *sql.DB) {
	db = d
}

// HubSpotStart simply redirects to the assembled authorize URL.
// Useful for /oauth/hubspot/start routes.
func HubSpotStart(c *gin.Context) {
	authURL := os.Getenv("HUBSPOT_AUTH_URL")
	if authURL == "" {
		// Fallback build (same logic used in ConnectPage)
		clientID := os.Getenv("HUBSPOT_CLIENT_ID")
		redirectURI := os.Getenv("HUBSPOT_REDIRECT_URI")
		scopes := os.Getenv("HUBSPOT_SCOPES")
		portalID := os.Getenv("HUBSPOT_PORTAL_ID")
		if clientID == "" || redirectURI == "" || scopes == "" {
			c.String(http.StatusBadRequest, "HubSpot env not configured")
			return
		}
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
		authURL = u.String()
	}
	c.Redirect(http.StatusFound, authURL)
}

// ---- Callback ----

type hubSpotTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func HubSpotCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "missing ?code")
		return
	}

	clientID := os.Getenv("HUBSPOT_CLIENT_ID")
	clientSecret := os.Getenv("HUBSPOT_CLIENT_SECRET")
	redirectURI := os.Getenv("HUBSPOT_REDIRECT_URI")

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		c.String(http.StatusBadRequest, "missing HubSpot env (client id/secret/redirect)")
		return
	}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("redirect_uri", redirectURI)
	form.Set("code", code)

	req, _ := http.NewRequest("POST", "https://api.hubapi.com/oauth/v1/token", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusBadGateway, "token exchange failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		var body bytes.Buffer
		body.ReadFrom(resp.Body)
		c.String(resp.StatusCode, "token exchange error: %s", body.String())
		return
	}

	var t hubSpotTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		c.String(http.StatusBadGateway, "decode token response failed: %v", err)
		return
	}

	// Persist the refresh_token. Adjust this to your app’s user logic.
	if db == nil {
		c.String(http.StatusInternalServerError, "db not initialized for hubspot handlers")
		return
	}

	if err := upsertHubSpotToken(db, t.RefreshToken); err != nil {
		c.String(http.StatusInternalServerError, "saving token failed: %v", err)
		return
	}

	// Success page
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, `<h2>HubSpot connected 🎉</h2><p>You can close this tab.</p>`)
}

func upsertHubSpotToken(db *sql.DB, refresh string) error {
	if strings.TrimSpace(refresh) == "" {
		return errors.New("empty refresh token")
	}

	// Strategy: attach it to the most-recent app_user row (for single-user demos).
	// If you have a session user, replace this with your own lookup.
	type row struct {
		ID string
	}
	var r row
	err := db.QueryRow(`SELECT id FROM app_user ORDER BY created_at DESC LIMIT 1`).Scan(&r.ID)
	if err != nil {
		return fmt.Errorf("load app_user: %w", err)
	}

	_, err = db.Exec(`UPDATE app_user SET hubspot_refresh_token = $1 WHERE id = $2`, refresh, r.ID)
	if err != nil {
		return fmt.Errorf("update token: %w", err)
	}
	_, _ = db.Exec(`UPDATE app_user SET updated_at = $1 WHERE id = $2`, time.Now(), r.ID)
	return nil
}
