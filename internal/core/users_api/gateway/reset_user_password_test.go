package gateway

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
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	providerI "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/stretchr/testify/assert"
)

type mockResetUserPasswordClient struct {
	providerI.CognitoIdentityProviderAPI
	serviceErr bool
}

func (m *mockResetUserPasswordClient) AdminResetUserPassword(
	*provider.AdminResetUserPasswordInput) (*provider.AdminResetUserPasswordOutput, error) {

	if m.serviceErr {
		return nil, errors.New("cognito does not exist")
	}
	return &provider.AdminResetUserPasswordOutput{}, nil
}

func TestResetUserPassword(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockResetUserPasswordClient{}}
	assert.NoError(t, gw.ResetUserPassword(aws.String("user123"), aws.String("userPoolId")))
}

func TestResetUserPasswordFailed(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockResetUserPasswordClient{serviceErr: true}}
	assert.Error(t, gw.ResetUserPassword(aws.String("user123"), aws.String("userPoolId")))
}
