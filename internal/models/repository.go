package models

type Repository struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Language      string   `json:"language"`
	Stars         int      `json:"stars"`
	Forks         int      `json:"forks"`
	OpenIssues    int      `json:"open_issues"`
	IsFork        bool     `json:"is_fork"`
	DefaultBranch string   `json:"default_branch"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	PushedAt      string   `json:"pushed_at"`
	HTMLURL       string   `json:"html_url"`
	HasWiki       bool     `json:"has_wiki"`
	HasProjects   bool     `json:"has_projects"`
	Archived      bool     `json:"archived"`
	LicenseName   string   `json:"license_name"`
	Topics        []string `json:"topics"`
}
