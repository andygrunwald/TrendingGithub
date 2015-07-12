package main

import (
	"github.com/ChimeraCoder/anaconda"
	"net/url"
)

type Twitter struct {
	API *anaconda.TwitterApi
}

func NewTwitterClient(config *TwitterConfiguration) *Twitter {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)

	client := Twitter{
		API: api,
	}

	return &client
}

func (client *Twitter) tweet(text string) (*anaconda.Tweet, error) {
	v := url.Values{}
	tweet, err := client.API.PostTweet(text, v)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}
