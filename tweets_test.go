package main

import (
	"net/url"
	"testing"

	"github.com/andygrunwald/go-trending"
	"github.com/google/go-github/github"
)

func TestTweets_IsProjectEmpty(t *testing.T) {
	ts := TweetSearch{}
	mock := []struct {
		Project trending.Project
		Result  bool
	}{
		{trending.Project{Name: ""}, true},
		{trending.Project{Name: "MyProject"}, false},
	}

	for _, item := range mock {
		res := ts.IsProjectEmpty(item.Project)
		if res != item.Result {
			t.Errorf("Failed for project \"%s\", got %v, expected %v", item.Project.Name, res, item.Result)
		}
	}
}

func TestTweets_BuildTweet(t *testing.T) {
	owner := "andygrunwald"
	repositoryName := "TrendingGithub"
	projectName := owner + "/" + repositoryName
	projectURL, _ := url.Parse("https://github.com/andygrunwald/TrendingGithub")
	projectDescription := "A twitter bot (@TrendingGithub) to tweet trending repositories and developers from GitHub"

	ts := TweetSearch{
		URLLength: 24,
	}

	stars := 123
	repository := &github.Repository{
		StargazersCount: &(stars),
	}

	mock := []struct {
		Project trending.Project
		Result  string
	}{
		//{trending.Project{Name: ""}, "true"},
		/*
			{trending.Project{
				Name:        "SuperDuperOwnerOrOrganisation/This-Is-A-Super-Long-Project-Name-That-Will-Maybe-Kill-My-Tweet-Generation-But-I-Think-It-Is-Useful-To-Test",
				Description: projectDescription + " and more and better and super duper text",
				Language:    "Go",
				URL:         projectURL,
			}, "andygrunwald/TrendingGithub - A twitter bot (@TrendingGithub) to tweet trending repositories and developers... https://github.com/andygrunwald/TrendingGithub #Go"},
		*/
		{trending.Project{
			Name:           "SuperDuperOwnerOrOrganisation/This-Is-A-Long-Project-Name-That-Will-Drop-The-Description-Of-The-Project",
			Owner:          "SuperDuperOwnerOrOrganisation",
			RepositoryName: "This-Is-A-Long-Project-Name-That-Will-Drop-The-Description-Of-The-Project",
			Description:    projectDescription + " and more and better and super duper text",
			Language:       "Go",
			URL:            projectURL,
		}, "SuperDuperOwnerOrOrganisation/This-Is-A-Long-Project-Name-That-Will-Drop-The-Description-Of-The-Project ★123 https://github.com/andygrunwald/TrendingGithub #Go"},
		{trending.Project{
			Name:           projectName + "-cool-super-project",
			Owner:          owner,
			RepositoryName: repositoryName + "-cool-super-project",
			Description:    projectDescription + " and more and better and super duper text",
			Language:       "Go",
			URL:            projectURL,
		}, "andygrunwald/TrendingGithub-cool-super-project: A twitter bot (@TrendingGithub) to tweet trending... ★123 https://github.com/andygrunwald/TrendingGithub #Go"},
		{trending.Project{
			Name:           projectName,
			Owner:          owner,
			RepositoryName: repositoryName,
			Description:    projectDescription,
			Language:       "Go",
			URL:            projectURL,
		}, "andygrunwald/TrendingGithub: A twitter bot (@TrendingGithub) to tweet trending repositories and developers... ★123 https://github.com/andygrunwald/TrendingGithub"},
		{trending.Project{
			Name:           projectName,
			Owner:          owner,
			RepositoryName: repositoryName,
			Description:    "Short description",
			Language:       "Go Lang",
			URL:            projectURL,
		}, "andygrunwald/TrendingGithub: Short description ★123 https://github.com/andygrunwald/TrendingGithub #GoLang"},
		{trending.Project{
			Name:           projectName,
			Owner:          owner,
			RepositoryName: repositoryName,
			Description:    "Project without a URL",
			Language:       "Go Lang",
		}, "andygrunwald/TrendingGithub: Project without a URL ★123 #GoLang"},
		{trending.Project{
			Name:           repositoryName + "/" + repositoryName,
			Owner:          repositoryName,
			RepositoryName: repositoryName,
			Description:    projectDescription,
			Language:       "Go",
			URL:            projectURL,
		}, "TrendingGithub: A twitter bot (@TrendingGithub) to tweet trending repositories and developers from GitHub ★123 https://github.com/andygrunwald/TrendingGithub #Go"},
	}

	for _, item := range mock {
		res := ts.BuildTweet(item.Project, repository)
		if res != item.Result {
			t.Errorf("Failed building a tweet for project \"%s\". Got \"%s\", expected \"%s\"", item.Project.Name, res, item.Result)
		}
	}
}
