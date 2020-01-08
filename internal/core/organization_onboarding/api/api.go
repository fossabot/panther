// Package api defines CRUD actions for the Cognito Api.
package api

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
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/panther-labs/panther/internal/core/organization_onboarding/gateway"
)

// The API has receiver methods for each of the handlers.
type API struct{}

var (
	awsSession                      = session.Must(session.NewSession())
	stepFunctionGateway gateway.API = gateway.New(awsSession)
)
