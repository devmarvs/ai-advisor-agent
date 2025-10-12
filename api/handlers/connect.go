package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ConnectPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
      <h1>Connect your accounts</h1>
      <p><a href="/oauth/google/start">Connect Google (Gmail + Calendar)</a></p>
      <p><a href="/oauth/hubspot/start">Connect HubSpot</a></p>
    `))
	}
}
