# GitHub Profile Reviewer

An API that will analyze a GitHub user's profile and return useful engineering insights.

## Run

Create a `.env` file:

```text
GITHUB_TOKEN=your_github_token_here
```

Start the server:

```bash
go run ./server
```

Then open:

```text
http://localhost:8080/hello
```

Fetch a GitHub profile:

```text
http://localhost:8080/profile/octocat
```

The response includes the user's profile and up to 100 recently updated public repositories.
