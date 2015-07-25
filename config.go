package main

import (
	"encoding/json"
	"io/ioutil"
)

// Configuration is the main data structure for the configuration
// This structure reflects the config.json file
type Configuration struct {
	Twitter TwitterConfiguration `json:"twitter"`
	Redis   RedisConfiguration   `json:"redis"`
}

// TwitterConfiguration is the configuration structure for the twitter connection.
type TwitterConfiguration struct {
	ConsumerKey       string `json:"consumer-key"`
	ConsumerSecret    string `json:"consumer-secret"`
	AccessToken       string `json:"access-token"`
	AccessTokenSecret string `json:"access-token-secret"`
}

// RedisConfiguration is the configuration structure for the redis connection.
type RedisConfiguration struct {
	URL  string `json:"url"`
	Auth string `json:"auth"`
}

// NewConfiguration will provide a new instance of Configuration
// configFile accepts a path to a configuration file based on config.json
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
