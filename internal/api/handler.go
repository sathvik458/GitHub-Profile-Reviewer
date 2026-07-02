package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github-profile-reviewer/internal/analyzer"
	"github-profile-reviewer/internal/github"
	"github-profile-reviewer/internal/models"
)

type GitHubClient interface {
	FetchUser(ctx context.Context, username string) (models.GitHubUser, error)
	FetchRepositories(ctx context.Context, username string) ([]models.Repository, error)
}

type Handler struct {
	githubClient GitHubClient
}

func NewHandler(githubClient GitHubClient) *Handler {
	return &Handler{
		githubClient: githubClient,
	}
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/profile/")
	username = strings.TrimSpace(username)

	if username == "" || strings.Contains(username, "/") {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}

	user, err := h.githubClient.FetchUser(r.Context(), username)
	if err != nil {
		if errors.Is(err, github.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "github user not found")
			return
		}

		log.Printf("failed to fetch github profile for %q: %v", username, err)
		writeError(w, http.StatusBadGateway, "failed to fetch github profile")
		return
	}

	repositories, err := h.githubClient.FetchRepositories(r.Context(), username)
	if err != nil {
		log.Printf("failed to fetch github repositories for %q: %v", username, err)
		writeError(w, http.StatusBadGateway, "failed to fetch github repositories")
		return
	}

	response := models.ProfileResponse{
		User:         user,
		Repositories: repositories,
		Analysis:     analyzer.Analyze(repositories),
	}

	writeJSON(w, http.StatusOK, response)
}
