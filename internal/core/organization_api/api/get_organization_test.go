package api

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
