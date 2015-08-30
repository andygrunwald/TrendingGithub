package storage

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	// GreyListTTL defined the TTL of a repository in seconds: 1 month and 15 days (~45 days)
	GreyListTTL = 60 * 60 * 24 * 45
	// OK is the standard response of a Redis server if everything went fine
	OK = "OK"
)

type RedisStorage struct{}

type RedisPool struct {
	pool *redis.Pool
}

type RedisConnection struct {
	conn redis.Conn
}

func (rs *RedisStorage) NewPool(url, auth string) Pool {
	rp := RedisPool{
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", url)
				if err != nil {
					return nil, err
				}
				if _, err := c.Do("AUTH", auth); err != nil {
					c.Close()
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}

	return rp
}

func (rp RedisPool) Close() error {
	return rp.pool.Close()
}

func (rp RedisPool) Get() Connection {
	rc := RedisConnection{
		conn: rp.pool.Get(),
	}
	return &rc
}

func (rc *RedisConnection) Close() error {
	return rc.conn.Close()
}

// Redis contains the connection to the redis server
/*
type Redis struct {
	Client redis.Conn
}
*/

/*
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
*/
// MarkRepositoryAsTweeted marks a single projects as "already tweeted".
// This information will be stored in Redis as a simple set with a TTL.
// The timestamp of the tweet will be used as value.
func (rc *RedisConnection) MarkRepositoryAsTweeted(projectName, score string) (bool, error) {
	result, err := redis.String(rc.conn.Do("SET", projectName, score, "EX", GreyListTTL, "NX"))
	if result == OK && err == nil {
		return true, err
	}
	return false, err
}

// IsRepositoryAlreadyTweeted checks if a project was already tweeted.
// If it is not available
//	a) the project was not tweeted yet
//	b) the project ttl expired and is ready to tweet again
func (rc *RedisConnection) IsRepositoryAlreadyTweeted(projectName string) (bool, error) {
	return redis.Bool(rc.conn.Do("EXISTS", projectName))
}
