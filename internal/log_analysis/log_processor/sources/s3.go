package sources

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
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
)

// ReadData reads incoming messages and returns a slice of DataStream items
func ReadData(messages []*string) (result []*common.DataStream, err error) {
	zap.L().Info("Reading data for messages", zap.Int("numMessages", len(messages)))
	for _, message := range messages {
		snsNotificationMessage := &SnsNotification{}
		if err := jsoniter.UnmarshalFromString(*message, snsNotificationMessage); err != nil {
			return nil, err
		}

		switch snsNotificationMessage.Type {
		case "Notification":
			streams, err := handleNotificationMessage(snsNotificationMessage)
			if err != nil {
				return nil, err
			}
			result = append(result, streams...)
		case "SubscriptionConfirmation":
			err := ConfirmSubscription(snsNotificationMessage)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("received unexpected message in SQS queue")
		}
	}
	return
}

// ConfirmSubscription will confirm the SNS->SQS subscription
func ConfirmSubscription(notification *SnsNotification) error {
	zap.L().Info("confirming sns subscription",
		zap.String("topicArn", notification.TopicArn))

	topicArn, err := arn.Parse(notification.TopicArn)
	if err != nil {
		return err
	}
	snsClient := sns.New(common.Session, aws.NewConfig().WithRegion(topicArn.Region))
	subscriptionConfiguration := &sns.ConfirmSubscriptionInput{
		Token:    notification.Token,
		TopicArn: aws.String(notification.TopicArn),
	}
	_, err = snsClient.ConfirmSubscription(subscriptionConfiguration)
	if err != nil {
		zap.L().Warn("failed to confirm subscription", zap.Error(err))
		return err
	}
	zap.L().Info("successfully confirmed subscription",
		zap.String("topicArn", notification.TopicArn))
	return nil
}

func handleNotificationMessage(notification *SnsNotification) (result []*common.DataStream, err error) {
	s3Objects, err := ParseNotification(notification.Message)
	if err != nil {
		return
	}

	for _, s3Object := range s3Objects {
		zap.L().Debug("going to load new object from S3",
			zap.String("bucket", *s3Object.S3Bucket),
			zap.String("key", *s3Object.S3ObjectKey),
		)

		s3Client, err := getS3Client(*s3Object.S3Bucket, notification.TopicArn)
		if err != nil {
			return nil, err
		}

		getObjectInput := &s3.GetObjectInput{
			Bucket: s3Object.S3Bucket,
			Key:    s3Object.S3ObjectKey,
		}
		output, err := s3Client.GetObject(getObjectInput)
		if err != nil {
			return nil, err
		}

		bufferedReader := bufio.NewReader(output.Body)

		// We peek into the file header to identify the content type
		// http.DetectContentType only uses up to the first 512 bytes
		headerBytes, readerErr := bufferedReader.Peek(512)
		if readerErr != nil {
			return nil, err
		}
		contentType := http.DetectContentType(headerBytes)

		var streamReader io.Reader

		// Checking for prefix because the returned type can have also charset used
		if strings.HasPrefix(contentType, "text/plain") {
			// if it's plain text, just return the buffered reader
			streamReader = bufferedReader
		} else if strings.HasPrefix(contentType, "application/x-gzip") {
			gzipReader, err := gzip.NewReader(bufferedReader)
			if err != nil {
				zap.L().Warn("failed to created gzip reader", zap.Error(err))
				return nil, err
			}
			streamReader = gzipReader
		}

		dataStream := &common.DataStream{
			Reader:  &streamReader,
			LogType: s3Object.LogType,
		}
		result = append(result, dataStream)
	}
	return result, err
}

// ParseNotification parses a message received
func ParseNotification(message string) ([]*S3ObjectInfo, error) {
	s3Objects, err := parseCloudTrailNotification(message)
	if err != nil {
		zap.L().Error("encountered issue when parsing CloudTrail notification", zap.Error(err))
		return nil, err
	}

	// If the input was not a CloudTrail notification, the result will be empty slice
	if len(s3Objects) == 0 {
		s3Objects, err = parseS3Event(message)
		if err != nil {
			zap.L().Error("encountered issue when parsing S3 notification", zap.Error(err))
			return nil, err
		}
	}

	if len(s3Objects) == 0 {
		zap.L().Error("notification is not of known type")
		return nil, errors.New("notification is not of known type")
	}
	return s3Objects, nil
}

// parseCloudTrailNotification will try to parse input as if it was a CloudTrail notification
// If the input was not a CloudTrail notification, it will return a empty slice
// The method returns error if it encountered some issue while trying to parse the notification
func parseCloudTrailNotification(message string) (result []*S3ObjectInfo, err error) {
	cloudTrailNotification := &cloudTrailNotification{}
	err = jsoniter.UnmarshalFromString(message, cloudTrailNotification)
	if err != nil {
		return nil, err
	}

	for _, s3Key := range cloudTrailNotification.S3ObjectKey {
		info := &S3ObjectInfo{
			S3Bucket:    cloudTrailNotification.S3Bucket,
			S3ObjectKey: s3Key,
		}
		result = append(result, info)
	}
	return result, nil
}

// parseS3Event will try to parse input as if it was an S3 Event (https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html)
// If the input was not an S3 Event  notification, it will return a empty slice
// The method returns error if it encountered some issue while trying to parse the notification
func parseS3Event(message string) (result []*S3ObjectInfo, err error) {
	notification := &events.S3Event{}
	err = jsoniter.UnmarshalFromString(message, notification)
	if err != nil {
		return nil, err
	}

	for _, record := range notification.Records {
		info := &S3ObjectInfo{
			S3Bucket:    aws.String(record.S3.Bucket.Name),
			S3ObjectKey: aws.String(record.S3.Object.Key),
		}
		result = append(result, info)
	}
	return result, nil
}

// cloudTrailNotification is the notification sent by CloudTrail whenever it delivers a new log file to S3
type cloudTrailNotification struct {
	S3Bucket    *string   `json:"s3Bucket"`
	S3ObjectKey []*string `json:"s3ObjectKey"`
}

//S3ObjectInfo contains information about the S3 object
type S3ObjectInfo struct {
	S3Bucket    *string
	S3ObjectKey *string
	LogType     *string
}

// SnsNotification struct represents an SNS message arriving to Panther SQS from a customer account.
// The message can either be of type 'Notification' or 'SubscriptionConfirmation'
// Since there is no AWS SDK-provided struct to represent both types
// we had to create this custom type to include fields from both types.
type SnsNotification struct {
	events.SNSEntity
	Token *string `json:"Token"`
}
