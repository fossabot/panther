package awstest

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

// AssumeRoleMock generates a set of fake credentials for testing.
func AssumeRoleMock(
	pollerInput *awsmodels.ResourcePollerInput,
	sess *session.Session,
) (*credentials.Credentials, error) {

	return &credentials.Credentials{}, nil
}
