// Package table manages all of the Dynamo calls (query, scan, get, write, etc).
package table

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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/panther-labs/panther/api/lambda/organization/models"
)

// API defines the interface for the table which can be used for mocking.
type API interface {
	Get() (*models.Organization, error)
	Put(*models.Organization) error
	Update(*models.Organization) (*models.Organization, error)
	AddActions(actions []*models.Action) (*models.Organization, error)
}

// OrganizationsTable encapsulates a connection to the Dynamo table.
type OrganizationsTable struct {
	Name   *string
	client dynamodbiface.DynamoDBAPI
}

// The OrganizationsTable must satisfy the API interface.
var _ API = (*OrganizationsTable)(nil)

// New creates a new Dynamo client which talks to the given table name.
func New(tableName string, sess *session.Session) *OrganizationsTable {
	return &OrganizationsTable{Name: aws.String(tableName), client: dynamodb.New(sess)}
}

// DynamoItem is a type alias for the item format expected by the Dynamo SDK.
type DynamoItem = map[string]*dynamodb.AttributeValue
