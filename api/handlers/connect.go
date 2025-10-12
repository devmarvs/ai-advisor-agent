package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
)

func ConnectPage(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("HUBSPOT_CLIENT_ID")
	redirectURI := os.Getenv("HUBSPOT_REDIRECT_URI")
	scopes := os.Getenv("HUBSPOT_SCOPES")
	authBase := os.Getenv("HUBSPOT_AUTH_BASE")

	// Build query parameters dynamically
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", scopes)
	q.Set("state", "connect-"+r.RemoteAddr) // small CSRF token

	authURL := fmt.Sprintf("%s?%s", authBase, q.Encode())

	// Simple HTML for the connect page
	tmpl := template.Must(template.New("connect").Parse(`
		<h1>Connect Your Accounts</h1>
		<p><a href="{{.HubSpotURL}}">Connect HubSpot</a></p>
		<p><a href="/oauth/google/start">Connect Google (Gmail + Calendar)</a></p>
	`))

	_ = tmpl.Execute(w, map[string]interface{}{
		"HubSpotURL": authURL,
	})
}
