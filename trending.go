package main

import (
	"errors"
	"math/rand"

	"github.com/andygrunwald/go-trending"
)

type TrendingClient interface {
	GetTrendingLanguages() ([]trending.Language, error)
	GetProjects(time, language string) ([]trending.Project, error)
}

// Trend is the datastructure to hold a github-trending client.
// This will be used to retrieve trending projects
type Trend struct {
	Client TrendingClient
}

// NewTrendingClient will provide a new instance of Trend.
func NewTrendingClient() *Trend {
	githubTrending := trending.NewTrending()

	t := &Trend{
		Client: githubTrending,
	}

	return t
}

// GetTimeFrames returns all available timeframes of go-trending package.
func (t *Trend) GetTimeFrames() []string {
	return []string{trending.TimeToday, trending.TimeWeek, trending.TimeMonth}
}

// GetTrendingLanguages returns all trending languages,
// but only the URLNames, because this is what we need.
// Errors are not important here.
func (t *Trend) GetTrendingLanguages() []string {
	languages, err := t.Client.GetTrendingLanguages()
	if err != nil {
		return []string{}
	}

	// I know. Slices with a predefined number of elements (0) is not a good idea.
	// But we are calling an external API and don`t know how many items will be there.
	// Furthere more we will filter some languages in the loop.
	// Does anyone got a better idea? Contact me!
	var trendingLanguages []string
	for _, language := range languages {
		if len(language.URLName) > 0 {
			trendingLanguages = append(trendingLanguages, language.URLName)
		}
	}

	return trendingLanguages
}

// GetRandomProjectGenerator returns a closure to retrieve a random project based on timeFrame.
// timeFrame is a string based on the timeframes provided by go-trending or GetTimeFrames.
// language is a (programing) language provided by go-trending. Can be empty as well.
func (t *Trend) GetRandomProjectGenerator(timeFrame, language string) func() (trending.Project, error) {
	var projects []trending.Project
	var err error

	// Get projects based on timeframe
	// This makes the initial HTTP call to github.
	projects, err = t.Client.GetProjects(timeFrame, language)
	if err != nil {
		return func() (trending.Project, error) {
			return trending.Project{}, err
		}
	}

	// Once we got the projects we will provide a closure
	// to retrieve random projects of this project list.
	return func() (trending.Project, error) {

		// Check the number of projects left in the list
		// If there are no more projects anymore, we will return an error.
		numOfProjects := len(projects)
		if numOfProjects == 0 {
			return trending.Project{}, errors.New("No projects found")
		}

		// If there are projects left, chose a random one ...
		randomNumber := rand.Intn(numOfProjects)
		randomProject := projects[randomNumber]

		// ... and delete the chosen project from our list.
		projects = append(projects[:randomNumber], projects[randomNumber+1:]...)

		return randomProject, nil
	}
}
