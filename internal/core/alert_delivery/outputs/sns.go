package outputs

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
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

// Sns sends an alert to an SNS Topic.
// nolint: dupl
func (client *OutputClient) Sns(alert *alertmodels.Alert, config *outputmodels.SnsConfig) *AlertDeliveryError {
	outputMessage := &snsOutputMessage{
		ID:          alert.PolicyID,
		Name:        alert.PolicyName,
		VersionID:   alert.PolicyVersionID,
		Description: alert.PolicyDescription,
		Runbook:     alert.Runbook,
		Severity:    alert.Severity,
		Tags:        alert.Tags,
	}

	serializedMessage, err := jsoniter.MarshalToString(outputMessage)
	if err != nil {
		zap.L().Error("Failed to serialize message", zap.Error(err))
		return &AlertDeliveryError{Message: "Failed to serialize message"}
	}

	snsMessageInput := &sns.PublishInput{
		TopicArn: config.TopicArn,
		Message:  aws.String(serializedMessage),
	}

	snsClient, err := client.getSnsClient(*config.TopicArn)
	if err != nil {
		return &AlertDeliveryError{Message: "Failed to create SNS client for topic", Permanent: true}
	}

	_, err = snsClient.Publish(snsMessageInput)
	if err != nil {
		zap.L().Error("Failed to send message to SNS topic", zap.Error(err))
		return &AlertDeliveryError{Message: "Failed to send message to SNS topic"}
	}
	return nil
}

//snsOutputMessage contains the fields that will be included in the SNS message
type snsOutputMessage struct {
	ID          *string   `json:"id"`
	Name        *string   `json:"name,omitempty"`
	VersionID   *string   `json:"versionId,omitempty"`
	Description *string   `json:"description,omitempty"`
	Runbook     *string   `json:"runbook,omitempty"`
	Severity    *string   `json:"severity"`
	Tags        []*string `json:"tags,omitempty"`
}

func (client *OutputClient) getSnsClient(topicArn string) (snsiface.SNSAPI, error) {
	parsedArn, err := arn.Parse(topicArn)
	if err != nil {
		zap.L().Error("failed to parse topic ARN", zap.Error(err))
		return nil, err
	}
	snsClient, ok := client.snsClients[parsedArn.Region]
	if !ok {
		snsClient = sns.New(client.session, aws.NewConfig().WithRegion(parsedArn.Region))
		client.snsClients[parsedArn.Region] = snsClient
	}
	return snsClient, nil
}
