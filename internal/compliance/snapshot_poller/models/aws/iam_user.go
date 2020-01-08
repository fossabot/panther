package aws

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
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	// IAMRootUserSchema is the schema identifier for IAMRootUser.
	IAMRootUserSchema = "AWS.IAM.RootUser"
	// IAMUserSchema is the schema identifier for IAMUser.
	IAMUserSchema = "AWS.IAM.User"
)

// IAMUser contains all information about an IAM User
type IAMUser struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from iam.User
	PasswordLastUsed    *time.Time
	Path                *string
	PermissionsBoundary *iam.AttachedPermissionsBoundary

	// Additional fields
	CredentialReport   *IAMCredentialReport
	Groups             []*iam.Group
	InlinePolicies     map[string]*string
	ManagedPolicyNames []*string
	VirtualMFA         *VirtualMFADevice
}

// IAMRootUser extends IAMUser, and contains some additional
// information only pertinent to the root account.
type IAMRootUser struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Additional fields
	CredentialReport *IAMCredentialReport
	VirtualMFA       *VirtualMFADevice
}

// IAMCredentialReport provides information on IAM credentials in an AWS Account.
//
// This includes status of credentials, console passwords, access keys, MFA devices, and more.
type IAMCredentialReport struct {
	UserName                  *string
	ARN                       *string
	UserCreationTime          *time.Time
	PasswordEnabled           *bool
	PasswordLastUsed          *time.Time
	PasswordLastChanged       *time.Time
	PasswordNextRotation      *time.Time
	MfaActive                 *bool
	AccessKey1Active          *bool
	AccessKey1LastRotated     *time.Time
	AccessKey1LastUsedDate    *time.Time
	AccessKey1LastUsedRegion  *string
	AccessKey1LastUsedService *string
	AccessKey2Active          *bool
	AccessKey2LastRotated     *time.Time
	AccessKey2LastUsedDate    *time.Time
	AccessKey2LastUsedRegion  *string
	AccessKey2LastUsedService *string
	Cert1Active               *bool
	Cert1LastRotated          *time.Time
	Cert2Active               *bool
	Cert2LastRotated          *time.Time
}

// VirtualMFADevice provides metadata about an IAM User's MFA device
type VirtualMFADevice struct {
	EnableDate   *time.Time
	SerialNumber *string
}
