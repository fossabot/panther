package destinations

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
	"bytes"
	"compress/gzip"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/glue/glueiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/registry"
)

// s3ObjectKeyFormat represents the format of the S3 object key
const s3ObjectKeyFormat = "%s%s-%s.gz"

var (
	maxFileSize = 100 * 1000 * 1000 // 100MB uncompressed file size, should result in ~10MB output file size
	// It should always be greater than the maximum expected event (log line) size
	maxDuration      = 1 * time.Minute // Holding events for maximum 1 minute in memory
	newLineDelimiter = []byte("\n")

	parserRegistry registry.Interface = registry.AvailableParsers() // initialize
)

// S3Destination sends normalized events to S3
type S3Destination struct {
	s3Client   s3iface.S3API
	snsClient  snsiface.SNSAPI
	glueClient glueiface.GlueAPI
	// s3Bucket is the s3Bucket where the data will be stored
	s3Bucket string
	// snsTopic is the SNS Topic ARN where we will send the notification
	// when we store new data in S3
	snsTopicArn string
	// used to track existing glue partitions, avoids excessive Glue API calls
	partitionExistsCache map[string]struct{}
}

// SendEvents stores events in S3.
// It continuously reads events from outputChannel, groups them in batches per log type
// and stores them in the appropriate S3 path. If the method encounters an error
// it stops reading from the outputChannel, writes an error to the errorChannel and terminates
func (destination *S3Destination) SendEvents(parsedEventChannel chan *common.ParsedEvent) error {
	logTypeToBuffer := make(map[string]*s3EventBuffer)
	eventsProcessed := 0
	zap.L().Info("starting to read events from channel")
	for event := range parsedEventChannel {
		eventsProcessed++
		data, err := jsoniter.Marshal(event.Event)
		if err != nil {
			zap.L().Warn("failed to marshall event", zap.Error(err))
			return err
		}

		buffer, ok := logTypeToBuffer[event.LogType]
		if !ok {
			buffer = &s3EventBuffer{}
			logTypeToBuffer[event.LogType] = buffer
		}

		canAdd, err := buffer.addEvent(data)
		if err != nil {
			return err
		}
		if !canAdd {
			if err = destination.sendData(event.LogType, buffer); err != nil {
				return err
			}

			canAdd, err = buffer.addEvent(data)
			if err != nil {
				return err
			}
			if !canAdd {
				// We will reach this point only if a single marshalled event is greater than maxFileSize
				// something that shouldn't happen normally
				zap.L().Error("event doesn't fit in single s3 object and will be dropped",
					zap.String("logtype", event.LogType))
			}
		}

		// Check if any buffers has data for longer than 1 minute
		if err = destination.sendExpiredData(logTypeToBuffer); err != nil {
			zap.L().Warn("failed to send data to S3", zap.Error(err))
			return err
		}
	}

	zap.L().Info("output channel closed, sending last events")
	// If the channel has been closed
	// send the buffered messages before terminating
	for logType, data := range logTypeToBuffer {
		if err := destination.sendData(logType, data); err != nil {
			return err
		}
	}
	zap.L().Info("Finished sending messages", zap.Int("events", eventsProcessed))
	return nil
}

func (destination *S3Destination) sendExpiredData(logTypeToEvents map[string]*s3EventBuffer) error {
	currentTime := time.Now().UTC()
	for logType, buffer := range logTypeToEvents {
		if currentTime.Sub(buffer.firstEventProcessedTime) > maxDuration {
			err := destination.sendData(logType, buffer)
			if err != nil {
				return err
			}
			// delete the entry after sending the data
			delete(logTypeToEvents, logType)
		}
	}
	return nil
}

