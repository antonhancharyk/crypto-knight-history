package utilities

import (
	"time"
)

func FromUnixToUTC(msec int64) time.Time {
	return time.Unix(msec/1000, (msec%1000)*1000000)
}

func SleepUntilNextHour() {
	now := time.Now().UTC()
	next := now.Truncate(time.Hour).Add(time.Hour)
	time.Sleep(time.Until(next))
}
