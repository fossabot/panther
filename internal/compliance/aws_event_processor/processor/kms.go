package processor

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/api/snapshot/aws"
)

func classifyKMS(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_awskeymanagementservice.html
	if strings.HasPrefix(eventName, "Decrypt") ||
		strings.HasPrefix(eventName, "GenerateDataKey") ||
		strings.HasPrefix(eventName, "Encrypt") {

		zap.L().Debug("kms: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	var keyARN string
	switch eventName {
	/*
		Missing (not sure if needed in all cases):
			(Connect/Create/Delete/Update)CustomKeyStore
			(Delete/Import)KeyMaterial
			(Retire/Revoke)Grant
	*/
	case "CancelKeyDeletion", "CreateGrant", "CreateKey", "DisableKey", "DisableKeyRotation", "EnableKey",
		"EnableKeyRotation", "PutKeyPolicy", "ScheduleKeyDeletion", "TagResource", "UntagResource",
		"UpdateAlias", "UpdateKeyDescription":
		keyARN = detail.Get("resources").Array()[0].Get("ARN").Str
	case "CreateAlias", "DeleteAlias":
		resources := detail.Get("resources").Array()
		for _, resource := range resources {
			resourceARN, err := arn.Parse(resource.Get("ARN").Str)
			if err != nil {
				zap.L().Error("kms: unable to extract ARN", zap.String("eventName", eventName))
				return nil
			}
			if strings.HasPrefix(resourceARN.Resource, "key/") {
				keyARN = resourceARN.String()
			}
		}
	default:
		zap.L().Warn("kms: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       eventName == "DeleteKey",
		EventName:    eventName,
		ResourceID:   keyARN,
		ResourceType: schemas.KmsKeySchema,
	}}
}
