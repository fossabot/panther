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
