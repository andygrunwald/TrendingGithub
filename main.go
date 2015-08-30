package main

import (
	"flag"
	"fmt"
	"github.com/andygrunwald/TrendingGithub/storage"
	"log"
	"time"
)

const (
	majorVersion             = 0
	minorVersion             = 1
	patchVersion             = 0
	tweetTimes               = 30 * time.Minute
	configurationRefreshTime = 24 * time.Hour
)

func main() {
	var (
		flagConfigFile = flag.String("config", "", "Path to configuration file.")
		flagVersion    = flag.Bool("version", false, "Outputs the version number and exits.")
		flagDebug      = flag.Bool("debug", false, "Outputs the tweet instead of tweet it. Useful for development.")
	)

	flag.Parse()

	// Output the version and exit
	if *flagVersion {
		fmt.Printf("@TrendingGithub v%d.%d.%d\n", majorVersion, minorVersion, patchVersion)
		return
	}

	// Check for configuration file if we are in production (non debug) mode
	if *flagDebug == false && len(*flagConfigFile) <= 0 {
		log.Fatal("No configuration file found. Please add the --config parameter")
	}

	if *flagDebug == true {
		*flagConfigFile = "./config.json.dist"
	}

	log.Println("Lets get ready to tweet trending content!")
	defer log.Println("Nice sesssion. A lot of knowledge was tweeted.")

	config, err := NewConfiguration(*flagConfigFile)
	if err != nil {
		log.Fatal("Configuration initialisation failed:", err)
	}

	twitter := GetTwitterClient(&config.Twitter, flagDebug)
	storageBackend := GetStorageBackend(&config.Redis, flagDebug)
	defer storageBackend.Close()

	StartTweeting(twitter, storageBackend)
}

// StartTweeting bundles the main logic of this bot.
// It schedules the times when we are looking for a new project to tweet.
// If we found a project, we will build the tweet and tweet it to our followers.
// Because we love our followers ;)
func StartTweeting(twitter *Twitter, storageBackend storage.Pool) {

	// Setup tweet scheduling
	ts := &TweetSearch{
		Channel:   make(chan *Tweet),
		Trending:  NewTrendingClient(),
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
		for _ = range time.Tick(tweetTimes) {
			go tweetSearch.GenerateNewTweet()
		}
	}()
}

func GetStorageBackend(config *RedisConfiguration, debug *bool) storage.Pool {
	var pool storage.Pool
	if *debug == false {
		storage := storage.RedisStorage{}
		pool = storage.NewPool(config.URL, config.Auth)
	} else {
		storage := storage.MemoryStorage{}
		pool = storage.NewPool("", "")
	}

	return pool
}

func GetTwitterClient(config *TwitterConfiguration, debug *bool) *Twitter {
	var twitter *Twitter
	// If we are running in debug mode, we won`t tweet the tweet.
	if *debug == false {
		twitter = NewTwitterClient(config)
		err := twitter.LoadConfiguration()
		if err != nil {
			log.Fatal("Twitter Configuration initialisation failed:", err)
		}
		// Refresh the configuration every day
		twitter.SetupConfigurationRefresh(configurationRefreshTime)
	} else {
		twitter = &Twitter{
			Configuration: GetDebugConfiguration(),
		}
	}

	return twitter
}
