package main

import (
	"flag"
)

var (
	flagConfigFile = flag.String("config", "", "Configuration file")
	flagPidFile    = flag.String("pidfile", "", "Write the process id into a given file")
	flagVersion    = flag.Bool("version", false, "Outputs the version number and exits")
	flagDebug      = flag.Bool("debug", false, "Outputs the tweet instead of tweet it")
)

func init() {
	flag.Parse()
}
