package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github-profile-reviewer/internal/models"
)

var ErrUserNotFound = errors.New("github user not found")

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func NewClient(token string) *Client {
	return &Client{
		baseURL: "https://api.github.com",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		token: token,
	}
}

func (c *Client) FetchUser(ctx context.Context, username string) (models.GitHubUser, error) {
	apiURL := c.baseURL + "/users/" + url.PathEscape(username)

	var user models.GitHubUser
	if err := c.get(ctx, apiURL, &user); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return models.GitHubUser{}, ErrUserNotFound
		}

		return models.GitHubUser{}, err
	}

	return user, nil
}

func (c *Client) FetchRepositories(ctx context.Context, username string) ([]models.Repository, error) {
	apiURL := c.baseURL + "/users/" + url.PathEscape(username) + "/repos?sort=updated&per_page=100"

	var repos []githubRepository
	if err := c.get(ctx, apiURL, &repos); err != nil {
		return nil, err
	}

	repositories := make([]models.Repository, 0, len(repos))
	for _, repo := range repos {
		repositories = append(repositories, repo.toModel())
	}

	return repositories, nil
}

func (c *Client) get(ctx context.Context, apiURL string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "github-profile-reviewer")

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrUserNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("github returned status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
