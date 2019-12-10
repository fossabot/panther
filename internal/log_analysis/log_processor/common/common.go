package common

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// ParsedEvent contains a single event that has already been processed
type ParsedEvent struct {
	Event   interface{} `json:"event"`
	LogType *string     `json:"logType"`
}

// BatchParsedEvents contains a list of parsed events
type BatchParsedEvents struct {
	Events []*ParsedEvent `json:"events" validate:"required"`
}

// DataStream represents a data stream that read by the processor
type DataStream struct {
	Reader *io.Reader
	// The log type if known
	// If it is nil, it means the log type hasn't been identified yet
	LogType *string
}

// Session AWS Session that can be used by components of the system
// Setting Max Retries to a higher number - we'd like to retry several times when reading from S3/pushing to Firehose before failing.
var Session = session.Must(session.NewSession(aws.NewConfig().WithMaxRetries(10)))
