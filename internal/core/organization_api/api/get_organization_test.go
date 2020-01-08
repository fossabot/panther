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
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/lambda/organization/models"
)

func (m *mockTable) Get() (*models.Organization, error) {
	args := m.Called()
	return args.Get(0).(*models.Organization), args.Error(1)
}

func TestGetOrganizationError(t *testing.T) {
	m := &mockTable{}
	m.On("Get").Return(
		(*models.Organization)(nil), errors.New(""))
	orgTable = m

	result, err := (API{}).GetOrganization(&models.GetOrganizationInput{})
	m.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestGetOrganization(t *testing.T) {
	testOrganization := &models.Organization{
		DisplayName: aws.String("panther-labs"),
		Email:       aws.String("contact@runpanther.io"),
		AwsConfig: &models.AwsConfig{
			UserPoolID:     aws.String("userPool"),
			AppClientID:    aws.String("appClient"),
			IdentityPoolID: aws.String("identityPool"),
		},
	}

	m := &mockTable{}
	m.On("Get").Return(testOrganization, nil)
	orgTable = m

	result, err := (API{}).GetOrganization(&models.GetOrganizationInput{})
	m.AssertExpectations(t)
	assert.NotNil(t, result)
	assert.Equal(t, &models.GetOrganizationOutput{Organization: testOrganization}, result)
	assert.NoError(t, err)
}
