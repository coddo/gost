package util

import (
	"time"
)

// Today retrieves the current date and time
func Today() time.Time {
	return time.Now().Local()
}

// NextDate gets the current date and time and then adds the given duration to it,
// then returns the result
func NextDate(duration time.Duration) time.Time {
	return Today().Add(duration)
}

// CompareDates checks whether two dates are identical or note.
// Both the date and the time are compared up to millisecond precision
func CompareDates(source, target time.Time) bool {
	date1 := source.Truncate(time.Millisecond)
	date2 := target.Truncate(time.Millisecond)

	return date1.Equal(date2)
}

// IsDateExpired tells if the current date and time are greater than the given one
func IsDateExpired(date time.Time) bool {
	today := Today()

	return date.Local().Before(today)
}
