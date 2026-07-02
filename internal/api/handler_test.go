package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github-profile-reviewer/internal/github"
	"github-profile-reviewer/internal/models"
)

type fakeGitHubClient struct {
	user             models.GitHubUser
	repositories     []models.Repository
	fetchUserErr     error
	fetchReposErr    error
	receivedUsername string
}

func (f *fakeGitHubClient) FetchUser(ctx context.Context, username string) (models.GitHubUser, error) {
	f.receivedUsername = username
	if f.fetchUserErr != nil {
		return models.GitHubUser{}, f.fetchUserErr
	}
	return f.user, nil
}

func (f *fakeGitHubClient) FetchRepositories(ctx context.Context, username string) ([]models.Repository, error) {
	f.receivedUsername = username
	if f.fetchReposErr != nil {
		return nil, f.fetchReposErr
	}
	return f.repositories, nil
}

func TestProfileReturnsProfileResponse(t *testing.T) {
	client := &fakeGitHubClient{
		user: models.GitHubUser{Login: "octocat"},
		repositories: []models.Repository{
			{
				Name:        "hello-world",
				Description: "Example repository",
				LicenseName: "MIT",
				Topics:      []string{"go"},
				PushedAt:    "2026-07-01T00:00:00Z",
			},
		},
	}

	handler := NewHandler(client)
	req := httptest.NewRequest(http.MethodGet, "/profile/octocat", nil)
	rec := httptest.NewRecorder()

	handler.Profile(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var response models.ProfileResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.User.Login != "octocat" {
		t.Fatalf("login = %q, want %q", response.User.Login, "octocat")
	}

	if len(response.Repositories) != 1 {
		t.Fatalf("repositories length = %d, want %d", len(response.Repositories), 1)
	}

	if response.Analysis.OverallScore == 0 {
		t.Fatal("expected analysis overall score")
	}

	if client.receivedUsername != "octocat" {
		t.Fatalf("username = %q, want %q", client.receivedUsername, "octocat")
	}
}

func TestProfileRejectsInvalidMethod(t *testing.T) {
	handler := NewHandler(&fakeGitHubClient{})
	req := httptest.NewRequest(http.MethodPost, "/profile/octocat", nil)
	rec := httptest.NewRecorder()

	handler.Profile(rec, req)

	assertErrorResponse(t, rec, http.StatusMethodNotAllowed, "method not allowed")

	if rec.Header().Get("Allow") != http.MethodGet {
		t.Fatalf("Allow header = %q, want %q", rec.Header().Get("Allow"), http.MethodGet)
	}
}

func TestProfileRequiresUsername(t *testing.T) {
	handler := NewHandler(&fakeGitHubClient{})
	req := httptest.NewRequest(http.MethodGet, "/profile/", nil)
	rec := httptest.NewRecorder()

	handler.Profile(rec, req)

	assertErrorResponse(t, rec, http.StatusBadRequest, "username is required")
}

func TestProfileReturnsNotFoundForMissingGitHubUser(t *testing.T) {
	handler := NewHandler(&fakeGitHubClient{
		fetchUserErr: github.ErrUserNotFound,
	})
	req := httptest.NewRequest(http.MethodGet, "/profile/missing-user", nil)
	rec := httptest.NewRecorder()

	handler.Profile(rec, req)

	assertErrorResponse(t, rec, http.StatusNotFound, "github user not found")
}

func TestProfileReturnsBadGatewayForRepositoryFailure(t *testing.T) {
	handler := NewHandler(&fakeGitHubClient{
		user:          models.GitHubUser{Login: "octocat"},
		fetchReposErr: errors.New("github unavailable"),
	})
	req := httptest.NewRequest(http.MethodGet, "/profile/octocat", nil)
	rec := httptest.NewRecorder()

	handler.Profile(rec, req)

	assertErrorResponse(t, rec, http.StatusBadGateway, "failed to fetch github repositories")
}

func assertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, status int, message string) {
	t.Helper()

	if rec.Code != status {
		t.Fatalf("status = %d, want %d", rec.Code, status)
	}

	var response models.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}

	if response.Error != message {
		t.Fatalf("error = %q, want %q", response.Error, message)
	}
}
