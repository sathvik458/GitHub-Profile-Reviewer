package analyzer

import (
	"time"

	"github-profile-reviewer/internal/models"
)

func Analyze(repositories []models.Repository) models.Analysis {
	return models.Analysis{
		DocumentationScore: documentationScore(repositories),
		RepositoryScore:    repositoryScore(repositories),
		ActivityScore:      activityScore(repositories, time.Now()),
	}
}

func documentationScore(repositories []models.Repository) int {
	if len(repositories) == 0 {
		return 0
	}

	total := 0
	for _, repo := range repositories {
		score := 0

		if repo.Description != "" {
			score += 40
		}

		if repo.LicenseName != "" {
			score += 30
		}

		if len(repo.Topics) > 0 {
			score += 30
		}

		total += score
	}

	return total / len(repositories)
}

func repositoryScore(repositories []models.Repository) int {
	if len(repositories) == 0 {
		return 0
	}

	total := 0
	for _, repo := range repositories {
		score := 50

		if !repo.IsFork {
			score += 20
		}

		if repo.Stars > 0 {
			score += 10
		}

		if repo.HasWiki || repo.HasProjects {
			score += 10
		}

		if repo.Archived {
			score -= 20
		}

		total += clamp(score)
	}

	return total / len(repositories)
}

func activityScore(repositories []models.Repository, now time.Time) int {
	if len(repositories) == 0 {
		return 0
	}

	total := 0
	for _, repo := range repositories {
		pushedAt, err := time.Parse(time.RFC3339, repo.PushedAt)
		if err != nil {
			continue
		}

		age := now.Sub(pushedAt)

		switch {
		case age <= 90*24*time.Hour:
			total += 100
		case age <= 180*24*time.Hour:
			total += 70
		case age <= 365*24*time.Hour:
			total += 40
		default:
			total += 10
		}
	}

	return total / len(repositories)
}

func clamp(score int) int {
	if score < 0 {
		return 0
	}

	if score > 100 {
		return 100
	}

	return score
}
