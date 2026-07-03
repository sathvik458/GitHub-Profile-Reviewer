package github

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestFetchUserSendsHeadersAndDecodesResponse(t *testing.T) {
	client := testClient("test-token", func(r *http.Request) *http.Response {
		if r.URL.Path != "/users/octocat" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/users/octocat")
		}

		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatalf("Authorization header = %q", r.Header.Get("Authorization"))
		}

		if r.Header.Get("User-Agent") != "github-profile-reviewer" {
			t.Fatalf("User-Agent header = %q", r.Header.Get("User-Agent"))
		}

		return jsonResponse(http.StatusOK, `{"login":"octocat","name":"The Octocat","public_repos":8}`)
	})

	user, err := client.FetchUser(context.Background(), "octocat")
	if err != nil {
		t.Fatalf("FetchUser() error = %v", err)
	}

	if user.Login != "octocat" {
		t.Fatalf("Login = %q, want %q", user.Login, "octocat")
	}

	if user.PublicRepos != 8 {
		t.Fatalf("PublicRepos = %d, want %d", user.PublicRepos, 8)
	}
}

func TestFetchUserReturnsNotFound(t *testing.T) {
	client := testClient("", func(r *http.Request) *http.Response {
		return jsonResponse(http.StatusNotFound, "")
	})

	_, err := client.FetchUser(context.Background(), "missing-user")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("error = %v, want ErrUserNotFound", err)
	}
}

func TestFetchRepositoriesDecodesAndMapsResponse(t *testing.T) {
	client := testClient("", func(r *http.Request) *http.Response {
		if r.URL.Path != "/users/octocat/repos" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/users/octocat/repos")
		}

		if r.URL.Query().Get("sort") != "updated" {
			t.Fatalf("sort query = %q, want %q", r.URL.Query().Get("sort"), "updated")
		}

		return jsonResponse(http.StatusOK, `[{
			"name":"github-profile-reviewer",
			"description":"Profile API",
			"language":"Go",
			"stargazers_count":5,
			"forks_count":2,
			"open_issues_count":1,
			"fork":false,
			"default_branch":"main",
			"license":{"name":"MIT"},
			"topics":["go","api"]
		}]`)
	})

	repos, err := client.FetchRepositories(context.Background(), "octocat")
	if err != nil {
		t.Fatalf("FetchRepositories() error = %v", err)
	}

	if len(repos) != 1 {
		t.Fatalf("len(repos) = %d, want %d", len(repos), 1)
	}

	if repos[0].Stars != 5 {
		t.Fatalf("Stars = %d, want %d", repos[0].Stars, 5)
	}

	if repos[0].LicenseName != "MIT" {
		t.Fatalf("LicenseName = %q, want %q", repos[0].LicenseName, "MIT")
	}
}

func TestGitHubStatusErrorIncludesMessage(t *testing.T) {
	client := testClient("", func(r *http.Request) *http.Response {
		return jsonResponse(http.StatusTooManyRequests, `{"message":"rate limit exceeded"}`)
	})

	_, err := client.FetchUser(context.Background(), "octocat")
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "429: rate limit exceeded") {
		t.Fatalf("error = %q, want status and message", err.Error())
	}
}

type roundTripFunc func(r *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r), nil
}

func testClient(token string, transport roundTripFunc) *Client {
	client := NewClient(token)
	client.baseURL = "https://api.test"
	client.httpClient = &http.Client{
		Transport: transport,
	}

	return client
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
