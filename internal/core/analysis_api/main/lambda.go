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

	"github.com/panther-labs/panther/internal/core/analysis_api/handlers"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var methodHandlers = map[string]gatewayapi.RequestHandler{
	// Policies only
	"GET /list":      handlers.ListPolicies,
	"GET /policy":    handlers.GetPolicy,
	"POST /policy":   handlers.CreatePolicy,
	"POST /suppress": handlers.Suppress,
	"POST /update":   handlers.ModifyPolicy,
	"POST /upload":   handlers.BulkUpload,

	// Rules only
	"GET /rule":         handlers.GetRule,
	"POST /rule":        handlers.CreateRule,
	"GET /rule/list":    handlers.ListRules,
	"POST /rule/update": handlers.ModifyRule,

	// Rules and Policies
	"POST /delete": handlers.DeletePolicies,
	"GET /enabled": handlers.GetEnabledPolicies,
	"POST /test":   handlers.TestPolicy,
}

func main() {
	handlers.Setup()
	lambda.Start(gatewayapi.LambdaProxy(methodHandlers))
}
