package utils

import (
	"time"
)

func MillisFromTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func TimeFromMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

func Yesterday() time.Time {
	return time.Now().AddDate(0, 0, -1)
}

// Returns a time.Time representing the next time it'll be the given hour (in 24 hour time). For example,
// passing in t=6:00 today and hour=14:00 would return 14:00 today, but passing in t=16:00 today and
// hour=14:00 would return 14:00 tomorrow.
func NextHour(t time.Time, hour int) time.Time {
	next := time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, t.Location())

	if next.Before(t) {
		next = next.AddDate(0, 0, 1)
	}

	return next
}
