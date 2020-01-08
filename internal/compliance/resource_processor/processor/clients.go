package processor

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
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/kelseyhightower/envconfig"

	analysisapi "github.com/panther-labs/panther/api/gateway/analysis/client"
	complianceapi "github.com/panther-labs/panther/api/gateway/compliance/client"
	resourceapi "github.com/panther-labs/panther/api/gateway/resources/client"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

const maxBackoff = 30 * time.Second

type envConfig struct {
	AlertQueueURL     string `required:"true" split_words:"true"`
	AnalysisAPIHost   string `required:"true" split_words:"true"`
	AnalysisAPIPath   string `required:"true" split_words:"true"`
	PolicyEngine      string `required:"true" split_words:"true"`
	ComplianceAPIHost string `required:"true" split_words:"true"`
	ComplianceAPIPath string `required:"true" split_words:"true"`
	ResourceAPIHost   string `required:"true" split_words:"true"`
	ResourceAPIPath   string `required:"true" split_words:"true"`
}

var (
	env envConfig

	awsSession   *session.Session
	lambdaClient lambdaiface.LambdaAPI
	sqsClient    sqsiface.SQSAPI

	httpClient       *http.Client
	complianceClient *complianceapi.PantherCompliance
	analysisClient   *analysisapi.PantherAnalysis
	resourceClient   *resourceapi.PantherResources
)

// Setup parses the environment and initializes AWS and API clients.
func Setup() {
	envconfig.MustProcess("", &env)

	awsSession = session.Must(session.NewSession())
	lambdaClient = lambda.New(awsSession)
	sqsClient = sqs.New(awsSession)

	httpClient = gatewayapi.GatewayClient(awsSession)
	complianceClient = complianceapi.NewHTTPClientWithConfig(
		nil, complianceapi.DefaultTransportConfig().
			WithHost(env.ComplianceAPIHost).WithBasePath("/"+env.ComplianceAPIPath))
	analysisClient = analysisapi.NewHTTPClientWithConfig(
		nil, analysisapi.DefaultTransportConfig().
			WithHost(env.AnalysisAPIHost).WithBasePath("/"+env.AnalysisAPIPath))
	resourceClient = resourceapi.NewHTTPClientWithConfig(
		nil, resourceapi.DefaultTransportConfig().
			WithHost(env.ResourceAPIHost).WithBasePath("/"+env.ResourceAPIPath))
}
