package seeder

import (
	"time"
)

func toBoolPtr(b bool) *bool {
	return &b
}

func toInt32Ptr(i int32) *int32 {
	return &i
}

func toFloat64Ptr(f float64) *float64 {
	return &f
}

func toStringPtr(s string) *string {
	return &s
}

func toDate(t time.Time) *time.Time {
	return &t
}

func toTime(t time.Time) *time.Time {
	return &t
}
