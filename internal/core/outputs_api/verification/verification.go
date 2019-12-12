package verification

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

var (
	emailVerificationTemplate = os.Getenv("EMAIL_VERIFICATION_TEMPLATE")
	sesConfigurationSet       = os.Getenv("SES_CONFIGURATION_SET")
	usersAPI                  = os.Getenv("USERS_API")
)

// OutputVerificationAPI defines the interface for the outputs table which can be used for mocking.
type OutputVerificationAPI interface {
	// GetVerificationStatus gets the verification status of an email
	GetVerificationStatus(output *models.AlertOutput) (*string, error)

	// VerifyOutput verifies an email address
	VerifyOutput(output *models.AlertOutput) (*models.AlertOutput, error)
}

// OutputVerification encapsulates a connection to the Dynamo rules table.
type OutputVerification struct {
	sesClient    sesiface.SESAPI
	lambdaClient lambdaiface.LambdaAPI
}

// NewVerification creates a new OutputVerification struct
func NewVerification(sess *session.Session) *OutputVerification {
	return &OutputVerification{
		sesClient:    ses.New(sess),
		lambdaClient: lambda.New(sess),
	}
}