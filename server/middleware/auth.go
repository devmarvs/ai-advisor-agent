package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := GetUserID(c)
		if uid == "" {
			c.Redirect(http.StatusFound, "/connect")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	uid, _ := GetSessionValue(c, "sid")
	return uid
}
