package auth

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const SessionCookie = "sid"

type User struct {
	ID    string
	Email string
}

// Very simple cookie session (production: replace with secure store)
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := c.Cookie(SessionCookie); err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/connect")
			c.Abort()
			return
		}
		c.Next()
	}
}

// helper to set a cookie once OAuth succeeds
func SetSession(c *gin.Context, userID string) {
	// In a real app, map sid->userID in DB/redis. For demo, store userID directly.
	// If you want DB-based, create a session table and store sid->userID.
	c.SetCookie(SessionCookie, userID, 60*60*24*30, "/", "", true, true)
}

func Logout(c *gin.Context) {
	c.SetCookie(SessionCookie, "", -1, "/", "", true, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// GetCurrentUser resolves the user ID stored in the sid cookie against the app_user table.
func GetCurrentUser(c *gin.Context, db *sql.DB) (*User, error) {
	sid, err := c.Cookie(SessionCookie)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var user User
	if err := db.QueryRowContext(ctx, `SELECT id, email FROM app_user WHERE id=$1`, sid).Scan(&user.ID, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}

// small helper
func BaseURL() string { return os.Getenv("APP_BASE_URL") }
