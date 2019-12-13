package destinations

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
)

type mockS3 struct {
	s3iface.S3API
	mock.Mock
}

func (m *mockS3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

type mockSns struct {
	snsiface.SNSAPI
	mock.Mock
}

// testEvent is a test event used for the purposes of this test
type testEvent struct {
	data string
}

func (m *mockSns) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

func TestSendDataToS3BeforeTerminating(t *testing.T) {
	mockSns := &mockSns{}
	mockS3 := &mockS3{}
	destination := &S3Destination{
		snsTopicArn: "arn:aws:sns:us-west-2:123456789012:test",
		s3Bucket:    "testbucket",
		snsClient:   mockSns,
		s3Client:    mockS3,
	}
	eventChannel := make(chan *common.ParsedEvent, 1)

	parsedEvent := &common.ParsedEvent{
		Event:   testEvent{data: "test"},
		LogType: "testtype",
	}

	// sending event to buffered channel
	eventChannel <- parsedEvent
	close(eventChannel)

	marshalledEvent, _ := jsoniter.Marshal(parsedEvent.Event)

	mockS3.On("PutObject", mock.Anything).Return(&s3.PutObjectOutput{}, nil)
	mockSns.On("Publish", mock.Anything).Return(&sns.PublishOutput{}, nil)

	require.NoError(t, destination.SendEvents(eventChannel))

	// There is no way to know the key of the S3 object since we are generating it based on time
	// I am fetching it from the actual request performed to S3 and:
	//1. Verifying the S3 object key is of the correct format
	//2. Verifying the rest of the fields are as expected
	putObjectInput := mockS3.Calls[0].Arguments.Get(0).(*s3.PutObjectInput)
	// Gzipping the test event
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)

	writer.Write(marshalledEvent) //nolint:errcheck
	writer.Write([]byte("\n"))    //nolint:errcheck
	writer.Close()                //nolint:errcheck

	bodyBytes, _ := ioutil.ReadAll(putObjectInput.Body)
	require.Equal(t, aws.String("testbucket"), putObjectInput.Bucket)
	require.Equal(t, buffer.Bytes(), bodyBytes)

	// Verifying Sns Publish payload
	publishInput := mockSns.Calls[0].Arguments.Get(0).(*sns.PublishInput)
	expectedS3Notification := &common.S3Notification{
		S3Bucket:    aws.String("testbucket"),
		S3ObjectKey: putObjectInput.Key,
		Events:      aws.Int(1),
		Bytes:       aws.Int(len(marshalledEvent) + len("\n")),
		Type:        aws.String(common.LogData),
		ID:          aws.String("testtype"),
	}
	marshalledExpectedS3Notification, _ := jsoniter.MarshalToString(expectedS3Notification)
	expectedSnsPublishInput := &sns.PublishInput{
		Message:  aws.String(marshalledExpectedS3Notification),
		TopicArn: aws.String("arn:aws:sns:us-west-2:123456789012:test"),
	}
	require.Equal(t, expectedSnsPublishInput, publishInput)
}

func TestSendDataIfSizeLimitHasBeenReached(t *testing.T) {
	mockSns := &mockSns{}
	mockS3 := &mockS3{}
	destination := &S3Destination{
		snsTopicArn: "arn:aws:sns:us-west-2:123456789012:test",
		s3Bucket:    "testbucket",
		snsClient:   mockSns,
		s3Client:    mockS3,
	}
	eventChannel := make(chan *common.ParsedEvent, 2)

	// sending 2 events to buffered channel
	// The second should already cause the S3 object size limits to be exceeded
	// so we expect two objects to be written to s3
	eventChannel <- &common.ParsedEvent{
		Event:   testEvent{data: "test"},
		LogType: "testtype",
	}
	eventChannel <- &common.ParsedEvent{
		Event:   testEvent{data: "test"},
		LogType: "testtype",
	}
	close(eventChannel)

	mockS3.On("PutObject", mock.Anything).Return(&s3.PutObjectOutput{}, nil).Twice()
	mockSns.On("Publish", mock.Anything).Return(&sns.PublishOutput{}, nil).Twice()

	// This is the size of a single event
	// We expect this to cause the S3Destination to create two objects in S3
	maxFileSize = 3

	require.NoError(t, destination.SendEvents(eventChannel))
}

