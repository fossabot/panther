package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/organization/models"
)

func (m *mockTable) AddActions(actions []*models.Action) (*models.Organization, error) {
	args := m.Called(actions)
	return args.Get(0).(*models.Organization), args.Error(1)
}

func TestCompleteActionError(t *testing.T) {
	m := &mockTable{}
	m.On("AddActions", mock.Anything, mock.Anything).Return(
		(*models.Organization)(nil), errors.New(""))
	orgTable = m

	action := models.VisitedOnboardingFlow
	result, err := (API{}).CompleteAction(&models.CompleteActionInput{
		CompletedActions: []*models.Action{&action},
	})
	m.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestCompleteAction(t *testing.T) {
	m := &mockTable{}
	action := models.VisitedOnboardingFlow
	output := &models.Organization{
		CompletedActions: []*models.Action{&action},
	}
	m.On("AddActions", mock.Anything, mock.Anything).Return(output, nil)
	orgTable = m

	result, err := (API{}).CompleteAction(&models.CompleteActionInput{
		CompletedActions: []*models.Action{&action},
	})
	m.AssertExpectations(t)
	assert.Equal(t, &models.CompleteActionOutput{CompletedActions: output.CompletedActions}, result)
	assert.NoError(t, err)
}
