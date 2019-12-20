package timestamp

import (
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/tools/cfngen/gluecf"
)

var (
	expectedString        = "2019-12-15 01:01:01 +0000 UTC" // from String()
	expectedMarshalString = `"2019-12-15 01:01:01.000000000"`
	expectedTime          = time.Date(2019, 12, 15, 1, 1, 1, 0, time.UTC)

	jsonUnmarshalString    = `"2019-12-15T01:01:01Z"`
	osqueryUnmarshalString = `"Sun Dec 15 01:01:01 2019 UTC"`
)

func TestTimestampRFC3339_String(t *testing.T) {
	ts := (RFC3339)(expectedTime)
	assert.Equal(t, expectedString, ts.String())
}

func TestTimestampRFC3339_Marshal(t *testing.T) {
	ts := (RFC3339)(expectedTime)
	jsonTS, err := jsoniter.Marshal(&ts)
	assert.NoError(t, err)
	assert.Equal(t, expectedMarshalString, string(jsonTS))
}

func TestTimestampRFC3339_Unmarshal(t *testing.T) {
	var ts RFC3339
	err := jsoniter.Unmarshal([]byte(jsonUnmarshalString), &ts)
	assert.NoError(t, err)
	assert.Equal(t, (RFC3339)(expectedTime), ts)
}

func TestTimestampANSICwithTZ_String(t *testing.T) {
	ts := (ANSICwithTZ)(expectedTime)
	assert.Equal(t, expectedString, ts.String())
}

func TestTimestampANSICwithTZ_Marshal(t *testing.T) {
	ts := (ANSICwithTZ)(expectedTime)
	jsonTS, err := jsoniter.Marshal(&ts)
	assert.NoError(t, err)
	assert.Equal(t, expectedMarshalString, string(jsonTS))
}

func TestTimestampANSICwithTZ_Unmarshal(t *testing.T) {
	var ts ANSICwithTZ
	err := jsoniter.Unmarshal([]byte(osqueryUnmarshalString), &ts)
	assert.NoError(t, err)
	assert.Equal(t, (ANSICwithTZ)(expectedTime), ts)
}

func TestTimestampRFC3339_SerializeToCF(t *testing.T) {
	col := struct {
		MyTime RFC3339
	}{
		(RFC3339)(expectedTime),
	}
	cfCol := gluecf.InferJSONColumns(col, GlueMappings...)
	assert.Equal(t, glueType, cfCol[0].Type)
}
