package main

import (
	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	Client redis.Conn
}

const (
	ProjectKey = "tweeted-repositories"
)

func NewRedisClient(config *RedisConfiguration) (*Redis, error) {
	redisClient, err := redis.Dial("tcp", config.URL)
	if err != nil {
		return nil, err
	}

	r := &Redis{
		Client: redisClient,
	}

	if len(config.Auth) == 0 {
		return r, nil
	}

	if _, err := r.Client.Do("AUTH", config.Auth); err != nil {
		r.Client.Close()
		return nil, err
	}

	return r, nil
}

func (r *Redis) AddRepositoryToTweetedList(projectName string) (int, error) {
	return redis.Int(r.Client.Do("SADD", ProjectKey, projectName))
}

func (r *Redis) IsRepositoryAlreadyTweeted(projectName string) (bool, error) {
	return redis.Bool(r.Client.Do("SISMEMBER", ProjectKey, projectName))
}
