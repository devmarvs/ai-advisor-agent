package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetSessionValue(c *gin.Context, name, value string) {
	secret := os.Getenv("SESSION_KEY")
	if secret == "" {
		secret = "dev-session-key-change-me"
	}
	data := value
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	sig := mac.Sum(nil)
	token := base64.RawURLEncoding.EncodeToString([]byte(data)) + "." + base64.RawURLEncoding.EncodeToString(sig)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetSessionValue(c *gin.Context, name string) (string, bool) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", false
	}
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return "", false
	}
	dataB, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", false
	}
	sigB, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", false
	}

	secret := os.Getenv("SESSION_KEY")
	if secret == "" {
		secret = "dev-session-key-change-me"
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(dataB))
	if !hmac.Equal(sigB, mac.Sum(nil)) {
		return "", false
	}
	return string(dataB), true
}
