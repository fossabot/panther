package processor

import (
	"strings"

	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

func classifyGuardDuty(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_amazonguardduty.html
	if eventName == "ArchiveFindings" ||
		eventName == "CreateIPSet" ||
		eventName == "CreateSampleFindings" ||
		eventName == "CreateThreatIntelSet" ||
		eventName == "DeclineInvitations" ||
		eventName == "DeleteFilter" ||
		eventName == "DeleteIPSet" ||
		eventName == "DeleteInvitations" ||
		eventName == "DeleteThreatIntelSet" ||
		eventName == "InviteMembers" ||
		eventName == "UnarchiveFindings" ||
		eventName == "UpdateFilter" ||
		eventName == "UpdateFindingsFeedback" ||
		eventName == "UpdateIPSet" ||
		eventName == "UpdateThreatIntelSet" ||
		eventName == "CreateFilter" {

		zap.L().Debug("guardduty: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	switch eventName {
	case "TagResource", "UntagResource", "UpdateDetector":
		// Single resource/region scan (only one detector can exist per region)
		return []*resourceChange{{
			AwsAccountID: accountID,
			EventName:    eventName,
			ResourceID: strings.Join([]string{
				accountID,
				detail.Get("awsRegion").Str,
				schemas.GuardDutySchema,
			}, ":"),
			ResourceType: schemas.GuardDutySchema,
		}}
	case "AcceptInvitation", "CreateDetector", "CreateMembers", "DeleteMembers", "DisassociateFromMasterAccount",
		"DisassociateMembers", "StartMonitoringMembers", "StopMonitoringMembers":
		// Full account scan
		return []*resourceChange{{
			AwsAccountID: accountID,
			EventName:    eventName,
			ResourceType: schemas.GuardDutySchema,
		}}
	case "DeleteDetector":
		// Special case where need to queue both a delete action and a meta re-scan
		return []*resourceChange{
			{
				AwsAccountID: accountID,
				Delete:       true,
				EventName:    eventName,
				ResourceID: strings.Join([]string{
					accountID,
					detail.Get("awsRegion").Str,
					schemas.GuardDutySchema,
				}, ":"),
				ResourceType: schemas.GuardDutySchema,
			},
			{
				AwsAccountID: accountID,
				EventName:    eventName,
				ResourceType: schemas.GuardDutySchema,
			}}
	default:
		zap.L().Warn("guardduty: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}
}
