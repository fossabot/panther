package registry

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/timestamp"
)

// Read CF
func readTestFile(filename string) ([]byte, error) {
	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	contents, err := ioutil.ReadAll(fd)
	return contents, err
}

type dummyEvent struct {
	FirstName   string
	LastName    string
	DOB         timestamp.RFC3339
	Anniversary timestamp.ANSICwithTZ
}

type dummyParser struct{}

func (dp *dummyParser) Parse(log string) []interface{} {
	events := []dummyEvent{
		{
			FirstName:   "Bob",
			LastName:    "Smith",
			DOB:         (timestamp.RFC3339)(time.Date(2019, 01, 1, 1, 1, 1, 1, time.UTC)),
			Anniversary: (timestamp.ANSICwithTZ)(time.Date(2040, 01, 1, 1, 1, 1, 1, time.UTC)),
		},
	}
	return []interface{}{events}
}
func (dp *dummyParser) LogType() string {
	return "dummy"
}

func TestGenerateGlueCloudFormation(t *testing.T) {
	expectedOutput, err := readTestFile("testdata/gluecf.json")
	require.NoError(t, err)

	// use simple consistent reference set of parsers
	testRegistry := Registry{
		(&dummyParser{}).LogType(): DefaultHourlyLogParser(&dummyParser{},
			&dummyEvent{}, "dummy"),
	}

	cf, err := generateGlueCloudFormation(testRegistry)
	require.NoError(t, err)

	// un-comment to see output
	// os.Stdout.Write(cf)

	assert.Equal(t, string(expectedOutput), string(cf))
}
