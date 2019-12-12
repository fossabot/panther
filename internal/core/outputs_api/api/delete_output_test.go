package api

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

var mockDeleteOutputInput = &models.DeleteOutputInput{
	OutputID: aws.String("outputId"),
}

func TestDeleteOutput(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	mockDefaultsTable.On("GetDefaults").Return(make([]*models.DefaultOutputsItem, 0), nil)
	mockOutputsTable.On("DeleteOutput", aws.String("outputId")).Return(nil)

	err := (API{}).DeleteOutput(mockDeleteOutputInput)

	assert.NoError(t, err)
	mockOutputsTable.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}

func TestDeleteOutputInUse(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	defaultOutputs := []*models.DefaultOutputsItem{{OutputIDs: aws.StringSlice([]string{"outputId"})}}

	mockDefaultsTable.On("GetDefaults").Return(defaultOutputs, nil)

	err := (API{}).DeleteOutput(mockDeleteOutputInput)

	require.Error(t, err)
	require.IsType(t, err, &genericapi.InUseError{})
	mockOutputsTable.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}

func TestDeleteOutputGetDefaultsFails(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	mockDefaultsTable.On("GetDefaults").Return(make([]*models.DefaultOutputsItem, 0), errors.New("Error"))

	err := (API{}).DeleteOutput(mockDeleteOutputInput)

	require.Error(t, err)
	mockOutputsTable.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}

func TestDeleteOutputDeleteFails(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	mockDefaultsTable.On("GetDefaults").Return(make([]*models.DefaultOutputsItem, 0), nil)
	mockOutputsTable.On("DeleteOutput", aws.String("outputId")).Return(errors.New("error"))

	err := (API{}).DeleteOutput(mockDeleteOutputInput)

	require.Error(t, err)
	mockOutputsTable.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}

func TestDeleteOutputRemovesDefaults(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	defaultOutputs := []*models.DefaultOutputsItem{{
		OutputIDs: aws.StringSlice([]string{"outputId1", "outputId2"}),
		Severity:  aws.String("INFO"),
	}}

	expectedPutDefaults := &models.DefaultOutputsItem{
		OutputIDs: aws.StringSlice([]string{"outputId2"}),
		Severity:  aws.String("INFO"),
	}

	mockDefaultsTable.On("GetDefaults").Return(defaultOutputs, nil)
	mockDefaultsTable.On("PutDefaults", expectedPutDefaults).Return(nil)
	mockOutputsTable.On("DeleteOutput", aws.String("outputId1")).Return(nil)

	mockDeleteOutputInput = &models.DeleteOutputInput{
		OutputID: aws.String("outputId1"),
		Force:    aws.Bool(true),
	}
	err := (API{}).DeleteOutput(mockDeleteOutputInput)

	assert.NoError(t, err)
	mockOutputsTable.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}

func TestDeleteOutputRemovesDefaultsFailsIfDefaultUpdateFails(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	defaultOutputs := []*models.DefaultOutputsItem{{
		OutputIDs: aws.StringSlice([]string{"outputId"}),
		Severity:  aws.String("INFO"),
	}}

	expectedPutDefaults := &models.DefaultOutputsItem{
		OutputIDs: make([]*string, 0),
		Severity:  aws.String("INFO"),
	}

	mockDefaultsTable.On("GetDefaults").Return(defaultOutputs, nil)
	mockDefaultsTable.On("PutDefaults", expectedPutDefaults).Return(errors.New("error"))

	mockDeleteOutputInput = &models.DeleteOutputInput{
		OutputID: aws.String("outputId"),
		Force:    aws.Bool(true),
	}
	err := (API{}).DeleteOutput(mockDeleteOutputInput)

	assert.Error(t, err)
	mockOutputsTable.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}
