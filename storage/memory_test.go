package storage

import (
	"testing"
)

var (
	testProjectName = "andygrunwald/TrendingGithub"
)

func TestMemory_MarkRepositoryAsTweeted(t *testing.T) {
	storage := MemoryStorage{}
	pool := storage.NewPool("", "")
	conn := pool.Get()

	res, err := conn.MarkRepositoryAsTweeted(testProjectName, "1440946305")
	if err != nil {
		t.Fatalf("Error of marking repository: \"%s\"", err)
	}

	if res == false {
		t.Fatal("Marking repositoriy failed, got false, expected true")
	}
}

func TestMemory_IsRepositoryAlreadyTweeted(t *testing.T) {
	storage := MemoryStorage{}
	pool := storage.NewPool("", "")
	conn := pool.Get()

	res, err := conn.IsRepositoryAlreadyTweeted(testProjectName)
	if err != nil {
		t.Fatalf("First already tweeted check throws an error: \"%s\"", err)
	}
	if res == true {
		t.Fatal("Repository was already tweeted, got true, expected false")
	}

	res, err = conn.MarkRepositoryAsTweeted(testProjectName, "1440946884")
	if err != nil {
		t.Fatalf("Error of marking repository: \"%s\"", err)
	}

	if res == false {
		t.Fatal("Marking repositoriy failed, got false, expected true")
	}

	res, err = conn.IsRepositoryAlreadyTweeted(testProjectName)
	if err != nil {
		t.Fatalf("Second already tweeted check throws an error: \"%s\"", err)
	}
	if res == false {
		t.Fatal("Repository was not already tweeted, got false, expected true")
	}
}
