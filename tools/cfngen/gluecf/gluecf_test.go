package gluecf

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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/timestamp"
	"github.com/panther-labs/panther/pkg/awsglue"
)

type dummyParserEvent struct {
	FirstName   string
	LastName    string
	DOB         timestamp.RFC3339
	Anniversary timestamp.ANSICwithTZ
}

func TestGenerateGlueCloudFormation(t *testing.T) {
	expectedOutput, err := readTestFile("testdata/gluecf.json")
	require.NoError(t, err)

	// use simple consistent reference set of parsers
	table, err := awsglue.NewGlueMetadata(awsglue.InternalDatabaseName, "dummy", "dummy",
		awsglue.GlueTableHourly, false, &dummyParserEvent{})
	require.NoError(t, err)
	tables := []*awsglue.GlueMetadata{table}

	cf, err := GenerateCloudFormation(tables)
	require.NoError(t, err)

	// un-comment to see output
	// os.Stdout.Write(cf)

	assert.Equal(t, expectedOutput, string(cf))
}
