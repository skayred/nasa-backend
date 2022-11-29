package utils

import (
	"time"
)

func DaysBetweenDates(from time.Time, to time.Time) []time.Time {
	results := []time.Time{}

	// Start - beginning of the start day
	start := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)

	// end - beginning of the next day
	end := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	end = end.AddDate(0, 0, 1)

	for day := start; day.After(end) == false; day = day.AddDate(0, 0, 1) {
		results = append(results, day)
	}

	return results
}
