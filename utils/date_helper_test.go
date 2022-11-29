package utils

import (
	"testing"
	"time"
)

func TestRegularDates(t *testing.T) {
	testDates(
		t,
		time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
		4,
		time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
	)
}

func TestSingleDay(t *testing.T) {
	testDates(
		t,
		time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
		2,
		time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
	)
}

func TestMonthBoundary(t *testing.T) {
	testDates(
		t,
		time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
		5,
		time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 3, 0, 0, 0, 0, time.UTC),
	)
}

func TestSpecialMonthBoundary(t *testing.T) {
	testDates(
		t,
		time.Date(2020, 2, 27, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 3, 2, 0, 0, 0, 0, time.UTC),
		6,
		time.Date(2020, 2, 27, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 3, 3, 0, 0, 0, 0, time.UTC),
	)
}

func TestSingleDayWithHours(t *testing.T) {
	testDates(
		t,
		time.Date(2020, 1, 3, 3, 4, 0, 0, time.UTC),
		time.Date(2020, 1, 3, 3, 4, 0, 0, time.UTC),
		2,
		time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
	)
}

func testDates(t *testing.T, from time.Time, to time.Time, expectedAmount int, expectedStart time.Time, expectedEnd time.Time) {
	days := DaysBetweenDates(from, to)

	if len(days) != expectedAmount {
		t.Fatalf("Expected %d days, got %d days", expectedAmount, len(days))
	}

	expectedStartFormatted := expectedStart.Format(time.RFC3339)
	expectedEndFormatted := expectedEnd.Format(time.RFC3339)

	actualStartFormatted := days[0].Format(time.RFC3339)
	actualEndFormatted := days[len(days)-1].Format(time.RFC3339)

	if expectedStartFormatted != actualStartFormatted {
		t.Fatalf("Expected %s as a first day, got %s", expectedStartFormatted, actualStartFormatted)
	}

	if expectedEndFormatted != actualEndFormatted {
		t.Fatalf("Expected %s as a first day, got %s", expectedEndFormatted, actualEndFormatted)
	}
}
