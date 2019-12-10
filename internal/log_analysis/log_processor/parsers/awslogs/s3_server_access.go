package awslogs

import (
	"encoding/csv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
)

// S3ServerAccess is an AWS S3 Access Log.
// Log format & samples can be seen here: https://docs.aws.amazon.com/AmazonS3/latest/dev/LogFormat.html
type S3ServerAccess struct {
	BucketOwner        *string    `json:"bucketowner,omitempty" validate:"required,len=64,alphanum"`
	Bucket             *string    `json:"bucket,omitempty"`
	Time               *time.Time `json:"time,omitempty"`
	RemoteIP           *string    `json:"remoteip,omitempty"`
	Requester          *string    `json:"requester,omitempty"`
	RequestID          *string    `json:"requestid,omitempty"`
	Operation          *string    `json:"operation,omitempty"`
	Key                *string    `json:"key,omitempty"`
	RequestURI         *string    `json:"requesturi,omitempty"`
	HTTPStatus         *int       `json:"httpstatus,omitempty" validate:"required,max=600,min=100"`
	ErrorCode          *string    `json:"errorcode,omitempty"`
	BytesSent          *int       `json:"bytessent,omitempty"`
	ObjectSize         *int       `json:"objectsize,omitempty"`
	TotalTime          *int       `json:"totaltime,omitempty"`
	TurnAroundTime     *int       `json:"turnaroundtime,omitempty"`
	Referrer           *string    `json:"referrer,omitempty"`
	UserAgent          *string    `json:"useragent,omitempty"`
	VersionID          *string    `json:"versionid,omitempty"`
	HostID             *string    `json:"hostid,omitempty"`
	SignatureVersion   *string    `json:"signatureversion,omitempty"`
	CipherSuite        *string    `json:"ciphersuite,omitempty"`
	AuthenticationType *string    `json:"authenticationtype,omitempty"`
	HostHeader         *string    `json:"hostheader,omitempty"`
	TLSVersion         *string    `json:"tlsVersion,omitempty"`
	AdditionalFields   []string   `json:"additionalFields,omitempty"`
}

// S3ServerAccessParser parses AWS S3 Server Access logs
type S3ServerAccessParser struct{}

// Parse returns the parsed events or nil if parsing failed
func (p *S3ServerAccessParser) Parse(log string) []interface{} {
	reader := csv.NewReader(strings.NewReader(log))
	reader.LazyQuotes = true
	reader.Comma = ' '

	records, err := reader.ReadAll()
	if err != nil {
		zap.L().Debug("failed to parse the log as csv")
		return nil
	}

	record := records[0]
	// The time in the logs is represented as [06/Feb/2019:00:00:38 +0000]
	// The CSV reader will break the above date to two different fields `[06/Feb/2019:00:00:38` and `+0000]`
	// We concatenate these fields before trying to parse them
	parsedTime, err := time.Parse("[2/Jan/2006:15:04:05-0700]", record[2]+record[3])
	if err != nil {
		zap.L().Debug("failed to parse timestamp log as csv")
		return nil
	}

	var additionalFields []string = nil
	if len(record) > 25 {
		additionalFields = record[25:]
	}

	event := &S3ServerAccess{
		BucketOwner:        csvStringToPointer(record[0]),
		Bucket:             csvStringToPointer(record[1]),
		Time:               aws.Time(parsedTime.In(time.UTC)),
		RemoteIP:           csvStringToPointer(record[4]),
		Requester:          csvStringToPointer(record[5]),
		RequestID:          csvStringToPointer(record[6]),
		Operation:          csvStringToPointer(record[7]),
		Key:                csvStringToPointer(record[8]),
		RequestURI:         csvStringToPointer(record[9]),
		HTTPStatus:         csvStringToIntPointer(record[10]),
		ErrorCode:          csvStringToPointer(record[11]),
		BytesSent:          csvStringToIntPointer(record[12]),
		ObjectSize:         csvStringToIntPointer(record[13]),
		TotalTime:          csvStringToIntPointer(record[14]),
		TurnAroundTime:     csvStringToIntPointer(record[15]),
		Referrer:           csvStringToPointer(record[16]),
		UserAgent:          csvStringToPointer(record[17]),
		VersionID:          csvStringToPointer(record[18]),
		HostID:             csvStringToPointer(record[19]),
		SignatureVersion:   csvStringToPointer(record[20]),
		CipherSuite:        csvStringToPointer(record[21]),
		AuthenticationType: csvStringToPointer(record[22]),
		HostHeader:         csvStringToPointer(record[23]),
		TLSVersion:         csvStringToPointer(record[24]),
		AdditionalFields:   additionalFields,
	}

	if err := parsers.Validator.Struct(event); err != nil {
		zap.L().Debug("failed to validate log", zap.Error(err))
		return nil
	}

	return []interface{}{event}
}

// LogType returns the log type supported by this parser
func (p *S3ServerAccessParser) LogType() string {
	return "AWS.S3ServerAccess"
}