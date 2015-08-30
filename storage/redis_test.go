package storage

import (
	"testing"
)

var (
	redisTestProjectName = "andygrunwald/TrendingGithub"
)

func TestRedis_Get_FailedConnection(t *testing.T) {
	storage := RedisStorage{}
	pool := storage.NewPool("no-host:0", "wrong-auth")
	defer pool.Close()
	conn := pool.Get()
	err := conn.Close()

	if err == nil {
		t.Fatal("No error thrown, but expected one, because no redis server available.")
	}
}

// TestRedis_Get_SuccessConnection will only succeed if there is a
// running redis server on localhost:6379
// At travis this is the case.
// @link http://docs.travis-ci.com/user/database-setup/#Redis
func TestRedis_Get_SuccessConnection(t *testing.T) {
	storage := RedisStorage{}
	pool := storage.NewPool("localhost:6379", "")
	defer pool.Close()
	conn := pool.Get()
	err := conn.Close()

	if err != nil {
		t.Fatalf("An error occured, but no one was expected: %s", err)
	}
}

func TestRedis_MarkRepositoryAsTweeted(t *testing.T) {
	storage := RedisStorage{}
	pool := storage.NewPool("localhost:6379", "")
	defer pool.Close()
	conn := pool.Get()
	defer conn.Close()

	res, err := conn.MarkRepositoryAsTweeted(redisTestProjectName, "1440946305")
	if err != nil {
		t.Fatalf("Error of marking repository: \"%s\"", err)
	}

	if res == false {
		t.Fatal("Marking repositoriy failed, got false, expected true")
	}
}

func TestRedis_IsRepositoryAlreadyTweeted(t *testing.T) {
	testProject := redisTestProjectName + "Foo"

	storage := RedisStorage{}
	pool := storage.NewPool("localhost:6379", "")
	defer pool.Close()
	conn := pool.Get()
	defer conn.Close()

	res, err := conn.IsRepositoryAlreadyTweeted(testProject)
	if err != nil {
		t.Fatalf("First already tweeted check throws an error: \"%s\"", err)
	}
	if res == true {
		t.Fatal("Repository was already tweeted, got true, expected false")
	}

	res, err = conn.MarkRepositoryAsTweeted(testProject, "1440946884")
	if err != nil {
		t.Fatalf("Error of marking repository: \"%s\"", err)
	}

	if res == false {
		t.Fatal("Marking repositoriy failed, got false, expected true")
	}

	res, err = conn.IsRepositoryAlreadyTweeted(testProject)
	if err != nil {
		t.Fatalf("Second already tweeted check throws an error: \"%s\"", err)
	}
	if res == false {
		t.Fatal("Repository was not already tweeted, got false, expected true")
	}
}
