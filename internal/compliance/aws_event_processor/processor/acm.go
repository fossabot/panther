package processor

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
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
