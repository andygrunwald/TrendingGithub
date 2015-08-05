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
	majorVersion = 0
	minorVersion = 0
	patchVersion = 2
	tweetTimes   = 30 * time.Minute
)

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
		// Schedule new tweet
		time.AfterFunc(tweetTimes, func() {
			generateNewTweet(tweetChan, config)
		})

		// Sometimes it happens that we won`t get a project.
		// In this situation we try to avoid empty tweets like ...
		//	* https://twitter.com/TrendingGithub/status/628714326564696064
		//	* https://twitter.com/TrendingGithub/status/628530032361795584
		//	* https://twitter.com/TrendingGithub/status/628348405790711808
		// we will return here
		// We do this check here and not in tweets.go, because otherwise
		// a new tweet won`t be scheduled
		if len(tweet.ProjectName) <= 0 {
			log.Print("No project found. No tweet sent.")
			continue
		}

		// If we are running in debug mode, we won`t tweet the tweet.
		// We will just output them.
		// This is a good development feature ;)
		if *debug {
			log.Printf("Tweet: %s (length: %d)", tweet.Tweet, len(tweet.Tweet))

		} else {
			twitter := NewTwitterClient(&config.Twitter)
			postedTweet, err := twitter.tweet(tweet.Tweet)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Tweet %s posted", postedTweet.IdStr)
			}
		}
		markTweetAsAlreadyTweeted(tweet.ProjectName, config)

	}
}
