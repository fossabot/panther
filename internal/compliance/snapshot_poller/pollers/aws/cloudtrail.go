package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	"go.uber.org/zap"

	apimodels "github.com/panther-labs/panther/api/gateway/resources/models"
	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/utils"
)

var (
	// CloudTrailClientFunc is the function it setup the CloudTrail client.
	CloudTrailClientFunc = setupCloudTrailClient
)

func setupCloudTrailClient(sess *session.Session, cfg *aws.Config) interface{} {
	cfg.MaxRetries = aws.Int(MaxRetries)
	return cloudtrail.New(sess, cfg)
}

// PollCloudTrailTrail polls a single CloudTrail trail resource
func PollCloudTrailTrail(
	pollerResourceInput *awsmodels.ResourcePollerInput,
	resourceARN arn.ARN,
	scanRequest *pollermodels.ScanEntry,
) interface{} {

	client := getClient(pollerResourceInput, "cloudtrail", resourceARN.Region).(cloudtrailiface.CloudTrailAPI)
	trail := getTrail(client, scanRequest.ResourceID)

	snapshot := buildCloudTrailSnapshot(client, trail)
	if snapshot == nil {
		return nil
	}
	snapshot.Region = aws.String(resourceARN.Region)
	snapshot.AccountID = aws.String(resourceARN.AccountID)
	return snapshot
}

// getTrail returns the specified cloudtrail
func getTrail(svc cloudtrailiface.CloudTrailAPI, trailARN *string) *cloudtrail.Trail {
	trail, err := svc.DescribeTrails(&cloudtrail.DescribeTrailsInput{
		TrailNameList: []*string{trailARN},
	})
	if err != nil {
		utils.LogAWSError("CloudTrail.DescribeTrails", err)
		return nil
	}

	if len(trail.TrailList) == 0 {
		zap.L().Warn("tried to scan non-existent resource",
			zap.String("resource", *trailARN),
			zap.String("resourceType", awsmodels.CloudTrailSchema))
		return nil
	}

	return trail.TrailList[0]
}

func describeTrails(svc cloudtrailiface.CloudTrailAPI) ([]*cloudtrail.Trail, error) {
	var in = &cloudtrail.DescribeTrailsInput{IncludeShadowTrails: aws.Bool(false)}
	var out *cloudtrail.DescribeTrailsOutput
	var err error

	if out, err = svc.DescribeTrails(in); err != nil {
		return nil, err
	}

	return out.TrailList, nil
}

func getTrailStatus(
	svc cloudtrailiface.CloudTrailAPI,
	trailARN *string,
) (*cloudtrail.GetTrailStatusOutput, error) {

	var in = &cloudtrail.GetTrailStatusInput{Name: trailARN}
	var out *cloudtrail.GetTrailStatusOutput
	var err error

	if out, err = svc.GetTrailStatus(in); err != nil {
		return nil, err
	}

	return out, nil
}

func listTagsCloudTrail(svc cloudtrailiface.CloudTrailAPI, trailArn *string) ([]*cloudtrail.Tag, error) {
	out, err := svc.ListTags(&cloudtrail.ListTagsInput{ResourceIdList: []*string{trailArn}})
	if err != nil {
		utils.LogAWSError("CloudTrail.ListTags", err)
		return nil, err
	}

	// Since we are only specifying one resource, this will always return one value.
	// Could optimize here by calling list-tags for all resources in the region, then looking them up
	// on a per resource limit.
	return out.ResourceTagList[0].TagsList, nil
}

func getEventSelectors(svc cloudtrailiface.CloudTrailAPI, trailARN *string) ([]*cloudtrail.EventSelector, error) {
	out, err := svc.GetEventSelectors(&cloudtrail.GetEventSelectorsInput{TrailName: trailARN})
	if err != nil {
		return nil, err
	}
	return out.EventSelectors, nil
}

// buildCloudTrailSnapshot builds a complete CloudTrail snapshot for a given trail
func buildCloudTrailSnapshot(svc cloudtrailiface.CloudTrailAPI, trail *cloudtrail.Trail) *awsmodels.CloudTrail {
	if trail == nil {
		return nil
	}
	cloudTrail := &awsmodels.CloudTrail{
		GenericResource: awsmodels.GenericResource{
			ResourceID:   trail.TrailARN,
			ResourceType: aws.String(awsmodels.CloudTrailSchema),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			ARN:  trail.TrailARN,
			Name: trail.Name,
		},
		CloudWatchLogsLogGroupArn:  trail.CloudWatchLogsLogGroupArn,
		CloudWatchLogsRoleArn:      trail.CloudWatchLogsRoleArn,
		HasCustomEventSelectors:    trail.HasCustomEventSelectors,
		HomeRegion:                 trail.HomeRegion,
		IncludeGlobalServiceEvents: trail.IncludeGlobalServiceEvents,
		IsMultiRegionTrail:         trail.IsMultiRegionTrail,
		IsOrganizationTrail:        trail.IsOrganizationTrail,
		KmsKeyId:                   trail.KmsKeyId,
		LogFileValidationEnabled:   trail.LogFileValidationEnabled,
		S3BucketName:               trail.S3BucketName,
		S3KeyPrefix:                trail.S3KeyPrefix,
		SnsTopicARN:                trail.SnsTopicARN,
		SnsTopicName:               trail.SnsTopicName, //nolint:staticcheck
	}

	status, err := getTrailStatus(svc, trail.TrailARN)
	if err != nil {
		utils.LogAWSError("CloudTrail.GetTrailStatus", err)
	} else {
		cloudTrail.Status = status
	}

	eventSelectors, err := getEventSelectors(svc, trail.TrailARN)
	if err != nil {
		utils.LogAWSError("CloudTrail.GetEventSelectors", err)
	} else {
		cloudTrail.EventSelectors = eventSelectors
	}

	tags, err := listTagsCloudTrail(svc, trail.TrailARN)
	if err == nil {
		cloudTrail.Tags = utils.ParseTagSlice(tags)
	}

	return cloudTrail
}

