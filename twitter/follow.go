package twitter

import (
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

const (
	// twitterStatusNone reflects the status "none" of GET friendships/lookup call
	// See https://dev.twitter.com/rest/reference/get/friendships/lookup
	twitterStatusNone = "none"
)

// SetupFollowNewPeopleScheduling sets up a scheduler and will search
// for a new person to follow every duration d. This is our growth hack.
func (client *Client) SetupFollowNewPeopleScheduling(d time.Duration) {
	go func() {
		for range time.Tick(d) {
			client.FollowNewPerson()
		}
	}()
	log.Printf("Growth hack: Enabled ✅  (every %s)\n", d.String())
}

// FollowNewPerson will follow a new person on twitter to raise the attraction for the bot.
// We will follow a new person who follow on random follower of @TrendingGithub
// Only persons who don`t have a relationship to the bot will be followed.
func (client *Client) FollowNewPerson() error {
	// Get all own followers
	c, err := client.API.GetFollowersIds(nil)
	if err != nil {
		return err
	}

	// We loop here, because we want to follow one person.
	// If we choose a random person it can be that we choose a person
	// that follows @TrendingGithub already.
	// We want to attract new persons ;)
	for {
		// Choose a random follower
		randomNumber := rand.Intn(len(c.Ids))

		// Request the follower from the random follower
		v := url.Values{}
		v.Add("user_id", strconv.FormatInt(c.Ids[randomNumber], 10))
		c, err = client.API.GetFollowersIds(v)
		if err != nil {
			return err
		}

		// Choose a random follower (again) from the random follower chosen before
		randomNumber = rand.Intn(len(c.Ids))

		// Get friendship details of @TrendingGithub and the chosen person
		v = url.Values{}
		v.Add("user_id", strconv.FormatInt(c.Ids[randomNumber], 10))
		friendships, err := client.API.GetFriendshipsLookup(v)
		if err != nil {
			return err
		}

		// Test if @TrendingGithub has a relationship to the new user.
		shouldIFollow := client.isThereARelationship(friendships)

		// If we got a relationship, we will repeat the process ...
		if !shouldIFollow {
			continue
		}

		// ... if not we will follow the new person
		// We drop the error and user here, because we got no logging yet ;)
		client.API.FollowUserId(c.Ids[randomNumber], nil)
		return nil
	}
}

// isThereARelationship will test if @TrendingGithub has a relationship to the new user.
// Only if there is no relationship ("none").
// Default wise we assume that we got a relationship already
// See https://dev.twitter.com/rest/reference/get/friendships/lookup
func (client *Client) isThereARelationship(friendships []anaconda.Friendship) bool {
	shouldIFollow := false
	for _, friend := range friendships {
		for _, status := range friend.Connections {
			if status == twitterStatusNone {
				shouldIFollow = true
				break
			}
		}
	}

	return shouldIFollow
}
