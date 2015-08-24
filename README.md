# TrendingGithub

[![Build Status](https://travis-ci.org/andygrunwald/TrendingGithub.svg?branch=master)](https://travis-ci.org/andygrunwald/TrendingGithub)

A twitter bot (**[@TrendingGithub](https://twitter.com/TrendingGithub)**) to tweet [trending repositories](https://github.com/trending) and [developers](https://github.com/trending/developers) from GitHub.
Follow us at **[@TrendingGithub](https://twitter.com/TrendingGithub)**.

[![@TrendingGithub twitter account](./img/TrendingGithub.png "@TrendingGithub twitter account")](https://twitter.com/TrendingGithub)

**Important note:** This is no official GitHub or Twitter product.

## Motivation

I love to discover new tools, new projects, new languages, new people who share the same passion like me, new coding best practices, new exciting ideas.
And [I use twitter a lot](https://twitter.com/andygrunwald) and have little time to check [trending repositories](https://github.com/trending) and [developers](https://github.com/trending/developers) on a daily basis.

Why not combine both to save time, favorite projects and developers and spread them by retweets?

## Installation

TODO

## Usage

```
$ TrendingGithub --help

Usage of ./TrendingGithub:
  --config="": Path to configuration file.
  --debug: Outputs the tweet instead of tweet it. Useful for development.
  --pidfile="": Write the process id into a given file.
  --version: Outputs the version number and exits.
```

The **--config** parameter is required. 
See [Configuration chapter](https://github.com/andygrunwald/TrendingGithub#configuration) for details.

**--debug** is quite useful for development.
It doesn`t output special information.
It only avoids using the Twitter API for tweet purposes and outputs the tweet on stdout.

**--pidfile** is useful to run this process in production and observe it with a monitoring tool like [Nagios](https://www.nagios.org/), [Icinga](https://www.icinga.org/), [Monit](https://mmonit.com/monit/), [supervisord](http://supervisord.org/) or similar.

## Configuration

The configuration is based on a JSON file.
You can use the [config.json.dist](./config.json.dist) file as base.

### Twitter

```
"twitter": {
    "consumer-key": "",
    "consumer-secret": "",
    "access-token": "",
    "access-token-secret": ""
  },
```

All these settings mentioned above are necessary to use the Twitter API and to set up a tweet by your application.
You can get those settings by [Twitter's application management](https://apps.twitter.com/).

### Redis

```
"redis": {
    "url": ":6379",
    "auth": ""
  }
```

*redis* contains the connection details to the [Redis](http://redis.io/) server.
This Redis server is used for blacklisting projects that were already tweeted.

**url** is the address of the Redis server in format *ip:port*. 
Example: *192.168.0.12:6379*.
If your server is running on localhost you can use *:6379* as a shortcut.

**auth** is the authentication string necessary for your Redis server if you use the [Authentication feature](http://redis.io/topics/security#authentication-feature).

## License

This project is released under the terms of the [MIT license](http://en.wikipedia.org/wiki/MIT_License).
