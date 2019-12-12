package processor

import (
	"regexp"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

const lambdaNameRegex = `(arn:(aws[a-zA-Z-]*)?:lambda:)?([a-z]{2}(-gov)?-[a-z]+-\d{1}:)?(\d{12}:)?` +
	`(function:)?([a-zA-Z0-9-_]+)(:(\$LATEST|[a-zA-Z0-9-_]+))?`

func classifyLambda(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_awslambda.html
	if eventName == "AddLayerVersionPermission" ||
		eventName == "InvokeAsync" ||
		eventName == "InvokeFunction" {

		zap.L().Debug("lambda: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	region := detail.Get("awsRegion").Str
	lambdaARN := arn.ARN{
		Partition: "aws",
		Service:   "lambda",
		Region:    region,
		AccountID: accountID,
		Resource:  "function:",
	}
	switch eventName {
	case "AddPermission", "CreateAlias", "CreateEventSourceMapping", "CreateFunction", "DeleteAlias", "DeleteFunction",
		"DeleteFunctionConcurrency", "PublishVersion", "PutFunctionConcurrency", "RemovePermission", "UpdateAlias",
		"UpdateAlias20150331", "UpdateEventSourceMapping", "UpdateFunctionCode", "UpdateFunctionConfiguration",
		"UpdateFunctionCode20150331v2", "PublishVersion20150331", "UpdateEventSourceMapping20150331", "CreateAlias20150331":
		functionName := detail.Get("requestParameters.functionName").Str
		// Lambda Fun! This will need to be updated once we support tracking multiple aliases.
		// Legal formats:
		// Function name - my-function (name-only), my-function:v1 (with alias).
		// Function ARN - arn:aws:lambda:us-west-2:123456789012:function:my-function.
		// Partial ARN - 123456789012:function:my-function.
		// Regex taken from lambda user documentation referenced above.
		re := regexp.MustCompile(lambdaNameRegex)
		lambdaARN.Resource += re.FindStringSubmatch(functionName)[7]
	case "DeleteEventSourceMapping":
		functionName := detail.Get("responseElements.functionArn").Str
		re := regexp.MustCompile(lambdaNameRegex)
		lambdaARN.Resource += re.FindStringSubmatch(functionName)[7]
	case "TagResource", "UntagResource", "TagResource20170331v2", "UntagResource20170331v2":
		var err error
		lambdaARN, err = arn.Parse(detail.Get("requestParameters.resource").Str)
		if err != nil {
			zap.L().Error("lambda: error parsing ARN", zap.String("eventName", eventName), zap.Error(err))
			return nil
		}
	default:
		zap.L().Warn("lambda: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       eventName == "DeleteFunction",
		EventName:    eventName,
		ResourceID:   lambdaARN.String(),
		ResourceType: schemas.LambdaFunctionSchema,
	}}
}
