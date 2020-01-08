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
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/panther-labs/panther/api/gateway/resources/client"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var (
	awsSession                 = session.Must(session.NewSession())
	sqsClient  sqsiface.SQSAPI = sqs.New(awsSession)
	queueURL                   = os.Getenv("SNAPSHOT_QUEUE_URL")

	transportConfig = client.DefaultTransportConfig().
			WithBasePath("/" + os.Getenv("RESOURCES_API_PATH")).
			WithHost(os.Getenv("RESOURCES_API_FQDN"))
	apiClient  = client.NewHTTPClientWithConfig(nil, transportConfig)
	httpClient = gatewayapi.GatewayClient(awsSession)
)
