package handler

import (
	"regexp"
	"strconv"
	"time"
)

var digitsRe = regexp.MustCompile("^[0-9\\.]+$")

func isDigit(s string) bool {
	return digitsRe.MatchString(s)
}

func parsePromTime(s string) (time.Time, error) {
	var zero time.Time
	if isDigit(s) {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return zero, err
		}
		return time.Unix(0, int64(f*1000*1000*1000)), nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return zero, err
	}
	return t, nil
}

func parsePromDuration(s string) (time.Duration, error) {
	var zero time.Duration
	if isDigit(s) {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return zero, err
		}
		return time.Duration(int64(f*1000*1000*1000)), nil
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return zero, err
	}
	return d, nil
}
