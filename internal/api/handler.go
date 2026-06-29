package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github-profile-reviewer/internal/github"
	"github-profile-reviewer/internal/models"
)

type Handler struct {
	githubClient *github.Client
}

func NewHandler(githubClient *github.Client) *Handler {
	return &Handler{
		githubClient: githubClient,
	}
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/profile/")
	username = strings.TrimSpace(username)

	if username == "" || strings.Contains(username, "/") {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	user, err := h.githubClient.FetchUser(r.Context(), username)
	if err != nil {
		if errors.Is(err, github.ErrUserNotFound) {
			http.Error(w, "github user not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to fetch github profile", http.StatusBadGateway)
		return
	}

	repositories, err := h.githubClient.FetchRepositories(r.Context(), username)
	if err != nil {
		http.Error(w, "failed to fetch github repositories", http.StatusBadGateway)
		return
	}

	response := models.ProfileResponse{
		User:         user,
		Repositories: repositories,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
