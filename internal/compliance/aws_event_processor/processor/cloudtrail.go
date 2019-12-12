package processor

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

func classifyCloudTrail(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_awscloudtrail.html
	if strings.HasPrefix(eventName, "Lookup") {
		zap.L().Debug("cloudtrail: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	trailARNBase := arn.ARN{
		Partition: "aws",
		Service:   "cloudtrail",
		Region:    detail.Get("awsRegion").Str,
		AccountID: accountID,
	}
	var err error

	switch eventName {
	case "AddTags", "RemoveTags":
		// This will always be an ARN
		trailARNBase, err = arn.Parse(detail.Get("requestParameters.resourceId").Str)
		if err != nil {
			zap.L().Error("cloudtrail: unable to parse ARN",
				zap.String("eventName", eventName),
				zap.String("resourceId", detail.Get("requestParameters.resourceId").Str))
			return nil
		}
	case "StartLogging", "StopLogging", "UpdateTrail":
		// The name requestParameter could be either the trail name or the trail ARN, so we try both
		trailARN, err := arn.Parse(detail.Get("requestParameters.name").Str)
		if err == nil {
			trailARNBase = trailARN
		} else {
			trailARNBase.Resource = "trail/" + detail.Get("requestParameters.name").Str
		}
	case "CreateTrail", "PutEventSelectors":
		// These events may effect the CloudTrail Meta resource, so must launch a full account scan
		return []*resourceChange{{
			AwsAccountID: accountID,
			EventName:    eventName,
			ResourceType: schemas.CloudTrailSchema,
		}}
	case "DeleteTrail":
		// The name requestParameter could be either the trail name or the trail ARN, so we try both
		// Special case of full account scan where we must also delete a resource
		trailARN, err := arn.Parse(detail.Get("requestParameters.name").Str)
		if err == nil {
			trailARNBase = trailARN
		} else {
			trailARNBase.Resource = "trail/" + detail.Get("requestParameters.name").Str
		}
		return []*resourceChange{
			{
				AwsAccountID: accountID,
				Delete:       true,
				EventName:    eventName,
				ResourceID:   trailARNBase.String(),
				ResourceType: schemas.CloudTrailSchema,
			},
			{
				AwsAccountID: accountID,
				EventName:    eventName,
				ResourceType: schemas.CloudTrailSchema,
			}}
	default:
		zap.L().Warn("cloudtrail: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	// This will only happen when the name parameter is an ARN and also the resource exists in a
	// different account than the account this event was logged in
	if accountID != trailARNBase.AccountID {
		zap.L().Info("cloudtrail: discarding resource from another account",
			zap.String("ResourceID", trailARNBase.String()),
			zap.String("AccountID", accountID))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       false,
		EventName:    eventName,
		ResourceID:   trailARNBase.String(),
		ResourceType: schemas.CloudTrailSchema,
	}}
}
