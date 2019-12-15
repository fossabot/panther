package api

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
