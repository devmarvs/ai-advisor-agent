package handler

import (
	"net/http"

	app "aiagentapi/app"
)

var router = app.SetupRouter()

// Handler is the Vercel serverless entrypoint.
func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
