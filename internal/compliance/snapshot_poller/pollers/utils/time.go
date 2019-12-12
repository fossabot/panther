package utils

import (
	"time"

	"github.com/go-openapi/strfmt"
)

var (
	// TimeNowFunc directs to the TimeNow function.
	// This is intended to be overridden for testing.
	TimeNowFunc = TimeNowRFC3339
)

// DateTimeFormat converts time.Time to strfmt.DateTime.
func DateTimeFormat(inputTime time.Time) *strfmt.DateTime {
	conv := strfmt.DateTime(inputTime)
	return &conv
}

// TimeNowRFC3339 returns the current time in RFC3339 format.
func TimeNowRFC3339() time.Time {
	return time.Now()
}

// ParseTimeRFC3339 parses a time string into a valid RFC3339 time.
func ParseTimeRFC3339(timeString string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

// StringToDateTime parses a time string into a strfmt.DateTime struct
func StringToDateTime(timeString string) *strfmt.DateTime {
	return DateTimeFormat(ParseTimeRFC3339(timeString))
}

// UnixTimeToDateTime parses an Int64 representing an epoch timestamp to a strfmt.DateTime struct
func UnixTimeToDateTime(epochTimeStamp int64) *strfmt.DateTime {
	return DateTimeFormat(time.Unix(epochTimeStamp, 0))
}
