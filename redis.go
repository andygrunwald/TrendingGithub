package main

import (
	"github.com/garyburd/redigo/redis"
)

// Redis contains the connection to the redis server
type Redis struct {
	Client redis.Conn
}

const (
	projectKey = "tweeted-repositories"
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

// AddRepositoryToTweetedList adds single projects to a "blacklist" in Redis.
// This list will be used to check if this project was already tweeted.
// The timestamp of the tweet will be used as score.
// @link http://redis.io/commands#sorted_set
func (r *Redis) AddRepositoryToTweetedList(projectName, score string) (int, error) {
	return redis.Int(r.Client.Do("ZADD", projectKey, score, projectName))
}

// IsRepositoryAlreadyTweeted is the "opposite" of AddRepositoryToTweetedList.
// It checks if a project is a member of our "blacklist" set.
// For this the score doesn`t matter.
func (r *Redis) IsRepositoryAlreadyTweeted(projectName string) (int, error) {
	val, err := r.Client.Do("ZSCORE", projectKey, projectName)

	if val == nil {
		return 0, err
	}

	return redis.Int(val, err)
}
