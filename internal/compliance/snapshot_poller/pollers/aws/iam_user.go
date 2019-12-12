package aws

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/cenkalti/backoff"
	"go.uber.org/zap"

	apimodels "github.com/panther-labs/panther/api/gateway/resources/models"
	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/utils"
)

const (
	rootAccountNameCredReport = "<root_account>"
	rootDeviceSerialSuffix    = ":mfa/root-account-mfa-device"
)

var (
	maxCredReportBackoff  = 1 * time.Minute
	userCredentialReports map[string]*awsmodels.IAMCredentialReport
	mfaDeviceMapping      map[string]*awsmodels.VirtualMFADevice
)

// PollIAMUser polls a single IAM User resource
func PollIAMUser(
	pollerResourceInput *awsmodels.ResourcePollerInput,
	resourceARN arn.ARN,
	scanRequest *pollermodels.ScanEntry,
) interface{} {

	client := getClient(pollerResourceInput, "iam", defaultRegion).(iamiface.IAMAPI)
	// See PollIAMRole for an explanation of this behavior
	resourceSplit := strings.Split(resourceARN.Resource, "/")
	user := getUser(client, aws.String(resourceSplit[len(resourceSplit)-1]))

	// Refresh the caches as needed
	var err error
	mfaDeviceMapping, err = listVirtualMFADevices(client)
	if err != nil {
		utils.LogAWSError("IAM.ListVirtualMFADevices", err)
		return nil
	}
	userCredentialReports, err = buildCredentialReport(client)
	if err != nil {
		zap.L().Error("failed to build credential report", zap.Error(err))
		return nil
	}

	snapshot := buildIAMUserSnapshot(client, user)
	if snapshot == nil {
		return nil
	}

	// If the user does not have a credential report, then continue on with the snapshot but
	// re-queue the user for a scan in fifteen minutes (the maximum delay time). The primary reason
	// a user would not have a credential report is if they were recently created and there has not
	// yet been time for a new credential report that includes them to have been generated.
	if snapshot.CredentialReport == nil {
		utils.Requeue(pollermodels.ScanMsg{
			Entries: []*pollermodels.ScanEntry{
				scanRequest,
			},
		}, utils.MaxRequeueDelaySeconds)
	}

	snapshot.AccountID = aws.String(resourceARN.AccountID)
	scanRequest.ResourceID = snapshot.ResourceID
	return snapshot
}

// PollIAMUser polls a single IAM User resource
func PollIAMRootUser(
	pollerResourceInput *awsmodels.ResourcePollerInput,
	_ arn.ARN,
	scanRequest *pollermodels.ScanEntry,
) interface{} {

	client := getClient(pollerResourceInput, "iam", defaultRegion).(iamiface.IAMAPI)
	// Refresh the caches as needed
	var err error
	mfaDeviceMapping, err = listVirtualMFADevices(client)
	if err != nil {
		utils.LogAWSError("IAM.ListVirtualMFADevices", err)
		return nil
	}
	userCredentialReports, err = buildCredentialReport(client)
	if err != nil {
		zap.L().Error("failed to build credential report", zap.Error(err))
		return nil
	}

	snapshot := buildIAMRootUserSnapshot()
	// If the root user does not have a credential report, then continue on with the snapshot but
	// re-queue the for a scan in fifteen minutes (the maximum delay time).
	if snapshot == nil || snapshot.CredentialReport == nil {
		utils.Requeue(pollermodels.ScanMsg{
			Entries: []*pollermodels.ScanEntry{
				scanRequest,
			},
		}, utils.MaxRequeueDelaySeconds)
		return nil
	}

	// Over ride this as it may be set incorrectly
	scanRequest.ResourceID = snapshot.ResourceID
	return snapshot
}

// getUser returns an individual IAM user
func getUser(svc iamiface.IAMAPI, userName *string) *iam.User {
	user, err := svc.GetUser(&iam.GetUserInput{
		UserName: userName,
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NoSuchEntity" {
				zap.L().Warn("tried to scan non-existent resource",
					zap.String("resource", *userName),
					zap.String("resourceType", awsmodels.IAMUserSchema))
				return nil
			}
		}
		utils.LogAWSError("IAM.GetUser", err)
		return nil
	}
	return user.User
}

