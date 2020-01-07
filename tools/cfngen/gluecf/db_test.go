package gluecf

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