// sendData puts data in S3 and sends notification to SNS
func (destination *S3Destination) sendData(logType string, buffer *s3EventBuffer) error {
	key := getS3ObjectKey(logType, buffer.firstEventProcessedTime)

	s3Notification := &common.S3Notification{
		S3Bucket:    aws.String(destination.s3Bucket),
		S3ObjectKey: aws.String(key),
		Events:      aws.Int(buffer.events),
		Bytes:       aws.Int(buffer.bytes),
		Type:        aws.String(common.LogData),
		ID:          aws.String(logType),
	}

	payload, err := buffer.getBytes()
	if err != nil {
		zap.L().Warn("failed to get buffer bytes", zap.Error(err))
		return err
	}

	if err := buffer.reset(); err != nil {
		zap.L().Warn("failed to reset buffer", zap.Error(err))
		return err
	}

	request := &s3.PutObjectInput{
		Bucket: aws.String(destination.s3Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(payload),
	}
	if _, err := destination.s3Client.PutObject(request); err != nil {
		zap.L().Warn("failed to put object in s3",
			zap.String("bucket", *request.Bucket),
			zap.String("key", *request.Key),
			zap.Error(err))
		return err
	}

	destination.createGluePartition(logType, buffer)

	marshalledNotification, marshallingError := jsoniter.MarshalToString(s3Notification)
	if marshallingError != nil {
		zap.L().Warn("failed to marshal notification", zap.Error(err))
		return marshallingError
	}

	input := &sns.PublishInput{
		TopicArn: aws.String(destination.snsTopicArn),
		Message:  aws.String(marshalledNotification),
	}
	if _, err := destination.snsClient.Publish(input); err != nil {
		zap.L().Warn("failed to send notification to topic",
			zap.String("topicArn", destination.snsTopicArn), zap.Error(err))
		return err
	}
	return nil
}

// create glue partition (best effort and log)
func (destination *S3Destination) createGluePartition(logType string, buffer *s3EventBuffer) {
	glueMetadata := parserRegistry.LookupParser(logType).Glue
	partitionPath := glueMetadata.PartitionPrefix(buffer.firstEventProcessedTime)
	if _, exists := destination.partitionExistsCache[partitionPath]; !exists {
		partitionErr := glueMetadata.CreateJSONPartition(destination.glueClient, destination.s3Bucket, buffer.firstEventProcessedTime)
		if partitionErr != nil {
			if awsErr, ok := partitionErr.(awserr.Error); ok {
				if awsErr.Code() == "AlreadyExistsException" {
					destination.partitionExistsCache[partitionPath] = struct{}{} // remember
					return
				}
			}
			zap.L().Error("failed to create glue partition",
				zap.String("bucket", destination.s3Bucket),
				zap.String("partition", partitionPath),
				zap.Error(partitionErr))
		} else {
			destination.partitionExistsCache[partitionPath] = struct{}{} // remember
			zap.L().Info("created glue partition",
				zap.String("bucket", destination.s3Bucket),
				zap.String("partition", partitionPath))
		}
	}
}

func getS3ObjectKey(logType string, timestamp time.Time) string {
	return fmt.Sprintf(s3ObjectKeyFormat,
		parserRegistry.LookupParser(logType).Glue.PartitionPrefix(timestamp.UTC()), // get the path used in Glue table
		timestamp.Format("20060102T150405Z"),
		uuid.New().String())
}

// s3EventBuffer is a group of events of the same type
// that will be stored in the same S3 object
type s3EventBuffer struct {
	buffer                  *bytes.Buffer
	writer                  *gzip.Writer
	bytes                   int
	events                  int
	firstEventProcessedTime time.Time
}

// addEvent adds new data to the s3EventBuffer
// If it returns true, the record was added successfully.
// If it returns false, the record couldn't be added because the group has exceeded file size limit
func (b *s3EventBuffer) addEvent(event []byte) (bool, error) {
	if b.buffer == nil {
		b.buffer = &bytes.Buffer{}
		b.writer = gzip.NewWriter(b.buffer)
		b.firstEventProcessedTime = time.Now().UTC()
	}

	// The size of the batch in bytes if the event is added
	projectedFileSize := b.bytes + len(event) + len(newLineDelimiter)
	if projectedFileSize > maxFileSize {
		return false, nil
	}

	_, err := b.writer.Write(event)
	if err != nil {
		zap.L().Warn("failed to add data to buffer", zap.Error(err))
		return false, err
	}

	// Adding new line delimiter
	_, err = b.writer.Write(newLineDelimiter)
	if err != nil {
		zap.L().Warn("failed to add data to buffer", zap.Error(err))
		return false, err
	}
	b.bytes = projectedFileSize
	b.events++
	return true, nil
}

// getBytes returns all bytes in the buffer
func (b *s3EventBuffer) getBytes() ([]byte, error) {
	if err := b.writer.Close(); err != nil {
		return nil, err
	}
	return b.buffer.Bytes(), nil
}

// reset resets the s3EventBuffer
func (b *s3EventBuffer) reset() error {
	b.bytes = 0
	b.events = 0
	if err := b.writer.Close(); err != nil {
		return err
	}
	b.writer = nil
	b.buffer = nil
	return nil
}
