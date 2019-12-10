package aws

import "github.com/aws/aws-sdk-go/service/elbv2"

const (
	Elbv2LoadBalancerSchema = "AWS.ELBV2.ApplicationLoadBalancer"
)

// Elbv2ApplicationLoadBalancer contains all information about an application load balancer
type Elbv2ApplicationLoadBalancer struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from elbv2.LoadBalancer
	AvailabilityZones      []*elbv2.AvailabilityZone
	CanonicalHostedZonedId *string
	DNSName                *string
	IpAddressType          *string
	Scheme                 *string
	SecurityGroups         []*string
	State                  *elbv2.LoadBalancerState
	Type                   *string
	VpcId                  *string

	// Additional fields
	WebAcl      *string
	Listeners   []*elbv2.Listener
	SSLPolicies map[string]*elbv2.SslPolicy
}
