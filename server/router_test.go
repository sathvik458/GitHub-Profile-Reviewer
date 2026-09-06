package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github-profile-reviewer/internal/api"
	"github-profile-reviewer/internal/models"
)

// stubGitHubClient satisfies api.GitHubClient with canned data. The router
// tests care about wiring, not about what GitHub returns.
type stubGitHubClient struct{}

func (stubGitHubClient) FetchUser(ctx context.Context, username string) (models.GitHubUser, error) {
	return models.GitHubUser{Login: username}, nil
}

func (stubGitHubClient) FetchRepositories(ctx context.Context, username string) ([]models.Repository, error) {
	return nil, nil
}

func testRouter() http.Handler {
	return newRouter(api.NewHandler(stubGitHubClient{}))
}

func TestRouterServesHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	testRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", contentType, "application/json")
	}

	var body healthBody
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode health response: %v", err)
	}

	if body.Status != "ok" {
		t.Fatalf("status = %q, want %q", body.Status, "ok")
	}
}

func TestRouterRoutesProfileRequests(t *testing.T) {
	rec := httptest.NewRecorder()
	testRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/profile/octocat", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRouterReturnsNotFoundForUnknownPath(t *testing.T) {
	rec := httptest.NewRecorder()
	testRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/not-a-route", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestPortPrefersEnvironment(t *testing.T) {
	t.Setenv("PORT", "9090")

	if got := port(); got != "9090" {
		t.Fatalf("port() = %q, want %q", got, "9090")
	}
}

func TestPortFallsBackToDefault(t *testing.T) {
	t.Setenv("PORT", "")

	if got := port(); got != defaultPort {
		t.Fatalf("port() = %q, want %q", got, defaultPort)
	}
}

type healthBody struct {
	Status string `json:"status"`
}
