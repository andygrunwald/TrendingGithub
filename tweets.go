package main

import (
	"github.com/andygrunwald/go-trending"
	"log"
	"strings"
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

	// First get the timeframes without any languages
	projectToTweet = timeframeLoopToSearchAProject(timeFrames, "", trendingClient, redisClient)

	// Check if we found a project. If yes tweet it.
	if isProjectEmpty(projectToTweet) == false {
		sendProject(tweetChan, projectToTweet)
		return
	}

	// If not, keep going and try to get some (trending) languages
	languages := trendingClient.GetTrendingLanguages()
	ShuffleStringSlice(languages)
	ShuffleStringSlice(timeFrames)

	for _, language := range languages {
		projectToTweet = timeframeLoopToSearchAProject(timeFrames, language, trendingClient, redisClient)

		// If we found a project, break this loop again.
		if isProjectEmpty(projectToTweet) == false {
			sendProject(tweetChan, projectToTweet)
			break
		}
	}

	// If we don`t find a tweet we need to send a empty project
	// This is necessary to reschedule the search for a new repository
	p := trending.Project{}
	sendProject(tweetChan, p)
}

// timeframeLoopToSearchAProject provides basicly a loop over incoming timeFrames (+ language)
// to try to find a new tweet.
// You can say that this is nearly the <3 of this bot.
func timeframeLoopToSearchAProject(timeFrames []string, language string, trendingClient *Trend, redisClient *Redis) trending.Project {
	var projectToTweet trending.Project

	for _, timeFrame := range timeFrames {
		if len(language) > 0 {
			log.Printf("Getting projects for timeframe %s and language %s", timeFrame, language)
		} else {
			log.Printf("Getting projects for timeframe %s", timeFrame)
		}

		getProject := trendingClient.GetRandomProjectGenerator(timeFrame, language)
		projectToTweet = findProjectWithRandomProjectGenerator(getProject, redisClient)

		// Check if we found a project.
		// If yes we can leave the loop and keep on rockin
		if isProjectEmpty(projectToTweet) == false {
			break
		}
	}

	return projectToTweet
}

// sendProject puts the project we want to tweet into the tweet queue
// If the queue is ready to receive a new project, this will be tweeted
func sendProject(tweetChan chan *Tweet, p trending.Project) {
	text := ""
	// Only build tweet if necessary
	if len(p.Name) > 0 {
		text = buildTweet(p)
	}

	tweet := &Tweet{
		Tweet:       text,
		ProjectName: p.Name,
	}
	tweetChan <- tweet
}

// isProjectEmpty checks if the incoming project is empty
func isProjectEmpty(p trending.Project) bool {
	if len(p.Name) > 0 {
		return false
	}

	return true
}

// findProjectWithRandomProjectGenerator retrieves a new project and
// checks if this was already tweeted.
func findProjectWithRandomProjectGenerator(getProject func() (trending.Project, error), redisClient *Redis) trending.Project {
	var projectToTweet trending.Project
	var project trending.Project
	var projectErr error

	for project, projectErr = getProject(); projectErr == nil; project, projectErr = getProject() {
		// Check if the project was already tweeted
		alreadyTweeted, err := redisClient.IsRepositoryAlreadyTweeted(project.Name)
		if err != nil {
			log.Println(err)
			continue
		}

		// If the project was already tweeted
		// we will skip this project and go to the next one
		if alreadyTweeted == true {
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

	// 25 letters for the url (+ 1 character for a whitespace)
	// TODO: Daily GET call to help/configuration to receive max length for URL
	// Only with this we are able to fill out 140 chars as max as possible.
	// @link https://dev.twitter.com/overview/t.co
	// @link https://dev.twitter.com/rest/reference/get/help/configuration
	//
	// Currently this value is hardcoded.
	// Anaconda (the twitter library we use) doesn`t support this yet.
	// There is a pull request waiting: https://github.com/ChimeraCoder/anaconda/pull/66
	// Today (2015-07-26) the values are
	//	"short_url_length": 22,
	//	"short_url_length_https": 23
	// We choose a few chars more to get some more time until Anaconda accepts the PR
	tweetLen -= 26

	// Check if the length of the project name is > 114 chars
	if nameLen := len(p.Name); nameLen < tweetLen {
		tweetLen -= len(p.Name)
		tweet += p.Name
	}

	// We only post a description if we got more than 20 charactes available
	// We have to add 3 chars more, because of the prefix " - "
	if tweetLen > 23 && len(p.Description) > 0 {
		tweetLen -= 3
		tweet += " - "

		projectDescription := ""
		if len(p.Description) < tweetLen {
			projectDescription = p.Description
		} else {
			projectDescription = crop(p.Description, (tweetLen - 4), "...", true)
		}

		tweetLen -= len(projectDescription)
		tweet += projectDescription
	}

	// Lets add the URL, but we don`t need to substract the chars
	// because we have done this before
	if p.URL != nil {
		tweet += " "
		tweet += p.URL.String()
	}

	// Lets check if we got space left to add the language as hashtag
	language := strings.Replace(p.Language, " ", "", -1)
	// len + 2, because of " #" in front of the hashtag
	hashTagLen := (len(language) + 2)
	if len(language) > 0 && tweetLen >= hashTagLen {
		tweet += " #" + language
		tweetLen -= hashTagLen
	}

	return tweet
}

// markTweetAsAlreadyTweeted adds a projectName to the global blacklist of already tweeted projects.
// For this we use a Sorted Set where the score is the timestamp of the tweet.
func markTweetAsAlreadyTweeted(projectName string, config *Configuration) (bool, error) {
	redisClient, err := NewRedisClient(&config.Redis)
	if err != nil {
		log.Fatal(err)
	}

	// Generate score in format YYYYMMDDHHiiss
	now := time.Now()
	score := now.Format("20060102150405")

	res, err := redisClient.MarkRepositoryAsTweeted(projectName, score)
	if err != nil || res != true {
		log.Printf("Error during adding project %s to tweeted list: %s (%d)", projectName, err, res)
	}

	return res, err
}
