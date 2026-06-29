package models

type ProfileResponse struct {
	User         GitHubUser   `json:"user"`
	Repositories []Repository `json:"repositories"`
}
