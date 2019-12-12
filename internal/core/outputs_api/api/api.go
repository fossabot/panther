// Package api defines CRUD actions for Panther alert outputs.
package api

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/panther-labs/panther/internal/core/outputs_api/encryption"
	"github.com/panther-labs/panther/internal/core/outputs_api/table"
	"github.com/panther-labs/panther/internal/core/outputs_api/verification"
)

// The API consists of receiver methods for each of the handlers.
type API struct{}

var (
	awsSession = session.Must(session.NewSession())

	encryptionKey encryption.API = encryption.New(os.Getenv("KEY_ID"), awsSession)

	outputsTable table.OutputsAPI = table.NewOutputs(
		os.Getenv("OUTPUTS_TABLE_NAME"),
		os.Getenv("OUTPUTS_DISPLAY_NAME_INDEX_NAME"),
		awsSession)
	defaultsTable table.DefaultsAPI = table.NewDefaults(
		os.Getenv("DEFAULTS_TABLE_NAME"),
		awsSession)

	outputVerification verification.OutputVerificationAPI = verification.NewVerification(awsSession)
)
