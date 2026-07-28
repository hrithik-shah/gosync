package timeutil

import (
	"regexp"
	"strings"
	"time"
)

var unitMap = map[string]time.Duration{
	"d": 24,
	"D": 24,
	"w": 7 * 24,
	"W": 7 * 24,
	"M": 30 * 24,
	"y": 365 * 24,
	"Y": 365 * 24,
}

// ParseDuration parses a duration string.
// Extends time.ParseDuration to support days, weeks, months, and years as defined in unitMap
func ParseDuration(s string) (time.Duration, error) {
	neg := false
	if len(s) > 0 && s[0] == '-' {
		neg = true
		s = s[1:]
	}

	re := regexp.MustCompile(`(\d*\.\d+|\d+)[^\d]*`)

	strs := re.FindAllString(s, -1)
	var sumDur time.Duration
	for _, str := range strs {
		var _hours time.Duration = 1
		for unit, hours := range unitMap {
			if strings.Contains(str, unit) {
				str = strings.ReplaceAll(str, unit, "h")
				_hours = hours
				break
			}
		}

		dur, err := time.ParseDuration(str)
		if err != nil {
			return 0, err
		}

		sumDur += dur * _hours
	}

	if neg {
		sumDur = -sumDur
	}
	return sumDur, nil
}
