package api

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

var mockSetDefaultsInput = &models.SetDefaultOutputsInput{
	Severity:  aws.String("INFO"),
	OutputIDs: aws.StringSlice([]string{"outputId1", "outputId2"}),
}

func TestSetDefaults(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	mockDefaultsTable := &mockDefaultsTable{}

	outputsTable = mockOutputsTable
	defaultsTable = mockDefaultsTable

	expectedDefaults := &models.DefaultOutputsItem{
		Severity:  aws.String("INFO"),
		OutputIDs: aws.StringSlice([]string{"outputId1", "outputId2"}),
	}

	expectedResult := &models.DefaultOutputs{
		Severity:  aws.String("INFO"),
		OutputIDs: aws.StringSlice([]string{"outputId1", "outputId2"}),
	}

	mockOutputsTable.On("GetOutput", aws.String("outputId1")).Return(&models.AlertOutputItem{}, nil)
	mockOutputsTable.On("GetOutput", aws.String("outputId2")).Return(&models.AlertOutputItem{}, nil)
	mockDefaultsTable.On("PutDefaults", expectedDefaults).Return(nil)

	result, err := (API{}).SetDefaultOutputs(mockSetDefaultsInput)

	require.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockOutputsTable.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}

func TestSetDefaultsFailureWhenOutputDoesntExist(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	mockDefaultsTable := &mockDefaultsTable{}

	outputsTable = mockOutputsTable
	defaultsTable = mockDefaultsTable

	mockOutputsTable.On("GetOutput", mock.Anything, mock.Anything).Return((*models.AlertOutputItem)(nil), errors.New("error"))

	result, err := (API{}).SetDefaultOutputs(mockSetDefaultsInput)

	require.Error(t, err)
	assert.Nil(t, result)
	mockOutputsTable.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}

func TestSetDefaultsFailureWhenPutFails(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	mockDefaultsTable := &mockDefaultsTable{}

	outputsTable = mockOutputsTable
	defaultsTable = mockDefaultsTable

	mockOutputsTable.On("GetOutput", mock.Anything, mock.Anything).Return(&models.AlertOutputItem{}, nil)
	mockDefaultsTable.On("PutDefaults", mock.Anything).Return(errors.New("error"))

	result, err := (API{}).SetDefaultOutputs(mockSetDefaultsInput)

	require.Error(t, err)
	assert.Nil(t, result)
	mockOutputsTable.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}
