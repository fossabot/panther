package awslogs

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
)

func TestS3AccessLogGetHttpOk(t *testing.T) {
	parser := &S3ServerAccessParser{}
	//nolint:lll
	log := "79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be awsexamplebucket [06/Feb/2019:00:00:38 +0000] 192.0.2.3 79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be 3E57427F3EXAMPLE REST.GET.VERSIONING - \"GET /awsexamplebucket?versioning HTTP/1.1\" 200 - 113 - 7 - \"-\" \"S3Console/0.4\" - s9lzHYrFp76ZVxRcpX9+5cjAnEH2ROuNkd2BHfIa6UkFVdtjf5mKR3/eTPFvsiP/XV/VLi31234= SigV2 ECDHE-RSA-AES128-GCM-SHA256 AuthHeader awsexamplebucket.s3.amazonaws.com TLSV1.1"

	date := time.Unix(1549411238, 0).In(time.UTC)
	expectedEvent := &S3ServerAccess{
		BucketOwner:        aws.String("79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be"),
		Bucket:             aws.String("awsexamplebucket"),
		Time:               aws.Time(date),
		RemoteIP:           aws.String("192.0.2.3"),
		Requester:          aws.String("79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be"),
		RequestID:          aws.String("3E57427F3EXAMPLE"),
		Operation:          aws.String("REST.GET.VERSIONING"),
		RequestURI:         aws.String("GET /awsexamplebucket?versioning HTTP/1.1"),
		HTTPStatus:         aws.Int(200),
		BytesSent:          aws.Int(113),
		TotalTime:          aws.Int(7),
		UserAgent:          aws.String("S3Console/0.4"),
		HostID:             aws.String("s9lzHYrFp76ZVxRcpX9+5cjAnEH2ROuNkd2BHfIa6UkFVdtjf5mKR3/eTPFvsiP/XV/VLi31234="),
		SignatureVersion:   aws.String("SigV2"),
		CipherSuite:        aws.String("ECDHE-RSA-AES128-GCM-SHA256"),
		AuthenticationType: aws.String("AuthHeader"),
		HostHeader:         aws.String("awsexamplebucket.s3.amazonaws.com"),
		TLSVersion:         aws.String("TLSV1.1"),
	}

	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestS3AccessLogGetHttpNotFound(t *testing.T) {
	parser := &S3ServerAccessParser{}
	//nolint:lll
	log := `79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be awsexamplebucket [06/Feb/2019:00:00:38 +0000] 192.0.2.3 79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be A1206F460EXAMPLE REST.GET.BUCKETPOLICY - "GET /awsexamplebucket?policy HTTP/1.1" 404 NoSuchBucketPolicy 297 - 38 - "-" "S3Console/0.4" - BNaBsXZQQDbssi6xMBdBU2sLt+Yf5kZDmeBUP35sFoKa3sLLeMC78iwEIWxs99CRUrbS4n11234= SigV2 ECDHE-RSA-AES128-GCM-SHA256 AuthHeader awsexamplebucket.s3.amazonaws.com TLSV1.1`

	date := time.Unix(1549411238, 0).In(time.UTC)
	expectedEvent := &S3ServerAccess{
		BucketOwner:        aws.String("79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be"),
		Bucket:             aws.String("awsexamplebucket"),
		Time:               aws.Time(date),
		RemoteIP:           aws.String("192.0.2.3"),
		Requester:          aws.String("79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be"),
		RequestID:          aws.String("A1206F460EXAMPLE"),
		Operation:          aws.String("REST.GET.BUCKETPOLICY"),
		RequestURI:         aws.String("GET /awsexamplebucket?policy HTTP/1.1"),
		HTTPStatus:         aws.Int(404),
		ErrorCode:          aws.String("NoSuchBucketPolicy"),
		BytesSent:          aws.Int(297),
		TotalTime:          aws.Int(38),
		UserAgent:          aws.String("S3Console/0.4"),
		HostID:             aws.String("BNaBsXZQQDbssi6xMBdBU2sLt+Yf5kZDmeBUP35sFoKa3sLLeMC78iwEIWxs99CRUrbS4n11234="),
		SignatureVersion:   aws.String("SigV2"),
		CipherSuite:        aws.String("ECDHE-RSA-AES128-GCM-SHA256"),
		AuthenticationType: aws.String("AuthHeader"),
		HostHeader:         aws.String("awsexamplebucket.s3.amazonaws.com"),
		TLSVersion:         aws.String("TLSV1.1"),
	}

	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestS3AccessLogPutHttpOK(t *testing.T) {
	parser := &S3ServerAccessParser{}
	//nolint:lll
	log := `79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be awsexamplebucket [06/Feb/2019:00:00:38 +0000] 192.0.2.3 79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be DD6CC733AEXAMPLE REST.PUT.OBJECT s3-dg.pdf "PUT /awsexamplebucket/s3-dg.pdf HTTP/1.1" 200 - - 4406583 41754 28 "-" "S3Console/0.4" - 10S62Zv81kBW7BB6SX4XJ48o6kpcl6LPwEoizZQQxJd5qDSCTLX0TgS37kYUBKQW3+bPdrg1234= SigV4 ECDHE-RSA-AES128-SHA AuthHeader awsexamplebucket.s3.amazonaws.com TLSV1.1`

	date := time.Unix(1549411238, 0).In(time.UTC)
	expectedEvent := &S3ServerAccess{
		BucketOwner:        aws.String("79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be"),
		Bucket:             aws.String("awsexamplebucket"),
		Key:                aws.String("s3-dg.pdf"),
		ObjectSize:         aws.Int(4406583),
		Time:               aws.Time(date),
		RemoteIP:           aws.String("192.0.2.3"),
		Requester:          aws.String("79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be"),
		RequestID:          aws.String("DD6CC733AEXAMPLE"),
		Operation:          aws.String("REST.PUT.OBJECT"),
		RequestURI:         aws.String("PUT /awsexamplebucket/s3-dg.pdf HTTP/1.1"),
		HTTPStatus:         aws.Int(200),
		TotalTime:          aws.Int(41754),
		TurnAroundTime:     aws.Int(28),
		UserAgent:          aws.String("S3Console/0.4"),
		HostID:             aws.String("10S62Zv81kBW7BB6SX4XJ48o6kpcl6LPwEoizZQQxJd5qDSCTLX0TgS37kYUBKQW3+bPdrg1234="),
		SignatureVersion:   aws.String("SigV4"),
		CipherSuite:        aws.String("ECDHE-RSA-AES128-SHA"),
		AuthenticationType: aws.String("AuthHeader"),
		HostHeader:         aws.String("awsexamplebucket.s3.amazonaws.com"),
		TLSVersion:         aws.String("TLSV1.1"),
	}

	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestS3AccessLogPutHttpOKExtraFields(t *testing.T) {
	parser := &S3ServerAccessParser{}
	//nolint:lll
	log := `79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be awsexamplebucket [06/Feb/2019:00:00:38 +0000] 192.0.2.3 79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be DD6CC733AEXAMPLE REST.PUT.OBJECT s3-dg.pdf "PUT /awsexamplebucket/s3-dg.pdf HTTP/1.1" 200 - - 4406583 41754 28 "-" "S3Console/0.4" - 10S62Zv81kBW7BB6SX4XJ48o6kpcl6LPwEoizZQQxJd5qDSCTLX0TgS37kYUBKQW3+bPdrg1234= SigV4 ECDHE-RSA-AES128-SHA AuthHeader awsexamplebucket.s3.amazonaws.com TLSV1.1 test1 test2`

	date := time.Unix(1549411238, 0).In(time.UTC)
	expectedEvent := &S3ServerAccess{
		BucketOwner:        aws.String("79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be"),
		Bucket:             aws.String("awsexamplebucket"),
		Key:                aws.String("s3-dg.pdf"),
		ObjectSize:         aws.Int(4406583),
		Time:               aws.Time(date),
		RemoteIP:           aws.String("192.0.2.3"),
		Requester:          aws.String("79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be"),
		RequestID:          aws.String("DD6CC733AEXAMPLE"),
		Operation:          aws.String("REST.PUT.OBJECT"),
		RequestURI:         aws.String("PUT /awsexamplebucket/s3-dg.pdf HTTP/1.1"),
		HTTPStatus:         aws.Int(200),
		TotalTime:          aws.Int(41754),
		TurnAroundTime:     aws.Int(28),
		UserAgent:          aws.String("S3Console/0.4"),
		HostID:             aws.String("10S62Zv81kBW7BB6SX4XJ48o6kpcl6LPwEoizZQQxJd5qDSCTLX0TgS37kYUBKQW3+bPdrg1234="),
		SignatureVersion:   aws.String("SigV4"),
		CipherSuite:        aws.String("ECDHE-RSA-AES128-SHA"),
		AuthenticationType: aws.String("AuthHeader"),
		HostHeader:         aws.String("awsexamplebucket.s3.amazonaws.com"),
		TLSVersion:         aws.String("TLSV1.1"),
		AdditionalFields:   []string{"test1", "test2"},
	}

	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestS3ServerAccessLogType(t *testing.T) {
	parser := &S3ServerAccessParser{}
	require.Equal(t, "AWS.S3ServerAccess", parser.LogType())
}