// getCredentialReport retrieves an existing credential report from AWS
func getCredentialReport(svc iamiface.IAMAPI) (*iam.GetCredentialReportOutput, error) {
	var getIn = &iam.GetCredentialReportInput{}
	var getOut *iam.GetCredentialReportOutput
	var getErr error

	if getOut, getErr = svc.GetCredentialReport(getIn); getErr != nil {
		return nil, getErr
	}

	return getOut, nil
}

// generateCredentialReport generates a credential report if one does not exist,
// and does not return until the report has been successfully generated.
func generateCredentialReport(svc iamiface.IAMAPI) (*iam.GenerateCredentialReportOutput, error) {
	var genIn = &iam.GenerateCredentialReportInput{}
	var genOut *iam.GenerateCredentialReportOutput
	var genErr error

	backoffOperation := func() error {
		if genOut, genErr = svc.GenerateCredentialReport(genIn); genErr != nil {
			return backoff.Permanent(genErr)
		}
		if *genOut.State != "COMPLETE" {
			return errors.New("report in progress")
		}
		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = maxCredReportBackoff
	backoffErr := backoff.Retry(backoffOperation, expBackoff)
	if backoffErr != nil {
		return nil, backoffErr
	}

	return genOut, nil
}

func parseCredReportBool(field string) bool {
	convertedBool, err := strconv.ParseBool(field)
	if err != nil {
		return false
	}

	return convertedBool
}

// extractCredentialReport converts a CSV credential report into a mapping of user to parsed report.
func extractCredentialReport(content []byte) (map[string]*awsmodels.IAMCredentialReport, error) {
	csvReader := csv.NewReader(bytes.NewReader(content))
	userCredReportMapping := make(map[string]*awsmodels.IAMCredentialReport)

	credReportRows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Iterate through all cred report rows, and skip the header row
	for _, credReportRow := range credReportRows[1:] {
		if len(credReportRow) != 22 {
			zap.L().Error("invalid credential report row")
			continue
		}

		credReport := &awsmodels.IAMCredentialReport{
			UserName:                  aws.String(credReportRow[0]),
			ARN:                       aws.String(credReportRow[1]),
			UserCreationTime:          aws.Time(utils.ParseTimeRFC3339(credReportRow[2])),
			PasswordEnabled:           aws.Bool(parseCredReportBool(credReportRow[3])),
			PasswordLastUsed:          aws.Time(utils.ParseTimeRFC3339(credReportRow[4])),
			PasswordLastChanged:       aws.Time(utils.ParseTimeRFC3339(credReportRow[5])),
			PasswordNextRotation:      aws.Time(utils.ParseTimeRFC3339(credReportRow[6])),
			MfaActive:                 aws.Bool(parseCredReportBool(credReportRow[7])),
			AccessKey1Active:          aws.Bool(parseCredReportBool(credReportRow[8])),
			AccessKey1LastRotated:     aws.Time(utils.ParseTimeRFC3339(credReportRow[9])),
			AccessKey1LastUsedDate:    aws.Time(utils.ParseTimeRFC3339(credReportRow[10])),
			AccessKey1LastUsedRegion:  aws.String(credReportRow[11]),
			AccessKey1LastUsedService: aws.String(credReportRow[12]),
			AccessKey2Active:          aws.Bool(parseCredReportBool(credReportRow[13])),
			AccessKey2LastRotated:     aws.Time(utils.ParseTimeRFC3339(credReportRow[14])),
			AccessKey2LastUsedDate:    aws.Time(utils.ParseTimeRFC3339(credReportRow[15])),
			AccessKey2LastUsedRegion:  aws.String(credReportRow[16]),
			AccessKey2LastUsedService: aws.String(credReportRow[17]),
			Cert1Active:               aws.Bool(parseCredReportBool(credReportRow[18])),
			Cert1LastRotated:          aws.Time(utils.ParseTimeRFC3339(credReportRow[19])),
			Cert2Active:               aws.Bool(parseCredReportBool(credReportRow[20])),
			Cert2LastRotated:          aws.Time(utils.ParseTimeRFC3339(credReportRow[21])),
		}
		userCredReportMapping[credReportRow[0]] = credReport
	}

	return userCredReportMapping, nil
}

// buildCredentialReport obtains an IAM Credential Report and generates a mapping from user to report.
func buildCredentialReport(
	iamSvc iamiface.IAMAPI) (map[string]*awsmodels.IAMCredentialReport, error) {

	var credentialReportRaw *iam.GetCredentialReportOutput
	var err error

	// Try to get the credential report
	credentialReportRaw, err = getCredentialReport(iamSvc)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case iam.ErrCodeCredentialReportNotPresentException, iam.ErrCodeCredentialReportExpiredException:
				zap.L().Debug("no credential report found, generating a new one")
				if _, err := generateCredentialReport(iamSvc); err != nil {
					return nil, err
				}
				credentialReportRaw, err = getCredentialReport(iamSvc)
				if err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}

	return extractCredentialReport(credentialReportRaw.Content)
}

