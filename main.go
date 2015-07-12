package main

import (
	"fmt"
	"github.com/andygrunwald/go-trending"
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

func StartTweeting(config *Configuration, debug *bool) {
	trend := NewTrendingClient()

	// TODO Maybe this can be better done by channels

	// TODO Make generation of trending project more intelligent
	// Steps to do:
	//		* Get all timeframes and check a random timeframe
	//		  if this is not a success remove this timeframe from slice
	//		  and continue.
	//		* If no project was chosen, lets get trending languages
	//		  and chose one randomly and request the timeframes again
	//		  and repeat the "slice from removes trick there"

	// Endless loop, because the bot should not stop tweeting :)
	for {
		redisClient, err := NewRedisClient(&config.Redis)
		if err != nil {
			log.Fatal(err)
		}
	NewTimeFrame:
		timeFrame := trend.getRandomTimeFrame()
		log.Printf("Getting projects for timeframe %s", timeFrame)
		getProject := trend.getRandomProjectGenerator(timeFrame)

		var getProjectError error
		var p trending.Project
		for getProjectError == nil {
		LoopStart:
			p, getProjectError = getProject()
			if getProjectError != nil {
				log.Println(getProjectError)
				goto NewTimeFrame
			}

			// We found a tweet
			alreadyTweeted, err := redisClient.IsRepositoryAlreadyTweeted(p.Name)
			if err != nil {
				log.Println(err)
				goto LoopStart
			}

			if alreadyTweeted > 0 {
				goto LoopStart
			}

			goto Tweet
		}
	Tweet:
		tweetProject(p, redisClient, config, debug)

		// Lets sleep for ~1h
		// Currently i think it is okay to tweet every hour,
		// because we don`t hit the rate limit of the Twitter API
		// and new projects must be trending ;)
		// With this we got 24 tweets per day.
		log.Println("Going to sleep now.")
		time.Sleep(1 * time.Hour)
	}
}

func tweetProject(p trending.Project, redisClient *Redis, config *Configuration, debug *bool) {
	tweet := buildTweet(p)

	// Generate score in format YYYYMMDDHHiiss
	now := time.Now()
	score := now.Format("20060102150405")

	// TODO Switch to sorted set and use timestamp as score
	res, err := redisClient.AddRepositoryToTweetedList(p.Name, score)
	if err != nil || res != 1 {
		log.Printf("Error during adding project %s to tweeted list: %s (%d)", p.Name, err, res)
	}

	if *debug {
		log.Printf("Tweet: %s (length: %d)", tweet, len(tweet))

	} else {
		twitter := NewTwitterClient(&config.Twitter)
		postedTweet, err := twitter.tweet(tweet)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Tweet %s posted", postedTweet.IdStr)
		}
	}
}

func buildTweet(p trending.Project) string {
	tweet := ""

	tweetLen := 140
	// 20 letters for the url
	// see https://dev.twitter.com/overview/t.co
	// TODO we have to replace this by an API call
	tweetLen -= 20

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

	if p.URL != nil {
		tweet += p.URL.String()
	}

	return tweet
}
