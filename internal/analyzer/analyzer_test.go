package analyzer

import (
	"testing"
	"time"

	"github-profile-reviewer/internal/models"
)

func TestDocumentationScore(t *testing.T) {
	repositories := []models.Repository{
		{
			Description: "API service",
			LicenseName: "MIT",
			Topics:      []string{"go", "api"},
		},
		{
			Description: "Small CLI",
		},
	}

	got := documentationScore(repositories)
	want := 70

	if got != want {
		t.Fatalf("documentationScore() = %d, want %d", got, want)
	}
}

func TestRepositoryScore(t *testing.T) {
	repositories := []models.Repository{
		{
			IsFork:  false,
			Stars:   3,
			HasWiki: true,
		},
		{
			IsFork:   true,
			Archived: true,
		},
	}

	got := repositoryScore(repositories)
	want := 60

	if got != want {
		t.Fatalf("repositoryScore() = %d, want %d", got, want)
	}
}

func TestActivityScore(t *testing.T) {
	now := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	repositories := []models.Repository{
		{PushedAt: now.AddDate(0, 0, -30).Format(time.RFC3339)},
		{PushedAt: now.AddDate(0, 0, -120).Format(time.RFC3339)},
		{PushedAt: now.AddDate(0, 0, -300).Format(time.RFC3339)},
		{PushedAt: now.AddDate(-2, 0, 0).Format(time.RFC3339)},
	}

	got := activityScore(repositories, now)
	want := 55

	if got != want {
		t.Fatalf("activityScore() = %d, want %d", got, want)
	}
}

func TestAnalyzeAddsRecommendations(t *testing.T) {
	analysis := Analyze([]models.Repository{
		{
			Name:     "old-project",
			PushedAt: "2020-01-01T00:00:00Z",
		},
	})

	if analysis.OverallScore == 0 {
		t.Fatal("expected overall score to be calculated")
	}

	if len(analysis.Recommendations) == 0 {
		t.Fatal("expected recommendations")
	}
}
