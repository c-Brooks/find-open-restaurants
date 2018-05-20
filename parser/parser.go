package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/gogap/logrus"

	"github.com/c-Brooks/find-open-restaurants/hours"
)

const (
	// ruleSeparatorToken separates rules into multiple sub-rules
	// note the whitespace, it is important
	ruleSeparatorToken = " / "

	// dayRangeCharacter denotes a range of days, e.g. Mon-Fri
	dayRangeCharacter = "-"
)

// We're using an array, not a map, because order is important for our case.
// Lucky us, time.Weekday has an integer value, and so can be used as an array index.
var weekdayLUT = [hours.DaysInWeek]string{
	time.Sunday:    "Sun",
	time.Monday:    "Mon",
	time.Tuesday:   "Tue",
	time.Wednesday: "Wed",
	time.Thursday:  "Thu",
	time.Friday:    "Fri",
	time.Saturday:  "Sat",
}

func findDay(day string) (time.Weekday, error) {
	for i, weekday := range weekdayLUT {
		if weekday == day {
			return time.Weekday(i), nil
		}
	}

	return 0, fmt.Errorf("invalid weekday: %s", day)
}

// apply these rules:
// 1) Mon-Wed -> Mon, Tue, Wed
// 2) Mon-Wed, Sat -> Mon, Tue, Wed, Sat
func normalizeDays(dayToken string) ([]time.Weekday, error) {
	normalizedDays := make([]time.Weekday, 0)
	days := strings.Split(dayToken, ", ")

	for _, day := range days {
		if strings.Contains(day, dayRangeCharacter) {
			// expand
			startDay, err := findDay(strings.Split(day, dayRangeCharacter)[0])
			endDay, err := findDay(strings.Split(day, dayRangeCharacter)[1])
			if err != nil {
				return nil, err
			}
			// go thru LUT (wrap around) and make an array of days from range (Mon-Wed)
			for i := startDay; ; i = ((i + 1) % hours.DaysInWeek) {
				normalizedDays = append(normalizedDays, time.Weekday(i))
				if i == endDay {
					break
				}
			}
		} else {
			// it's just one day, not a range
			i, err := findDay(day)
			if err != nil {
				return nil, err
			}
			normalizedDays = append(normalizedDays, time.Weekday(i))
		}
	}
	return normalizedDays, nil
}

// apply these rules:
// 1) 9 pm -> 09:00pm
// 2) 9:30 pm -> 09:30pm
func normalizeTime(time, period string) string {
	hasSingleDigitHour := len(strings.Split(time, ":")[0]) == 1
	hasNoMinutes := !strings.Contains(time, ":")

	if hasSingleDigitHour {
		time = "0" + time
	}
	if hasNoMinutes {
		time = time + ":00"
	}
	return time + period
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 // not found
}

// Turns a start and end time (strings) into a hours.TimeRange
func parseTimeTokens(startTimeStr, endTimeStr string) hours.TimeRange {
	startTime, err := time.Parse("03:04pm", startTimeStr)
	if err != nil {
		logrus.Fatalf("could not parse startTime %s: %v", startTimeStr, err)
	}
	endTime, err := time.Parse("03:04pm", endTimeStr)
	if err != nil {
		logrus.Fatalf("could not parse endTime %s: %v", endTimeStr, err)
	}

	return hours.TimeRange{hours.TimeOfDay{startTime}, hours.TimeOfDay{endTime}}
}

// HoursFromString parses the CSV rows into Hours
// rows have this format:
//
// Mon-Thu, Sun 11:30 am - 10 pm  / Sat 5:30 pm - 10 pm
//
// tokens:
// day(s) - can be one day or a range (if there's a hyphen separating), or a mix
// timeRange - always two times of day, separated by spaces and a hyphen
// slash - denotes different rules separator
func HoursFromString(row string) (hours.Hours, error) {
	hrs := make(hours.Hours, len(row))

	// handle case when there are multiple rules
	// call itself with each row and return combined Hours
	if strings.Contains(row, ruleSeparatorToken) {
		rules := strings.Split(row, ruleSeparatorToken)

		for _, rule := range rules {
			hoursForRule, err := HoursFromString(rule)
			if err != nil {
				return nil, err
			}

			hrs.Combine(hoursForRule)
		}
		return hrs, nil
	}

	fields := strings.Fields(row)
	timeSeparator := indexOf(dayRangeCharacter, fields)
	startTimeStr := normalizeTime(fields[timeSeparator-2], fields[timeSeparator-1])
	endTimeStr := normalizeTime(fields[timeSeparator+1], fields[timeSeparator+2])

	dayToken := fields[:timeSeparator-2]
	t := parseTimeTokens(startTimeStr, endTimeStr)

	days, err := normalizeDays(strings.Join(dayToken, " "))
	if err != nil {
		return nil, err
	}

	for _, wkday := range days {
		hrs[wkday] = t
	}
	return hrs, nil
}
