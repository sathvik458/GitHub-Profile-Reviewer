package github

import "github-profile-reviewer/internal/models"

type githubRepository struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Language      string   `json:"language"`
	Stars         int      `json:"stargazers_count"`
	Forks         int      `json:"forks_count"`
	OpenIssues    int      `json:"open_issues_count"`
	IsFork        bool     `json:"fork"`
	DefaultBranch string   `json:"default_branch"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	PushedAt      string   `json:"pushed_at"`
	HTMLURL       string   `json:"html_url"`
	HasWiki       bool     `json:"has_wiki"`
	HasProjects   bool     `json:"has_projects"`
	Archived      bool     `json:"archived"`
	Topics        []string `json:"topics"`
	License       *struct {
		Name string `json:"name"`
	} `json:"license"`
}

func (r githubRepository) toModel() models.Repository {
	repo := models.Repository{
		Name:          r.Name,
		Description:   r.Description,
		Language:      r.Language,
		Stars:         r.Stars,
		Forks:         r.Forks,
		OpenIssues:    r.OpenIssues,
		IsFork:        r.IsFork,
		DefaultBranch: r.DefaultBranch,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
		PushedAt:      r.PushedAt,
		HTMLURL:       r.HTMLURL,
		HasWiki:       r.HasWiki,
		HasProjects:   r.HasProjects,
		Archived:      r.Archived,
		Topics:        r.Topics,
	}

	if r.License != nil {
		repo.LicenseName = r.License.Name
	}

	return repo
}
