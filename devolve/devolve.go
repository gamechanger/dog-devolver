package devolve

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var dogstatsdEventPattern = regexp.MustCompile(`_e{\d+,\d+}`)

// Convert a DogStatsD formatted string to its equivalent
// in the StatsD protocol. If this is not possible, e.g.
// if it is a DataDog event, we return "".
func Devolve(in string) (string, error) {
	if dogstatsdEventPattern.FindString(in) != "" {
		return "", errors.New("Received DataDog event, not StatsD compatible")
	}

	split := strings.SplitN(in, "|", 3)
	if len(split) == 1 {
		return "", errors.New(fmt.Sprintf("String is not any valid *StatsD format: %s", in))
	}
	nameValue, metricType := split[0], split[1]

	return nameValue + "|" + metricType, nil
}
