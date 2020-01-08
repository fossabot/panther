package osquerylogs

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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/timestamp"
)

func TestSnapshotLog(t *testing.T) {
	//nolint:lll
	log := `{"action": "snapshot","snapshot": [{"parent": "0","path": "/sbin/launchd","pid": "1"}],"name": "process_snapshot","hostIdentifier": "hostname.local","calendarTime": "Tue Nov 5 06:08:26 2018 UTC","unixTime": "1462228052","epoch": "314159265","counter": "1","numerics": false}`

	expectedTime := time.Unix(1541398106, 0).UTC()
	expectedEvent := &Snapshot{
		Action:         aws.String("snapshot"),
		Name:           aws.String("process_snapshot"),
		Epoch:          aws.Int(314159265),
		HostIdentifier: aws.String(("hostname.local")),
		UnixTime:       aws.Int(1462228052),
		CalendarTime:   (*timestamp.ANSICwithTZ)(&expectedTime),
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
