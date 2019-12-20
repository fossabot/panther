package timestamp

import (
	"time"

	"github.com/panther-labs/panther/tools/cfngen/gluecf"
)

// These objects are used to read timestamps and ensure a consistent JSON output for timestamps.

// NOTE: prefix the name of all objects with Timestamp so schema generation can automatically understand these.
// NOTE: the suffix of the names is meant to reflect the time format being read (unmarshal)

// We want our output JSON timestamps to be: YYYY-MM-DD HH:MM:SS.fffffffff
// https://aws.amazon.com/premiumsupport/knowledge-center/query-table-athena-timestamp-empty/
const (
	jsonMarshalLayout = `"2006-01-02 15:04:05.000000000"`

	ansicWithTZUnmarshalLayout = `"Mon Jan 2 15:04:05 2006 MST"` // similar to time.ANSIC but with MST

	glueType = "timestamp" // type in Glue tables for timestamps
)

var (
	// GlueMappings for timestamps. Reference this when generating CF
	GlueMappings = []gluecf.CustomMapping{
		{
			From: "timestamp.RFC3339",
			To:   glueType,
		},
		{
			From: "timestamp.ANSICwithTZ",
			To:   glueType,
		},
	}
)

// use these functions to parse all incoming dates to ensure UTC consistency
func Parse(layout, value string) (RFC3339, error) {
	t, err := time.Parse(layout, value)
	return (RFC3339)(t.UTC()), err
}

func Unix(sec int64, nsec int64) RFC3339 {
	return (RFC3339)(time.Unix(sec, nsec).UTC())
}

type RFC3339 time.Time

func (ts *RFC3339) String() string {
	return (*time.Time)(ts).UTC().String() // ensure UTC
}

func (ts *RFC3339) MarshalJSON() ([]byte, error) {
	return []byte((*time.Time)(ts).UTC().Format(jsonMarshalLayout)), nil // ensure UTC
}

func (ts *RFC3339) UnmarshalJSON(jsonBytes []byte) (err error) {
	return (*time.Time)(ts).UnmarshalJSON(jsonBytes)
}

// Like time.ANSIC but with MST
type ANSICwithTZ time.Time

func (ts *ANSICwithTZ) String() string {
	return (*time.Time)(ts).UTC().String() // ensure UTC
}

func (ts *ANSICwithTZ) MarshalJSON() ([]byte, error) {
	return []byte((*time.Time)(ts).UTC().Format(jsonMarshalLayout)), nil // ensure UTC
}

func (ts *ANSICwithTZ) UnmarshalJSON(text []byte) (err error) {
	t, err := time.Parse(ansicWithTZUnmarshalLayout, string(text))
	if err != nil {
		return
	}
	*ts = (ANSICwithTZ)(t)
	return
}
