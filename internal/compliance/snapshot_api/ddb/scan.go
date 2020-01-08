package ddb

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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// ScanEnabledIntegrations returns all enabled integrations based on type.
// It performs a DDB scan of the entire table with a filter expression.
func (ddb *DDB) ScanEnabledIntegrations(input *models.ListIntegrationsInput) ([]*models.SourceIntegration, error) {
	filt := expression.And(
		expression.Name("scanEnabled").Equal(expression.Value(true)),
		expression.Name("integrationType").Equal(expression.Value(input.IntegrationType)),
	)
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, &genericapi.InternalError{Message: "failed to build dynamodb expression"}
	}

	output, err := ddb.Client.Scan(&dynamodb.ScanInput{
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String(ddb.TableName),
	})
	if err != nil {
		return nil, &genericapi.AWSError{Err: err, Method: "Dynamodb.Scan"}
	}

	var enabledIntegrations []*models.SourceIntegration
	if err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &enabledIntegrations); err != nil {
		return nil, err
	}

	if enabledIntegrations == nil {
		enabledIntegrations = make([]*models.SourceIntegration, 0)
	}
	return enabledIntegrations, nil
}
