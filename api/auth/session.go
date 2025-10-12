package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const SessionCookie = "sid"

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

// small helper
func BaseURL() string { return os.Getenv("APP_BASE_URL") }
