package main

import (
	"net/http"

	"github-profile-reviewer/internal/api"
)

// newRouter builds the route table. Keeping it out of main is what makes the
// wiring testable: main itself can never be exercised by a test.
//
// The static file handler sits last so that /healthz and /profile/ always take
// priority over any file that happens to share that name in web/.
func newRouter(handler *api.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", api.Health)
	mux.HandleFunc("/profile/", handler.Profile)

	// Serve the browser UI from the web/ directory.
	// http.FileServer sets correct Content-Type headers and handles
	// ETags/Last-Modified so the browser can cache assets properly.
	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/", fs)

	return mux
}
