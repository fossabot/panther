package osquerylogs

import (
	"time"
)

const timeLayout = `"Mon Jan 2 15:04:05 2006 MST"`

// Time struct is a wrapper for time.Time
// with custom unmarshalling logic
type Time struct {
	time.Time
}

// UnmarshalJSON unmarshals JSON field using OsQuery time layout
func (t *Time) UnmarshalJSON(text []byte) error {
	parsedTime, err := time.Parse(timeLayout, string(text))
	if err != nil {
		return err
	}
	*t = Time{parsedTime}
	return nil
}
