package storage

import (
	"time"
)

// MemoryStorageContainer is the backend of the "in memory" storage engine.
// It supports a key (string) and a duration as time.
// The time duration can act as a TTL.
type MemoryStorageContainer map[string]time.Time

// MemoryStorage represents the in memory storage engine.
// This storage can be useful for debugging / development
type MemoryStorage struct{}

// MemoryPool is the pool of connections to your local memory ;)
type MemoryPool struct {
	storage MemoryStorageContainer
}

// MemoryConnection represents a in memory connection
type MemoryConnection struct {
	storage MemoryStorageContainer
}

// NewPool returns a pool to communicate with your in memory
func (ms *MemoryStorage) NewPool(url, auth string) Pool {
	return MemoryPool{
		storage: make(MemoryStorageContainer),
	}
}

// Close closes a in memory pool
func (mp MemoryPool) Close() error {
	return nil
}

// Get returns you a connection to your in memory storage
func (mp MemoryPool) Get() Connection {
	return &MemoryConnection{
		storage: mp.storage,
	}
}

// Err will return an error once one occurred
func (mc *MemoryConnection) Err() error {
	return nil
}

// Close shuts down a in memory connection
func (mc *MemoryConnection) Close() error {
	return nil
}

// MarkRepositoryAsTweeted marks a single projects as "already tweeted".
// This information will be stored as a hashmap with a TTL.
// The timestamp of the tweet will be used as value.
func (mc *MemoryConnection) MarkRepositoryAsTweeted(projectName, score string) (bool, error) {
	// Add grey listing to current time
	now := time.Now()
	future := now.Add(time.Second * BlackListTTL)

	mc.storage[projectName] = future

	return true, nil
}

// IsRepositoryAlreadyTweeted checks if a project was already tweeted.
// If it is not available
//	a) the project was not tweeted yet
//	b) the project ttl expired and is ready to tweet again
func (mc *MemoryConnection) IsRepositoryAlreadyTweeted(projectName string) (bool, error) {
	if val, ok := mc.storage[projectName]; ok {
		if res := val.Before(time.Now()); res {
			delete(mc.storage, projectName)
			return false, nil
		}

		return true, nil
	}

	return false, nil
}
