package gluecf

import (
	"bytes"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"testing"

	"github.com/panther-labs/panther/tools/cfngen"
)

func TestTables(t *testing.T) {
	expectedOutput, err := readTestFile("testdata/tables.template.json")
	require.NoError(t, err)

	// pass in bucket name
	parameters := make(map[string]interface{})
	parameters["Bucket"] = &cfngen.Parameter{
		Type:        "String",
		Description: "Bucket to hold data for table",
	}

	resources := make(map[string]interface{})

	catalogID := "12345"
	dbName := "db1"

	db := NewDatabase(catalogID, dbName, "Test database")

	resources[dbName] = db

	// same for both tables
	columns := []Column{
		{Name: "c1", Type: "int", Comment: "foo"},
		{Name: "c2", Type: "varchar", Comment: "bar"},
	}

	partitionKeys := []Column{
		{Name: "year", Type: "int", Comment: "year"},
		{Name: "month", Type: "int", Comment: "month"},
		{Name: "day", Type: "int", Comment: "day"},
	}

	tableName := "parquetTable"
	description := "Test table"
	location := cfngen.Sub{Sub: "s3//${Bucket}/" + dbName + "/" + tableName}
	table := NewParquetTable(&NewTableInput{
		CatalogID:     catalogID,
		DatabaseName:  dbName,
		Name:          tableName,
		Description:   description,
		Location:      location,
		Columns:       columns,
		PartitionKeys: partitionKeys,
	})
	table.DependsOn = []string{dbName} // table depends on db resource
	resources[tableName] = table

	tableName = "jsonlTable"
	description = "Test table"
	location = cfngen.Sub{Sub: "s3//${Bucket}/" + dbName + "/" + tableName}
	table = NewJSONLTable(&NewTableInput{
		CatalogID:     catalogID,
		DatabaseName:  dbName,
		Name:          tableName,
		Description:   description,
		Location:      location,
		Columns:       columns,
		PartitionKeys: partitionKeys,
	})
	table.DependsOn = []string{dbName} // table depends on db resource
	resources[tableName] = table

	cfTemplate := cfngen.NewTemplate("Test template", parameters, resources)

	cf := &bytes.Buffer{}

	require.NoError(t, cfTemplate.WriteCloudFormation(cf))

	// uncomment to see output
	// os.Stdout.Write(cf.Bytes())

	assert.Equal(t, expectedOutput, cf.String())
}
