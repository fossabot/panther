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
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

func classifyELBV2(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_elasticloadbalancingv2.html
	if eventName == "DeleteTargetGroup" ||
		eventName == "CreateTargetGroup" ||
		eventName == "ModifyTargetGroup" ||
		eventName == "ModifyTargetGroupAttributes" ||
		eventName == "RegisterTargets" ||
		eventName == "DeregisterTargets" {

		zap.L().Debug("elbv2: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	var parseErr error
	region := detail.Get("awsRegion").Str
	lbARN := arn.ARN{
		Partition: "aws",
		Service:   "elasticloadbalancing",
		Region:    region,
		AccountID: accountID,
	}

	// We don't have a separate resource for listeners or listener rules yet, but they're built into
	// the load balancer so we need to update it. Fortunately, the load balancer ARN can be exactly
	// determined from the the ARNs of its components:
	// arn:aws:elasticloadbalancing:region:account-id:loadbalancer/[app|net]/lb-name/lb-id
	// arn:aws:elasticloadbalancing:region:account-id:listener/[app|net]/lb-name/lb-id/listener-id
	// arn:aws:elasticloadbalancing:region:account-id:listener-rule/[app|net]/lb-name/lb-id/listener-id/rule-id
	// So if we split the resource on the '/' character, we always need the elements at indices one,
	// two and three.
	switch eventName {
	case "AddListenerCertificates", "CreateRule", "DeleteListener", "ModifyListener", "RemoveListenerCertificates":
		listenerARN := detail.Get("requestParameters.listenerArn").Str
		arnComponents := strings.Split(listenerARN, "/")
		lbARN.Resource = strings.Join([]string{
			"loadbalancer",
			arnComponents[1],
			arnComponents[2],
			arnComponents[3],
		}, "/")
	case "DeleteRule", "ModifyRule":
		ruleARN := detail.Get("requestParameters.ruleArn").Str
		arnComponents := strings.Split(ruleARN, "/")
		lbARN.Resource = strings.Join([]string{
			"loadbalancer",
			arnComponents[1],
			arnComponents[2],
			arnComponents[3],
		}, "/")
	case "AddTags", "RemoveTags":
		var changes []*resourceChange
		for _, resource := range detail.Get("requestParameters.resourceArns").Array() {
			resourceARN, err := arn.Parse(resource.Str)
			if err != nil {
				zap.L().Error("elbv2: error parsing ARN", zap.String("eventName", eventName), zap.Error(err))
				return changes
			}
			if strings.HasPrefix(resourceARN.Resource, "targetgroup/") {
				continue
			}
			changes = append(changes, &resourceChange{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   resourceARN.String(),
				ResourceType: schemas.Elbv2LoadBalancerSchema,
			})
		}
		return changes
	case "CreateListener", "DeleteLoadBalancer", "ModifyLoadBalancerAttributes", "SetIpAddressType", "SetSecurityGroups", "SetSubnets":
		lbARN, parseErr = arn.Parse(detail.Get("requestParameters.loadBalancerArn").Str)
	case "CreateLoadBalancer":
		var changes []*resourceChange
		for _, lb := range detail.Get("responseElements.loadBalancers").Array() {
			lbARN, err := arn.Parse(lb.Get("loadBalancerArn").Str)
			if err != nil {
				zap.L().Error("elbv2: error parsing ARN", zap.String("eventName", eventName), zap.Error(err))
				return changes
			}
			changes = append(changes, &resourceChange{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   lbARN.String(),
				ResourceType: schemas.Elbv2LoadBalancerSchema,
			})
		}
		return changes
	case "SetRulePriorities":
		var changes []*resourceChange
		for _, rule := range detail.Get("requestParameters.rulePriorities").Array() {
			ruleARN := rule.Get("ruleArn").Str
			arnComponents := strings.Split(ruleARN, "/")
			lbARN.Resource = strings.Join([]string{
				"loadbalancer",
				arnComponents[1],
				arnComponents[2],
				arnComponents[3],
			}, "/")
			changes = append(changes, &resourceChange{
				AwsAccountID: accountID,
				Delete:       false,
				EventName:    eventName,
				ResourceID:   lbARN.String(),
				ResourceType: schemas.Elbv2LoadBalancerSchema,
			})
		}
		return changes
	default:
		zap.L().Error("elbv2: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	if parseErr != nil {
		zap.L().Warn("elbv2: error parsing ARN", zap.String("eventName", eventName), zap.Error(parseErr))
	}
	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       eventName == "DeleteLoadBalancer",
		EventName:    eventName,
		ResourceID:   lbARN.String(),
		ResourceType: schemas.Elbv2LoadBalancerSchema,
	}}
}
