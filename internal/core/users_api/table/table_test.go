package users

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDynamoClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func TestNew(t *testing.T) {
	assert.NotNil(t, New("table", session.Must(session.NewSession())))
}
