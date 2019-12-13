package delivery

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
