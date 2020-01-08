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
