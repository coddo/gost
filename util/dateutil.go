package util

import (
	"time"
)

func CompareDates(source, target time.Time) bool {
	date1 := source.Truncate(time.Millisecond)
	date2 := target.Truncate(time.Millisecond)

	return date1.Equal(date2)
}
