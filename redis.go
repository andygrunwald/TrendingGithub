package main

import (
	"github.com/garyburd/redigo/redis"
)

// Redis contains the connection to the redis server
type Redis struct {
	Client redis.Conn
}

const (
	// GreyListTTL defined the TTL of a repository in seconds: 1 month and 15 days (~45 days)
	GreyListTTL = 60 * 60 * 24 * 45
	// OK is the standard response of a Redis server if everything went fine
	OK = "OK"
)

// NewRedisClient provides a new instance of a Redis connection
func NewRedisClient(config *RedisConfiguration) (*Redis, error) {
	redisClient, err := redis.Dial("tcp", config.URL)
	if err != nil {
		return nil, err
	}

	r := &Redis{
		Client: redisClient,
	}

	// If the redis server isn`t protected by an Auth, we are done here.
	if len(config.Auth) == 0 {
		return r, nil
	}

	if _, err := r.Client.Do("AUTH", config.Auth); err != nil {
		r.Client.Close()
		return nil, err
	}

	return r, nil
}

// MarkRepositoryAsTweeted marks a single projects as "already tweeted".
// This information will be stored in Redis as a simple set with a TTL.
// The timestamp of the tweet will be used as value.
func (r *Redis) MarkRepositoryAsTweeted(projectName, score string) (bool, error) {
	result, err := redis.String(r.Client.Do("SET", projectName, score, "EX", GreyListTTL, "NX"))
	if result == OK && err == nil {
		return true, err
	}
	return false, err
}

// IsRepositoryAlreadyTweeted checks if a project was already tweeted.
// If it is not available
//	a) the project was not tweeted yet
//	b) the project ttl expired and is ready to tweet again
func (r *Redis) IsRepositoryAlreadyTweeted(projectName string) (bool, error) {
	return redis.Bool(r.Client.Do("EXISTS", projectName))
}
