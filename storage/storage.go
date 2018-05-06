package storage

import (
	"io"
)

const (
	// GreyListTTL defined the TTL of a repository in seconds: 30 days (~30 days)
	GreyListTTL = 60 * 60 * 24 * 30
)

// Storage represents a new storage type.
// Examples as storages are Redis or in memory.
type Storage interface {
	NewPool(url, auth string) Pool
}

// Pool is the implementation of a specific storage type.
// It should handle a pool of connections to communicate with the storage type.
type Pool interface {
	io.Closer
	Get() Connection
}

// Connection represents a single connection out of a pool from a storage type.
type Connection interface {
	io.Closer

	// Err will return an error once one occured
	Err() error

	// MarkRepositoryAsTweeted marks a single projects as "already tweeted".
	// This information will be stored in Redis as a simple set with a TTL.
	// The timestamp of the tweet will be used as value.
	MarkRepositoryAsTweeted(projectName, score string) (bool, error)

	// IsRepositoryAlreadyTweeted checks if a project was already tweeted.
	// If it is not available
	//	a) the project was not tweeted yet
	//	b) the project ttl expired and is ready to tweet again
	IsRepositoryAlreadyTweeted(projectName string) (bool, error)
}

// NewBackend returns a new connection pool based on the requested storage engine.
func NewBackend(storageURL string, storageAuth string) Pool {
	storageBackend := RedisStorage{}
	pool := storageBackend.NewPool(storageURL, storageAuth)

	return pool
}

// NewDebugBackend returns a new connection pool for in memory.
func NewDebugBackend() Pool {
	storageBackend := MemoryStorage{}
	pool := storageBackend.NewPool("", "")

	return pool
}
