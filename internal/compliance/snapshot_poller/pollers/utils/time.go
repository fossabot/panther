package utils

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