// listUsers returns all the users in the account, excluding the root account.
func listUsers(iamSvc iamiface.IAMAPI) (users []*iam.User) {
	err := iamSvc.ListUsersPages(
		&iam.ListUsersInput{},
		func(page *iam.ListUsersOutput, lastPage bool) bool {
			users = append(users, page.Users...)
			return true
		},
	)
	if err != nil {
		utils.LogAWSError("IAM.ListUsersPages", err)
	}
	return
}

// getUserPolicies aggregates all the policies assigned to a user by polling both
// the ListUserPolicies and ListAttachedUserPolicies APIs.
func getUserPolicies(iamSvc iamiface.IAMAPI, userName *string) (
	inlinePolicies []*string, managedPolicies []*string, err error) {

	err = iamSvc.ListUserPoliciesPages(
		&iam.ListUserPoliciesInput{UserName: userName},
		func(page *iam.ListUserPoliciesOutput, lastPage bool) bool {
			inlinePolicies = append(inlinePolicies, page.PolicyNames...)
			return true
		},
	)
	if err != nil {
		utils.LogAWSError("IAM.ListUserPolicies", err)
	}

	err = iamSvc.ListAttachedUserPoliciesPages(
		&iam.ListAttachedUserPoliciesInput{UserName: userName},
		func(page *iam.ListAttachedUserPoliciesOutput, lastPage bool) bool {
			for _, attachedPolicy := range page.AttachedPolicies {
				managedPolicies = append(managedPolicies, attachedPolicy.PolicyName)
			}
			return true
		},
	)
	if err != nil {
		utils.LogAWSError("IAM.ListAttachedUserPolicies", err)
	}

	return
}

