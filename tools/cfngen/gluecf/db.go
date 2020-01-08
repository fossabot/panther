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
