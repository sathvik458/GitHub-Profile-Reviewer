package models

type Analysis struct {
	DocumentationScore int      `json:"documentation_score"`
	RepositoryScore    int      `json:"repository_score"`
	ActivityScore      int      `json:"activity_score"`
	OverallScore       int      `json:"overall_score"`
	Recommendations    []string `json:"recommendations"`
}
