package twitter

import (
	"net/url"
	"sync"

	"github.com/ChimeraCoder/anaconda"
)

// API is the interface to decouple the twitter API
type API interface {
	GetConfiguration(v url.Values) (conf anaconda.Configuration, err error)
	PostTweet(status string, v url.Values) (tweet anaconda.Tweet, err error)
	GetFollowersIds(v url.Values) (c anaconda.Cursor, err error)
	GetFriendshipsLookup(v url.Values) (friendships []anaconda.Friendship, err error)
	FollowUserId(userID int64, v url.Values) (user anaconda.User, err error)
}

// Client is the data structure to reflect the twitter client
type Client struct {
	API           API
	Configuration *anaconda.Configuration
	Mutex         *sync.Mutex
}

// NewClient returns a new client to communicate with twitter.
// To auth against twitter, credentials like consumer key and access tokens
// are necessary. Those can be retrieved via https://apps.twitter.com/.
func NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *Client {
	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerSecret)
	client := &Client{
		API:   api,
		Mutex: &sync.Mutex{},
	}

	return client
}

// NewDebugClient returns a new debug client to interact with.
// This client will not communicate via twitter.
// This client should be used for debugging or developing purpose.
// This client loads a debug configuration.
func NewDebugClient() *Client {
	client := &Client{
		// TODO Create a debugging client
		// API:   api,
		Configuration: &anaconda.Configuration{
			ShortUrlLength:      24,
			ShortUrlLengthHttps: 25,
		},
	}
	return client
}
