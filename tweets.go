package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	githubwrap "github.com/andygrunwald/TrendingGithub/github"
	"github.com/andygrunwald/TrendingGithub/storage"
	trendingwrap "github.com/andygrunwald/TrendingGithub/trending"
	"github.com/andygrunwald/go-trending"
	"github.com/google/go-github/github"
	"github.com/andygrunwald/TrendingGithub/twitter"
)

type TweetSearch struct {
	Channel   chan *Tweet
	Trending  *trendingwrap.Trend
	Storage   storage.Pool
	URLLength int
}

// Tweet is a structure to store the tweet and the project name based on the tweet.
type Tweet struct {
	Tweet       string
	ProjectName string
}

// GenerateNewTweet is responsible to search a new project / repository
// and build a new tweet based on this.
// The generated tweet will be sent to tweetChan.
func (ts *TweetSearch) GenerateNewTweet() {
	var projectToTweet trending.Project

	// Get timeframes and randomize them
	timeFrames := ts.Trending.GetTimeFrames()
	ShuffleStringSlice(timeFrames)

	// First get the timeframes without any languages
	projectToTweet = ts.TimeframeLoopToSearchAProject(timeFrames, "")

	// Check if we found a project. If yes tweet it.
	if ts.IsProjectEmpty(projectToTweet) == false {
		ts.SendProject(projectToTweet)
		return
	}

	// If not, keep going and try to get some (trending) languages
	languages := ts.Trending.GetTrendingLanguages()
	ShuffleStringSlice(languages)
	ShuffleStringSlice(timeFrames)

	for _, language := range languages {
		projectToTweet = ts.TimeframeLoopToSearchAProject(timeFrames, language)

		// If we found a project, break this loop again.
		if ts.IsProjectEmpty(projectToTweet) == false {
			ts.SendProject(projectToTweet)
			break
		}
	}
}

// timeframeLoopToSearchAProject provides basicly a loop over incoming timeFrames (+ language)
// to try to find a new tweet.
// You can say that this is nearly the <3 of this bot.
func (ts *TweetSearch) TimeframeLoopToSearchAProject(timeFrames []string, language string) trending.Project {
	var projectToTweet trending.Project

	for _, timeFrame := range timeFrames {
		if len(language) > 0 {
			log.Printf("Getting projects for timeframe %s and language %s", timeFrame, language)
		} else {
			log.Printf("Getting projects for timeframe %s", timeFrame)
		}

		getProject := ts.Trending.GetRandomProjectGenerator(timeFrame, language)
		projectToTweet = ts.FindProjectWithRandomProjectGenerator(getProject)

		// Check if we found a project.
		// If yes we can leave the loop and keep on rockin
		if ts.IsProjectEmpty(projectToTweet) == false {
			break
		}
	}

	return projectToTweet
}

// sendProject puts the project we want to tweet into the tweet queue
// If the queue is ready to receive a new project, this will be tweeted
func (ts *TweetSearch) SendProject(p trending.Project) {
	text := ""
	// Only build tweet if necessary
	if len(p.Name) > 0 {
		// This is a really hack here ...
		// We have to abstract this a little bit.
		// Eieieieiei
		repository, _ := githubwrap.GetRepositoryDetails(p.Owner, p.RepositoryName)

		text = ts.BuildTweet(p, repository)
	}

	tweet := &Tweet{
		Tweet:       text,
		ProjectName: p.Name,
	}
	ts.Channel <- tweet
}

// isProjectEmpty checks if the incoming project is empty
func (ts *TweetSearch) IsProjectEmpty(p trending.Project) bool {
	if len(p.Name) > 0 {
		return false
	}

	return true
}

// FindProjectWithRandomProjectGenerator retrieves a new project and checks if this was already tweeted.
func (ts *TweetSearch) FindProjectWithRandomProjectGenerator(getProject func() (trending.Project, error)) trending.Project {
	var projectToTweet trending.Project
	var project trending.Project
	var projectErr error

	storageConn := ts.Storage.Get()
	defer storageConn.Close()

	// TODO Lets throw an error, when we dont get a project at all
	// This happened in the past and the bot tweeted nothing.

	for project, projectErr = getProject(); projectErr == nil; project, projectErr = getProject() {
		// Check if the project was already tweeted
		alreadyTweeted, err := storageConn.IsRepositoryAlreadyTweeted(project.Name)
		if err != nil {
			log.Println(err)
			continue
		}

		// If the project was already tweeted
		// we will skip this project and go to the next one
		if alreadyTweeted == true {
			continue
		}

		// This project wasn`t tweeted yet, so we will take over this job
		projectToTweet = project
		break
	}

	return projectToTweet
}

