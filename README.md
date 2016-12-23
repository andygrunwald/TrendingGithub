# [@TrendingGithub](https://twitter.com/TrendingGithub)

[![Build Status](https://travis-ci.org/andygrunwald/TrendingGithub.svg?branch=master)](https://travis-ci.org/andygrunwald/TrendingGithub)
[![GoDoc](https://godoc.org/github.com/andygrunwald/TrendingGithub?status.svg)](https://godoc.org/github.com/andygrunwald/TrendingGithub)
[![Coverage Status](https://coveralls.io/repos/andygrunwald/TrendingGithub/badge.svg?branch=master&service=github)](https://coveralls.io/github/andygrunwald/TrendingGithub?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/andygrunwald/TrendingGithub)](https://goreportcard.com/report/github.com/andygrunwald/TrendingGithub)

A twitter bot (**[@TrendingGithub](https://twitter.com/TrendingGithub)**) to tweet [trending repositories](https://github.com/trending) and [developers](https://github.com/trending/developers) from GitHub.

> Follow us at **[@TrendingGithub](https://twitter.com/TrendingGithub)**.

[![@TrendingGithub twitter account](./img/TrendingGithub.png "@TrendingGithub twitter account")](https://twitter.com/TrendingGithub)

**Important:** This is no official GitHub or Twitter product.

## Features

* Tweets trending projects every 30 minutes
* Refreshes the configuration of twitters URL shortener t.co every 24 hours
* Greylisting of repositories for 30 days (to avoid tweeting a project multiple times in a short timeframe)
* Maximum use of 140 chars per tweet to fill up with information
* Debug / development mode
* Multiple storage backends (currently [Redis](http://redis.io/) and in memory)

## Installation

1. Download the [latest release](https://github.com/andygrunwald/TrendingGithub/releases/latest)
2. Extract the archive (zip / tar.gz)
3. Start the bot via `./TrendingGithub -debug`

For linux this can look like:

```sh
curl -L  https://github.com/andygrunwald/TrendingGithub/releases/download/v0.3.1/TrendingGithub-v0.3.1-linux-amd64.tar.gz -o TrendingGithub-v0.3.1-linux-amd64.tar.gz
tar xzvf TrendingGithub-v0.3.1-linux-amd64.tar.gz
cd TrendingGithub-v0.3.1-linux-amd64
./TrendingGithub -debug
```

## Usage

```
$ ./TrendingGithub -help
Usage of ./TrendingGithub:
  -debug
    	Outputs the tweet instead of tweet it (useful for development). Default: false. Env var: TRENDINGGITHUB_DEBUG
  -storage-auth string
    	Storage Auth (e.g. myPassword or <empty>). Default: empty.  Env var: TRENDINGGITHUB_STORAGE_AUTH
  -storage-url string
    	Storage URL (e.g. 1.2.3.4:6379 or :6379). Default: :6379.  Env var: TRENDINGGITHUB_STORAGE_URL (default ":6379")
  -twitter-access-token string
    	Twitter-API: Access token. Default: empty. Env var: TRENDINGGITHUB_TWITTER_ACCESS_TOKEN
  -twitter-access-token-secret string
    	Twitter-API: Access token secret. Default: empty. Env var: TRENDINGGITHUB_TWITTER_ACCESS_TOKEN_SECRET
  -twitter-consumer-key string
    	Twitter-API: Consumer key. Default: empty. Env var: TRENDINGGITHUB_TWITTER_CONSUMER_KEY
  -twitter-consumer-secret string
    	Twitter-API: Consumer secret. Default: empty. Env var: TRENDINGGITHUB_TWITTER_CONSUMER_SECRET
  -twitter-follow-new-person
    	Twitter: Follows a friend of one of our followers. Default: false. Env var: TRENDINGGITHUB_TWITTER_FOLLOW_NEW_PERSON
  -version
    	Outputs the version number and exit. Default: false. Env var: TRENDINGGITHUB_VERSION
```

**Every parameter can be set by environment variable as well.**

**Twitter-API settings** (`twitter-access-token`, `twitter-access-token-secret`, `twitter-consumer-key` and `twitter-consumer-secret`) are necessary to use the Twitter API and to set up a tweet by your application.
You can get those settings by [Twitter's application management](https://apps.twitter.com/).

If you want to play around or develop this bot, use the `debug` setting.
It avoids using the Twitter API for tweet purposes and outputs the tweet on stdout.

The Redis url (`storage-url`)is the address of the Redis server in format *ip:port* (e.g. *192.168.0.12:6379*).
If your server is running on localhost you can use *:6379* as a shortcut.
`storage-auth` is the authentication string necessary for your Redis server if you use the [Authentication feature](http://redis.io/topics/security#authentication-feature).

## Storage backends

Why is a storage backend needed at all?

We are looking for popular projects in a regular interval.
To avoid tweeting a project or developer multiple times after another we add those records to a blacklist for a specific time.

At the moment there are two backends implemented:

* Memory (used in development)
* Redis (used in production)

## Growth hack

We implemented a small growth hack to get a few followers.
This hack was suggested by my colleague [@mre](https://github.com/mre).
It works like described:

* Get all followers from [@TrendingGithub](https://twitter.com/TrendingGithub)
* Choose a random one and get the followers of the choosen person
* Check if this person follows us already
* If yes, repeat
* If no, follow this person

This feature can be activated via the `twitter-follow-new-person` flag.

## Motivation

I love to discover new tools, new projects, new languages, new coding best practices, new exciting ideas and new people who share the same passion like me.
[I use twitter a lot](https://twitter.com/andygrunwald) and have little time to check [trending repositories](https://github.com/trending) and [developers](https://github.com/trending/developers) on a daily basis.

Why not combine both to save time and spread favorite projects and developers via tweets?

## TODO

* Code: Extend logging
* Code: Add expvar support
* Code: Documentation

## License

This project is released under the terms of the [MIT license](http://en.wikipedia.org/wiki/MIT_License).
