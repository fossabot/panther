package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/aws/aws-sdk-go/service/wafregional/wafregionaliface"
	"go.uber.org/zap"

	apimodels "github.com/panther-labs/panther/api/resources/models"
	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/utils"
)

// Set as variables to be overridden in testing
var (
	Elbv2ClientFunc = setupElbv2Client
	sslPolicies     = map[string]*elbv2.SslPolicy{}
)

func setupElbv2Client(sess *session.Session, cfg *aws.Config) interface{} {
	cfg.MaxRetries = aws.Int(MaxRetries)
	return elbv2.New(sess, cfg)
}

// PollELBV2 LoadBalancer polls a single ELBV2 Application Load Balancer resource
func PollELBV2LoadBalancer(
	pollerResourceInput *awsmodels.ResourcePollerInput,
	resourceARN arn.ARN,
	scanRequest *pollermodels.ScanEntry,
) interface{} {

	elbv2Client := getClient(pollerResourceInput, "elbv2", resourceARN.Region).(elbv2iface.ELBV2API)
	wafClient := getClient(pollerResourceInput, "waf-regional", resourceARN.Region).(wafregionaliface.WAFRegionalAPI)
	loadBalancer := getApplicationLoadBalancer(elbv2Client, scanRequest.ResourceID)

	snapshot := buildElbv2ApplicationLoadBalancerSnapshot(elbv2Client, wafClient, loadBalancer)
	if snapshot == nil {
		return nil
	}

	snapshot.AccountID = aws.String(resourceARN.AccountID)
	snapshot.Region = aws.String(resourceARN.Region)
	return snapshot
}

// getApplicationLoadBalancer returns a specifc ELBV2 application load balancer
func getApplicationLoadBalancer(svc elbv2iface.ELBV2API, loadBalancerARN *string) *elbv2.LoadBalancer {
	loadBalancer, err := svc.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []*string{loadBalancerARN},
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "LoadBalancerNotFound" {
				zap.L().Warn("tried to scan non-existent resource",
					zap.String("resource", *loadBalancerARN),
					zap.String("resourceType", awsmodels.Elbv2LoadBalancerSchema))
				return nil
			}
		}
		utils.LogAWSError("ELBV2.DescribeLoadBalancers", err)
		return nil
	}
	return loadBalancer.LoadBalancers[0]
}

// describeLoadBalancers returns a list of all Load Balancers in the account in the current region
func describeLoadBalancers(elbv2Svc elbv2iface.ELBV2API) (loadBalancers []*elbv2.LoadBalancer) {
	err := elbv2Svc.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{},
		func(page *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			loadBalancers = append(loadBalancers, page.LoadBalancers...)
			return true
		})
	if err != nil {
		utils.LogAWSError("ELBV2.DescribeLoadBalancersPages", err)
	}
	return
}

// describeListeners returns all the listeners for a given ELBV2 load balancer
func describeListeners(elbv2Svc elbv2iface.ELBV2API, arn *string) (listeners []*elbv2.Listener) {
	err := elbv2Svc.DescribeListenersPages(&elbv2.DescribeListenersInput{LoadBalancerArn: arn},
		func(page *elbv2.DescribeListenersOutput, lastPage bool) bool {
			listeners = append(listeners, page.Listeners...)
			return true
		})
	if err != nil {
		utils.LogAWSError("ELBV2.DescribeListenersPages", err)
	}
	return
}

// describeTags returns all the tags associated to the given load balancer
func describeTags(svc elbv2iface.ELBV2API, arn *string) ([]*elbv2.Tag, error) {
	tags, err := svc.DescribeTags(&elbv2.DescribeTagsInput{ResourceArns: []*string{arn}})
	if err != nil {
		utils.LogAWSError("ELBV2.DescribeTags", err)
		return nil, err
	}

	return tags.TagDescriptions[0].Tags, nil
}

// describeSSLPolicies returns all the SSL policies in the current region
func describeSSLPolicies(svc elbv2iface.ELBV2API) ([]*elbv2.SslPolicy, error) {
	sslPoliciesDescription, err := svc.DescribeSSLPolicies(&elbv2.DescribeSSLPoliciesInput{})
	if err != nil {
		utils.LogAWSError("ELBV2.DescribeSSLPolicies", err)
		return nil, err
	}
	return sslPoliciesDescription.SslPolicies, nil
}

// getWebACLForResource returns the web ACL ID for the given application load balancer
func getWebACLForResource(wafRegionalSvc wafregionaliface.WAFRegionalAPI, arn *string) (*string, error) {
	out, err := wafRegionalSvc.GetWebACLForResource(
		&wafregional.GetWebACLForResourceInput{ResourceArn: arn},
	)
	if err != nil {
		return nil, err
	}

	if out.WebACLSummary == nil {
		return nil, nil
	}

	return out.WebACLSummary.WebACLId, nil
}

// generateSSLPolices sets up the sslPolicies map for reference
func generateSSLPolicies(svc elbv2iface.ELBV2API) {
	policies, err := describeSSLPolicies(svc)
	if err == nil {
		sslPolicies = make(map[string]*elbv2.SslPolicy, len(policies))
		for _, policy := range policies {
			sslPolicies[*policy.Name] = policy
		}
	}
}

