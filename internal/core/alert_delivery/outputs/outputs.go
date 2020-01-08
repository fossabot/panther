package outputs

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
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var (
	policyURLPrefix = os.Getenv("POLICY_URL_PREFIX")
	alertURLPrefix  = os.Getenv("ALERT_URL_PREFIX")
)

// HTTPWrapper encapsulates the Golang's http client
type HTTPWrapper struct {
	httpClient HTTPiface
}

// PostInput type
type PostInput struct {
	url     *string
	body    map[string]interface{}
	headers map[string]*string
}

// HTTPWrapperiface is the interface for our wrapper around Golang's http client
type HTTPWrapperiface interface {
	post(*PostInput) *AlertDeliveryError
}

// HTTPiface is an interface for http.Client to simplify unit testing.
type HTTPiface interface {
	Do(*http.Request) (*http.Response, error)
}

// API is the interface for output delivery that can be used for mocks in tests.
type API interface {
	Slack(*alertmodels.Alert, *outputmodels.SlackConfig) *AlertDeliveryError
	Email(*alertmodels.Alert, *outputmodels.EmailConfig) *AlertDeliveryError
	PagerDuty(*alertmodels.Alert, *outputmodels.PagerDutyConfig) *AlertDeliveryError
	Github(*alertmodels.Alert, *outputmodels.GithubConfig) *AlertDeliveryError
	Jira(*alertmodels.Alert, *outputmodels.JiraConfig) *AlertDeliveryError
	Opsgenie(*alertmodels.Alert, *outputmodels.OpsgenieConfig) *AlertDeliveryError
	MsTeams(*alertmodels.Alert, *outputmodels.MsTeamsConfig) *AlertDeliveryError
	Sqs(*alertmodels.Alert, *outputmodels.SqsConfig) *AlertDeliveryError
	Sns(*alertmodels.Alert, *outputmodels.SnsConfig) *AlertDeliveryError
	getSnsClient(topicArn string) (snsiface.SNSAPI, error)
	getSqsClient(queueURL string) (sqsiface.SQSAPI, error)
}

// OutputClient encapsulates the clients that allow sending alerts to multiple outputs
type OutputClient struct {
	session     *session.Session
	httpWrapper HTTPWrapperiface
	sesClient   sesiface.SESAPI
	// Map from region -> client
	sqsClients map[string]sqsiface.SQSAPI
	snsClients map[string]snsiface.SNSAPI
	mailFrom   *string
}

// OutputClient must satisfy the API interface.
var _ API = (*OutputClient)(nil)

// New creates a new client for alert delivery.
func New(sess *session.Session) *OutputClient {
	return &OutputClient{
		session:     sess,
		httpWrapper: &HTTPWrapper{httpClient: &http.Client{}},
		// TODO Lazy initialization of clients
		sesClient:  ses.New(sess),
		sqsClients: make(map[string]sqsiface.SQSAPI),
		snsClients: make(map[string]snsiface.SNSAPI),
		mailFrom:   aws.String(os.Getenv("MAIL_FROM")),
	}
}

func generateAlertMessage(alert *alertmodels.Alert) *string {
	if aws.StringValue(alert.Type) == alertmodels.RuleType {
		return aws.String(getDisplayName(alert) + " failed")
	}
	return aws.String(getDisplayName(alert) + " failed on new resources")
}

func generateAlertTitle(alert *alertmodels.Alert) *string {
	if aws.StringValue(alert.Type) == alertmodels.RuleType {
		return aws.String("New Alert: " + getDisplayName(alert))
	}
	return aws.String("Policy Failure: " + getDisplayName(alert))
}

func getDisplayName(alert *alertmodels.Alert) string {
	if alert.PolicyName != nil && *alert.PolicyName != "" {
		return *alert.PolicyName
	}
	return *alert.PolicyID
}

func generateURL(alert *alertmodels.Alert) string {
	if aws.StringValue(alert.Type) == alertmodels.RuleType {
		return alertURLPrefix + *alert.AlertID
	}
	return policyURLPrefix + *alert.PolicyID
}
