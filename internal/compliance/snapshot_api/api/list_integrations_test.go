package api

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
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/internal/compliance/snapshot_api/ddb"
	"github.com/panther-labs/panther/internal/compliance/snapshot_api/ddb/modelstest"
)

func TestListIntegrations(t *testing.T) {
	lastScanEndTime, err := time.Parse(time.RFC3339, "2019-04-10T23:00:00Z")
	require.NoError(t, err)

	lastScanStartTime, err := time.Parse(time.RFC3339, "2019-04-10T22:59:00Z")
	require.NoError(t, err)

	db = &ddb.DDB{
		Client: &modelstest.MockDDBClient{
			MockScanAttributes: []map[string]*dynamodb.AttributeValue{
				{
					"awsAccountId":         {S: aws.String("123456789012")},
					"eventStatus":          {S: aws.String(models.StatusOK)},
					"integrationId":        {S: aws.String(testIntegrationID)},
					"integrationLabel":     {S: aws.String(testIntegrationLabel)},
					"integrationType":      {S: aws.String(testIntegrationType)},
					"lastScanEndTime":      {S: aws.String(lastScanEndTime.Format(time.RFC3339))},
					"lastScanErrorMessage": {S: aws.String("")},
					"lastScanStartTime":    {S: aws.String(lastScanStartTime.Format(time.RFC3339))},
					"scanEnabled":          {BOOL: aws.Bool(true)},
					"scanIntervalMins":     {N: aws.String(strconv.Itoa(1440))},
					"scanStatus":           {S: aws.String(models.StatusOK)},
				},
			},
			TestErr: false,
		},
		TableName: "test",
	}

	expected := &models.SourceIntegration{
		SourceIntegrationMetadata: &models.SourceIntegrationMetadata{
			AWSAccountID:     aws.String("123456789012"),
			IntegrationID:    aws.String(testIntegrationID),
			IntegrationLabel: aws.String(testIntegrationLabel),
			IntegrationType:  aws.String(testIntegrationType),
			ScanEnabled:      aws.Bool(true),
			ScanIntervalMins: aws.Int(1440),
		},
		SourceIntegrationStatus: &models.SourceIntegrationStatus{
			ScanStatus:  aws.String(models.StatusOK),
			EventStatus: aws.String(models.StatusOK),
		},
		SourceIntegrationScanInformation: &models.SourceIntegrationScanInformation{
			LastScanEndTime:      &lastScanEndTime,
			LastScanErrorMessage: aws.String(""),
			LastScanStartTime:    &lastScanStartTime,
		},
	}
	out, err := apiTest.ListIntegrations(&models.ListIntegrationsInput{})

	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.Len(t, out, 1)
	assert.Equal(t, expected, out[0])
}

// An empty list of integrations is returned instead of null
func TestListIntegrationsEmpty(t *testing.T) {
	db = &ddb.DDB{
		Client: &modelstest.MockDDBClient{
			MockScanAttributes: []map[string]*dynamodb.AttributeValue{},
			TestErr:            false,
		},
		TableName: "test",
	}

	out, err := apiTest.ListIntegrations(&models.ListIntegrationsInput{})

	require.NoError(t, err)
	assert.Equal(t, []*models.SourceIntegration{}, out)
}

func TestHandleListIntegrationsScanError(t *testing.T) {
	db = &ddb.DDB{
		Client: &modelstest.MockDDBClient{
			MockScanAttributes: []map[string]*dynamodb.AttributeValue{},
			TestErr:            true,
		},
		TableName: "test",
	}

	out, err := apiTest.ListIntegrations(&models.ListIntegrationsInput{})

	require.NotNil(t, err)
	assert.Nil(t, out)
}
