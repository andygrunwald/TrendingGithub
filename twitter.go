package main

import (
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	"sync"
	"time"
)

const (
	// Set for debugging purpose
	DebugURLLength = 25
)

type TwitterAPI interface {
	GetConfiguration(v url.Values) (conf anaconda.Configuration, err error)
	PostTweet(status string, v url.Values) (tweet anaconda.Tweet, err error)
}

// Twitter is the datastructure to store the twitter client
type Twitter struct {
	API           TwitterAPI
	Configuration *anaconda.Configuration
	Mutex         *sync.Mutex
}

// NewTwitterClient returns a new client to communicate with twitter (obvious, right?)
func NewTwitterClient(config *TwitterConfiguration) *Twitter {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)

	client := Twitter{
		API:   api,
		Mutex: &sync.Mutex{},
	}

	return &client
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
