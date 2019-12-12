package api

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

var mockGetOrganizationOutputsInput = &models.GetOrganizationOutputsInput{}

var alertOutputItem = &models.AlertOutputItem{
	OutputID:           aws.String("outputId"),
	DisplayName:        aws.String("displayName"),
	CreatedBy:          aws.String("createdBy"),
	CreationTime:       aws.String("creationTime"),
	LastModifiedBy:     aws.String("lastModifiedBy"),
	LastModifiedTime:   aws.String("lastModifiedTime"),
	OutputType:         aws.String("slack"),
	VerificationStatus: aws.String(models.VerificationStatusSuccess),
	EncryptedConfig:    make([]byte, 1),
}

func TestGetOrganizationOutputs(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockEncryptionKey := new(mockEncryptionKey)
	encryptionKey = mockEncryptionKey
	mockOutputVerification := &mockOutputVerification{}
	outputVerification = mockOutputVerification
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	mockOutputsTable.On("GetOutputs").Return([]*models.AlertOutputItem{alertOutputItem}, nil)
	mockEncryptionKey.On("DecryptConfig", make([]byte, 1), mock.Anything).Return(nil)
	mockDefaultsTable.On("GetDefaults", mock.Anything).Return([]*models.DefaultOutputsItem{}, nil)

	expectedAlertOutput := &models.AlertOutput{
		OutputID:           aws.String("outputId"),
		OutputType:         aws.String("slack"),
		CreatedBy:          aws.String("createdBy"),
		CreationTime:       aws.String("creationTime"),
		DisplayName:        aws.String("displayName"),
		LastModifiedBy:     aws.String("lastModifiedBy"),
		LastModifiedTime:   aws.String("lastModifiedTime"),
		VerificationStatus: aws.String(models.VerificationStatusSuccess),
		OutputConfig:       &models.OutputConfig{},
		DefaultForSeverity: []*string{},
	}

	result, err := (API{}).GetOrganizationOutputs(mockGetOrganizationOutputsInput)

	assert.NoError(t, err)
	assert.Equal(t, []*models.AlertOutput{expectedAlertOutput}, result)
	mockOutputsTable.AssertExpectations(t)
	mockEncryptionKey.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}

func TestGetOrganizationOutputsDdbError(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable

	mockOutputsTable.On("GetOutputs").Return([]*models.AlertOutputItem{}, errors.New("fake error"))

	_, err := (API{}).GetOrganizationOutputs(mockGetOrganizationOutputsInput)

	assert.Error(t, errors.New("fake error"), err)
	mockOutputsTable.AssertExpectations(t)
}

func TestGetOrganizationDecryptionError(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockEncryptionKey := new(mockEncryptionKey)
	encryptionKey = mockEncryptionKey

	mockOutputsTable.On("GetOutputs").Return([]*models.AlertOutputItem{alertOutputItem}, nil)
	mockEncryptionKey.On("DecryptConfig", make([]byte, 1), mock.Anything).Return(errors.New("fake error"))

	_, err := (API{}).GetOrganizationOutputs(mockGetOrganizationOutputsInput)

	assert.Error(t, errors.New("fake error"), err)
	mockOutputsTable.AssertExpectations(t)
	mockEncryptionKey.AssertExpectations(t)
}