// buildTweet is responsible to build a 140 length string based on the project we found.
func (ts *TweetSearch) BuildTweet(p trending.Project, repo *github.Repository) string {
	tweet := ""
	// Base length of a tweet
	tweetLen := 140

	// Number of letters for the url (+ 1 character for a whitespace)
	// As URL shortener t.co from twitter is used
	// URLLength will be constantly refreshed
	tweetLen -= ts.URLLength + 1

	// Sometimes the owner name is the same as the repository name
	// Like FreeCodeCamp / FreeCodeCamp, docker/docker or flarum / flarum
	// In such cases we will drop the owner name and just use the repository name.
	usedName := p.Name
	if p.Owner == p.RepositoryName {
		usedName = p.RepositoryName
	}

	// Check if the length of the project name is bigger than the space in the tweet
	// Max length of a project name on github is 100 chars
	if nameLen := len(usedName); nameLen < tweetLen {
		tweetLen -= nameLen
		tweet += usedName
	}

	// We only post a description if we got more than 20 charactes available
	// We have to add 2 chars more, because of the prefix ": "
	if tweetLen > 22 && len(p.Description) > 0 {
		tweetLen -= 2
		tweet += ": "

		projectDescription := ""
		if len(p.Description) < tweetLen {
			projectDescription = p.Description
		} else {
			projectDescription = Crop(p.Description, (tweetLen - 4), "...", true)
		}

		tweetLen -= len(projectDescription)
		tweet += projectDescription
	}

	stars := strconv.Itoa(*repo.StargazersCount)
	if starsLen := len(stars) + 2; tweetLen >= starsLen {
		tweet += " ★" + stars
		tweetLen -= starsLen
	}

	// Lets add the URL, but we don`t need to substract the chars
	// because we have done this before
	if p.URL != nil {
		tweet += " "
		tweet += p.URL.String()
	}

	// Lets check if we got space left to add the language as hashtag
	language := strings.Replace(p.Language, " ", "", -1)
	// len + 2, because of " #" in front of the hashtag
	hashTagLen := (len(language) + 2)
	if len(language) > 0 && tweetLen >= hashTagLen {
		tweet += " #" + language
		tweetLen -= hashTagLen
	}

	return tweet
}

// markTweetAsAlreadyTweeted adds a projectName to the global blacklist of already tweeted projects.
// For this we use a Sorted Set where the score is the timestamp of the tweet.
func (ts *TweetSearch) MarkTweetAsAlreadyTweeted(projectName string) (bool, error) {
	storageConn := ts.Storage.Get()
	defer storageConn.Close()

	// Generate score in format YYYYMMDDHHiiss
	now := time.Now()
	score := now.Format("20060102150405")

	res, err := storageConn.MarkRepositoryAsTweeted(projectName, score)
	if err != nil || res != true {
		log.Printf("Error during adding project %s to tweeted list: %s (%v)", projectName, err, res)
	}

	return res, err
}

// StartTweeting bundles the main logic of this bot.
// It schedules the times when we are looking for a new project to tweet.
// If we found a project, we will build the tweet and tweet it to our followers.
// Because we love our followers ;)
func StartTweeting(twitter *twitter.Twitter, storageBackend storage.Pool) {

	// Setup tweet scheduling
	ts := &TweetSearch{
		Channel:   make(chan *Tweet),
		Trending:  trendingwrap.NewTrendingClient(),
		Storage:   storageBackend,
		URLLength: twitter.Configuration.ShortUrlLengthHttps,
	}
	SetupRegularTweetSearchProcess(ts)

	// Waiting for tweets ...
	for tweet := range ts.Channel {
		// Sometimes it happens that we won`t get a project.
		// In this situation we try to avoid empty tweets like ...
		//	* https://twitter.com/TrendingGithub/status/628714326564696064
		//	* https://twitter.com/TrendingGithub/status/628530032361795584
		//	* https://twitter.com/TrendingGithub/status/628348405790711808
		// we will return here
		// We do this check here and not in tweets.go, because otherwise
		// a new tweet won`t be scheduled
		if len(tweet.ProjectName) <= 0 {
			log.Print("No project found. No tweet sent.")
			continue
		}

		// In debug mode the twitter variable is not available, so we won`t tweet the tweet.
		// We will just output them.
		// This is a good development feature ;)
		if twitter.API == nil {
			log.Printf("Tweet: %s (length: %d)", tweet.Tweet, len(tweet.Tweet))

		} else {
			postedTweet, err := twitter.Tweet(tweet.Tweet)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Tweet %s posted", postedTweet.IdStr)
			}
		}
		ts.MarkTweetAsAlreadyTweeted(tweet.ProjectName)
	}
}

func SetupRegularTweetSearchProcess(tweetSearch *TweetSearch) {
	go func() {
		for _ = range time.Tick(tweetTime) {
			go tweetSearch.GenerateNewTweet()
		}
	}()
}

// ShuffleStringSlice will randomize a string slice.
// I know that is a really bad shuffle logic (i won`t call this an algorithm, why? because i wrote and understand it :D)
// But this is YOUR chance to contribute to an open source project.
// Replace this by a cool one!
func ShuffleStringSlice(a []string) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

// Crop is a modified "sub string" function allowing to limit a string length to a certain number of chars (from either start or end of string) and having a pre/postfix applied if the string really was cropped.
// content is the string to perform the operation on
// chars is the max number of chars of the string. Negative value means cropping from end of string.
// afterstring is the pre/postfix string to apply if cropping occurs.
// crop2space is true, then crop will be applied at nearest space. False otherwise.
//
// This function is a port from the TYPO3 CMS (written in PHP)
// @link https://github.com/TYPO3/TYPO3.CMS/blob/aae88a565bdbbb69032692f2d20da5f24d285cdc/typo3/sysext/frontend/Classes/ContentObject/ContentObjectRenderer.php#L4065
func Crop(content string, chars int, afterstring string, crop2space bool) string {
	if chars == 0 {
		return content
	}

	if len(content) < chars || (chars < 0 && len(content) < (chars*-1)) {
		return content
	}

	var cropedContent string
	truncatePosition := -1

	if chars < 0 {
		cropedContent = content[len(content)+chars:]
		if crop2space == true {
			truncatePosition = strings.Index(cropedContent, " ")
		}
		if truncatePosition >= 0 {
			cropedContent = cropedContent[truncatePosition+1:]
		}
		cropedContent = afterstring + cropedContent

	} else {
		cropedContent = content[:chars-1]
		if crop2space == true {
			truncatePosition = strings.LastIndex(cropedContent, " ")
		}
		if truncatePosition >= 0 {
			cropedContent = cropedContent[0:truncatePosition]
		}
		cropedContent += afterstring
	}

	return cropedContent
}
