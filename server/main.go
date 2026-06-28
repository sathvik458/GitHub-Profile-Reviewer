package main

import (
	"fmt"
	"log"
	"net/http"

	"github-profile-reviewer/internal/api"
	"github-profile-reviewer/internal/config"
	"github-profile-reviewer/internal/github"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.GitHubToken == "" {
		log.Println("GITHUB_TOKEN is not set")
	}

	githubClient := github.NewClient(cfg.GitHubToken)
	handler := api.NewHandler(githubClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/profile/", handler.Profile)

	addr := ":8080"
	log.Printf("server listening on http://localhost%s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from GitHub Profile Reviewer")
}
