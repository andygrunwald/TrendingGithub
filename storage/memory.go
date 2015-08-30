package storage

import (
	"time"
)

type MemoryStorageContainer map[string]time.Time

type MemoryStorage struct{}

type MemoryPool struct {
	storage map[string]time.Time
}

type MemoryConnection struct {
	storage map[string]time.Time
}

func (ms *MemoryStorage) NewPool(url, auth string) Pool {
	return MemoryPool{
		storage: make(map[string]time.Time),
	}
}

func (mp MemoryPool) Close() error {
	return nil
}

func (mp MemoryPool) Get() Connection {
	return &MemoryConnection{
		storage: mp.storage,
	}
}

func (mc *MemoryConnection) Close() error {
	return nil
}

func (mc *MemoryConnection) MarkRepositoryAsTweeted(projectName, score string) (bool, error) {
	// Add greylisting to current time
	now := time.Now()
	future := now.Add(time.Second * GreyListTTL)

	mc.storage[projectName] = future

	return true, nil
}

func (mc *MemoryConnection) IsRepositoryAlreadyTweeted(projectName string) (bool, error) {
	if val, ok := mc.storage[projectName]; ok {
		if res := val.Before(time.Now()); res == true {
			delete(mc.storage, projectName)
			return false, nil
		}

		return true, nil
	}

	return false, nil
}
