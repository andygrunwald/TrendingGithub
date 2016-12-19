package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"math/rand"
	"net/url"
	"strconv"
	"sync"
	"time"
	"log"
)

const (
	// Set for debugging purpose
	DebugURLLength = 25
	// Status "none" of GET friendships/lookup call
	// @link https://dev.twitter.com/rest/reference/get/friendships/lookup
	TwitterStatusNone = "none"
)

type TwitterAPI interface {
	GetConfiguration(v url.Values) (conf anaconda.Configuration, err error)
	PostTweet(status string, v url.Values) (tweet anaconda.Tweet, err error)
	GetFollowersIds(v url.Values) (c anaconda.Cursor, err error)
	GetFriendshipsLookup(v url.Values) (friendships []anaconda.Friendship, err error)
	FollowUserId(userId int64, v url.Values) (user anaconda.User, err error)
}

// Twitter is the datastructure to store the twitter client
type Twitter struct {
	API           TwitterAPI
	Configuration *anaconda.Configuration
	Mutex         *sync.Mutex
}

// NewClient returns a new client to communicate with twitter (obvious, right?)
func NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string, debug *bool) *Twitter {
	var client *Twitter
	// If we are running in debug mode, we won`t tweet the tweet.
	if *debug == false {

		// Create anaconda client
		anaconda.SetConsumerKey(consumerKey)
		anaconda.SetConsumerSecret(consumerSecret)
		api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
		client = Twitter{
			API:   api,
			Mutex: &sync.Mutex{},
		}
		err := client.LoadConfiguration()
		if err != nil {
			log.Fatal("Twitter Configuration initialisation failed:", err)
		}
	} else {
		client = &Twitter{
			Configuration: GetDebugConfiguration(),
		}
	}

	return client
}

func (client *Twitter) LoadConfiguration() error {
	v := url.Values{}
	conf, err := client.API.GetConfiguration(v)
	if err != nil {
		return nil
	}

	client.Mutex.Lock()
	client.Configuration = &conf
	client.Mutex.Unlock()

	return nil
}

func (client *Twitter) SetupConfigurationRefresh(d time.Duration) {
	go func() {
		for _ = range time.Tick(d) {
			client.LoadConfiguration()
		}
	}()
}

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
func (client *Twitter) FollowNewPerson() {
	// Get own followers
	c, err := client.API.GetFollowersIds(nil)
	if err != nil {
		return
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
			return
		}

		// Choose a random follower (again) from the random follower chosen before
		randomNumber = rand.Intn(len(c.Ids))

		// Get friendship details of @TrendingGithub and the chosen person
		v = url.Values{}
		v.Add("user_id", strconv.FormatInt(c.Ids[randomNumber], 10))
		friendships, err := client.API.GetFriendshipsLookup(v)
		if err != nil {
			return
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
		return
	}
}

// isThereARelationship will test if @TrendingGithub has a relationship to the new user.
// Only if there is no relationship ("none")
// Default wise we assume that we got a relationship already
// @link https://dev.twitter.com/rest/reference/get/friendships/lookup
func (client *Twitter) isThereARelationship(friendships []anaconda.Friendship) bool {
	shouldIFollow := false
	for _, friend := range friendships {
		for _, status := range friend.Connections {
			if status == TwitterStatusNone {
				shouldIFollow = true
				break
			}
		}
	}

	return shouldIFollow
}

// Tweet will .... tweet the text :D ... Badumts
func (client *Twitter) Tweet(text string) (*anaconda.Tweet, error) {
	v := url.Values{}
	tweet, err := client.API.PostTweet(text, v)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}

func GetDebugConfiguration() *anaconda.Configuration {
	return &anaconda.Configuration{
		ShortUrlLength:      24,
		ShortUrlLengthHttps: 25,
	}
}
