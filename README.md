# GitHub Profile Reviewer

[![CI](https://github.com/sathvik458/GitHub-Profile-Reviewer/actions/workflows/ci.yml/badge.svg)](https://github.com/sathvik458/GitHub-Profile-Reviewer/actions/workflows/ci.yml)

An API that analyzes a GitHub user's profile and returns useful engineering insights.

## Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/healthz` | Liveness probe. Returns `{"status":"ok"}`. |
| `GET` | `/profile/{username}` | Profile, repositories, and analysis scores. |

`/profile/{username}` returns the user's profile, up to 100 recently updated public
repositories, and analysis scores for documentation, repository quality, activity, an
overall score, and recommendations.

Errors are always JSON:

```json
{
  "error": "github user not found"
}
```

| Status | Meaning |
|---|---|
| `400` | Username missing or malformed |
| `404` | No such GitHub user |
| `405` | Wrong HTTP method (with an `Allow` header) |
| `502` | GitHub itself failed |

## Configuration

| Variable | Default | Description |
|---|---|---|
| `GITHUB_TOKEN` | _(none)_ | Raises the GitHub rate limit from 60 to 5,000 requests/hour. |
| `PORT` | `8080` | Listen port. Hosting platforms set this for you. |

Real environment variables take precedence over `.env`, and a missing `.env` is not an
error, so the same binary works locally, in CI, and in a container.

## Run locally

Create a `.env` file:

```text
GITHUB_TOKEN=your_github_token_here
```

Then:

```bash
go run ./server
curl localhost:8080/healthz
curl localhost:8080/profile/octocat
```

## Run with Docker

The image is a multi-stage build: the Go toolchain compiles a static binary, and the
final stage is distroless with nothing but that binary — roughly 7MB, no shell, running
as a non-root user.

```bash
docker build -t github-profile-reviewer .
docker run --rm -p 8080:8080 -e GITHUB_TOKEN=your_token github-profile-reviewer
```

The server handles `SIGTERM`, so `docker stop` drains in-flight requests instead of
killing them.

## Tests

```bash
go test -race -cover ./...
```

Nothing in the suite touches the network. Three strategies are used, one per layer:

- **`internal/analyzer`** — white-box unit tests over the scoring functions, with the
  clock injected so scores are deterministic.
- **`internal/api`** — a hand-written fake satisfying the `GitHubClient` interface,
  covering every error-to-status branch.
- **`internal/github`** — a custom `http.RoundTripper`, so outbound request headers,
  paths, and query parameters are asserted without a socket.
