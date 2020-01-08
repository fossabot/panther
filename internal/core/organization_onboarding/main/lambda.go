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
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/panther-labs/panther/api/lambda/onboarding/models"
	"github.com/panther-labs/panther/internal/core/organization_onboarding/api"
	"github.com/panther-labs/panther/pkg/genericapi"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

var router = genericapi.NewRouter(nil, &api.API{})

func lambdaHandler(ctx context.Context, input *models.LambdaInput) (interface{}, error) {
	_, logger := lambdalogger.ConfigureGlobal(ctx, nil)
	defer func() { _ = logger.Sync() }()
	return router.Handle(input)
}

func main() {
	lambda.Start(lambdaHandler)
}
