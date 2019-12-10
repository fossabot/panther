package processor

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/api/snapshot/aws"
)

func classifyWAF(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// These cases are tough because they don't link these resources back to any attached Web ACLs,
	// of which there could be several. Just scan all web ACLs for now until there is a link table
	// or a sub-resource for each of these. This catches 11 API calls to WAF non Web ACL resources.
	if strings.HasPrefix(eventName, "Update") && eventName != "UpdateWebACL" {
		return []*resourceChange{{
			AwsAccountID: accountID,
			EventName:    eventName,
			Region:       schemas.GlobalRegion,
			ResourceType: schemas.WafWebAclSchema,
		}}
	}

	// All the API calls we don't care about (until we build resources for them)
	if strings.HasSuffix(eventName, "Set") || // 11
		strings.HasSuffix(eventName, "Rule") || // 6
		strings.HasSuffix(eventName, "RuleGroup") || // 3
		// Permission policies affect rule groups
		eventName == "DeletePermissionPolicy" ||
		eventName == "PutPermissionPolicy" {

		zap.L().Debug("waf: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_awswaf.html
	var wafARN string
	switch eventName {
	case "CreateWebACL":
		wafARN = detail.Get("responseElements.webACL.webACLArn").Str
	case "DeleteLoggingConfiguration":
		wafARN = detail.Get("requestParameters.resourceArn").Str
	case "DeleteWebACL", "UpdateWebACL":
		// arn:aws:waf::account-id:resource-type/resource-id
		wafARN = strings.Join([]string{
			"arn",
			"aws",     // Partition
			"waf",     // Service
			"",        // Region (global service so no region)
			accountID, // Account ID
			"webacl/" + detail.Get("requestParameters.webACLId").Str, // Resource-type/id
		}, ":")
	case "PutLoggingConfiguration":
		wafARN = detail.Get("requestParameters.loggingConfiguration.resourceArn").Str
	default:
		zap.L().Error("waf: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	parsedARN, err := arn.Parse(wafARN)
	if err != nil {
		zap.L().Warn("waf: error parsing ARN", zap.String("eventName", eventName), zap.Error(err))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       eventName == "DeleteWebACL",
		EventName:    eventName,
		ResourceID:   parsedARN.String(),
		ResourceType: schemas.WafWebAclSchema,
	}}
}
