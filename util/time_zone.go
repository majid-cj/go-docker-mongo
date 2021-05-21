package util

import (
	"os"
	"time"
)

// GetTimeNow ...
func GetTimeNow() time.Time {
	loc, err := time.LoadLocation(os.Getenv("TIME_ZONE"))
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}

// TimeAfter ...
func TimeAfter(timebefore time.Time, duration time.Duration) time.Time {
	return timebefore.Add(duration)
}
