package main

import (
	"github.com/andygrunwald/go-trending"
	"log"
	"time"
)

// Tweet is a structure to store the tweet and the project name based on the tweet.
type Tweet struct {
	Tweet       string
	ProjectName string
}

// generateNewTweet is responsible to search a new project / repository
// and build a new tweet based on this.
// The generated tweet will be sent to tweetChan.
func generateNewTweet(tweetChan chan *Tweet, config *Configuration) {
	// TODO Make generation of trending project more intelligent.
	// Currently we do:
	// 		* Get all timeframes and get a random timeframe.
	//		  Get projects based on this timeframe and check if we can tweet them.
	//		  If this is not a success remove this timeframe from slice and continue.
	//
	// Steps to improve the area:
	//		* If no project was chosen, lets get trending languages
	//		  and chose one randomly and request the timeframes again
	//		  and repeat the "slice from removes trick there"
	var projectToTweet trending.Project
	trendingClient := NewTrendingClient()
	redisClient, err := NewRedisClient(&config.Redis)
	if err != nil {
		log.Fatal(err)
	}

	// Get timeframes and randomize them
	timeFrames := trendingClient.GetTimeFrames()
	ShuffleStringSlice(timeFrames)

	for _, timeFrame := range timeFrames {
		log.Printf("Getting projects for timeframe %s", timeFrame)
		getProject := trendingClient.GetRandomProjectGenerator(timeFrame)
		projectToTweet = findProjectWithRandomProjectGenerator(getProject, redisClient)

		// Check if we found a project.
		// If yes we can leave the loop and keep on rockin
		if len(projectToTweet.Name) > 0 {
			break
		}
	}

	tweet := &Tweet{
		Tweet:       buildTweet(projectToTweet),
		ProjectName: projectToTweet.Name,
	}
	tweetChan <- tweet
}

// findProjectWithRandomProjectGenerator retrieves a new project and
// checks if this was already tweeted.
func findProjectWithRandomProjectGenerator(getProject func() (trending.Project, error), redisClient *Redis) trending.Project {
	var projectToTweet trending.Project

	for project, err := getProject(); err == nil; {
		// Check if the project was already tweeted
		alreadyTweeted, err := redisClient.IsRepositoryAlreadyTweeted(project.Name)
		if err != nil {
			log.Println(err)
			continue
		}

		// If the project was already tweeted
		// we will skip this project and go to the next one
		if alreadyTweeted > 0 {
			continue
		}

		// This project wasn`t tweeted yet, so we will take over this job
		projectToTweet = project
		break
	}

	return projectToTweet
}

// buildTweet is responsible to build a 140 length string based on the project we found.
func buildTweet(p trending.Project) string {
	// TODO: We need to improve this.
	// 		* Get the real number of chars for a URL for t.co.
	//		* Get the length of the description and check if we need to trim this.
	// 		* Try to put more information into it (like Language)
	tweet := ""

	// Base length of a tweet
	tweetLen := 140

	// 20 letters for the url
	// see https://dev.twitter.com/overview/t.co
	// TODO we have to replace this by an API call
	tweetLen -= 20

	// Check if the length of the project name is > 120 chars
	// We substract 3 chars, because we will add a suffix " - "
	if nameLen := len(p.Name); nameLen < (tweetLen - 3) {
		tweetLen -= len(p.Name)
		tweet += p.Name

		// Add name suffix " - "
		tweetLen -= 3
		tweet += " - "
	}

	// We only post descriptions if we got more than 20 charactes available
	// + 1 character for a whitespace
	if tweetLen > 21 {
		if len(p.Description) < tweetLen {
			tweet += p.Description
		} else {
			tweet += p.Description[0:(tweetLen - 1)]
		}
		tweet += " "
	}

	// Lets add the URL, but we don`t need to substract the chars
	// because we have done this before
	if p.URL != nil {
		tweet += p.URL.String()
	}

	return tweet
}

// markTweetAsAlreadyTweeted adds a projectName to the global blacklist of already tweeted projects.
// For this we use a Sorted Set where the score is the timestamp of the tweet.
func markTweetAsAlreadyTweeted(projectName string, config *Configuration) (int, error) {
	redisClient, err := NewRedisClient(&config.Redis)
	if err != nil {
		log.Fatal(err)
	}

	// Generate score in format YYYYMMDDHHiiss
	now := time.Now()
	score := now.Format("20060102150405")

	res, err := redisClient.AddRepositoryToTweetedList(projectName, score)
	if err != nil || res != 1 {
		log.Printf("Error during adding project %s to tweeted list: %s (%d)", projectName, err, res)
	}

	return res, err
}
