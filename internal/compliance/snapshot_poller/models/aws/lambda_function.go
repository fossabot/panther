package aws

import "github.com/aws/aws-sdk-go/service/lambda"

const (
	LambdaFunctionSchema = "AWS.Lambda.Function"
)

// LambdaFunction contains all the information about an Lambda Function
type LambdaFunction struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from lambda.FunctionConfiguration
	CodeSha256       *string
	CodeSize         *int64
	DeadLetterConfig *lambda.DeadLetterConfig
	Description      *string
	Environment      *lambda.EnvironmentResponse
	Handler          *string
	KMSKeyArn        *string
	LastModified     *string
	Layers           []*lambda.Layer
	MasterArn        *string
	MemorySize       *int64
	RevisionId       *string
	Role             *string
	Runtime          *string
	Timeout          *int64
	TracingConfig    *lambda.TracingConfigResponse
	Version          *string
	VpcConfig        *lambda.VpcConfigResponse

	// Additional fields
	Policy *lambda.GetPolicyOutput
}
