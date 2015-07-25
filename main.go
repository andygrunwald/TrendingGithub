package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	majorVersion = 1
	minorVersion = 0
	patchVersion = 0
	tweetTimes   = 5 * time.Second
)

// TODO: Daily GET call to help/configuration to receive max length for URL
// Only with this we are able to fill out 140 chars as max as possible.
// @link https://dev.twitter.com/overview/t.co
// @link https://dev.twitter.com/rest/reference/get/help/configuration

func main() {

	// Output the version and exit
	if *flagVersion {
		fmt.Printf("@TrendingGithub v%d.%d.%d\n", majorVersion, minorVersion, patchVersion)
		return
	}

	// Check for configuration file
	if len(*flagConfigFile) <= 0 {
		log.Fatal("No configuration file found. Please add the --config parameter")
	}

	// PID-File
	if len(*flagPidFile) > 0 {
		ioutil.WriteFile(*flagPidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
	}

	log.Println("Lets get ready to tweet trending content!")
	defer log.Println("Nice sesssion. A lot of knowledge was tweeted.")

	config, err := NewConfiguration(flagConfigFile)
	if err != nil {
		log.Fatal("Configuration initialisation failed:", err)
	}

	StartTweeting(config, flagDebug)
}

// StartTweeting bundles the main logic of this bot.
// It schedules the times when we are looking for a new project to tweet.
// If we found a project, we will build the tweet and tweet it to our followers.
// Because we love our followers ;)
func StartTweeting(config *Configuration, debug *bool) {
	tweetChan := make(chan *Tweet)

	// Schedule first tweet
	time.AfterFunc(tweetTimes, func() {
		generateNewTweet(tweetChan, config)
	})

	// Waiting for tweets ...
	for tweet := range tweetChan {

		// If we are running in debug mode, we won`t tweet the tweet.
		// We will just output them.
		// This is a good development feature ;)
		if *debug {
			log.Printf("Tweet: %s (length: %d)", tweet, len(tweet.Tweet))

		} else {
			markTweetAsAlreadyTweeted(tweet.ProjectName, config)

			twitter := NewTwitterClient(&config.Twitter)
			postedTweet, err := twitter.tweet(tweet.Tweet)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Tweet %s posted", postedTweet.IdStr)
			}
		}

		// Schedule new tweet
		time.AfterFunc(tweetTimes, func() {
			generateNewTweet(tweetChan, config)
		})
	}
}
