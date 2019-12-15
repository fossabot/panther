package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/lambda/organization/models"
)

// The handler signatures must match those in the LambdaInput struct.
func TestRouter(t *testing.T) {
	assert.NoError(t, router.VerifyHandlers(&models.LambdaInput{}))
}
