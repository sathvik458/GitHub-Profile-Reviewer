package main

import (
	"net/http"

	"github-profile-reviewer/internal/api"
)

// newRouter builds the route table. Keeping it out of main is what makes the
// wiring testable: main itself can never be exercised by a test.
func newRouter(handler *api.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", api.Health)
	mux.HandleFunc("/profile/", handler.Profile)

	return mux
}
