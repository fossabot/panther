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
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/organization/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

func (m *mockDynamoClient) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

func TestUpdateDoestNotExist(t *testing.T) {
	mockClient := &mockDynamoClient{}
	returnErr := awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "", nil)
	mockClient.On("UpdateItem", mock.Anything).Return(
		(*dynamodb.UpdateItemOutput)(nil), returnErr)
	table := &OrganizationsTable{client: mockClient, Name: aws.String("table-name")}

	result, err := table.Update(&models.Organization{})
	mockClient.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.DoesNotExistError{}, err)
}

func TestUpdateServiceError(t *testing.T) {
	mockClient := &mockDynamoClient{}
	mockClient.On("UpdateItem", mock.Anything).Return(
		(*dynamodb.UpdateItemOutput)(nil), errors.New("service unavailable"))
	table := &OrganizationsTable{client: mockClient, Name: aws.String("table-name")}

	result, err := table.Update(&models.Organization{})
	mockClient.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.AWSError{}, err)
}

func TestUpdateUnmarshalError(t *testing.T) {
	mockClient := &mockDynamoClient{}
	// output has wrong type for one of the fields
	output := &dynamodb.UpdateItemOutput{
		Attributes: DynamoItem{"awsConfig": {SS: aws.StringSlice([]string{"panther", "labs"})}},
	}
	mockClient.On("UpdateItem", mock.Anything).Return(output, nil)
	table := &OrganizationsTable{client: mockClient, Name: aws.String("test-table")}

	result, err := table.Update(&models.Organization{})
	mockClient.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.InternalError{}, err)
}

func TestUpdate(t *testing.T) {
	mockClient := &mockDynamoClient{}
	org := &models.Organization{}

	output := &dynamodb.UpdateItemOutput{
		Attributes: DynamoItem{"id": {S: aws.String("1")}},
	}

	expectedUpdate := expression.
		Set(expression.Name("alertReportFrequency"), expression.Value(org.AlertReportFrequency)).
		Set(expression.Name("awsConfig"), expression.Value(org.AwsConfig)).
		Set(expression.Name("displayName"), expression.Value(org.DisplayName)).
		Set(expression.Name("email"), expression.Value(org.Email)).
		Set(expression.Name("phone"), expression.Value(org.Phone)).
		Set(expression.Name("remediationConfig"), expression.Value(org.RemediationConfig))
	expectedCondition := expression.AttributeExists(expression.Name("id"))
	expectedExpression, _ := expression.NewBuilder().WithCondition(expectedCondition).WithUpdate(expectedUpdate).Build()

	expectedUpdateItemInput := &dynamodb.UpdateItemInput{
		ConditionExpression:       expectedExpression.Condition(),
		ExpressionAttributeNames:  expectedExpression.Names(),
		ExpressionAttributeValues: expectedExpression.Values(),
		Key:                       DynamoItem{"id": {S: aws.String("1")}},
		ReturnValues:              aws.String("ALL_NEW"),
		TableName:                 aws.String("test-table"),
		UpdateExpression:          expectedExpression.Update(),
	}

	mockClient.On("UpdateItem", expectedUpdateItemInput).Return(output, nil)
	table := &OrganizationsTable{client: mockClient, Name: aws.String("test-table")}

	result, err := table.Update(org)
	mockClient.AssertExpectations(t)
	require.NoError(t, err)
	expected := &models.Organization{}
	assert.Equal(t, expected, result)
}

func TestAddActions(t *testing.T) {
	mockClient := &mockDynamoClient{}
	output := &dynamodb.UpdateItemOutput{
		Attributes: DynamoItem{"id": {S: aws.String("1")}},
	}
	mockClient.On("UpdateItem", mock.Anything).Return(output, nil)
	table := &OrganizationsTable{client: mockClient, Name: aws.String("test-table")}
	action := models.VisitedOnboardingFlow
	result, err := table.AddActions([]*models.Action{&action})
	mockClient.AssertExpectations(t)
	require.NoError(t, err)
	expected := &models.Organization{}
	assert.Equal(t, expected, result)
}
