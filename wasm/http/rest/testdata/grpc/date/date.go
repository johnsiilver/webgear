package date

import (
	"time"
)

// WeekOf takes the current time and goes backwards till it finds a Sunday, then returns that date.
func WeekOf(t time.Time) time.Time {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	for {
		if t.Weekday() != time.Sunday {
			t = t.Add(-24 * time.Hour)
			continue
		}
		break
	}
	return t
}

// SafeUnixNano calls WeekOf() and then zeros out the hour, minute, second, nsec fields 
// and returns the UnixNano() value.
func SafeUnixNano(t time.Time) int64 {
	t = WeekOf(t)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).UnixNano()
}