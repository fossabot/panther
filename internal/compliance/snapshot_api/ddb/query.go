package ddb

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// QueryIntegrations returns all snapshot integrations
func (ddb *DDB) QueryIntegrations() ([]*models.SourceIntegration, error) {
	builder := expression.NewBuilder()
	expr, err := builder.Build()
	if err != nil {
		return nil, &genericapi.InternalError{Message: "failed to build dynamodb expression"}
	}

	output, err := ddb.Client.Query(&dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String(ddb.TableName),
	})
	if err != nil {
		return nil, &genericapi.AWSError{Err: err, Method: "Dynamodb.Scan"}
	}

	var integrations []*models.SourceIntegration
	if err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &integrations); err != nil {
		return nil, err
	}

	if integrations == nil {
		integrations = make([]*models.SourceIntegration, 0)
	}

	return integrations, nil
}
