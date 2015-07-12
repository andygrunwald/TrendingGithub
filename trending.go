package main

import (
	"errors"
	"github.com/andygrunwald/go-trending"
	"math/rand"
)

type Trend struct {
	Client *trending.Trending
}

func NewTrendingClient() *Trend {
	githubTrending := trending.NewTrending()

	t := &Trend{
		Client: githubTrending,
	}

	return t
}

func (t *Trend) getRandomTimeFrame() string {
	elements := t.getTimeFrames()
	return elements[rand.Intn(len(elements))]
}

func (t *Trend) getTimeFrames() []string {
	return []string{trending.TimeToday, trending.TimeWeek, trending.TimeMonth}
}

func (t *Trend) getRandomProjectGenerator(timeFrame string) func() (trending.Project, error) {
	var projects []trending.Project
	var err error

	githubTrending := trending.NewTrending()
	projects, err = githubTrending.GetProjects(timeFrame, "")
	if err != nil {
		return func() (trending.Project, error) {
			return trending.Project{}, err
		}
	}

	return func() (trending.Project, error) {

		numOfProjects := len(projects)
		if numOfProjects == 0 {
			return trending.Project{}, errors.New("No projects found")
		}

		randomNumber := rand.Intn(numOfProjects)
		randomProject := projects[randomNumber]

		// Delete project from list
		projects = append(projects[:randomNumber], projects[randomNumber+1:]...)

		return randomProject, nil
	}
}
