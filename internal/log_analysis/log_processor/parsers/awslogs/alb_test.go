package awslogs

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/timestamp"
)

func TestHTTPLog(t *testing.T) {
	log := "http 2018-08-26T14:17:23.186641Z app/my-loadbalancer/50dc6c495c0c9188 192.168.131.39:2817 " +
		"10.0.0.1:80 0.000 0.001 0.000 200 200 34 366 \"GET http://www.example.com:80/ HTTP/1.1\" " +
		"\"curl/7.46.0\" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 " +
		"\"Root=1-58337262-36d228ad5d99923122bbe354\" \"-\" \"-\" 0 2018-08-26T14:17:23.186641Z \"forward\" \"-\" \"-\""

	expectedTime := time.Unix(1535293043, 186641000).UTC()

	expectedEvent := &ALB{
		Type:                   aws.String("http"),
		Timestamp:              (*timestamp.RFC3339)(&expectedTime),
		ELB:                    aws.String("app/my-loadbalancer/50dc6c495c0c9188"),
		ClientIP:               aws.String("192.168.131.39"),
		ClientPort:             aws.Int(2817),
		TargetIP:               aws.String("10.0.0.1"),
		TargetPort:             aws.Int(80),
		RequestProcessingTime:  aws.Float64(0.0),
		TargetProcessingTime:   aws.Float64(0.001),
		ResponseProcessingTime: aws.Float64(0.000),
		ELBStatusCode:          aws.Int(200),
		TargetStatusCode:       aws.Int(200),
		ReceivedBytes:          aws.Int(34),
		SentBytes:              aws.Int(366),
		RequestHTTPMethod:      aws.String("GET"),
		RequestHTTPVersion:     aws.String("HTTP/1.1"),
		RequestURL:             aws.String("http://www.example.com:80/"),
		UserAgent:              aws.String("curl/7.46.0"),
		SSLCipher:              nil,
		SSLProtocol:            nil,
		TargetGroupARN:         aws.String("arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067"),
		TraceID:                aws.String("Root=1-58337262-36d228ad5d99923122bbe354"),
		DomainName:             nil,
		ChosenCertARN:          nil,
		MatchedRulePriority:    aws.Int(0),
		RequestCreationTime:    (*timestamp.RFC3339)(&expectedTime),
		ActionsExecuted:        []string{"forward"},
		RedirectURL:            nil,
		ErrorReason:            nil,
	}

	parser := &ALBParser{}
	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestHTTPSLog(t *testing.T) {
	log := "https 2018-08-26T14:17:23.186641Z app/my-loadbalancer/50dc6c495c0c9188 192.168.131.39:2817 10.0.0.1:80 " +
		"0.086 0.048 0.037 200 200 0 57 \"GET https://www.example.com:443/ HTTP/1.1\" \"curl/7.46.0\" " +
		"ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 " +
		"\"Root=1-58337281-1d84f3d73c47ec4e58577259\" \"www.example.com\" " +
		"\"arn:aws:acm:us-east-2:123456789012:certificate/12345678-1234-1234-1234-123456789012\" " +
		"1 2018-08-26T14:17:23.186641Z \"authenticate,forward\" \"-\" \"-\""

	expectedTime := time.Unix(1535293043, 186641000).UTC()

	expectedEvent := &ALB{
		Type:                   aws.String("https"),
		Timestamp:              (*timestamp.RFC3339)(&expectedTime),
		ELB:                    aws.String("app/my-loadbalancer/50dc6c495c0c9188"),
		ClientIP:               aws.String("192.168.131.39"),
		ClientPort:             aws.Int(2817),
		TargetIP:               aws.String("10.0.0.1"),
		TargetPort:             aws.Int(80),
		RequestProcessingTime:  aws.Float64(0.086),
		TargetProcessingTime:   aws.Float64(0.048),
		ResponseProcessingTime: aws.Float64(0.037),
		ELBStatusCode:          aws.Int(200),
		TargetStatusCode:       aws.Int(200),
		ReceivedBytes:          aws.Int(0),
		SentBytes:              aws.Int(57),
		RequestHTTPMethod:      aws.String("GET"),
		RequestHTTPVersion:     aws.String("HTTP/1.1"),
		RequestURL:             aws.String("https://www.example.com:443/"),
		UserAgent:              aws.String("curl/7.46.0"),
		SSLCipher:              aws.String("ECDHE-RSA-AES128-GCM-SHA256"),
		SSLProtocol:            aws.String("TLSv1.2"),
		TargetGroupARN:         aws.String("arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067"),
		TraceID:                aws.String("Root=1-58337281-1d84f3d73c47ec4e58577259"),
		DomainName:             aws.String("www.example.com"),
		ChosenCertARN:          aws.String("arn:aws:acm:us-east-2:123456789012:certificate/12345678-1234-1234-1234-123456789012"),
		MatchedRulePriority:    aws.Int(1),
		RequestCreationTime:    (*timestamp.RFC3339)(&expectedTime),
		ActionsExecuted:        []string{"authenticate", "forward"},
		RedirectURL:            nil,
		ErrorReason:            nil,
	}

	parser := &ALBParser{}
	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestHTTP2Log(t *testing.T) {
	//nolint:lll
	log := "h2 2018-08-26T14:17:23.186641Z app/my-loadbalancer/50dc6c495c0c9188 10.0.1.252:48160 " +
		"10.0.0.66:9000 0.000 0.002 0.000 200 200 5 257 \"GET https://10.0.2.105:773/ HTTP/2.0\" " +
		"\"curl/7.46.0\" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 " +
		"\"Root=1-58337327-72bd00b0343d75b906739c42\" \"-\" \"-\" 1 2018-08-26T14:17:23.186641Z " +
		"\"redirect\" \"https://example.com:80/\" \"-\""

	expectedTime := time.Unix(1535293043, 186641000).UTC()

	expectedEvent := &ALB{
		Type:                   aws.String("h2"),
		Timestamp:              (*timestamp.RFC3339)(&expectedTime),
		ELB:                    aws.String("app/my-loadbalancer/50dc6c495c0c9188"),
		ClientIP:               aws.String("10.0.1.252"),
		ClientPort:             aws.Int(48160),
		TargetIP:               aws.String("10.0.0.66"),
		TargetPort:             aws.Int(9000),
		RequestProcessingTime:  aws.Float64(0.000),
		TargetProcessingTime:   aws.Float64(0.002),
		ResponseProcessingTime: aws.Float64(0.000),
		ELBStatusCode:          aws.Int(200),
		TargetStatusCode:       aws.Int(200),
		ReceivedBytes:          aws.Int(5),
		SentBytes:              aws.Int(257),
		RequestHTTPMethod:      aws.String("GET"),
		RequestHTTPVersion:     aws.String("HTTP/2.0"),
		RequestURL:             aws.String("https://10.0.2.105:773/"),
		UserAgent:              aws.String("curl/7.46.0"),
		SSLCipher:              aws.String("ECDHE-RSA-AES128-GCM-SHA256"),
		SSLProtocol:            aws.String("TLSv1.2"),
		TargetGroupARN:         aws.String("arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067"),
		TraceID:                aws.String("Root=1-58337327-72bd00b0343d75b906739c42"),
		DomainName:             nil,
		ChosenCertARN:          nil,
		MatchedRulePriority:    aws.Int(1),
		RequestCreationTime:    (*timestamp.RFC3339)(&expectedTime),
		ActionsExecuted:        []string{"redirect"},
		RedirectURL:            aws.String("https://example.com:80/"),
		ErrorReason:            nil,
	}

	parser := &ALBParser{}
	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestAlbLogType(t *testing.T) {
	parser := &ALBParser{}
	require.Equal(t, "AWS.ALB", parser.LogType())
}
