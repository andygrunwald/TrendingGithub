package main

import (
	"flag"
)

var (
	flagConfigFile = flag.String("config", "", "Path to configuration file.")
	flagVersion    = flag.Bool("version", false, "Outputs the version number and exits.")
	flagDebug      = flag.Bool("debug", false, "Outputs the tweet instead of tweet it. Useful for development.")
)

func init() {
	flag.Parse()
}
