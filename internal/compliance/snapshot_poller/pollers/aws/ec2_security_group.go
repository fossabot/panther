package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"go.uber.org/zap"

	apimodels "github.com/panther-labs/panther/api/resources/models"
	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/utils"
)

// PollEC2SecurityGroup polls a single EC2 Security Group resource
func PollEC2SecurityGroup(
	pollerResourceInput *awsmodels.ResourcePollerInput,
	resourceARN arn.ARN,
	scanRequest *pollermodels.ScanEntry,
) interface{} {

	client := getClient(pollerResourceInput, "ec2", resourceARN.Region).(ec2iface.EC2API)
	sgID := strings.Replace(resourceARN.Resource, "security-group/", "", 1)
	securityGroup := getSecurityGroup(client, aws.String(sgID))

	snapshot := buildEc2SecurityGroupSnapshot(client, securityGroup)
	if snapshot == nil {
		return nil
	}
	snapshot.ResourceID = scanRequest.ResourceID
	snapshot.AccountID = aws.String(resourceARN.AccountID)
	snapshot.Region = aws.String(resourceARN.Region)
	snapshot.ARN = scanRequest.ResourceID
	return snapshot
}

// getSecurityGroup returns a specific EC2 security group
func getSecurityGroup(svc ec2iface.EC2API, securityGroupID *string) *ec2.SecurityGroup {
	securityGroup, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{securityGroupID},
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "InvalidGroup.NotFound" {
				zap.L().Warn("tried to scan non-existent resource",
					zap.String("resource", *securityGroupID),
					zap.String("resourceType", awsmodels.Ec2SecurityGroupSchema))
				return nil
			}
		}
		utils.LogAWSError("EC2.DescribeSecurityGroups", err)
		return nil
	}

	return securityGroup.SecurityGroups[0]
}

// describeSecurityGroupsPages returns all Security Groups for a given region
func describeSecurityGroups(ec2Svc ec2iface.EC2API) (securityGroups []*ec2.SecurityGroup) {
	err := ec2Svc.DescribeSecurityGroupsPages(&ec2.DescribeSecurityGroupsInput{},
		func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			securityGroups = append(securityGroups, page.SecurityGroups...)
			return true
		})

	if err != nil {
		utils.LogAWSError("EC2.DescribeSecurityGroupsPages", err)
	}
	return
}

func buildEc2SecurityGroupSnapshot(_ ec2iface.EC2API, securityGroup *ec2.SecurityGroup) *awsmodels.Ec2SecurityGroup {
	if securityGroup == nil {
		return nil
	}
	ec2SecurityGroupSnapshot := &awsmodels.Ec2SecurityGroup{
		GenericResource: awsmodels.GenericResource{
			ResourceType: aws.String(awsmodels.Ec2SecurityGroupSchema),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			Name: securityGroup.GroupName,
			ID:   securityGroup.GroupId,
			Tags: utils.ParseTagSlice(securityGroup.Tags),
		},
		Description:         securityGroup.Description,
		IpPermissions:       securityGroup.IpPermissions,
		IpPermissionsEgress: securityGroup.IpPermissionsEgress,
		OwnerId:             securityGroup.OwnerId,
		VpcId:               securityGroup.VpcId,
	}

	return ec2SecurityGroupSnapshot
}

// PollEc2SecurityGroups gathers information on each Security Group in an AWS account.
func PollEc2SecurityGroups(pollerInput *awsmodels.ResourcePollerInput) ([]*apimodels.AddResourceEntry, error) {
	zap.L().Debug("starting EC2 Security Group resource poller")
	ec2SecurityGroupSnapshots := make(map[string]*awsmodels.Ec2SecurityGroup)

	for _, regionID := range utils.GetServiceRegions(pollerInput.Regions, "ec2") {
		sess := session.Must(session.NewSession(&aws.Config{Region: regionID}))
		creds, err := AssumeRoleFunc(pollerInput, sess)
		if err != nil {
			return nil, err
		}

		ec2Svc := EC2ClientFunc(sess, &aws.Config{Credentials: creds}).(ec2iface.EC2API)

		// Start with generating a list of all Security Groups
		securityGroups := describeSecurityGroups(ec2Svc)
		if len(securityGroups) == 0 {
			zap.L().Debug("no EC2 Security Groups found", zap.String("region", *regionID))
			continue
		}

		// For each Security Group, build out a full snapshot
		for _, securityGroup := range securityGroups {
			ec2SecurityGroupSnapshot := buildEc2SecurityGroupSnapshot(ec2Svc, securityGroup)

			// arn:aws:ec2:region:account-id:security-group/sg-id
			resourceID := strings.Join(
				[]string{
					"arn",
					pollerInput.AuthSourceParsedARN.Partition,
					"ec2",
					*regionID,
					*ec2SecurityGroupSnapshot.OwnerId,
					"security-group/" + *ec2SecurityGroupSnapshot.ID,
				},
				":",
			)

			// Populate generic fields
			ec2SecurityGroupSnapshot.ResourceID = aws.String(resourceID)

			// Populate AWS generic fields
			ec2SecurityGroupSnapshot.AccountID = aws.String(pollerInput.AuthSourceParsedARN.AccountID)
			ec2SecurityGroupSnapshot.Region = regionID
			ec2SecurityGroupSnapshot.ARN = aws.String(resourceID)

			if _, ok := ec2SecurityGroupSnapshots[resourceID]; !ok {
				ec2SecurityGroupSnapshots[resourceID] = ec2SecurityGroupSnapshot
			} else {
				zap.L().Info("overwriting existing EC2 Security Group snapshot",
					zap.String("resourceId", resourceID))
				ec2SecurityGroupSnapshots[resourceID] = ec2SecurityGroupSnapshot
			}
		}
	}

	resources := make([]*apimodels.AddResourceEntry, 0, len(ec2SecurityGroupSnapshots))
	for resourceID, ec2SecurityGroupSnapshot := range ec2SecurityGroupSnapshots {
		resources = append(resources, &apimodels.AddResourceEntry{
			Attributes:      ec2SecurityGroupSnapshot,
			ID:              apimodels.ResourceID(resourceID),
			IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
			IntegrationType: apimodels.IntegrationTypeAws,
			Type:            awsmodels.Ec2SecurityGroupSchema,
		})
	}

	return resources, nil
}
