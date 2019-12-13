package gluecf

// Generate CF for a gluecf table: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-glue-table.html

// Below structs match CF structure

type Column struct {
	Name    string // required
	Type    string // required
	Comment string `json:",omitempty"`
}

type SerdeInfo struct {
	SerializationLibrary string                 `json:",omitempty"`
	Parameters           map[string]interface{} `json:",omitempty"`
}

type StorageDescriptor struct { // nolint
	InputFormat            string      // required
	OutputFormat           string      // required
	Compressed             bool        `json:",omitempty"`
	Location               interface{} // required
	BucketColumns          []Column    `json:",omitempty"`
	SortColumns            []Column    `json:",omitempty"`
	StoredAsSubDirectories bool        `json:",omitempty"`
	SerdeInfo              SerdeInfo   // required
	Columns                []Column    // required
}

type TableInput struct {
	TableType         string            // required
	Name              string            // required
	Description       string            `json:",omitempty"`
	StorageDescriptor StorageDescriptor // required
	PartitionKeys     []Column          `json:",omitempty"`
}

type TableProperties struct {
	CatalogID    interface{} `json:"CatalogId"` // required,  string or Ref{}, need json tag to keep linter happy
	DatabaseName string      // required
	TableInput   TableInput  // required
}

type Table struct {
	Type       string   // required
	DependsOn  []string `json:",omitempty"`
	Properties TableProperties
}

// Core function to create a table
func newExternalTable(catalogID interface{}, databaseName, name, description string, sd *StorageDescriptor, pks []Column) (db *Table) {
	db = &Table{
		Type: "AWS::Glue::Table",
		Properties: TableProperties{
			CatalogID:    catalogID,
			DatabaseName: databaseName,
			TableInput: TableInput{
				TableType:         "EXTERNAL_TABLE",
				Name:              name,
				Description:       description,
				StorageDescriptor: *sd,
				PartitionKeys:     pks,
			},
		},
	}

	return
}

// inputs to table generation functions
type NewTableInput struct {
	CatalogID     interface{} // type interface{} allows strings and structs
	DatabaseName  string
	Name          string
	Description   string
	Location      interface{} // type interface{} allows strings and structs
	Columns       []Column
	PartitionKeys []Column
}

func NewParquetTable(input *NewTableInput) (db *Table) {
	sd := &StorageDescriptor{
		InputFormat:  "org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat",
		OutputFormat: "org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat",
		SerdeInfo: SerdeInfo{
			SerializationLibrary: "org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe",
			Parameters: map[string]interface{}{
				"serialization.format": "1",
			},
		},
		Location: input.Location,
		Columns:  input.Columns,
	}

	return newExternalTable(input.CatalogID, input.DatabaseName, input.Name, input.Description, sd, input.PartitionKeys)
}

func NewJSONLTable(input *NewTableInput) (db *Table) {
	sd := &StorageDescriptor{
		InputFormat:  "org.apache.hadoop.mapred.TextInputFormat",
		OutputFormat: "org.apache.hadoop.hive.ql.io.IgnoreKeyTextOutputFormat",
		SerdeInfo: SerdeInfo{
			SerializationLibrary: "org.openx.data.jsonserde.JsonSerDe",
			Parameters: map[string]interface{}{
				"serialization.format": "1",
				"case.insensitive":     "FALSE", // preserve case!
			},
		},
		Location: input.Location,
		Columns:  input.Columns,
	}

	return newExternalTable(input.CatalogID, input.DatabaseName, input.Name, input.Description, sd, input.PartitionKeys)
}
