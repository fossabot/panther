package table

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	dynamoAwsConfig = &dynamodb.AttributeValue{M: map[string]*dynamodb.AttributeValue{
		"appClientId":    {S: aws.String("111")},
		"identityPoolId": {S: aws.String("us-west-2:1234")},
		"userPoolId":     {S: aws.String("us-west-2_1234")},
	}}
)

type mockDynamoClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func TestNew(t *testing.T) {
	assert.NotNil(t, New("table", session.Must(session.NewSession())))
}
