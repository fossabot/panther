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
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/tools/cfngen"
)

func TestDatabase(t *testing.T) {
	expectedOutput, err := readTestFile("testdata/db.template.json")
	require.NoError(t, err)

	catalogID := "12345"
	dbName := "db1"
	description := "Test db"

	db := NewDatabase(catalogID, dbName, description)

	resources := make(map[string]interface{})

	resources[dbName] = db

	cfTemplate := cfngen.NewTemplate("Test template", nil, resources, nil)

	cf := &bytes.Buffer{}

	require.NoError(t, cfTemplate.WriteCloudFormation(cf))

	// uncomment to see output
	// os.Stdout.Write(cf.Bytes())

	assert.Equal(t, expectedOutput, cf.String())
}