// buildCloudTrails combines the output of each required API call to build the CloudTrailSnapshot.
//
// It returns a mapping of CloudTrailARN to CloudTrailSnapshot.
func buildCloudTrails(
	cloudtrailSvc cloudtrailiface.CloudTrailAPI,
) awsmodels.CloudTrails {

	cloudTrails := make(awsmodels.CloudTrails)

	zap.L().Debug("describing CloudTrails")
	trails, err := describeTrails(cloudtrailSvc)
	if err != nil {
		utils.LogAWSError("CloudTrail.Describe", err)
		// Return early since there are no CloudTrails
		return cloudTrails
	}

	// Build each CloudTrail's snapshot by requesting additional context from CloudTrail/S3 APIs
	for _, trail := range trails {
		cloudTrail := buildCloudTrailSnapshot(cloudtrailSvc, trail)
		cloudTrails[*trail.TrailARN] = cloudTrail
	}

	return cloudTrails
}

// PollCloudTrails gathers information on all CloudTrails in an AWS account.
func PollCloudTrails(pollerInput *awsmodels.ResourcePollerInput) ([]*apimodels.AddResourceEntry, error) {
	zap.L().Debug("starting CloudTrail resource poller")
	cloudTrailSnapshots := make(awsmodels.CloudTrails)

	for _, regionID := range utils.GetServiceRegions(pollerInput.Regions, "cloudtrail") {
		zap.L().Debug("building CloudTrail snapshots", zap.String("region", *regionID))
		sess := session.Must(session.NewSession(&aws.Config{Region: regionID}))
		creds, err := AssumeRoleFunc(pollerInput, sess)
		if err != nil {
			return nil, err
		}

		cloudTrailSvc := CloudTrailClientFunc(sess, &aws.Config{Credentials: creds}).(cloudtrailiface.CloudTrailAPI)

		// Build the list of all CloudTrails for the given region
		regionTrails := buildCloudTrails(cloudTrailSvc)

		// Insert each trail into the master list of CloudTrails (if it is not there already)
		for trailARN, trail := range regionTrails {
			trail.Region = regionID
			trail.AccountID = aws.String(pollerInput.AuthSourceParsedARN.AccountID)
			if _, ok := cloudTrailSnapshots[trailARN]; !ok {
				cloudTrailSnapshots[trailARN] = trail
			} else {
				zap.L().Info(
					"overwriting existing CloudTrail snapshot",
					zap.String("resourceId", trailARN),
				)
				cloudTrailSnapshots[trailARN] = trail
			}
		}
	}

	zap.L().Debug("finished polling CloudTrail", zap.Int("count", len(cloudTrailSnapshots)))

	metaResourceID := utils.GenerateResourceID(
		pollerInput.AuthSourceParsedARN.AccountID,
		"",
		awsmodels.CloudTrailMetaSchema,
	)

	// Handle the case where there is not a CloudTrail to return
	if len(cloudTrailSnapshots) == 0 {
		return []*apimodels.AddResourceEntry{{
			Attributes: &awsmodels.CloudTrailMeta{
				GenericResource: awsmodels.GenericResource{
					ResourceID:   aws.String(metaResourceID),
					ResourceType: aws.String(awsmodels.CloudTrailMetaSchema),
				},
				GenericAWSResource: awsmodels.GenericAWSResource{
					AccountID: aws.String(pollerInput.AuthSourceParsedARN.AccountID),
					Name:      aws.String(awsmodels.CloudTrailMetaSchema),
					Region:    aws.String("global"),
				},
				Trails: []*string{},
			},
			ID:              apimodels.ResourceID(metaResourceID),
			IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
			IntegrationType: apimodels.IntegrationTypeAws,
			Type:            awsmodels.CloudTrailMetaSchema,
		}}, nil
	}

	accountSnapshot := &awsmodels.CloudTrailMeta{
		GenericResource: awsmodels.GenericResource{
			ResourceID:   aws.String(metaResourceID),
			ResourceType: aws.String(awsmodels.CloudTrailMetaSchema),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			AccountID: aws.String(pollerInput.AuthSourceParsedARN.AccountID),
			Name:      aws.String(awsmodels.CloudTrailMetaSchema),
			Region:    aws.String("global"),
		},
	}

	resources := make([]*apimodels.AddResourceEntry, 0, len(cloudTrailSnapshots)+1)
	for _, trail := range cloudTrailSnapshots {
		resources = append(resources, &apimodels.AddResourceEntry{
			Attributes:      trail,
			ID:              apimodels.ResourceID(*trail.ARN),
			IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
			IntegrationType: apimodels.IntegrationTypeAws,
			Type:            awsmodels.CloudTrailSchema,
		})
		accountSnapshot.Trails = append(accountSnapshot.Trails, trail.ResourceID)
		if *trail.IsMultiRegionTrail && *trail.Status.IsLogging {
			accountSnapshot.GlobalEventSelectors = append(
				accountSnapshot.GlobalEventSelectors,
				trail.EventSelectors...,
			)
		}
	}

	resources = append(resources, &apimodels.AddResourceEntry{
		Attributes:      accountSnapshot,
		ID:              apimodels.ResourceID(metaResourceID),
		IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
		IntegrationType: apimodels.IntegrationTypeAws,
		Type:            awsmodels.CloudTrailMetaSchema,
	})

	return resources, nil
}
