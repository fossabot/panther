package lambdalogger

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
)

var testContext = lambdacontext.NewContext(
	context.Background(), &lambdacontext.LambdaContext{AwsRequestID: "test-request-id"})

func TestConfigureGlobalDebug(t *testing.T) {
	DebugEnabled = true
	lc, logger := ConfigureGlobal(testContext, nil)
	assert.NotNil(t, lc)
	assert.NotNil(t, logger)
}

func TestConfigureGlobalProd(t *testing.T) {
	DebugEnabled = false
	lc, logger := ConfigureGlobal(testContext, nil)
	assert.NotNil(t, lc)
	assert.NotNil(t, logger)
}

func TestConfigureExtraFields(t *testing.T) {
	lc, logger := ConfigureGlobal(testContext, map[string]interface{}{"panther": "labs"})
	assert.NotNil(t, lc)
	assert.NotNil(t, logger)
}

func TestConfigureGlobalError(t *testing.T) {
	assert.Panics(t, func() { ConfigureGlobal(context.Background(), nil) })
}
