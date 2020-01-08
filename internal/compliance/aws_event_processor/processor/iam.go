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
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

const rootUserName = "AWS ROOT USER"

func classifyIAM(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_identityandaccessmanagement.html
	if eventName == "ChangePassword" ||
		eventName == "ResetServiceSpecificCredential" ||
		eventName == "GenerateCredentialReport" ||
		eventName == "CreateVirtualMFADevice" || // MFA device creation/deletion is not related to
		eventName == "DeleteVirtualMFADevice" || // users. See (Enable/Disable)MFADevice for that.
		eventName == "CreateInstanceProfile" {

		zap.L().Debug("iam: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	var resourceType string
	var err error
	resourceDelete := false
	iamARN := arn.ARN{
		Partition: "aws",
		Service:   "iam",
		Region:    "",
		AccountID: accountID,
	}
	switch eventName {
	case "AddRoleToInstanceProfile", "DeleteRolePermissionsBoundary", "DeleteRolePolicy":
		resourceType = aws.IAMRoleSchema
		iamARN.Resource = "role/" + detail.Get("requestParameters.roleName").Str
	case "AddUserToGroup", "RemoveUserFromGroup":
		userName := detail.Get("requestParameters.userName").Str
		var userType string

		// Not sure if it is actually possible to add the root user to a group
		if userName == rootUserName {
			userName = "root"
			userType = aws.IAMRootUserSchema
		} else {
			userName = "user/" + userName
			userType = aws.IAMUserSchema
		}

		return []*resourceChange{
			{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   iamARN.String() + "group/" + detail.Get("requestParameters.groupName").Str,
				ResourceType: aws.IAMGroupSchema,
			},
			{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   iamARN.String() + userName,
				ResourceType: userType,
			},
		}
	case "AttachGroupPolicy", "DetachGroupPolicy":
		return []*resourceChange{
			{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   detail.Get("requestParameters.policyArn").Str,
				ResourceType: aws.IAMPolicySchema,
			},
			{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   iamARN.String() + "group/" + detail.Get("requestParameters.groupName").Str,
				ResourceType: aws.IAMGroupSchema,
			},
		}
	case "AttachRolePolicy", "DetachRolePolicy":
		return []*resourceChange{
			{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   detail.Get("requestParameters.policyArn").Str,
				ResourceType: aws.IAMPolicySchema,
			},
			{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   iamARN.String() + "role/" + detail.Get("requestParameters.roleName").Str,
				ResourceType: aws.IAMRoleSchema,
			},
		}
	case "AttachUserPolicy", "DetachUserPolicy":
		userName := detail.Get("requestParameters.userName").Str
		var userType string

		if userName == rootUserName {
			userName = "root"
			userType = aws.IAMRootUserSchema
		} else {
			userName = "user/" + userName
			userType = aws.IAMUserSchema
		}

		return []*resourceChange{
			{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   detail.Get("requestParameters.policyArn").Str,
				ResourceType: aws.IAMPolicySchema,
			},
			{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   iamARN.String() + userName,
				ResourceType: userType,
			},
		}
	case "CreateAccessKey", "CreateLoginProfile", "CreateServiceSpecificCredential", "CreateUser", "DeactivateMFADevice",
		"DeleteLoginProfile", "DeleteSSHPublicKey", "DeleteServiceSpecificCredential", "DeleteSigningCertificate",
		"DeleteUserPermissionsBoundary", "DeleteUserPolicy", "EnableMFADevice", "PutUserPermissionsBoundary",
		"PutUserPolicy", "TagUser", "UntagUser", "UpdateAccessKey", "UpdateLoginProfile", "UpdateSSHPublicKey",
		"UpdateServiceSpecificCredential", "UpdateSigningCertificate", "UploadSSHPublicKey", "UploadSigningCertificate",
		"DeleteAccessKey":
		// Attempt to get the userName the easy way
		userName := detail.Get("requestParameters.userName").Str

		if userName == rootUserName {
			// We got root, set the appropriate fields and break
			resourceType = aws.IAMRootUserSchema
			iamARN.Resource = "root"
			break
		} else if userName != "" {
			// We got a normal user, set the appropriate fields and break
			resourceType = aws.IAMUserSchema
			iamARN.Resource = "user/" + detail.Get("requestParameters.userName").Str
			break
		}

		// Often, the userName may be omitted. If it is, AWS implicitly assumes it based on the AWS
		// access Key ID signing the request. This is most frequently used to manage the root account
		// credentials for an account that has no other users.
		iamARN, err = arn.Parse(detail.Get("userIdentity.arn").Str)

		// Error check here because we are about the use iamARN
		if err != nil {
			zap.L().Error("iam: error handling iam user event", zap.String("eventName", eventName), zap.Error(err))
			return nil
		}

		if iamARN.Resource == "root" {
			resourceType = aws.IAMRootUserSchema
		} else {
			resourceType = aws.IAMUserSchema
		}
	case "CreateGroup":
		iamARN, err = arn.Parse(detail.Get("responseElements.group.arn").Str)
		resourceType = aws.IAMGroupSchema
	case "CreatePolicy":
		iamARN, err = arn.Parse(detail.Get("responseElements.policy.arn").Str)
		resourceType = aws.IAMPolicySchema
	case "CreatePolicyVersion", "DeletePolicyVersion", "SetDefaultPolicyVersion":
		iamARN, err = arn.Parse(detail.Get("requestParameters.policyArn").Str)
		resourceType = aws.IAMPolicySchema
	case "CreateRole", "CreateServiceLinkedRole":
		iamARN, err = arn.Parse(detail.Get("responseElements.role.arn").Str)
		resourceType = aws.IAMRoleSchema
	case "DeleteAccountPasswordPolicy", "UpdateAccountPasswordPolicy":
		return []*resourceChange{
			{
				AwsAccountID: accountID,
				// We don't actually allow an account to not have a PasswordPolicy resource, it's
				// just marked as 'AnyExist = False' in it's attributes
				Delete:       false,
				EventName:    eventName,
				ResourceID:   accountID + "::" + aws.PasswordPolicySchema,
				ResourceType: aws.PasswordPolicySchema,
			},
		}
	case "DeleteGroup":
		resourceType = aws.IAMGroupSchema
		resourceDelete = true
		iamARN.Resource = "group/" + detail.Get("requestParameters.groupName").Str
	case "DeleteGroupPolicy", "PutGroupPolicy":
		// This is deleting an inline policy, which we're not tracking as a separate resource
		resourceType = aws.IAMGroupSchema
		iamARN.Resource = "group/" + detail.Get("requestParameters.groupName").Str
	case "DeletePolicy":
		resourceType = aws.IAMPolicySchema
		resourceDelete = true
		iamARN, err = arn.Parse(detail.Get("requestParameters.policyArn").Str)
	case "DeleteRole", "DeleteServiceLinkedRole":
		resourceType = aws.IAMRoleSchema
		resourceDelete = true
		iamARN.Resource = "role/" + detail.Get("requestParameters.roleName").Str
	case "DeleteUser":
		resourceType = aws.IAMUserSchema
		resourceDelete = true
		iamARN.Resource = "user/" + detail.Get("requestParameters.userName").Str
	case "PutRolePermissionsBoundary", "PutRolePolicy", "RemoveRoleFromInstanceProfile", "TagRole",
		"UntagRole", "UpdateAssumeRolePolicy", "UpdateRole", "UpdateRoleDescription":
		resourceType = aws.IAMRoleSchema
		iamARN.Resource = "role/" + detail.Get("requestParameters.roleName").Str
	case "UpdateGroup":
		// Special case cause the name could change, thus changing the ARN. We handle this by creating
		// a new group resource, and allowing the old one to eventually time out
		resourceType = aws.IAMGroupSchema
		newName := detail.Get("requestParameters.newGroupName").Str
		if newName != "" {
			iamARN.Resource = "group/" + newName
		} else {
			iamARN.Resource = "group/" + detail.Get("requestParameters.groupName").Str
		}
	case "UpdateUser":
		// Same special case as UpdateGroup
		newName := detail.Get("requestParameters.newUserName").Str
		resourceType = aws.IAMUserSchema
		if newName != "" {
			iamARN.Resource = "user/" + newName
		} else {
			userName := detail.Get("requestParameters.userName").Str
			if userName == rootUserName {
				resourceType = aws.IAMRootUserSchema
				iamARN.Resource = "root"
			} else {
				iamARN.Resource = "user/" + userName
			}
		}
	default:
		zap.L().Warn("iam: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	if err != nil {
		zap.L().Error(
			"iam: error occurred during event processing",
			zap.String("eventName", eventName),
			zap.Error(err))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: iamARN.AccountID,
		Delete:       resourceDelete,
		EventName:    eventName,
		ResourceID:   iamARN.String(),
		ResourceType: resourceType,
	}}
}
