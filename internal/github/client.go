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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return models.GitHubUser{}, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "github-profile-reviewer")

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.GitHubUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return models.GitHubUser{}, ErrUserNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return models.GitHubUser{}, fmt.Errorf("github returned status %d", resp.StatusCode)
	}

	var user models.GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return models.GitHubUser{}, err
	}

	return user, nil
}
