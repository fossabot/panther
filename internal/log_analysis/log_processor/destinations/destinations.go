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
	"os"

	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
)

// Destination defines the interface that all Destinations should follow
type Destination interface {
	SendEvents(parsedEventChannel chan *common.ParsedEvent, errChan chan error)
}

//CreateDestination the method returns the appropriate Destination based on configuration
func CreateDestination() Destination {
	zap.L().Debug("creating S3 destination")
	s3BucketName := os.Getenv("S3_BUCKET")

	if s3BucketName != "" {
		return createS3Destination(s3BucketName)
	}
	return createFirehoseDestination()
}

func createFirehoseDestination() Destination {
	client := firehose.New(common.Session)
	zap.L().Debug("created Firehose destination")
	return &FirehoseDestination{
		client:         client,
		firehosePrefix: "panther",
	}
}

func createS3Destination(s3BucketName string) Destination {
	return &S3Destination{
		s3Client:             s3.New(common.Session),
		snsClient:            sns.New(common.Session),
		glueClient:           glue.New(common.Session),
		s3Bucket:             s3BucketName,
		snsTopicArn:          os.Getenv("SNS_TOPIC_ARN"),
		partitionExistsCache: make(map[string]struct{}),
	}
}
