package config

import (
	"os"
	"strconv"
)

func defaultValue(fromEnv string, defaultIfMissing string) string {
	if fromEnv == "" {
		return defaultIfMissing
	}
	return fromEnv
}

// Given a prefix, e.g. PORT_SPEC, iterate through values starting at
// PORT_SPEC_0, then PORT_SPEC_1, etc. Appends each non-empty value to
// the slice until it finds an empty one, then returns.
func configSlice(prefix string) []string {
	var result []string
	index := 0
	for {
		currentValue := os.Getenv(prefix + "_" + strconv.Itoa(index))
		if currentValue == "" {
			break
		} else {
			result = append(result, currentValue)
			index += 1
		}
	}
	return result
}

var LOG_TO = os.Getenv("DOG_DEVOLVER_LOG_TO")

var LISTEN_IP = defaultValue(os.Getenv("DOG_DEVOLVER_LISTEN_IP"), "127.0.0.1")
var LISTEN_PORT = defaultValue(os.Getenv("DOG_DEVOLVER_LISTEN_PORT"), "9000")

// Specify a single DogStatsD target to which data will be proxied unchanged
var DOGSTATSD_HOST = defaultValue(os.Getenv("DOG_DEVOLVER_DOGSTATSD_HOST"), "127.0.0.1")
var DOGSTATSD_PORT = defaultValue(os.Getenv("DOG_DEVOLVER_DOGSTATSD_PORT"), "8135")

// Specify multiple StatsD targets to which "devolved" DogStatsD data will be proxied
// These are specified as environment variable with base-0 integer suffixes, e.g.
// DOG_DEVOLVER_STATSD_HOST_0=127.0.0.1
// DOG_DEVOLVER_STATSD_PORT_0=8135
// DOG_DEVOLVER_STATSD_HOST_1=statsd.gc.com
// DOG_DEVOLVER_STATSD_PORT_1=8135
var STATSD_HOSTS = configSlice("DOG_DEVOLVER_STATSD_HOST")
var STATSD_PORTS = configSlice("DOG_DEVOLVER_STATSD_PORT")
