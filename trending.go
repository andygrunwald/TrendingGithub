package main

import (
	"errors"
	"github.com/andygrunwald/go-trending"
	"math/rand"
)

// Trend is the datastructure to hold a github-trending client.
// This will be used to retrieve trending projects
type Trend struct {
	Client *trending.Trending
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

// GetRandomProjectGenerator returns a closure to retrieve a random project based on timeFrame.
// timeFrame is a string based on the timeframes provided by go-trending or GetTimeFrames.
// language is a (programing) language provided by go-trending. Can be empty as well.
func (t *Trend) GetRandomProjectGenerator(timeFrame, language string) func() (trending.Project, error) {
	var projects []trending.Project
	var err error

	// Get projects based on timeframe
	// This makes the initial HTTP call to github.
	githubTrending := trending.NewTrending()
	projects, err = githubTrending.GetProjects(timeFrame, language)
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
