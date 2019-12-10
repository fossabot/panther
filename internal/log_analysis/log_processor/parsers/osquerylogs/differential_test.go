package osquerylogs

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
)

func TestDifferentialLog(t *testing.T) {
	//nolint:lll
	log := `{"name":"pack_incident-response_mounts","hostIdentifier":"Quans-MacBook-Pro-2.local","calendarTime":"Tue Nov 5 06:08:26 2018 UTC","unixTime":"1572934106","epoch":"0","counter":"62","logNumericsAsNumbers":"false","decorations":{"host_uuid":"F919E9BF-0BF1-5456-8F6C-335243AEA537"},"columns":{"blocks":"61202533"},"action":"added","log_type":"result"}`

	expectedDate := Time{time.Unix(1541398106, 0).In(time.UTC)}
	expectedEvent := &Differential{
		Action:               aws.String("added"),
		Name:                 aws.String("pack_incident-response_mounts"),
		Epoch:                aws.Int(0),
		HostIdentifier:       aws.String(("Quans-MacBook-Pro-2.local")),
		UnixTime:             aws.Int(1572934106),
		LogNumericsAsNumbers: aws.Bool(false),
		LogType:              aws.String("result"),
		CalendarTime:         &expectedDate,
		Columns: map[string]string{
			"blocks": "61202533",
		},
		Counter: aws.Int(62),
		Decorations: map[string]string{
			"host_uuid": "F919E9BF-0BF1-5456-8F6C-335243AEA537",
		},
	}

	parser := &DifferentialParser{}
	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestDifferentialLogWithoutLogNumericAsNumbers(t *testing.T) {
	//nolint:lll
	log := `{"action":"added","calendarTime":"Tue Nov 5 06:08:26 2018 UTC","columns":{"build_distro":"10.12"},"counter":"255","decorations":{"host_uuid":"37821E12-CC8A-5AA3-A90C-FAB28A5BF8F9" },"epoch":"0","hostIdentifier":"host.lan","log_type":"result","name":"pack_osquery-monitoring_osquery_info","unixTime":"1536682461"}`
	expectedDate := Time{time.Unix(1541398106, 0).In(time.UTC)}

	expectedEvent := &Differential{
		Action:         aws.String("added"),
		Name:           aws.String("pack_osquery-monitoring_osquery_info"),
		Epoch:          aws.Int(0),
		HostIdentifier: aws.String(("host.lan")),
		UnixTime:       aws.Int(1536682461),
		LogType:        aws.String("result"),
		CalendarTime:   &expectedDate,
		Columns: map[string]string{
			"build_distro": "10.12",
		},
		Counter: aws.Int(255),
		Decorations: map[string]string{
			"host_uuid": "37821E12-CC8A-5AA3-A90C-FAB28A5BF8F9",
		},
	}

	parser := &DifferentialParser{}
	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestOsQueryDifferentialLogType(t *testing.T) {
	parser := &DifferentialParser{}
	require.Equal(t, "Osquery.Differential", parser.LogType())
}
