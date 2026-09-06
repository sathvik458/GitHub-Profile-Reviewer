package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github-profile-reviewer/internal/api"
	"github-profile-reviewer/internal/config"
	"github-profile-reviewer/internal/github"
)

const (
	defaultPort     = "8080"
	shutdownTimeout = 10 * time.Second
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.GitHubToken == "" {
		log.Println("GITHUB_TOKEN is not set: using the unauthenticated GitHub API (60 requests/hour)")
	}

	githubClient := github.NewClient(cfg.GitHubToken)
	handler := api.NewHandler(githubClient)

	srv := &http.Server{
		Addr:    ":" + port(),
		Handler: newRouter(handler),

		// ReadHeaderTimeout is the one that matters for safety: without it a
		// client can hold a connection open by trickling headers (Slowloris).
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	if err := run(srv); err != nil {
		log.Fatal(err)
	}
}

// run starts srv and blocks until the process is signalled, then drains
// in-flight requests before returning. Containers and orchestrators stop a
// process with SIGTERM; without this the runtime waits out its grace period
// and SIGKILLs us, dropping whatever was still being served.
func run(srv *http.Server) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("server listening on http://localhost%s", srv.Addr)

		// Shutdown makes ListenAndServe return ErrServerClosed, which is the
		// expected path out rather than a failure.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}

		serverErr <- nil
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
	}

	log.Println("shutdown signal received, draining connections")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Println("server stopped cleanly")

	return nil
}

// port reads PORT because hosting platforms assign the port rather than
// letting the process choose it.
func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}

	return defaultPort
}
