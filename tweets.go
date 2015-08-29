package main

import (
	"github.com/andygrunwald/go-trending"
	"log"
	"strings"
	"time"
)

type TweetSearch struct {
	Channel       chan *Tweet
	Configuration *Configuration
	URLLength     int
}

// Tweet is a structure to store the tweet and the project name based on the tweet.
type Tweet struct {
	Tweet       string
	ProjectName string
}

// GenerateNewTweet is responsible to search a new project / repository
// and build a new tweet based on this.
// The generated tweet will be sent to tweetChan.
func (ts *TweetSearch) GenerateNewTweet() {
	var projectToTweet trending.Project
	trendingClient := NewTrendingClient()
	redisClient, err := NewRedisClient(&ts.Configuration.Redis)
	if err != nil {
		log.Fatal(err)
	}

	// Get timeframes and randomize them
	timeFrames := trendingClient.GetTimeFrames()
	ShuffleStringSlice(timeFrames)

	// First get the timeframes without any languages
	projectToTweet = ts.TimeframeLoopToSearchAProject(timeFrames, "", trendingClient, redisClient)

	// Check if we found a project. If yes tweet it.
	if ts.IsProjectEmpty(projectToTweet) == false {
		ts.SendProject(projectToTweet)
		return
	}

	// If not, keep going and try to get some (trending) languages
	languages := trendingClient.GetTrendingLanguages()
	ShuffleStringSlice(languages)
	ShuffleStringSlice(timeFrames)

	for _, language := range languages {
		projectToTweet = ts.TimeframeLoopToSearchAProject(timeFrames, language, trendingClient, redisClient)

		// If we found a project, break this loop again.
		if ts.IsProjectEmpty(projectToTweet) == false {
			ts.SendProject(projectToTweet)
			break
		}
	}
}

// timeframeLoopToSearchAProject provides basicly a loop over incoming timeFrames (+ language)
// to try to find a new tweet.
// You can say that this is nearly the <3 of this bot.
func (ts *TweetSearch) TimeframeLoopToSearchAProject(timeFrames []string, language string, trendingClient *Trend, redisClient *Redis) trending.Project {
	var projectToTweet trending.Project

	for _, timeFrame := range timeFrames {
		if len(language) > 0 {
			log.Printf("Getting projects for timeframe %s and language %s", timeFrame, language)
		} else {
			log.Printf("Getting projects for timeframe %s", timeFrame)
		}

		getProject := trendingClient.GetRandomProjectGenerator(timeFrame, language)
		projectToTweet = ts.FindProjectWithRandomProjectGenerator(getProject, redisClient)

		// Check if we found a project.
		// If yes we can leave the loop and keep on rockin
		if ts.IsProjectEmpty(projectToTweet) == false {
			break
		}
	}

	return projectToTweet
}

// sendProject puts the project we want to tweet into the tweet queue
// If the queue is ready to receive a new project, this will be tweeted
func (ts *TweetSearch) SendProject(p trending.Project) {
	text := ""
	// Only build tweet if necessary
	if len(p.Name) > 0 {
		text = ts.BuildTweet(p)
	}

	tweet := &Tweet{
		Tweet:       text,
		ProjectName: p.Name,
	}
	ts.Channel <- tweet
}

// isProjectEmpty checks if the incoming project is empty
func (ts *TweetSearch) IsProjectEmpty(p trending.Project) bool {
	if len(p.Name) > 0 {
		return false
	}

	return true
}

// findProjectWithRandomProjectGenerator retrieves a new project and
// checks if this was already tweeted.
func (ts *TweetSearch) FindProjectWithRandomProjectGenerator(getProject func() (trending.Project, error), redisClient *Redis) trending.Project {
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
func (ts *TweetSearch) BuildTweet(p trending.Project) string {
	tweet := ""

	// Base length of a tweet
	tweetLen := 140

	// Number of letters for the url (+ 1 character for a whitespace)
	// As URL shortener t.co from twitter is used
	// @link https://dev.twitter.com/overview/t.co
	tweetLen -= ts.URLLength + 1

	// Check if the length of the project name is bigger than the space in the tweet
	// BUG(andygrunwald): If a name of a project got more chars as a tweet, we will generate a tweet without project name
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
			projectDescription = Crop(p.Description, (tweetLen - 4), "...", true)
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
func (ts *TweetSearch) MarkTweetAsAlreadyTweeted(projectName string, config *Configuration) (bool, error) {
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