// buildElbv2ApplicationLoadBalancerSnapshot makes all the calls to build up a snapshot of a given
// application load balancer
func buildElbv2ApplicationLoadBalancerSnapshot(
	elbv2Svc elbv2iface.ELBV2API,
	wafRegionalSvc wafregionaliface.WAFRegionalAPI,
	lb *elbv2.LoadBalancer,
) *awsmodels.Elbv2ApplicationLoadBalancer {

	if lb == nil {
		return nil
	}

	applicationLoadBalancer := &awsmodels.Elbv2ApplicationLoadBalancer{
		GenericResource: awsmodels.GenericResource{
			ResourceID:   lb.LoadBalancerArn,
			TimeCreated:  utils.DateTimeFormat(*lb.CreatedTime),
			ResourceType: aws.String(awsmodels.Elbv2LoadBalancerSchema),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			ARN:  lb.LoadBalancerArn,
			Name: lb.LoadBalancerName,
		},
		AvailabilityZones:      lb.AvailabilityZones,
		CanonicalHostedZonedId: lb.CanonicalHostedZoneId,
		DNSName:                lb.DNSName,
		IpAddressType:          lb.IpAddressType,
		Scheme:                 lb.Scheme,
		SecurityGroups:         lb.SecurityGroups,
		State:                  lb.State,
		Type:                   lb.Type,
		VpcId:                  lb.VpcId,
	}

	tags, err := describeTags(elbv2Svc, lb.LoadBalancerArn)
	if err == nil {
		applicationLoadBalancer.Tags = utils.ParseTagSlice(tags)
	}

	// Build the list of listeners and associated SSL Policies for the load balancer
	listeners := describeListeners(elbv2Svc, lb.LoadBalancerArn)
	if len(listeners) != 0 {
		applicationLoadBalancer.Listeners = listeners
		applicationLoadBalancer.SSLPolicies = make(map[string]*elbv2.SslPolicy)
		for _, listener := range listeners {
			if listener.SslPolicy == nil {
				continue
			}
			if sslPolicies == nil {
				generateSSLPolicies(elbv2Svc)
			}
			if policy, ok := sslPolicies[*listener.SslPolicy]; ok {
				applicationLoadBalancer.SSLPolicies[*listener.SslPolicy] = policy
			}
		}
	}

	// Try to find a webACL ID
	webACL, err := getWebACLForResource(wafRegionalSvc, lb.LoadBalancerArn)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "AccessDeniedException" {
				zap.L().Debug(
					"AccessDeniedException, additional privileges were not granted.",
					zap.String("API", "WAF.Regional.GetWebAclForResource"))
			}
		} else {
			utils.LogAWSError("WAF.Regional.GetWebAclForResource", err)
		}
	} else {
		applicationLoadBalancer.WebAcl = webACL
	}

	return applicationLoadBalancer
}

// PollElbv2ApplicationLoadBalancers gathers information on each application load balancer for an AWS account.
func PollElbv2ApplicationLoadBalancers(pollerInput *awsmodels.ResourcePollerInput) ([]*apimodels.AddResourceEntry, error) {
	zap.L().Debug("starting ELBV2 Application Load Balancer resource poller")
	elbv2LoadBalancerSnapshots := make(map[string]*awsmodels.Elbv2ApplicationLoadBalancer)

	for _, regionID := range utils.GetServiceRegions(pollerInput.Regions, "elasticloadbalancing") {
		sess := session.Must(session.NewSession(&aws.Config{Region: regionID}))
		creds, err := AssumeRoleFunc(pollerInput, sess)
		if err != nil {
			return nil, err
		}

		config := &aws.Config{Credentials: creds}
		elbv2Svc := Elbv2ClientFunc(sess, config).(elbv2iface.ELBV2API)
		wafRegionalSvc := WafRegionalClientFunc(sess, config).(wafregionaliface.WAFRegionalAPI)

		// Start with generating a list of all load balancers
		loadBalancers := describeLoadBalancers(elbv2Svc)
		if len(loadBalancers) == 0 {
			zap.L().Debug(
				"No application load balancers found.",
				zap.String("region", *regionID),
			)
			continue
		}

		// Next generate a list of SSL policies to be shared by the load balancer snapshots
		generateSSLPolicies(elbv2Svc)

		for _, loadBalancer := range loadBalancers {
			elbv2LoadBalancer := buildElbv2ApplicationLoadBalancerSnapshot(
				elbv2Svc,
				wafRegionalSvc,
				loadBalancer,
			)
			if elbv2LoadBalancer == nil {
				continue
			}
			elbv2LoadBalancer.AccountID = aws.String(pollerInput.AuthSourceParsedARN.AccountID)
			elbv2LoadBalancer.Region = regionID

			if _, ok := elbv2LoadBalancerSnapshots[*elbv2LoadBalancer.ARN]; !ok {
				elbv2LoadBalancerSnapshots[*elbv2LoadBalancer.ARN] = elbv2LoadBalancer
			} else {
				zap.L().Info(
					"overwriting existing ELB v2 Load Balancer snapshot",
					zap.String("resourceId", *elbv2LoadBalancer.ARN),
				)
				elbv2LoadBalancerSnapshots[*elbv2LoadBalancer.ARN] = elbv2LoadBalancer
			}
		}
	}

	resources := make([]*apimodels.AddResourceEntry, 0, len(elbv2LoadBalancerSnapshots))
	for resourceID, elbv2LoadBalancerSnapshot := range elbv2LoadBalancerSnapshots {
		resources = append(resources, &apimodels.AddResourceEntry{
			Attributes:      elbv2LoadBalancerSnapshot,
			ID:              apimodels.ResourceID(resourceID),
			IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
			IntegrationType: apimodels.IntegrationTypeAws,
			Type:            awsmodels.Elbv2LoadBalancerSchema,
		})
	}

	return resources, nil
}
