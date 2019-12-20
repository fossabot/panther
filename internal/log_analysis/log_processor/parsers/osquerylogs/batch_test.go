package osquerylogs

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/timestamp"
)

func TestBatchLog(t *testing.T) {
	//nolint:lll
	log := `{"diffResults": {"added": [ { "name": "osqueryd", "path": "/usr/local/bin/osqueryd", "pid": "97830" } ],"removed": [ { "name": "osqueryd", "path": "/usr/local/bin/osqueryd", "pid": "97650" } ] },"name": "processes", "hostname": "hostname.local", "calendarTime": "Tue Nov 5 06:08:26 2018 UTC","unixTime": "1412123850", "epoch": "314159265", "counter": "1" }`

	expectedTime := time.Unix(1541398106, 0).UTC()
	expectedEvent := &Batch{
		CalendarTime: (*timestamp.ANSICwithTZ)(&expectedTime),
		Name:         aws.String("processes"),
		Epoch:        aws.Int(314159265),
		Hostname:     aws.String(("hostname.local")),
		UnixTime:     aws.Int(1412123850),
		Counter:      aws.Int(1),
		DiffResults: &BatchDiffResults{
			Added: []map[string]string{
				{
					"name": "osqueryd",
					"path": "/usr/local/bin/osqueryd",
					"pid":  "97830",
				},
			},
			Removed: []map[string]string{
				{
					"name": "osqueryd",
					"path": "/usr/local/bin/osqueryd",
					"pid":  "97650",
				},
			},
		},
	}

	parser := &BatchParser{}
	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestOsQueryBatchLogType(t *testing.T) {
	parser := &BatchParser{}
	require.Equal(t, "Osquery.Batch", parser.LogType())
}