func TestSendDataIfTimeLimitHasBeenReached(t *testing.T) {
	mockSns := &mockSns{}
	mockS3 := &mockS3{}
	destination := &S3Destination{
		snsTopicArn: "arn:aws:sns:us-west-2:123456789012:test",
		s3Bucket:    "testbucket",
		snsClient:   mockSns,
		s3Client:    mockS3,
	}
	eventChannel := make(chan *common.ParsedEvent, 2)

	// sending 2 events to buffered channel
	// The second should already cause the S3 object size limits to be exceeded
	// so we expect two objects to be written to s3
	eventChannel <- &common.ParsedEvent{
		Event:   testEvent{data: "test"},
		LogType: "testtype",
	}
	eventChannel <- &common.ParsedEvent{
		Event:   testEvent{data: "test"},
		LogType: "testtype",
	}
	close(eventChannel)

	mockS3.On("PutObject", mock.Anything).Return(&s3.PutObjectOutput{}, nil).Twice()
	mockSns.On("Publish", mock.Anything).Return(&sns.PublishOutput{}, nil).Twice()

	// We expect this to cause the S3Destination to create two objects in S3
	maxDuration = 1 * time.Nanosecond

	require.NoError(t, destination.SendEvents(eventChannel))
}

func TestSendDataToS3FromMultipleLogTypesBeforeTerminating(t *testing.T) {
	mockSns := &mockSns{}
	mockS3 := &mockS3{}
	destination := &S3Destination{
		snsTopicArn: "arn:aws:sns:us-west-2:123456789012:test",
		s3Bucket:    "testbucket",
		snsClient:   mockSns,
		s3Client:    mockS3,
	}
	eventChannel := make(chan *common.ParsedEvent, 2)

	eventChannel <- &common.ParsedEvent{
		Event:   testEvent{data: "test"},
		LogType: "testtype1",
	}
	eventChannel <- &common.ParsedEvent{
		Event:   testEvent{data: "test"},
		LogType: "testtype2",
	}
	close(eventChannel)

	mockS3.On("PutObject", mock.Anything).Return(&s3.PutObjectOutput{}, nil).Twice()
	mockSns.On("Publish", mock.Anything).Return(&sns.PublishOutput{}, nil).Twice()

	require.NoError(t, destination.SendEvents(eventChannel))
}

func TestSendDataFailsIfS3Fails(t *testing.T) {
	mockSns := &mockSns{}
	mockS3 := &mockS3{}
	destination := &S3Destination{
		snsTopicArn: "arn:aws:sns:us-west-2:123456789012:test",
		s3Bucket:    "testbucket",
		snsClient:   mockSns,
		s3Client:    mockS3,
	}
	eventChannel := make(chan *common.ParsedEvent, 1)

	eventChannel <- &common.ParsedEvent{
		Event:   testEvent{data: "test"},
		LogType: "testtype",
	}
	close(eventChannel)

	mockS3.On("PutObject", mock.Anything).Return(&s3.PutObjectOutput{}, errors.New("")).Twice()

	require.Error(t, destination.SendEvents(eventChannel))
}

func TestSendDataFailsIfSnsFails(t *testing.T) {
	mockSns := &mockSns{}
	mockS3 := &mockS3{}
	destination := &S3Destination{
		snsTopicArn: "arn:aws:sns:us-west-2:123456789012:test",
		s3Bucket:    "testbucket",
		snsClient:   mockSns,
		s3Client:    mockS3,
	}
	eventChannel := make(chan *common.ParsedEvent, 1)

	eventChannel <- &common.ParsedEvent{
		Event:   testEvent{data: "test"},
		LogType: "testtype",
	}
	close(eventChannel)

	mockS3.On("PutObject", mock.Anything).Return(&s3.PutObjectOutput{}, nil)
	mockSns.On("Publish", mock.Anything).Return(&sns.PublishOutput{}, errors.New("test"))

	require.Error(t, destination.SendEvents(eventChannel))
}
