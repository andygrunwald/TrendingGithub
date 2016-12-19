package twitter

import (
	"math/rand"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

const (
	// twitterStatusNone reflects the status "none" of GET friendships/lookup call
	// See https://dev.twitter.com/rest/reference/get/friendships/lookup
	twitterStatusNone = "none"
)

// TwitterAPI is the interface to decouple the twitter API.
type TwitterAPI interface {
	GetConfiguration(v url.Values) (conf anaconda.Configuration, err error)
	PostTweet(status string, v url.Values) (tweet anaconda.Tweet, err error)
	GetFollowersIds(v url.Values) (c anaconda.Cursor, err error)
	GetFriendshipsLookup(v url.Values) (friendships []anaconda.Friendship, err error)
	FollowUserId(userId int64, v url.Values) (user anaconda.User, err error)
}

// Twitter is the data structure to reflect the twitter client
type Twitter struct {
	API           TwitterAPI
	Configuration *anaconda.Configuration
	Mutex         *sync.Mutex
}

// NewClient returns a new client to communicate with twitter.
// If debug is enabled, we will load a debug configuration for the twitter client.
func NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string, debug bool) *Twitter {
	var client *Twitter

	// If we are running in debug mode, we won`t tweet the tweet.
	if !debug {
		// Create anaconda client
		anaconda.SetConsumerKey(consumerKey)
		anaconda.SetConsumerSecret(consumerSecret)
		api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
		client = &Twitter{
			API:   api,
			Mutex: &sync.Mutex{},
		}
	} else {
		client = &Twitter{
			// Debug configuration
			Configuration: &anaconda.Configuration{
				ShortUrlLength:      24,
				ShortUrlLengthHttps: 25,
			},
		}
	}

	return client
}

// LoadConfiguration loads the configuration of the twitter client from twitter.
// See https://dev.twitter.com/rest/reference/get/help/configuration
func (client *Twitter) LoadConfiguration() error {
	v := url.Values{}
	conf, err := client.API.GetConfiguration(v)
	if err != nil {
		return err
	}

	client.Mutex.Lock()
	client.Configuration = &conf
	client.Mutex.Unlock()

	return nil
}

// SetupConfigurationRefresh sets up a scheduler and will refresh the configuration from twitter every duration d.
// The reason for this is that these values can change over time.
// See https://dev.twitter.com/rest/reference/get/help/configuration
func (client *Twitter) SetupConfigurationRefresh(d time.Duration) {
	go func() {
		for _ = range time.Tick(d) {
			client.LoadConfiguration()
		}
	}()
}

// SetupFollowNewPeopleScheduling sets up a scheduler and will search for a new person to follow every duration d.
// This is our growth hack feature.
func (client *Twitter) SetupFollowNewPeopleScheduling(d time.Duration) {
	go func() {
		for _ = range time.Tick(d) {
			client.FollowNewPerson()
		}
	}()
}

// FollowNewPerson will follow a new person on twitter to raise the attraction for the bot.
// We will follow a new person who follow on random follower of @TrendingGithub
// Only persons who don`t have a relationship to the bot will be followed.
func (client *Twitter) FollowNewPerson() error{
	// Get own followers
	c, err := client.API.GetFollowersIds(nil)
	if err != nil {
		return err
	}

	// We loop here, because we want to follow one person.
	// If we choose a random person it can be that we choose a person
	// that follows @TrendingGithub already.
	// We want to attract new persons ;)
	for {
		// Choose a random follower
		randomNumber := rand.Intn(len(c.Ids))

		// Request the follower from the random follower
		v := url.Values{}
		v.Add("user_id", strconv.FormatInt(c.Ids[randomNumber], 10))
		c, err = client.API.GetFollowersIds(v)
		if err != nil {
			return err
		}

		// Choose a random follower (again) from the random follower chosen before
		randomNumber = rand.Intn(len(c.Ids))

		// Get friendship details of @TrendingGithub and the chosen person
		v = url.Values{}
		v.Add("user_id", strconv.FormatInt(c.Ids[randomNumber], 10))
		friendships, err := client.API.GetFriendshipsLookup(v)
		if err != nil {
			return err
		}

		// Test if @TrendingGithub has a relationship to the new user.
		shouldIFollow := client.isThereARelationship(friendships)

		// If we got a relationship, we will repeat the process ...
		if shouldIFollow == false {
			continue
		}

		// ... if not we will follow the new person
		// We drop the error and user here, because we got no logging yet ;)
		client.API.FollowUserId(c.Ids[randomNumber], nil)
		return nil
	}
}

// isThereARelationship will test if @TrendingGithub has a relationship to the new user.
// Only if there is no relationship ("none").
// Default wise we assume that we got a relationship already
// See https://dev.twitter.com/rest/reference/get/friendships/lookup
func (client *Twitter) isThereARelationship(friendships []anaconda.Friendship) bool {
	shouldIFollow := false
	for _, friend := range friendships {
		for _, status := range friend.Connections {
			if status == twitterStatusNone {
				shouldIFollow = true
				break
			}
		}
	}

	return shouldIFollow
}

// Tweet will ... tweet the text :D ... Badum ts
func (client *Twitter) Tweet(text string) (*anaconda.Tweet, error) {
	v := url.Values{}
	tweet, err := client.API.PostTweet(text, v)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}