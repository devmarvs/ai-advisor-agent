package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"aiagentapi/auth"

	"github.com/gin-gonic/gin"
)

func HubSpotStart() gin.HandlerFunc {
	return func(c *gin.Context) {
		base := os.Getenv("OAUTH_REDIRECT_BASE_URL")
		redirect := base + "/oauth/hubspot/callback"
		params := url.Values{}
		params.Set("client_id", os.Getenv("HUBSPOT_CLIENT_ID"))
		params.Set("redirect_uri", redirect)
		params.Set("scope", "crm.objects.contacts.read crm.objects.contacts.write crm.objects.owners.read crm.schemas.contacts.read crm.objects.contacts.settings.read crm.objects.deals.read crm.objects.deals.write crm.lists.read crm.objects.notes.read crm.objects.notes.write")
		params.Set("response_type", "code")
		c.Redirect(http.StatusTemporaryRedirect, "https://app.hubspot.com/oauth/authorize?"+params.Encode())
	}
}

func HubSpotCallback(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.String(400, "missing code")
			return
		}

		base := os.Getenv("OAUTH_REDIRECT_BASE_URL")
		redirect := base + "/oauth/hubspot/callback"

		form := url.Values{}
		form.Set("grant_type", "authorization_code")
		form.Set("client_id", os.Getenv("HUBSPOT_CLIENT_ID"))
		form.Set("client_secret", os.Getenv("HUBSPOT_CLIENT_SECRET"))
		form.Set("redirect_uri", redirect)
		form.Set("code", code)

		req, _ := http.NewRequest("POST", "https://api.hubapi.com/oauth/v1/token", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		defer resp.Body.Close()

		var tok struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int64  `json:"expires_in"`
			TokenType    string `json:"token_type"`
		}
		json.NewDecoder(resp.Body).Decode(&tok)
		if tok.RefreshToken == "" {
			c.String(500, "no refresh token from HubSpot")
			return
		}

		// get current user’s email from session (demo assumption); in prod, map cookie->userId in DB
		sid, _ := c.Cookie(auth.SessionCookie)
		userID := sid
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = db.ExecContext(ctx, `UPDATE app_user SET hubspot_refresh_token=$2 WHERE id=$1`, userID, tok.RefreshToken)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
