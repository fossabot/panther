package gluecf

// Generate CF for a gluecf database:  https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-glue-database.html

// NOTE: the use of type interface{} allows strings and structs (e.g., cfngen.Ref{} and cfngen.Sub{} )

// Matches CF structure
type DatabaseInput struct {
	Name        interface{}
	Description string `json:",omitempty"`
}

type DatabaseProperties struct {
	CatalogID     interface{} `json:"CatalogId"` // required, string or Ref{}, need json tag to keep linter happy
	DatabaseInput DatabaseInput
}

type Database struct {
	Type       string
	DependsOn  []string `json:",omitempty"`
	Properties DatabaseProperties
}

func NewDatabase(catalogID interface{}, name, description string) (db *Database) {
	db = &Database{
		Type: "AWS::Glue::Database",
		Properties: DatabaseProperties{
			CatalogID: catalogID,
			DatabaseInput: DatabaseInput{
				Name:        name,
				Description: description,
			},
		},
	}

	return
}
