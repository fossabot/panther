package handlers

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// Queue a policy for re-analysis (evaluate against all applicable resources).
//
// This ensures policy changes are reflected almost immediately (instead of waiting for daily scan).
func queuePolicy(policy *tableItem) error {
	body, err := jsoniter.MarshalToString(policy.Policy(""))
	if err != nil {
		zap.L().Error("failed to marshal policy", zap.Error(err))
		return err
	}

	zap.L().Info("queueing policy for analysis",
		zap.String("policyId", string(policy.ID)),
		zap.String("resourceQueueURL", env.ResourceQueueURL))
	_, err = sqsClient.SendMessage(
		&sqs.SendMessageInput{MessageBody: &body, QueueUrl: &env.ResourceQueueURL})
	return err
}
