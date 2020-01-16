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
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
)

// ReadSQSMessages reads incoming messages containing SNS notifications and returns a slice of DataStream items
func ReadSQSMessages(messages []events.SQSMessage) (result []*common.DataStream, err error) {
	zap.L().Debug("reading data for messages", zap.Int("numMessages", len(messages)))
	for _, message := range messages {
		snsNotificationMessage := &SnsNotification{}
		if err := jsoniter.UnmarshalFromString(message.Body, snsNotificationMessage); err != nil {
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
	return result, nil
}

// ConfirmSubscription will confirm the SNS->SQS subscription
func ConfirmSubscription(notification *SnsNotification) (err error) {
	operation := common.OpLogManager.Start("ConfirmSubscription", common.OpLogSNSServiceDim)
	defer func() {
		operation.Stop()
		// sns dim info
		operation.Log(err, zap.String("topicArn", notification.TopicArn))
	}()

	topicArn, err := arn.Parse(notification.TopicArn)
	if err != nil {
		return errors.Wrap(err, "failed to parse topic arn: "+notification.TopicArn)
	}
	snsClient := sns.New(common.Session, aws.NewConfig().WithRegion(topicArn.Region))
	subscriptionConfiguration := &sns.ConfirmSubscriptionInput{
		Token:    notification.Token,
		TopicArn: aws.String(notification.TopicArn),
	}
	_, err = snsClient.ConfirmSubscription(subscriptionConfiguration)
	if err != nil {
		err = errors.Wrap(err, "failed to confirm subscription for: "+notification.TopicArn)
		return err
	}
	return nil
}

func handleNotificationMessage(notification *SnsNotification) (result []*common.DataStream, err error) {
	s3Objects, err := ParseNotification(notification.Message)
	if err != nil {
		return nil, err
	}
	for _, s3Object := range s3Objects {
		var dataStream *common.DataStream
		dataStream, err = readS3Object(s3Object, notification.TopicArn)
		if err != nil {
			return
		}
		result = append(result, dataStream)
	}
	return result, err
}

func readS3Object(s3Object *S3ObjectInfo, topicArn string) (dataStream *common.DataStream, err error) {
	operation := common.OpLogManager.Start("readS3Object", common.OpLogS3ServiceDim)
	defer func() {
		operation.Stop()
		operation.Log(err,
			// s3 dim info
			zap.String("bucket", s3Object.S3Bucket),
			zap.String("key", s3Object.S3ObjectKey))
	}()

	s3Client, err := getS3Client(s3Object.S3Bucket, topicArn)
	if err != nil {
		err = errors.Wrapf(err, "failed to get S3 client for s3://%s/%s",
			s3Object.S3Bucket, s3Object.S3ObjectKey)
		return nil, err
	}

	getObjectInput := &s3.GetObjectInput{
		Bucket: &s3Object.S3Bucket,
		Key:    &s3Object.S3ObjectKey,
	}
	output, err := s3Client.GetObject(getObjectInput)
	if err != nil {
		err = errors.Wrapf(err, "GetObject() failed for s3://%s/%s",
			s3Object.S3Bucket, s3Object.S3ObjectKey)
		return nil, err
	}

	bufferedReader := bufio.NewReader(output.Body)

	// We peek into the file header to identify the content type
	// http.DetectContentType only uses up to the first 512 bytes
	headerBytes, err := bufferedReader.Peek(512)
	if err != nil {
		if err != bufio.ErrBufferFull && err != io.EOF { // EOF or ErrBufferFull means file is shorter than n
			err = errors.Wrapf(err, "failed to Peek() in S3 payload for s3://%s/%s",
				s3Object.S3Bucket, s3Object.S3ObjectKey)
			return nil, err
		}
		err = nil // not really an error
	}
	contentType := http.DetectContentType(headerBytes)

	var streamReader io.Reader

	// Checking for prefix because the returned type can have also charset used
	if strings.HasPrefix(contentType, "text/plain") {
		// if it's plain text, just return the buffered reader
		streamReader = bufferedReader
	} else if strings.HasPrefix(contentType, "application/x-gzip") {
		var gzipReader *gzip.Reader
		gzipReader, err = gzip.NewReader(bufferedReader)
		if err != nil {
			err = errors.Wrapf(err, "failed to created gzip reader for s3://%s/%s",
				s3Object.S3Bucket, s3Object.S3ObjectKey)
			return nil, err
		}
		streamReader = gzipReader
	}

	dataStream = &common.DataStream{
		Reader: streamReader,
		Hints: common.DataStreamHints{
			S3: &common.S3DataStreamHints{
				Bucket:      s3Object.S3Bucket,
				Key:         s3Object.S3ObjectKey,
				ContentType: contentType,
			},
		},
	}
	return dataStream, err
}

// ParseNotification parses a message received
func ParseNotification(message string) ([]*S3ObjectInfo, error) {
	s3Objects, err := parseCloudTrailNotification(message)
	if err != nil {
		return nil, err
	}

	// If the input was not a CloudTrail notification, the result will be empty slice
	if len(s3Objects) == 0 {
		s3Objects, err = parseS3Event(message)
		if err != nil {
			return nil, err
		}
	}

	if len(s3Objects) == 0 {
		return nil, errors.New("notification is not of known type: " + message)
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
		return nil, errors.Wrap(err, "failed to parse CloudTrail event")
	}

	for _, s3Key := range cloudTrailNotification.S3ObjectKey {
		info := &S3ObjectInfo{
			S3Bucket:    *cloudTrailNotification.S3Bucket,
			S3ObjectKey: *s3Key,
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
		return nil, errors.Wrap(err, "failed to parse S3 event")
	}

	for _, record := range notification.Records {
		info := &S3ObjectInfo{
			S3Bucket:    record.S3.Bucket.Name,
			S3ObjectKey: record.S3.Object.Key,
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

// S3ObjectInfo contains information about the S3 object
type S3ObjectInfo struct {
	S3Bucket    string
	S3ObjectKey string
}

// SnsNotification struct represents an SNS message arriving to Panther SQS from a customer account.
// The message can either be of type 'Notification' or 'SubscriptionConfirmation'
// Since there is no AWS SDK-provided struct to represent both types
// we had to create this custom type to include fields from both types.
type SnsNotification struct {
	events.SNSEntity
	Token *string `json:"Token"`
}
