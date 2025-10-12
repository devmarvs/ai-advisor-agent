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

func GoogleStart() gin.HandlerFunc {
	return func(c *gin.Context) {
		base := os.Getenv("OAUTH_REDIRECT_BASE_URL")
		redirect := base + "/oauth/google/callback"
		params := url.Values{}
		params.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
		params.Set("redirect_uri", redirect)
		params.Set("response_type", "code")
		params.Set("access_type", "offline")
		params.Set("prompt", "consent")
		params.Set("scope", strings.Join([]string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/gmail.modify",
			"https://www.googleapis.com/auth/calendar",
		}, " "))
		c.Redirect(http.StatusTemporaryRedirect, "https://accounts.google.com/o/oauth2/v2/auth?"+params.Encode())
	}
}

func GoogleCallback(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.String(400, "missing code")
			return
		}
		base := os.Getenv("OAUTH_REDIRECT_BASE_URL")
		redirect := base + "/oauth/google/callback"

		// exchange token
		form := url.Values{}
		form.Set("code", code)
		form.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
		form.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
		form.Set("redirect_uri", redirect)
		form.Set("grant_type", "authorization_code")

		req, _ := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
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
			IdToken      string `json:"id_token"`
			TokenType    string `json:"token_type"`
		}
		json.NewDecoder(resp.Body).Decode(&tok)
		if tok.RefreshToken == "" {
			c.String(500, "no refresh_token; ensure prompt=consent & access_type=offline")
			return
		}

		// get email
		rq, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
		rq.Header.Set("Authorization", "Bearer "+tok.AccessToken)
		r2, err := http.DefaultClient.Do(rq)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		defer r2.Body.Close()
		var ui struct {
			Email string `json:"email"`
		}
		json.NewDecoder(r2.Body).Decode(&ui)

		// upsert app_user
		var userID string
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = db.QueryRowContext(ctx, `
      INSERT INTO app_user(email, google_refresh_token)
      VALUES ($1,$2)
      ON CONFLICT (email) DO UPDATE SET google_refresh_token=EXCLUDED.google_refresh_token
      RETURNING id
    `, ui.Email, tok.RefreshToken).Scan(&userID)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		auth.SetSession(c, userID)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
