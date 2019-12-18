package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/lambda/alerts/models"
)

// The handler signatures must match those in the LambdaInput struct.
func TestRouter(t *testing.T) {
	assert.Nil(t, router.VerifyHandlers(&models.LambdaInput{}))
}
