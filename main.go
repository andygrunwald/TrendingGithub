package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/andygrunwald/TrendingGithub/flags"
	"github.com/andygrunwald/TrendingGithub/storage"
	"github.com/andygrunwald/TrendingGithub/twitter"
)

const (
	// Version of @TrendingGithub
	Version = "0.2.0"

	tweetTime = 5 * time.Second
	// TODO: Reset
	//tweetTime                = 30 * time.Minute
	configurationRefreshTime = 24 * time.Hour
	followNewPersonTime      = 45 * time.Minute
)

func main() {
	var (
		// Twitter
		twitterConsumerKey       = flags.String("twitter-consumer-key", "TRENDINGGITHUB_TWITTER_CONSUMER_KEY", "", "Twitter-API: Consumer key. Env var: TRENDINGGITHUB_TWITTER_CONSUMER_KEY")
		twitterConsumerSecret    = flags.String("twitter-consumer-secret", "TRENDINGGITHUB_TWITTER_CONSUMER_SECRET", "", "Twitter-API: Consumer secret. Env var: TRENDINGGITHUB_TWITTER_CONSUMER_SECRET")
		twitterAccessToken       = flags.String("twitter-access-token", "TRENDINGGITHUB_TWITTER_ACCESS_TOKEN", "", "Twitter-API: Access token. Env var: TRENDINGGITHUB_TWITTER_ACCESS_TOKEN")
		twitterAccessTokenSecret = flags.String("twitter-access-token-secret", "TRENDINGGITHUB_TWITTER_ACCESS_TOKEN_SECRET", "", "Twitter-API: Access token secret. Env var: TRENDINGGITHUB_TWITTER_ACCESS_TOKEN_SECRET")
		twitterFollowNewPerson   = flags.Bool("twitter-follow-new-person", "TRENDINGGITHUB_TWITTER_FOLLOW_NEW_PERSON", false, "Twitter: Follows a friend of one of our followers. Env var: TRENDINGGITHUB_TWITTER_FOLLOW_NEW_PERSON")

		// Redis storage
		storageURL  = flags.String("storage-url", "TRENDINGGITHUB_STORAGE_URL", "", "Storage URL (e.g. 1.2.3.4:6379 or :6379). Env var: TRENDINGGITHUB_STORAGE_URL")
		storageAuth = flags.String("storage-auth", "TRENDINGGITHUB_STORAGE_AUTH", "", "Storage Auth (e.g. myPassword or <empty>). Env var: TRENDINGGITHUB_STORAGE_AUTH")

		showVersion = flags.Bool("version", "TRENDINGGITHUB_VERSION", false, "Outputs the version number and exit. Env var: TRENDINGGITHUB_VERSION")
		debugMode = flags.Bool("debug", "TRENDINGGITHUB_DEBUG", false, "Outputs the tweet instead of tweet it (useful for development). Env var: TRENDINGGITHUB_DEBUG")
	)
	flag.Parse()

	// Output the version and exit
	if *showVersion {
		fmt.Printf("@TrendingGithub v%s\n", Version)
		return
	}

	log.Println("Hey, nice to meet you. My name is @TrendingGithub. Lets get ready to tweet some trending content!")
	defer log.Println("Nice sesssion. A lot of knowledge was tweeted. Good work and see you next time!")

	// Prepare the twitter client
	twitterClient := twitter.NewClient(*twitterConsumerKey, *twitterConsumerSecret, *twitterAccessToken, *twitterAccessTokenSecret, debugMode)
	twitterClient.SetupConfigurationRefresh(configurationRefreshTime)

	// Activate our growth hack feature
	// Checkout the README for details or read the code (suggested).
	if *twitterFollowNewPerson {
		twitterClient.SetupFollowNewPeopleScheduling(followNewPersonTime)
	}

	// Request a storage backend
	storageBackend := storage.GetBackend(*storageURL, *storageAuth, *debugMode)
	defer storageBackend.Close()

	// Let the party begin
	StartTweeting(twitterClient, storageBackend)
}