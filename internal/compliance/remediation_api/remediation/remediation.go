package remediation

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"

	"github.com/panther-labs/panther/api/gateway/remediation/models"
)

//InvokerAPI is the interface for the Invoker,
// the component that is responsible for invoking Remediation Lambda
type InvokerAPI interface {
	Remediate(*models.RemediateResource) error
	GetRemediations() (*models.Remediations, error)
}

//Invoker is responsible for invoking Remediation Lambda
type Invoker struct {
	lambdaClient lambdaiface.LambdaAPI
	awsSession   *session.Session
}

//NewInvoker method returns a new instance of Invoker
func NewInvoker(sess *session.Session) *Invoker {
	return &Invoker{
		lambdaClient: lambda.New(sess),
		awsSession:   sess,
	}
}
