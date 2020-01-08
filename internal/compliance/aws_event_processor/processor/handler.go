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
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/sqs"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/resources/client/operations"
	api "github.com/panther-labs/panther/api/gateway/resources/models"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	"github.com/panther-labs/panther/pkg/awsbatch/sqsbatch"
)

const maxBackoffSeconds = 30

// Handle is the entry point for the event stream analysis.
//
// Do not make any assumptions about the correctness of the incoming data.
func Handle(batch *events.SQSEvent) error {
	// De-duplicate all updates and deletes before delivering them.
	// At most one change will be reported per resource (update or delete).
	//
	// For example, if a bucket is Deleted, Created, then Modified all in this batch,
	// we will send a single update request (i.e. queue a bucket scan).
	changes := make(map[string]*resourceChange, len(batch.Records)) // keyed by resourceID

	// Get the most recent integrations to map Account ID to IntegrationID
	if err := refreshAccounts(); err != nil {
		return err
	}

	for _, record := range batch.Records {
		// Using gjson to get only the fields we need is > 10x faster than running json.Unmarshal multiple times
		switch gjson.Get(record.Body, "Type").Str {
		case "Notification": // sns wrapped message
			zap.L().Debug("wrapped sns message - assuming cloudtrail is in Message field")
			handleCloudtrail(gjson.Get(record.Body, "Message").Str, changes)

		case "SubscriptionConfirmation": // sns confirmation message
			topicArn, err := arn.Parse(gjson.Get(record.Body, "TopicArn").Str)
			if err != nil {
				zap.L().Warn("invalid confirmation arn", zap.Error(err))
				continue
			}

			token := gjson.Get(record.Body, "Token").Str
			if err = handleSnsConfirmation(topicArn, &token); err != nil {
				return err
			}

		case "": // raw data
			handleCloudtrail(record.Body, changes)
		default: // Unexpected type
			zap.L().Warn("unexpected record type",
				zap.String("type", gjson.Get(record.Body, "Type").Str),
				zap.String("body", record.Body),
			)
			continue
		}
	}

	return submitChanges(changes)
}

func handleCloudtrail(body string, changes map[string]*resourceChange) {
	// this event potentially requires a change to some number of resources
	// TODO - store the raw event somewhere
	for _, summary := range classifyCloudTrailLog(body) {
		zap.L().Info("resource change required", zap.Any("changeDetail", summary))
		// TODO - Update this to not overwrite scan requests of different types
		// More details here: https://panther-labs.atlassian.net/browse/ENG-1113
		if entry, ok := changes[summary.ResourceID]; !ok || summary.EventTime > entry.EventTime {
			changes[summary.ResourceID] = summary // the newest event for this resource we've seen so far
		}
	}
}

func submitChanges(changes map[string]*resourceChange) error {
	var deleteRequest api.DeleteResources
	requestsByDelay := make(map[int64]*poller.ScanMsg)

	for _, change := range changes {
		if change.Delete {
			deleteRequest.Resources = append(deleteRequest.Resources, &api.DeleteEntry{
				ID: api.ResourceID(change.ResourceID),
			})
		} else {
			// Possible configurations:
			// ID = “”, region =“”:				Account wide service scan
			// ID = “”, region =“west”:			Region wide service scan
			// ID = “abc-123”, region =“”:		Single resource scan
			// ID = “abc-123”, region =“west”:	Undefined in spec, treated as single resource scan downstream
			var resourceID *string
			var region *string
			if change.ResourceID != "" {
				resourceID = &change.ResourceID
			}
			if change.Region != "" {
				region = &change.Region
			}

			if _, ok := requestsByDelay[change.Delay]; !ok {
				requestsByDelay[change.Delay] = &poller.ScanMsg{}
			}

			// Group all changes together by their delay time. This will maintain our ability to
			// group together changes that happened close together in time. I imagine in cases where
			// we set a delay it will be a fairly uniform delay.
			requestsByDelay[change.Delay].Entries = append(requestsByDelay[change.Delay].Entries, &poller.ScanEntry{
				AWSAccountID:     &change.AwsAccountID,
				IntegrationID:    &change.IntegrationID,
				Region:           region,
				ResourceID:       resourceID,
				ResourceType:     &change.ResourceType,
				ScanAllResources: aws.Bool(false),
			})
		}
	}

	// Send deletes to resources-api
	if len(deleteRequest.Resources) > 0 {
		zap.L().Info("deleting resources", zap.Any("deleteRequest", &deleteRequest))
		_, err := apiClient.Operations.DeleteResources(
			&operations.DeleteResourcesParams{Body: &deleteRequest, HTTPClient: httpClient})

		if err != nil {
			zap.L().Error("resource deletion failed", zap.Error(err))
			return err
		}
	}

	if len(requestsByDelay) > 0 {
		batchInput := &sqs.SendMessageBatchInput{QueueUrl: &queueURL}
		// Send resource scan requests to the poller queue
		for delay, request := range requestsByDelay {
			zap.L().Info("queueing resource scans", zap.Any("updateRequest", request))
			body, err := jsoniter.MarshalToString(request)
			if err != nil {
				zap.L().Error("resource queueing failed: json marshal", zap.Error(err))
				return err
			}

			batchInput.Entries = append(batchInput.Entries, &sqs.SendMessageBatchRequestEntry{
				Id:           aws.String(strconv.FormatInt(delay, 10)),
				MessageBody:  aws.String(body),
				DelaySeconds: aws.Int64(delay),
			})
		}

		if err := sqsbatch.SendMessageBatch(sqsClient, maxBackoffSeconds, batchInput); err != nil {
			return err
		}
	}

	return nil
}
