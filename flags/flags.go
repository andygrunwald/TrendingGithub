package flags

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

// Bool registers a flag and returns the pointer to the resulting boolean.
// The default value is passed as fallback and env sets the env variable
// that can override the default.
func Bool(name, env string, fallback bool, help string) *bool {
	value := fallback
	if v := os.Getenv(env); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			value = b
		}
	}

	return flag.Bool(name, value, help)
}

// String registers a flag and returns the pointer to the resulting string.
// The default value is passed as fallback and env sets the env variable
// that can override the default.
func String(name, env, fallback, help string) *string {
	value := fallback
	if v := os.Getenv(env); v != "" {
		value = v
	}

	return flag.String(name, value, help)
}

// Duration registers a flag and returns the pointer to the resulting duration.
// The default value is passed as fallback and env sets the env variable
// that can override the default.
func Duration(name, env string, fallback time.Duration, help string) *time.Duration {
	value := fallback
	if v := os.Getenv(env); v != "" {
		vv, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("Error parsing duration from env variable %s: %s", env, v)
		}

		value = vv
	}

	return flag.Duration(name, value, help)
}
