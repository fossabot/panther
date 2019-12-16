package apihandlers

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/gateway/remediation/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

func TestGetRemediations(t *testing.T) {
	mockInvoker := &mockInvoker{}
	invoker = mockInvoker

	request := &events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{},
	}

	remediationsParameters := map[string]interface{}{
		"KMSMasterKeyID": "",
		"SSEAlgorithm":   "AES256",
	}
	remediations := &models.Remediations{
		"AWS.S3.EnableBucketEncryption": remediationsParameters,
	}

	mockInvoker.On("GetRemediations").Return(remediations, nil)

	expectedResponseBody := map[string]interface{}{"AWS.S3.EnableBucketEncryption": remediationsParameters}
	response := GetRemediations(request)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	var responseBody map[string]interface{}
	assert.NoError(t, jsoniter.UnmarshalFromString(response.Body, &responseBody))
	assert.Equal(t, expectedResponseBody, responseBody)
	mockInvoker.AssertExpectations(t)
}

func TestGetRemediationsLambdaDoesntExist(t *testing.T) {
	mockInvoker := &mockInvoker{}
	invoker = mockInvoker

	request := &events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{},
	}

	mockInvoker.On("GetRemediations").Return(
		nil, &genericapi.DoesNotExistError{Message: "there is no aws remediation lambda configured for organization"})

	expectedResponseBody := &models.Error{Message: aws.String("Remediation Lambda not found or misconfigured")}
	response := GetRemediations(request)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	responseBody := &models.Error{}
	assert.NoError(t, jsoniter.UnmarshalFromString(response.Body, responseBody))
	assert.Equal(t, expectedResponseBody, responseBody)
	mockInvoker.AssertExpectations(t)
}
