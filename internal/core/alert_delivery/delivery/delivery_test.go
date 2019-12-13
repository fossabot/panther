package delivery

import (
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/stretchr/testify/mock"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
	"github.com/panther-labs/panther/internal/core/alert_delivery/outputs"
)

type mockOutputsClient struct {
	outputs.API
	mock.Mock
}

func (m *mockOutputsClient) Slack(alert *alertmodels.Alert, config *outputmodels.SlackConfig) *outputs.AlertDeliveryError {
	args := m.Called(alert, config)
	return args.Get(0).(*outputs.AlertDeliveryError)
}

type mockLambdaClient struct {
	lambdaiface.LambdaAPI
	mock.Mock
}

func (m *mockLambdaClient) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*lambda.InvokeOutput), args.Error(1)
}
