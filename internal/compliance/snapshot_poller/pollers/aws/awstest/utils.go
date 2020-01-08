package awstest

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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/go-openapi/strfmt"
)

// Example output from the AWS API
var (
	ExampleTimeParsed, _       = time.Parse(time.RFC3339, "2019-04-02T17:16:30+00:00")
	ExampleTime                = strfmt.DateTime(ExampleTimeParsed)
	ExampleDate                = aws.Time(ExampleTimeParsed)
	ExampleIntegrationID       = aws.String("8e39aa9d-9823-4872-a1bd-40fd8795634b")
	ExampleAuthSource          = "arn:aws:iam::123456789012:role/PantherAuditRole"
	ExampleAuthSource2         = "arn:aws:iam::210987654321:role/PantherAuditRole"
	ExampleAuthSourceParsedARN = ParseExampleAuthSourceARN(ExampleAuthSource)
	ExampleAccountId           = aws.String("123456789012")

	ExampleRegions = []*string{
		aws.String("ap-southeast-2"),
		aws.String("eu-central-1"),
		aws.String("us-west-2"),
	}
)

// ParseExampleAuthSourceARN returns a parsed Auth Source ARN
func ParseExampleAuthSourceARN(arnToParse string) arn.ARN {
	parsedArn, err := arn.Parse(arnToParse)
	if err != nil {
		return arn.ARN{}
	}

	return parsedArn
}
