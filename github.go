package main

import (
	"github.com/google/go-github/github"
)

func GetRepositoryDetails(owner, repo string) (*github.Repository, error) {
	client := github.NewClient(nil)
	repository, _, err := client.Repositories.Get(owner, repo)

	return repository, err
}
