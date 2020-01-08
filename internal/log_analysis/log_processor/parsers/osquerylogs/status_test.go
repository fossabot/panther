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

func TestStatusLog(t *testing.T) {
	//nolint:lll
	log := `{"hostIdentifier":"jacks-mbp.lan","calendarTime":"Tue Nov 5 06:08:26 2018 UTC","unixTime":"1535731040","severity":"0","filename":"scheduler.cpp","line":"83","message":"Executing scheduled query pack_incident-response_arp_cache: select * from arp_cache;","version":"3.2.6","decorations":{"host_uuid":"37821E12-CC8A-5AA3-A90C-FAB28A5BF8F9","username":"user"},"log_type":"status"}`

	expectedTime := time.Unix(1541398106, 0).UTC()
	expectedEvent := &Status{
		HostIdentifier: aws.String("jacks-mbp.lan"),
		CalendarTime:   (*timestamp.ANSICwithTZ)(&expectedTime),
		UnixTime:       aws.Int(1535731040),
		Severity:       aws.Int(0),
		Filename:       aws.String("scheduler.cpp"),
		Line:           aws.Int(83),
		Message:        aws.String("Executing scheduled query pack_incident-response_arp_cache: select * from arp_cache;"),
		Version:        aws.String("3.2.6"),
		LogType:        aws.String("status"),
		Decorations: map[string]string{
			"host_uuid": "37821E12-CC8A-5AA3-A90C-FAB28A5BF8F9",
			"username":  "user",
		},
	}

	parser := &StatusParser{}
	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestOsQueryStatusLogType(t *testing.T) {
	parser := &StatusParser{}
	require.Equal(t, "Osquery.Status", parser.LogType())
}