// listVirtualMFADevices returns a mapping of UserID to VirtualMFADeviceSnapshot.
func listVirtualMFADevices(
	iamSvc iamiface.IAMAPI) (map[string]*awsmodels.VirtualMFADevice, error) {

	vmfaDevicesInput := &iam.ListVirtualMFADevicesInput{
		// We only want MFA devices associated with a user.
		AssignmentStatus: aws.String("Assigned"),
	}
	var vmfaDevices []*iam.VirtualMFADevice
	err := iamSvc.ListVirtualMFADevicesPages(
		vmfaDevicesInput,
		func(page *iam.ListVirtualMFADevicesOutput, lastPage bool) bool {
			vmfaDevices = append(vmfaDevices, page.VirtualMFADevices...)
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	mfaDeviceMapping := make(map[string]*awsmodels.VirtualMFADevice)
	for _, vmfaDevice := range vmfaDevices {
		if vmfaDevice.User != nil && vmfaDevice.User.UserId != nil {
			mfaDeviceMapping[*vmfaDevice.User.UserId] = &awsmodels.VirtualMFADevice{
				EnableDate:   vmfaDevice.EnableDate,
				SerialNumber: vmfaDevice.SerialNumber,
			}
		}
	}

	return mfaDeviceMapping, nil
}

// listGroupsForUser returns all the IAM Groups a given IAM User belongs to
func listGroupsForUser(iamSvc iamiface.IAMAPI, userName *string) (groups []*iam.Group) {
	err := iamSvc.ListGroupsForUserPages(&iam.ListGroupsForUserInput{UserName: userName},
		func(page *iam.ListGroupsForUserOutput, lastPage bool) bool {
			groups = append(groups, page.Groups...)
			return true
		})
	if err != nil {
		utils.LogAWSError("IAM.ListGroupsForUserPages", err)
	}
	return
}

// getUserPolicy gets the inline policy documents for a given IAM user and inline policy name
func getUserPolicy(svc iamiface.IAMAPI, userName *string, policyName *string) *string {
	policy, err := svc.GetUserPolicy(&iam.GetUserPolicyInput{
		UserName:   userName,
		PolicyName: policyName,
	})

	if err != nil {
		utils.LogAWSError("IAM.GetUserPolicy", err)
		return nil
	}

	decodedPolicy, err := url.QueryUnescape(*policy.PolicyDocument)
	if err != nil {
		zap.L().Error("IAM: unable to url decode inline policy document",
			zap.String("policy document", *policy.PolicyDocument),
			zap.String("policy name", *policyName),
			zap.String("user", *userName),
		)
		return nil
	}

	return aws.String(decodedPolicy)
}

// buildIAMUserSnapshot builds an IAMUserSnapshot for a given IAM User
func buildIAMUserSnapshot(iamSvc iamiface.IAMAPI, user *iam.User) *awsmodels.IAMUser {
	if user == nil {
		return nil
	}
	iamUserSnapshot := &awsmodels.IAMUser{
		GenericResource: awsmodels.GenericResource{
			ResourceID:   user.Arn,
			TimeCreated:  utils.DateTimeFormat(*user.CreateDate),
			ResourceType: aws.String(awsmodels.IAMUserSchema),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			ARN:    user.Arn,
			ID:     user.UserId,
			Name:   user.UserName,
			Region: aws.String(awsmodels.GlobalRegion),
			Tags:   utils.ParseTagSlice(user.Tags),
		},
		PasswordLastUsed:    user.PasswordLastUsed,
		Path:                user.Path,
		PermissionsBoundary: user.PermissionsBoundary,
	}

	// Get IAM Policies associated to the user.
	// There is no error logging here because it is logged in getUserPolicies.
	inlinePolicyNames, managedPolicies, err := getUserPolicies(iamSvc, user.UserName)
	if err == nil {
		iamUserSnapshot.ManagedPolicyNames = managedPolicies
		if inlinePolicyNames != nil {
			iamUserSnapshot.InlinePolicies = make(map[string]*string, len(inlinePolicyNames))
			for _, inlinePolicy := range inlinePolicyNames {
				iamUserSnapshot.InlinePolicies[*inlinePolicy] = getUserPolicy(iamSvc, user.UserName, inlinePolicy)
			}
		}
	}

	// Build the credential report for all users if they don't exist already
	if userCredentialReports != nil {
		if userCredentialReport, ok := userCredentialReports[*user.UserName]; ok {
			iamUserSnapshot.CredentialReport = userCredentialReport
		}
	}

	// Look up any virtual MFA devices attached to the user
	if mfaDeviceMapping != nil {
		if mfaSnapshot, ok := mfaDeviceMapping[*user.UserId]; ok {
			iamUserSnapshot.VirtualMFA = mfaSnapshot
		}
	}

	// Look up any groups the user is a member of
	iamUserSnapshot.Groups = listGroupsForUser(iamSvc, user.UserName)

	return iamUserSnapshot
}

func buildIAMRootUserSnapshot() *awsmodels.IAMRootUser {
	rootCredReport, ok := userCredentialReports[rootAccountNameCredReport]
	if !ok {
		zap.L().Error("unable to find credential report for root user",
			zap.Any("credential report", userCredentialReports))
		return nil
	}

	rootARN, err := arn.Parse(*rootCredReport.ARN)
	if err != nil {
		zap.L().Error(
			"unable to extract root user account ID",
			zap.String("root arn", *rootCredReport.ARN),
			zap.Error(err))
		return nil
	}
	rootSnapshot := &awsmodels.IAMRootUser{
		GenericResource: awsmodels.GenericResource{
			ResourceID:   rootCredReport.ARN,
			TimeCreated:  utils.DateTimeFormat(*rootCredReport.UserCreationTime),
			ResourceType: aws.String(awsmodels.IAMRootUserSchema),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			AccountID: aws.String(rootARN.AccountID),
			ARN:       rootCredReport.ARN,
			Name:      rootCredReport.UserName,
			Region:    aws.String(awsmodels.GlobalRegion),
		},
		CredentialReport: rootCredReport,
	}

	// Add final MFA and UserID fields to Root Snapshot
	for userID, vMFADeviceSnapshot := range mfaDeviceMapping {
		if strings.HasSuffix(*vMFADeviceSnapshot.SerialNumber, rootDeviceSerialSuffix) {
			rootSnapshot.ID = aws.String(userID)
			rootSnapshot.VirtualMFA = vMFADeviceSnapshot
		}
	}

	return rootSnapshot
}

// PollIAMUsers generates a snapshot for each IAM User.
//
// This function returns a slice of Events.
func PollIAMUsers(pollerInput *awsmodels.ResourcePollerInput) ([]*apimodels.AddResourceEntry, error) {
	zap.L().Debug("starting IAM User resource poller")
	sess := session.Must(session.NewSession(&aws.Config{}))
	creds, err := AssumeRoleFunc(pollerInput, sess)
	if err != nil {
		return nil, err
	}

	iamSvc := IAMClientFunc(sess, &aws.Config{Credentials: creds}).(iamiface.IAMAPI)

	// List all IAM Users in the account
	users := listUsers(iamSvc)
	if len(users) == 0 {
		zap.L().Debug("no IAM users found")
	}

	// Build the credential report for all users
	userCredentialReports, err = buildCredentialReport(iamSvc)
	if err != nil {
		zap.L().Error("failed to build credential report", zap.Error(err))
		return nil, err
	}

	// Get all VMFA snapshots
	mfaDeviceMapping, err = listVirtualMFADevices(iamSvc)
	if err != nil {
		utils.LogAWSError("IAM.ListVirtualMFADevices", err)
		return nil, err
	}

	// Create IAM User snapshots
	var resources []*apimodels.AddResourceEntry
	for _, user := range users {
		// The IAM.User struct has a Tags field, indicating what tags the User has
		// The API call IAM.GetUser returns an IAM.User struct, with all appropriate fields set
		// The API call IAM.ListUsers returns a slice of IAM.User structs, but does not set the tags
		// field for any of these structs regardless of whether the corresponding user has tags set
		// This patches that gap
		fullUser := getUser(iamSvc, user.UserName)
		iamUserSnapshot := buildIAMUserSnapshot(iamSvc, fullUser)
		if iamUserSnapshot == nil {
			continue
		}
		iamUserSnapshot.AccountID = aws.String(pollerInput.AuthSourceParsedARN.AccountID)
		// If the user does not have a credential report, then continue on with the snapshot but
		// re-queue the user for a scan in fifteen minutes (the maximum delay time). The primary reason
		// a user would not have a credential report is if they were recently created and there has not
		// yet been time for a new credential report that includes them to have been generated.
		if iamUserSnapshot.CredentialReport == nil {
			utils.Requeue(pollermodels.ScanMsg{
				Entries: []*pollermodels.ScanEntry{
					{
						AWSAccountID:  iamUserSnapshot.AccountID,
						IntegrationID: pollerInput.IntegrationID,
						ResourceID:    iamUserSnapshot.ResourceID,
						ResourceType:  iamUserSnapshot.ResourceType,
					},
				},
			}, utils.MaxRequeueDelaySeconds)
		}

		resources = append(resources, &apimodels.AddResourceEntry{
			Attributes:      iamUserSnapshot,
			ID:              apimodels.ResourceID(*user.Arn),
			IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
			IntegrationType: apimodels.IntegrationTypeAws,
			Type:            awsmodels.IAMUserSchema,
		})
	}

	rootSnapshot := buildIAMRootUserSnapshot()
	if rootSnapshot == nil {
		// Re-scan the root user if there was any error building it
		utils.Requeue(pollermodels.ScanMsg{
			Entries: []*pollermodels.ScanEntry{
				{
					AWSAccountID:  aws.String(pollerInput.AuthSourceParsedARN.AccountID),
					IntegrationID: pollerInput.IntegrationID,
					ResourceType:  aws.String(awsmodels.IAMRootUserSchema),
				},
			},
		}, utils.MaxRequeueDelaySeconds)
	} else {
		// Create the IAM Root User snapshot
		resources = append(resources, &apimodels.AddResourceEntry{
			Attributes:      rootSnapshot,
			ID:              apimodels.ResourceID(*rootSnapshot.ARN),
			IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
			IntegrationType: apimodels.IntegrationTypeAws,
			Type:            awsmodels.IAMRootUserSchema,
		})
	}

	return resources, nil
}
