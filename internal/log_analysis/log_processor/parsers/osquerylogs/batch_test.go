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
