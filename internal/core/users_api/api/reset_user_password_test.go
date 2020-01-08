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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/internal/core/users_api/gateway"
	"github.com/panther-labs/panther/pkg/genericapi"
)

type mockGatewayResetUserPasswordClient struct {
	gateway.API
	gatewayErr bool
}

func (m *mockGatewayResetUserPasswordClient) ResetUserPassword(*string, *string) error {
	if m.gatewayErr {
		return &genericapi.AWSError{}
	}
	return nil
}

func TestResetUserPasswordGatewayErr(t *testing.T) {
	userGateway = &mockGatewayResetUserPasswordClient{gatewayErr: true}
	input := &models.ResetUserPasswordInput{
		ID:         aws.String("user123"),
		UserPoolID: aws.String("fakePoolId"),
	}
	assert.Error(t, (API{}).ResetUserPassword(input))
}

func TestResetUserPasswordHandle(t *testing.T) {
	userGateway = &mockGatewayResetUserPasswordClient{}
	input := &models.ResetUserPasswordInput{
		ID:         aws.String("user123"),
		UserPoolID: aws.String("fakePoolId"),
	}
	assert.NoError(t, (API{}).ResetUserPassword(input))
}
