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
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"

	"github.com/panther-labs/panther/internal/core/users_api/gateway"
	users "github.com/panther-labs/panther/internal/core/users_api/table"
)

// The API has receiver methods for each of the handlers.
type API struct{}

var (
	organizationAPI = os.Getenv("ORGANIZATION_API")
	awsSession      = session.Must(session.NewSession())

	lambdaClient lambdaiface.LambdaAPI = lambda.New(awsSession)
	userGateway  gateway.API           = gateway.New(awsSession)
	userTable    users.API             = users.New(os.Getenv("USERS_TABLE_NAME"), awsSession)
)
