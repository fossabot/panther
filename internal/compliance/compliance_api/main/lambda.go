package main

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
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"

	"github.com/panther-labs/panther/internal/compliance/compliance_api/handlers"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var methodHandlers = map[string]gatewayapi.RequestHandler{
	"GET /describe-org":      handlers.DescribeOrg,
	"GET /describe-policy":   handlers.DescribePolicy,
	"GET /describe-resource": handlers.DescribeResource,
	"GET /org-overview":      handlers.GetOrgOverview,
	"GET /status":            handlers.GetStatus,

	"POST /delete": handlers.DeleteStatus,
	"POST /status": handlers.SetStatus,
	"POST /update": handlers.UpdateMetadata,
}

func main() {
	envconfig.MustProcess("", &handlers.Env)
	lambda.Start(gatewayapi.LambdaProxy(methodHandlers))
}
