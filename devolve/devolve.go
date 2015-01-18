package devolve

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var dogstatsdEventPattern = regexp.MustCompile(`_e{\d+,\d+}`)
var dogstatsdMetricPattern = regexp.MustCompile(`(.+?:.+?)\|([cghs]|ms)\|?(@.+?)?(\|.*)?$`)

func removeEmptyStrings(slice []string) []string {
	var result []string
	for _, s := range slice {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

// Convert a DogStatsD formatted string to its equivalent
// in the StatsD protocol. If this is not possible, e.g.
// if it is a DataDog event, we return "".
func Devolve(in string) (string, error) {
	if dogstatsdEventPattern.FindString(in) != "" {
		return "", errors.New("Received DataDog event, not StatsD compatible")
	}

	matches := dogstatsdMetricPattern.FindStringSubmatch(in)
	if len(matches) < 3 {
		return "", errors.New(fmt.Sprintf("String is not any valid *StastsD format: %s", in))
	}
	return strings.Join(removeEmptyStrings(matches[1:4]), "|"), nil

}
