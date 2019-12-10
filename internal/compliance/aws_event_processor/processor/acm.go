package processor

import (
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/api/snapshot/aws"
)

func classifyACM(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_awscertificatemanager.html
	if eventName == "ExportCertificate" ||
		eventName == "ResendValidationEmail" {

		zap.L().Debug("acm: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	var certARN string
	switch eventName {
	case "AddTagsToCertificate", "DeleteCertificate", "RemoveTags", "RenewCertificate", "UpdateCertificateOptions":
		certARN = detail.Get("requestParameters.certificateArn").Str
	case "ImportCertificate", "RequestCertificate":
		certARN = detail.Get("responseElements.certificateArn").Str
	default:
		zap.L().Warn("acm: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       eventName == "DeleteCertificate",
		EventName:    eventName,
		ResourceID:   certARN,
		ResourceType: schemas.AcmCertificateSchema,
	}}
}
