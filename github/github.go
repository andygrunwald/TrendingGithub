package github

import (
	"context"

	"github.com/google/go-github/github"
)

// Repository represents a single repository from github.
// This struct is a stripped down version of github.Repository.
// We only return the values we need here.
type Repository struct {
	StargazersCount *int `json:"stargazers_count,omitempty"`
}

// GetRepositoryDetails will retrieve details about the repository owner/repo from github.
func GetRepositoryDetails(owner, repo string) (*Repository, error) {
	client := github.NewClient(nil)
	repository, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if repository == nil {
		return nil, err
	}

	r := &Repository{
		StargazersCount: repository.StargazersCount,
	}
	return r, err
}
