package hours

import (
	"time"
)

const (
	OneDayInSeconds = 60 * 60 * 24
	DaysInWeek      = 7
)

// Hours is map of weekday (Sunday = 0, Sat = 6) to timeRange
// note that if a weekday is not present, TODO
type Hours map[time.Weekday]TimeRange

// Combine combines one Hours map with another, overwriting the
// method receiver if there are conflicts
func (h Hours) Combine(hours Hours) {
	for day, timerange := range hours {
		h[day] = timerange
	}
}

// TimeRange is the range of hours a restaurant is open for
type TimeRange struct {
	StartTime TimeOfDay
	EndTime   TimeOfDay
}

// TimeOfDay is a wrapper around time.Time that receives the secondsOfDay method
type TimeOfDay struct {
	time.Time
}

// secondsOfDay returns the number of seconds since 12am
func (tod TimeOfDay) secondsOfDay() int {
	return tod.Hour()*60*60 + tod.Minute()*60 + tod.Second()
}

// Contains checks if a time is contained in a timeRange
//
// take the time of day (in seconds) from time and compare it
// with the times in the timeRange
func (tr TimeRange) Contains(tod TimeOfDay) bool {
	startTimeSeconds := tr.StartTime.secondsOfDay()
	endTimeSeconds := tr.EndTime.secondsOfDay()

	// Special case where a restaurant is open until later than midnight
	// subtract a day from the start-time
	if endTimeSeconds < startTimeSeconds {
		startTimeSeconds -= OneDayInSeconds
	}

	tSeconds := tod.secondsOfDay()
	isAfterStartTime := tSeconds > startTimeSeconds
	isBeforeEndTime := tSeconds < endTimeSeconds
	return isAfterStartTime && isBeforeEndTime
}
