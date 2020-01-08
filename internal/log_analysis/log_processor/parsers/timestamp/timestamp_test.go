package timestamp

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

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
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
