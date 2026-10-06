package utils

import (
	"fmt"
	"time"
)

// ParseDate parses a string in "2006-01-02" format. An empty string yields a
// zero time.Time.
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %v", err)
	}

	return t, nil
}

// ParseTime parses a string in "15:04:05" format. An empty string yields a
// zero time.Time.
func ParseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %v", err)
	}

	return t, nil
}
