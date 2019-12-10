package osquerylogs

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
)

func TestSnapshotLog(t *testing.T) {
	//nolint:lll
	log := `{"action": "snapshot","snapshot": [{"parent": "0","path": "/sbin/launchd","pid": "1"}],"name": "process_snapshot","hostIdentifier": "hostname.local","calendarTime": "Tue Nov 5 06:08:26 2018 UTC","unixTime": "1462228052","epoch": "314159265","counter": "1","numerics": false}`

	expectedDate := Time{time.Unix(1541398106, 0).In(time.UTC)}
	expectedEvent := &Snapshot{
		Action:         aws.String("snapshot"),
		Name:           aws.String("process_snapshot"),
		Epoch:          aws.Int(314159265),
		HostIdentifier: aws.String(("hostname.local")),
		UnixTime:       aws.Int(1462228052),
		CalendarTime:   &expectedDate,
		Counter:        aws.Int(1),
		Snapshot: []map[string]string{
			{
				"parent": "0",
				"path":   "/sbin/launchd",
				"pid":    "1",
			},
		},
	}

	parser := &SnapshotParser{}
	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestOsQuerySnapshotLogType(t *testing.T) {
	parser := &SnapshotParser{}
	require.Equal(t, "Osquery.Snapshot", parser.LogType())
}
