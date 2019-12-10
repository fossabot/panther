package processor

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"go.uber.org/zap"
)

// SNS client factory that can be replaced for unit tests.
var snsClientBuilder = buildSnsClient

// Confirm the SNS subscription if the source account is a registered customer.
//
// Returns an error only if the confirmation failed and needs to be retried.
func handleSnsConfirmation(topicArn arn.ARN, token *string) error {
	if _, ok := accounts[topicArn.AccountID]; !ok {
		zap.L().Warn("refusing sns confirmation from unknown account",
			zap.String("accountId", topicArn.AccountID))
		return nil
	}

	if aws.StringValue(token) == "" {
		zap.L().Warn("no sns confirmation token", zap.String("topicArn", topicArn.String()))
		return nil
	}

	zap.L().Info("confirming sns subscription", zap.String("topicArn", topicArn.String()))
	snsClient, err := snsClientBuilder(&topicArn.Region)
	if err != nil {
		zap.L().Error("sns client creation failed", zap.Error(err))
		return err // retry session creation
	}

	response, err := snsClient.ConfirmSubscription(
		&sns.ConfirmSubscriptionInput{Token: token, TopicArn: aws.String(topicArn.String())})
	if err != nil {
		zap.L().Error("sns confirmation failed", zap.Error(err))
		return err // retry confirmation
	}

	zap.L().Info("sns subscription confirmed successfully",
		zap.String("subscriptionArn", aws.StringValue(response.SubscriptionArn)))
	return nil
}

func buildSnsClient(region *string) (snsiface.SNSAPI, error) {
	sess, err := session.NewSession(&aws.Config{Region: region})
	if err != nil {
		return nil, err
	}

	return sns.New(sess), nil
}
