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
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/organization/models"
)

func (m *mockTable) Update(input *models.Organization) (*models.Organization, error) {
	args := m.Called(input)
	return args.Get(0).(*models.Organization), args.Error(1)
}

func TestUpdateOrganizationError(t *testing.T) {
	m := &mockTable{}
	m.On("Update", mock.Anything).Return(
		(*models.Organization)(nil), errors.New(""))
	orgTable = m

	result, err := (API{}).UpdateOrganization(&models.UpdateOrganizationInput{})
	m.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestUpdateOrganization(t *testing.T) {
	m := &mockTable{}
	output := &models.Organization{
		DisplayName:          aws.String("panther-labs"),
		AlertReportFrequency: aws.String("P1W"),
		Email:                aws.String("fake@email.com"),
		AwsConfig: &models.AwsConfig{
			UserPoolID:     aws.String("userPool"),
			AppClientID:    aws.String("appClient"),
			IdentityPoolID: aws.String("identityPool"),
		},
	}
	m.On("Update", mock.Anything).Return(output, nil)
	orgTable = m

	result, err := (API{}).UpdateOrganization(&models.UpdateOrganizationInput{})
	m.AssertExpectations(t)
	assert.Equal(t, &models.UpdateOrganizationOutput{Organization: output}, result)
	assert.NoError(t, err)
}
