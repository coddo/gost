package util

import (
	"time"
)

// Now retrieves the current date and time
func Now() time.Time {
	return time.Now().Local()
}

// NextDate uses the selected date and time to add a given duration to it,
// then returns the result
func NextDate(date time.Time, duration time.Duration) time.Time {
	return date.Local().Add(duration)
}

// NextDateFromNow gets the current date and time and then adds the given duration to it,
// then returns the result
func NextDateFromNow(duration time.Duration) time.Time {
	return Now().Add(duration)
}

// CompareDates checks whether two dates are identical or note.
// Both the date and the time are compared up to millisecond precision
func CompareDates(source, target time.Time) bool {
	date1 := source.Local().Truncate(time.Millisecond)
	date2 := target.Local().Truncate(time.Millisecond)

	return date1.Equal(date2)
}

// IsDateExpired tells if the limit date and time are greater than the given one
func IsDateExpired(date, limit time.Time) bool {
	date = date.Local()
	limit = limit.Local()

	return date.Local().Before(limit)
}

// IsDateExpiredFromNow tells if the current date and time are greater than the given one
func IsDateExpiredFromNow(date time.Time) bool {
	today := Now()

	return date.Local().Before(today)
}
