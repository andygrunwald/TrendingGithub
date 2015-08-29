package main

import (
	"github.com/andygrunwald/go-trending"
	"testing"
)

func TestTweets_IsProjectEmpty(t *testing.T) {
	ts := TweetSearch{}
	mock := []struct {
		Project trending.Project
		Result  bool
	}{
		{trending.Project{Name: ""}, true},
		{trending.Project{Name: "MyProject"}, false},
	}

	for _, item := range mock {
		res := ts.IsProjectEmpty(item.Project)
		if res != item.Result {
			t.Errorf("Failed for project \"%s\", got %v, expected %v", item.Project.Name, res, item.Result)
		}
	}
}
