package genericapi

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestAlreadyExistsError(t *testing.T) {
	err := &AlreadyExistsError{Route: "Do", Message: "name=panther"}
	assert.Equal(t, "Do failed: already exists: name=panther", err.Error())
}

func TestAWSError(t *testing.T) {
	err := &AWSError{Route: "Do", Method: "dynamodb.PutItem", Err: errors.New("not authorized")}
	assert.Equal(t, "Do failed: AWS dynamodb.PutItem error: not authorized", err.Error())
}

func TestDoesNotExistError(t *testing.T) {
	err := &DoesNotExistError{Route: "Do", Message: "name=panther"}
	assert.Equal(t, "Do failed: does not exist: name=panther", err.Error())
}

func TestInternalError(t *testing.T) {
	err := &InternalError{Route: "Do", Message: "can't marshal to JSON"}
	assert.Equal(t, "Do failed: internal error: can't marshal to JSON", err.Error())
}

func TestInUseError(t *testing.T) {
	err := &InUseError{Route: "Do", Message: "name=panther"}
	assert.Equal(t, "Do failed: still in use: name=panther", err.Error())
}

func TestInvalidInputError(t *testing.T) {
	err := &InvalidInputError{Route: "Do", Message: "you forgot something"}
	assert.Equal(t, "Do failed: invalid input: you forgot something", err.Error())
}

func TestLambdaErrorEmpty(t *testing.T) {
	err := &LambdaError{}
	assert.Equal(t, "lambda error returned: (nil)", err.Error())
}

func TestLambdaError(t *testing.T) {
	err := &LambdaError{
		Route: "Do", FunctionName: "rules-api", ErrorMessage: aws.String("task timed out")}
	assert.Equal(t, "Do failed: lambda error returned: rules-api: task timed out", err.Error())
}
