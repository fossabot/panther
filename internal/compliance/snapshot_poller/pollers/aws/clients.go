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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

const MaxRetries = 6

// Key used for the client cache to neatly encapsulate an integration, service, and region
type clientKey struct {
	IntegrationID string
	Service       string
	Region        string
}

type cachedClient struct {
	Client      interface{}
	Credentials *credentials.Credentials
}

var clientCache = make(map[clientKey]cachedClient)

// Functions used to build clients, keyed by service
var clientFuncs = map[string]func(session2 *session.Session, config *aws.Config) interface{}{
	"acm":                    AcmClientFunc,
	"applicationautoscaling": ApplicationAutoScalingClientFunc,
	"cloudformation":         CloudFormationClientFunc,
	"cloudtrail":             CloudTrailClientFunc,
	"cloudwatchlogs":         CloudWatchLogsClientFunc,
	"configservice":          ConfigServiceClientFunc,
	"dynamodb":               DynamoDBClientFunc,
	"ec2":                    EC2ClientFunc,
	"elbv2":                  Elbv2ClientFunc,
	"guardduty":              GuardDutyClientFunc,
	"iam":                    IAMClientFunc,
	"kms":                    KmsClientFunc,
	"lambda":                 LambdaClientFunc,
	"rds":                    RDSClientFunc,
	"redshift":               RedshiftClientFunc,
	"s3":                     S3ClientFunc,
	"waf":                    WafClientFunc,
	"waf-regional":           WafRegionalClientFunc,
}

// getClient returns a valid client for a given integration, service, and region using caching.
func getClient(pollerInput *awsmodels.ResourcePollerInput, service string, region string) interface{} {
	cacheKey := clientKey{
		IntegrationID: *pollerInput.IntegrationID,
		Service:       service,
		Region:        region,
	}

	// Return the cached client if the credentials used to build it are not expired
	if cachedClient, exists := clientCache[cacheKey]; exists {
		if !cachedClient.Credentials.IsExpired() {
			if cachedClient.Client != nil {
				return cachedClient.Client
			}
			zap.L().Warn("nil client was cached", zap.Any("cache key", cacheKey))
		}
	}

	// Build a new client on cache miss OR if the client in the cache has expired credentials

	// Build the new session and credentials
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	creds, err := AssumeRoleFunc(pollerInput, sess)
	if err != nil {
		zap.L().Error("unable to assume role to build client cache", zap.Error(err))
		return nil
	}

	// Build the new client and cache it with the credentials used to build it
	if clientFunc, ok := clientFuncs[service]; ok {
		client := clientFunc(sess, &aws.Config{Credentials: creds})
		clientCache[cacheKey] = cachedClient{
			Client:      client,
			Credentials: creds,
		}
		return client
	}

	zap.L().Error("cannot build client for unsupported service",
		zap.String("service", service),
	)
	return nil
}
