package main

import (
	"github.com/ChimeraCoder/anaconda"
	"net/url"
)

// Twitter is the datastructure to store the twitter client
type Twitter struct {
	API *anaconda.TwitterApi
}

// NewTwitterClient returns a new client to communicate with twitter (obvious, right?)
func NewTwitterClient(config *TwitterConfiguration) *Twitter {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)

	client := Twitter{
		API: api,
	}

	return &client
}

// tweet will .... tweet the text :D ... Badumts
func (client *Twitter) tweet(text string) (*anaconda.Tweet, error) {
	v := url.Values{}
	tweet, err := client.API.PostTweet(text, v)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}
