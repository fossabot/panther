package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

var mockUpdateOutputInput = &models.UpdateOutputInput{
	OutputID:           aws.String("outputId"),
	DisplayName:        aws.String("displayName"),
	UserID:             aws.String("userId"),
	OutputConfig:       &models.OutputConfig{Sns: &models.SnsConfig{}},
	DefaultForSeverity: aws.StringSlice([]string{"CRITICAL", "HIGH"}),
}

func TestUpdateOutput(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockEncryptionKey := &mockEncryptionKey{}
	encryptionKey = mockEncryptionKey
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	alertOutputItem := &models.AlertOutputItem{
		OutputID:           aws.String("outputId"),
		DisplayName:        aws.String("displayName"),
		CreatedBy:          aws.String("createdBy"),
		CreationTime:       aws.String("createdTime"),
		OutputType:         aws.String("sns"),
		VerificationStatus: aws.String(models.VerificationStatusSuccess),
		EncryptedConfig:    make([]byte, 1),
	}

	mockDefaultsTable.On("GetDefaults", mock.Anything).Return([]*models.DefaultOutputsItem{}, nil)
	mockDefaultsTable.On("GetDefault", mock.Anything).Return(&models.DefaultOutputsItem{}, nil)
	mockDefaultsTable.On("PutDefaults", mock.Anything).Return(nil)
	mockOutputsTable.On("UpdateOutput", mock.Anything).Return(alertOutputItem, nil)
	mockOutputsTable.On("GetOutputByName", aws.String("displayName")).Return(nil, nil)
	mockEncryptionKey.On("EncryptConfig", mock.Anything).Return(make([]byte, 1), nil)
	mockEncryptionKey.On("DecryptConfig", mock.Anything, mock.Anything).Return(nil)

	result, err := (API{}).UpdateOutput(mockUpdateOutputInput)

	assert.NoError(t, err)
	assert.Equal(t, aws.String("outputId"), result.OutputID)
	assert.Equal(t, aws.String("displayName"), result.DisplayName)
	assert.Equal(t, aws.String("createdBy"), result.CreatedBy)
	assert.Equal(t, aws.String("userId"), result.LastModifiedBy)
	assert.Equal(t, aws.String("sns"), result.OutputType)

	mockOutputsTable.AssertExpectations(t)
}

func TestUpdateOutputOtherItemExists(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable

	preExistingAlertItem := &models.AlertOutputItem{
		OutputID: aws.String("outputId-2"),
	}

	mockOutputsTable.On("GetOutputByName", aws.String("displayName")).Return(preExistingAlertItem, nil)

	result, err := (API{}).UpdateOutput(mockUpdateOutputInput)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockOutputsTable.AssertExpectations(t)
}

func TestUpdateSameOutpuOutput(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockEncryptionKey := &mockEncryptionKey{}
	encryptionKey = mockEncryptionKey
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	alertOutputItem := &models.AlertOutputItem{
		OutputID:           aws.String("outputId"),
		DisplayName:        aws.String("displayName"),
		CreatedBy:          aws.String("createdBy"),
		CreationTime:       aws.String("createdTime"),
		OutputType:         aws.String("sns"),
		VerificationStatus: aws.String(models.VerificationStatusSuccess),
		EncryptedConfig:    make([]byte, 1),
	}

	preExistingAlertItem := &models.AlertOutputItem{
		OutputID: mockUpdateOutputInput.OutputID,
	}

	mockDefaultsTable.On("GetDefaults", mock.Anything).Return([]*models.DefaultOutputsItem{}, nil)
	mockDefaultsTable.On("GetDefault", mock.Anything).Return(&models.DefaultOutputsItem{}, nil)
	mockDefaultsTable.On("PutDefaults", mock.Anything).Return(nil)
	mockOutputsTable.On("UpdateOutput", mock.Anything).Return(alertOutputItem, nil)
	mockOutputsTable.On("GetOutputByName", aws.String("displayName")).Return(preExistingAlertItem, nil)
	mockEncryptionKey.On("EncryptConfig", mock.Anything).Return(make([]byte, 1), nil)
	mockEncryptionKey.On("DecryptConfig", mock.Anything, mock.Anything).Return(nil)

	result, err := (API{}).UpdateOutput(mockUpdateOutputInput)

	assert.NoError(t, err)
	assert.Equal(t, aws.String("outputId"), result.OutputID)
	assert.Equal(t, aws.String("displayName"), result.DisplayName)
	assert.Equal(t, aws.String("createdBy"), result.CreatedBy)
	assert.Equal(t, aws.String("userId"), result.LastModifiedBy)
	assert.Equal(t, aws.String("sns"), result.OutputType)
	assert.Equal(t, aws.String(models.VerificationStatusSuccess), result.VerificationStatus)

	mockOutputsTable.AssertExpectations(t)
}
