package aws

import (
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"go.uber.org/zap"

	apimodels "github.com/panther-labs/panther/api/gateway/resources/models"
	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/utils"
)

const (
	localPolicyScope = "Local"
)

// PollIAMPolicy polls a single IAM Policy resource
func PollIAMPolicy(
	pollerResourceInput *awsmodels.ResourcePollerInput,
	resourceARN arn.ARN,
	scanRequest *pollermodels.ScanEntry,
) interface{} {

	client := getClient(pollerResourceInput, "iam", defaultRegion).(iamiface.IAMAPI)
	policy := getIAMPolicy(client, scanRequest.ResourceID)

	snapshot := buildIAMPolicySnapshot(client, policy)
	if snapshot == nil {
		return nil
	}
	snapshot.AccountID = aws.String(resourceARN.AccountID)
	return snapshot
}

// getPolicy returns a specific IAM policy
func getIAMPolicy(svc iamiface.IAMAPI, policyARN *string) *iam.Policy {
	policy, err := svc.GetPolicy(&iam.GetPolicyInput{
		PolicyArn: policyARN,
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NoSuchEntity" {
				zap.L().Warn("tried to scan non-existent resource",
					zap.String("resource", *policyARN),
					zap.String("resourceType", awsmodels.IAMPolicySchema))
				return nil
			}
		}
		utils.LogAWSError("IAM.GetPolicy", err)
		return nil
	}
	return policy.Policy
}

// listPolicies returns all IAM policies in the account
func listPolicies(iamSvc iamiface.IAMAPI) (policies []*iam.Policy, err error) {
	err = iamSvc.ListPoliciesPages(
		// We only want to scan Customer managed policies
		&iam.ListPoliciesInput{Scope: aws.String(localPolicyScope)},
		func(page *iam.ListPoliciesOutput, lastPage bool) bool {
			policies = append(policies, page.Policies...)
			return true
		},
	)

	return
}

// listEntitiesForPolicy returns the entities that have the given policy
func listEntitiesForPolicy(
	iamSvc iamiface.IAMAPI, arn *string) (entities *awsmodels.IAMPolicyEntities) {

	entities = &awsmodels.IAMPolicyEntities{}
	err := iamSvc.ListEntitiesForPolicyPages(
		&iam.ListEntitiesForPolicyInput{PolicyArn: arn},
		func(page *iam.ListEntitiesForPolicyOutput, lastPage bool) bool {
			entities.PolicyGroups = append(entities.PolicyGroups, page.PolicyGroups...)
			entities.PolicyRoles = append(entities.PolicyRoles, page.PolicyRoles...)
			entities.PolicyUsers = append(entities.PolicyUsers, page.PolicyUsers...)
			return true
		},
	)
	if err != nil {
		utils.LogAWSError("IAM.ListEntitiesForPolicyPages", err)
	}
	return
}

// getPolicyVersion returns a specific policy document given a policy ARN and version number
func getPolicyVersion(
	iamSvc iamiface.IAMAPI, arn *string, version *string) (policyDoc string, err error) {

	policy, err := iamSvc.GetPolicyVersion(
		&iam.GetPolicyVersionInput{PolicyArn: arn, VersionId: version},
	)
	if err != nil {
		return
	}

	policyDoc, err = url.QueryUnescape(*policy.PolicyVersion.Document)
	return
}

// buildIAMPolicySnapshot builds a complete IAMPolicySnapshot
func buildIAMPolicySnapshot(iamSvc iamiface.IAMAPI, policy *iam.Policy) *awsmodels.IAMPolicy {
	if policy == nil {
		return nil
	}

	policySnapshot := &awsmodels.IAMPolicy{
		GenericResource: awsmodels.GenericResource{
			ResourceID:   policy.Arn,
			TimeCreated:  utils.DateTimeFormat(*policy.CreateDate),
			ResourceType: aws.String(awsmodels.IAMPolicySchema),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			ARN:    policy.Arn,
			Name:   policy.PolicyName,
			ID:     policy.PolicyId,
			Region: aws.String(awsmodels.GlobalRegion),
		},
		AttachmentCount:               policy.AttachmentCount,
		DefaultVersionId:              policy.DefaultVersionId,
		Description:                   policy.Description,
		IsAttachable:                  policy.IsAttachable,
		Path:                          policy.Path,
		PermissionsBoundaryUsageCount: policy.PermissionsBoundaryUsageCount,
		UpdateDate:                    policy.UpdateDate,
	}

	policySnapshot.Entities = listEntitiesForPolicy(iamSvc, policy.Arn)

	policyDocument, err := getPolicyVersion(iamSvc, policy.Arn, policy.DefaultVersionId)
	if err != nil {
		utils.LogAWSError("IAM.GetPolicyVersion", err)
	} else {
		policySnapshot.PolicyDocument = aws.String(policyDocument)
	}

	return policySnapshot
}

// PollIamPolicies gathers information on each IAM policy for an AWS account.
func PollIamPolicies(pollerInput *awsmodels.ResourcePollerInput) ([]*apimodels.AddResourceEntry, error) {
	zap.L().Debug("starting IAM Policy resource poller")
	sess := session.Must(session.NewSession(&aws.Config{}))
	creds, err := AssumeRoleFunc(pollerInput, sess)
	if err != nil {
		return nil, err
	}

	iamSvc := IAMClientFunc(sess, &aws.Config{Credentials: creds}).(iamiface.IAMAPI)

	// Start with generating a list of all policies
	policies, listErr := listPolicies(iamSvc)
	if listErr != nil {
		utils.LogAWSError("IAM.ListPolicies", listErr)
		return nil, nil
	}

	var resources []*apimodels.AddResourceEntry
	for _, policy := range policies {
		iamPolicySnapshot := buildIAMPolicySnapshot(iamSvc, policy)
		if iamPolicySnapshot == nil {
			continue
		}
		iamPolicySnapshot.AccountID = aws.String(pollerInput.AuthSourceParsedARN.AccountID)

		resources = append(resources, &apimodels.AddResourceEntry{
			Attributes:      iamPolicySnapshot,
			ID:              apimodels.ResourceID(*iamPolicySnapshot.ARN),
			IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
			IntegrationType: apimodels.IntegrationTypeAws,
			Type:            awsmodels.IAMPolicySchema,
		})
	}

	return resources, nil
}