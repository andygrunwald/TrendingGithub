package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/andygrunwald/TrendingGithub/flags"
	"github.com/andygrunwald/TrendingGithub/storage"
	trendingwrap "github.com/andygrunwald/TrendingGithub/trending"
	"github.com/andygrunwald/TrendingGithub/twitter"
)

const (
	// Version of @TrendingGithub
	Version = "0.2.0"

	tweetTime = 5 * time.Second
	//tweetTime                = 30 * time.Minute
	configurationRefreshTime = 24 * time.Hour
	followNewPersonTime      = 45 * time.Minute
)

func main() {
	var (
		// Twitter
		twitterConsumerKey       = flags.String("twitter-consumer-key", "TRENDINGGITHUB_TWITTER_CONSUMER_KEY", "", "Twitter-API: Consumer key")
		twitterConsumerSecret    = flags.String("twitter-consumer-secret", "TRENDINGGITHUB_TWITTER_CONSUMER_SECRET", "", "Twitter-API: Consumer secret")
		twitterAccessToken       = flags.String("twitter-access-token", "TRENDINGGITHUB_TWITTER_ACCESS_TOKEN", "", "Twitter-API: Access token")
		twitterAccessTokenSecret = flags.String("twitter-access-token-secret", "TRENDINGGITHUB_TWITTER_ACCESS_TOKEN_SECRET", "", "Twitter-API: Access token secret")
		twitterFollowNewPerson   = flags.Bool("twitter-follow-new-person", "TRENDINGGITHUB_TWITTER_FOLLOW_NEW_PERSON", false, "Twitter: Follows a friend of one of our followers")

		// Redis storage
		storageURL  = flags.String("storage-url", "TRENDINGGITHUB_STORAGE_URL", "", "Storage URL (e.g. 1.2.3.4:6379 or :6379")
		storageAuth = flags.String("storage-auth", "TRENDINGGITHUB_STORAGE_AUTH", "", "Storage Auth (e.g. myPassword or <empty>")

		flagVersion = flags.Bool("version", "TRENDINGGITHUB_VERSION", false, "Outputs the version number and exit")
		flagDebug   = flags.Bool("debug", "TRENDINGGITHUB_DEBUG", false, "Outputs the tweet instead of tweet it. Useful for development.")
	)
	flag.Parse()

	// Output the version and exit
	if *flagVersion {
		fmt.Printf("@TrendingGithub v%s\n", Version)
		return
	}

	log.Println("Hey, nice to meet you. My name is @TrendingGithub. Lets get ready to tweet some trending content!")
	defer log.Println("Nice sesssion. A lot of knowledge was tweeted. Good work and see you next time!")

	twitterClient := twitter.NewClient(*twitterConsumerKey, *twitterConsumerSecret, *twitterAccessToken, *twitterAccessTokenSecret, flagDebug)

	// Refresh the configuration every day
	twitterClient.SetupConfigurationRefresh(configurationRefreshTime)

	// Activate our growth hack feature
	if *twitterFollowNewPerson {
		twitterClient.SetupFollowNewPeopleScheduling(followNewPersonTime)
	}

	storageBackend := storage.GetBackend(*storageURL, *storageAuth, flagDebug)
	defer storageBackend.Close()

	StartTweeting(twitterClient, storageBackend)
}

// StartTweeting bundles the main logic of this bot.
// It schedules the times when we are looking for a new project to tweet.
// If we found a project, we will build the tweet and tweet it to our followers.
// Because we love our followers ;)
func StartTweeting(twitter *twitter.Twitter, storageBackend storage.Pool) {

	// Setup tweet scheduling
	ts := &TweetSearch{
		Channel:   make(chan *Tweet),
		Trending:  trendingwrap.NewTrendingClient(),
		Storage:   storageBackend,
		URLLength: twitter.Configuration.ShortUrlLengthHttps,
	}
	SetupRegularTweetSearchProcess(ts)

	// Waiting for tweets ...
	for tweet := range ts.Channel {
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

		// In debug mode the twitter variable is not available, so we won`t tweet the tweet.
		// We will just output them.
		// This is a good development feature ;)
		if twitter.API == nil {
			log.Printf("Tweet: %s (length: %d)", tweet.Tweet, len(tweet.Tweet))

		} else {
			postedTweet, err := twitter.Tweet(tweet.Tweet)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Tweet %s posted", postedTweet.IdStr)
			}
		}
		ts.MarkTweetAsAlreadyTweeted(tweet.ProjectName)
	}
}

func SetupRegularTweetSearchProcess(tweetSearch *TweetSearch) {
	go func() {
		for _ = range time.Tick(tweetTime) {
			go tweetSearch.GenerateNewTweet()
		}
	}()
}
