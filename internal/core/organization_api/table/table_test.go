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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	dynamoAwsConfig = &dynamodb.AttributeValue{M: map[string]*dynamodb.AttributeValue{
		"appClientId":    {S: aws.String("111")},
		"identityPoolId": {S: aws.String("us-west-2:1234")},
		"userPoolId":     {S: aws.String("us-west-2_1234")},
	}}
)

type mockDynamoClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func TestNew(t *testing.T) {
	assert.NotNil(t, New("table", session.Must(session.NewSession())))
}
