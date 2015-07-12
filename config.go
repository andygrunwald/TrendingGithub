package main

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Twitter TwitterConfiguration `json:"twitter"`
	Redis   RedisConfiguration   `json:"redis"`
}

type TwitterConfiguration struct {
	ConsumerKey       string `json:"consumer-key"`
	ConsumerSecret    string `json:"consumer-secret"`
	AccessToken       string `json:"access-token"`
	AccessTokenSecret string `json:"access-token-secret"`
}

type RedisConfiguration struct {
	URL  string `json:"url"`
	Auth string `json:"auth"`
}

func NewConfiguration(configFile *string) (*Configuration, error) {
	fileContent, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	var config Configuration
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
