package destinations

import (
	"os"

	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
)

// Destination defines the interface that all Destinations should follow
type Destination interface {
	SendEvents(parsedEventChannel chan *common.ParsedEvent) error
}

//CreateDestination the method returns the appropriate Destination based on configuration
func CreateDestination() Destination {
	zap.L().Info("creating destination")
	s3BucketName := os.Getenv("S3_BUCKET")

	if s3BucketName != "" {
		return createS3Destination(s3BucketName)
	}
	return createFirehoseDestination()
}

func createFirehoseDestination() Destination {
	client := firehose.New(common.Session)
	zap.L().Info("created Firehose destination")
	return &FirehoseDestination{
		client:         client,
		firehosePrefix: "panther",
	}
}

func createS3Destination(s3BucketName string) Destination {
	return &S3Destination{
		s3Client:    s3.New(common.Session),
		snsClient:   sns.New(common.Session),
		s3Bucket:    s3BucketName,
		snsTopicArn: os.Getenv("SNS_TOPIC_ARN"),
	}
}
