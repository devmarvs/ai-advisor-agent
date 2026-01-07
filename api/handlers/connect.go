package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ConnectPage renders a centered, minimal page with a link to begin Google OAuth.
func ConnectPage(c *gin.Context) {
	googleURL := "/oauth/google/start"

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, connectPageHTML(googleURL))
}

func connectPageHTML(googleURL string) string {
	const header = `<!doctype html>
<html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Connect accounts</title>
<style>
body{font-family:Inter,system-ui,-apple-system,Segoe UI,Roboto,sans-serif;background:#0b0c10;color:#e5e7eb;margin:0}
.center{min-height:100vh;display:grid;place-items:center}
.card{padding:28px; border:1px solid #1f2937; background:#111827; border-radius:14px; text-align:center; max-width:560px}
a.btn{display:inline-block;margin:8px 6px;padding:12px 16px;border-radius:10px;background:#0ea5e9;color:#fff;text-decoration:none}
h1{margin:0 0 12px 0}
p{color:#9ca3af}
</style></head>
<body><div class="center">
  <div class="card">
    <h1>Connect your accounts</h1>
`
	const footer = `  </div>
</div></body></html>`

	var b strings.Builder
	b.WriteString(header)
	b.WriteString(`    <p><a class="btn" href="`)
	b.WriteString(googleURL)
	b.WriteString(`">Connect Google (Gmail + Calendar)</a></p>` + "\n")
	b.WriteString(footer)
	return b.String()
}
