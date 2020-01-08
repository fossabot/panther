package utils

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
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"go.uber.org/zap"
)

// GetRegions returns all the active AWS regions for a given account.
func GetRegions(ec2Svc ec2iface.EC2API) (regions []*string) {
	regionsOutput, err := ec2Svc.DescribeRegions(&ec2.DescribeRegionsInput{})

	if err != nil {
		LogAWSError("EC2.DescribeRegions", err)
		return nil
	}

	for _, region := range regionsOutput.Regions {
		regions = append(regions, region.RegionName)
	}
	return
}

// GetServiceRegions returns the intersection of the active regions passed in by the poller input
// and the regions specific to the given service
func GetServiceRegions(activeRegions []*string, serviceID string) (regions []*string) {
	serviceRegions, exists := endpoints.RegionsForService(
		endpoints.DefaultPartitions(),
		endpoints.AwsPartitionID,
		serviceID,
	)
	if !exists {
		zap.L().Error("no regions found for service", zap.String("service", serviceID))
		return nil
	}

	for _, region := range activeRegions {
		if _, ok := serviceRegions[*region]; ok {
			regions = append(regions, region)
		}
	}
	if len(regions) == 0 {
		activeRegionsS := make([]string, len(activeRegions))
		for _, region := range regions {
			activeRegionsS = append(activeRegionsS, *region)
		}
		zap.L().Debug(
			"no shared regions found between service regions and active regions",
			zap.String("service", serviceID),
			zap.Strings("activeRegions", activeRegionsS),
		)
	}
	return
}
