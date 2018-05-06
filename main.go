package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/andygrunwald/TrendingGithub/flags"
	"github.com/andygrunwald/TrendingGithub/storage"
	"github.com/andygrunwald/TrendingGithub/twitter"
)

const (
	// Name of the application
	Name = "@TrendingGithub"
)

var (
	// Version of @TrendingGithub
	version = "dev"

	// Build commit of @TrendingGithub
	commit = "none"

	// Build date of @TrendingGithub
	date = "unknown"
)

func main() {
	var (
		// Twitter
		twitterConsumerKey       = flags.String("twitter-consumer-key", "TRENDINGGITHUB_TWITTER_CONSUMER_KEY", "", "Twitter-API: Consumer key. Env var: TRENDINGGITHUB_TWITTER_CONSUMER_KEY")
		twitterConsumerSecret    = flags.String("twitter-consumer-secret", "TRENDINGGITHUB_TWITTER_CONSUMER_SECRET", "", "Twitter-API: Consumer secret. Env var: TRENDINGGITHUB_TWITTER_CONSUMER_SECRET")
		twitterAccessToken       = flags.String("twitter-access-token", "TRENDINGGITHUB_TWITTER_ACCESS_TOKEN", "", "Twitter-API: Access token. Env var: TRENDINGGITHUB_TWITTER_ACCESS_TOKEN")
		twitterAccessTokenSecret = flags.String("twitter-access-token-secret", "TRENDINGGITHUB_TWITTER_ACCESS_TOKEN_SECRET", "", "Twitter-API: Access token secret. Env var: TRENDINGGITHUB_TWITTER_ACCESS_TOKEN_SECRET")
		twitterFollowNewPerson   = flags.Bool("twitter-follow-new-person", "TRENDINGGITHUB_TWITTER_FOLLOW_NEW_PERSON", false, "Twitter: Follows a friend of one of our followers. Env var: TRENDINGGITHUB_TWITTER_FOLLOW_NEW_PERSON")

		// Timings
		tweetTime                = flags.Duration("twitter-tweet-time", "TRENDINGGITHUB_TWITTER_TWEET_TIME", 30*time.Minute, "Twitter: Time interval to search a new project and tweet it. Env var: TRENDINGGITHUB_TWITTER_TWEET_TIME")
		configurationRefreshTime = flags.Duration("twitter-conf-refresh-time", "TRENDINGGITHUB_TWITTER_CONF_REFRESH_TIME", 24*time.Hour, "Twitter: Time interval to refresh the configuration of twitter (e.g. char length for short url). Env var: TRENDINGGITHUB_TWITTER_CONF_REFRESH_TIME")
		followNewPersonTime      = flags.Duration("twitter-follow-new-person-time", "TRENDINGGITHUB_TWITTER_FOLLOW_NEW_PERSON_TIME", 45*time.Minute, "Growth hack: Time interval to search for a new person to follow. Env var: TRENDINGGITHUB_TWITTER_FOLLOW_NEW_PERSON_TIME")

		// Redis storage
		storageURL  = flags.String("storage-url", "TRENDINGGITHUB_STORAGE_URL", ":6379", "Storage URL (e.g. 1.2.3.4:6379 or :6379). Env var: TRENDINGGITHUB_STORAGE_URL")
		storageAuth = flags.String("storage-auth", "TRENDINGGITHUB_STORAGE_AUTH", "", "Storage Auth (e.g. myPassword or <empty>). Env var: TRENDINGGITHUB_STORAGE_AUTH")

		expVarPort  = flags.Int("expvar-port", "TRENDINGGITHUB_EXPVAR_PORT", 8123, "Port which will be used for the expvar TCP server. Env var: TRENDINGGITHUB_EXPVAR_PORT")
		showVersion = flags.Bool("version", "TRENDINGGITHUB_VERSION", false, "Outputs the version number and exit. Env var: TRENDINGGITHUB_VERSION")
		debugMode   = flags.Bool("debug", "TRENDINGGITHUB_DEBUG", false, "Outputs the tweet instead of tweet it (useful for development). Env var: TRENDINGGITHUB_DEBUG")
	)
	flag.Parse()

	// Output the version and exit
	if *showVersion {
		fmt.Printf("%s v%v, commit %v, built at %v", Name, version, commit, date)
		return
	}

	log.Printf("Hey, my name is %s (v%s). Lets get ready to tweet some trending content!\n", Name, version)
	defer log.Println("Nice session. A lot of knowledge was spreaded. Good work. See you next time!")

	twitterClient := initTwitterClient(*twitterConsumerKey, *twitterConsumerSecret, *twitterAccessToken, *twitterAccessTokenSecret, *debugMode, *configurationRefreshTime)

	// Activate the growth hack feature
	if *twitterFollowNewPerson {
		log.Println("Growth hack \"Follow a friend of a friend\": Enabled ✅ ")
		twitterClient.SetupFollowNewPeopleScheduling(*followNewPersonTime)
	}

	storageBackend := initStorageBackend(*storageURL, *storageAuth, *debugMode)
	initExpvarServer(*expVarPort)

	// Let the party begin
	StartTweeting(twitterClient, storageBackend, *tweetTime)
}

// initTwitterClient prepares and initializes the twitter client
func initTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string, debugMode bool, confRefreshTime time.Duration) *twitter.Client {
	var twitterClient *twitter.Client

	if debugMode {
		// When we are running in a debug mode, we are running with a debug configuration.
		// So we don`t need to load the configuration from twitter here.
		twitterClient = twitter.NewDebugClient()

	} else {
		twitterClient = twitter.NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret)
		err := twitterClient.LoadConfiguration()
		if err != nil {
			log.Fatalf("Twitter Configuration: Initialisation ❌  (%s)", err)
		}
		log.Println("Twitter Configuration: Initialisation ✅")
		twitterClient.SetupConfigurationRefresh(confRefreshTime)
	}

	return twitterClient
}

// initExpvarServer will start a small tcp server for the expvar package.
// This server is only available via localhost on localhost:port/debug/vars
func initExpvarServer(port int) {
	sock, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Expvar: Initialisation ❌  (%s)", err)
	}

	go func() {
		log.Printf("Expvar: Available at http://localhost:%d/debug/vars", port)
		http.Serve(sock, nil)
	}()

	log.Println("Expvar: Initialisation ✅")
}

// initStorageBackend will start the storage backend
func initStorageBackend(address, auth string, debug bool) storage.Pool {
	var storageBackend storage.Pool

	if debug {
		storageBackend = storage.NewDebugBackend()
	} else {
		storageBackend = storage.NewBackend(address, auth)
	}

	defer storageBackend.Close()
	log.Println("Storage backend: Initialisation ✅")

	return storageBackend
}
