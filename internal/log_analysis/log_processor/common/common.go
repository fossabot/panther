package common

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
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// ParsedEvent contains a single event that has already been processed
type ParsedEvent struct {
	Event   interface{} `json:"event"`
	LogType string      `json:"logType"`
}

// DataStream represents a data stream that read by the processor
type DataStream struct {
	Reader io.Reader
	Hints  DataStreamHints
	// The log type if known
	// If it is nil, it means the log type hasn't been identified yet
	LogType *string
}

// Used in a DataStream as meta data to describe the data
type DataStreamHints struct {
	S3 *S3DataStreamHints // if nil, no hint
}

// Used in a DataStreamHints as meta data to describe the S3 object backing the stream
type S3DataStreamHints struct {
	Bucket      string
	Key         string
	ContentType string
}

// S3Notification is sent when new data is available in S3
type S3Notification struct {
	// S3Bucket is name of the S3 Bucket where data is available
	S3Bucket *string `json:"s3Bucket"`
	// S3ObjectKey is the key of the S3 object that contains the new data
	S3ObjectKey *string `json:"s3ObjectKey"`
	// Events is the number of events in the S3 object
	Events *int `json:"events"`
	// Bytes is the uncompressed size in bytes of the S3 object
	Bytes *int `json:"bytes"`
	// Type is the type of data available in the S3 object (LogData,RuleOutput)
	Type *string `json:"type"`
	// ID is an identified for the data in the S3 object. In case of LogData this will be
	// the Log Type, in case of RuleOutput data this will be the RuleID
	ID *string `json:"id"`
}

const (
	LogData    = "LogData"
	RuleOutput = "RuleOutput"
)

// Session AWS Session that can be used by components of the system
// Setting Max Retries to a higher number - we'd like to retry several times when reading from S3/pushing to Firehose before failing.
var Session = session.Must(session.NewSession(aws.NewConfig().WithMaxRetries(10)))
