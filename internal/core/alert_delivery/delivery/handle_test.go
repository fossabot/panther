package delivery

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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/internal/core/alert_delivery/models"
	"github.com/panther-labs/panther/internal/core/alert_delivery/outputs"
)

func TestMustParseIntPanic(t *testing.T) {
	assert.Panics(t, func() { mustParseInt("") })
}

func TestHandleAlerts(t *testing.T) {
	mockClient := &mockOutputsClient{}
	outputClient = mockClient
	mockClient.On("Slack", mock.Anything, mock.Anything).Return((*outputs.AlertDeliveryError)(nil))
	setCaches()
	alerts := []*models.Alert{sampleAlert(), sampleAlert(), sampleAlert()}
	assert.NotPanics(t, func() { HandleAlerts(alerts) })
}

func TestHandleAlertsPermanentlyFailed(t *testing.T) {
	createdAtTime, _ := time.Parse(time.RFC3339, "2019-05-03T11:40:13Z")
	mockClient := &mockOutputsClient{}
	outputClient = mockClient
	mockClient.On("Slack", mock.Anything, mock.Anything).Return(&outputs.AlertDeliveryError{})
	sqsClient = &mockSQSClient{}
	setCaches()
	os.Setenv("ALERT_RETRY_DURATION_MINS", "5")
	alert := sampleAlert()
	alert.CreatedAt = &createdAtTime
	alerts := []*models.Alert{alert, alert, alert}
	sqsMessages = 0

	HandleAlerts(alerts)
	assert.Equal(t, 0, sqsMessages)
}

func TestHandleAlertsTemporarilyFailed(t *testing.T) {
	createdAtTime := time.Now()
	mockClient := &mockOutputsClient{}
	outputClient = mockClient
	mockClient.On("Slack", mock.Anything, mock.Anything).Return(&outputs.AlertDeliveryError{})
	sqsClient = &mockSQSClient{}
	setCaches()
	os.Setenv("ALERT_RETRY_DURATION_MINS", "5")
	os.Setenv("ALERT_QUEUE_URL", "sqs.url")
	os.Setenv("MIN_RETRY_DELAY_SECS", "10")
	os.Setenv("MAX_RETRY_DELAY_SECS", "30")
	alert := sampleAlert()
	alert.CreatedAt = &createdAtTime
	alerts := []*models.Alert{alert, alert, alert}
	sqsMessages = 0

	HandleAlerts(alerts)
	assert.Equal(t, 3, sqsMessages)
}
