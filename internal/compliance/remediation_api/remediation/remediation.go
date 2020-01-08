package remediation

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
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"

	"github.com/panther-labs/panther/api/gateway/remediation/models"
)

//InvokerAPI is the interface for the Invoker,
// the component that is responsible for invoking Remediation Lambda
type InvokerAPI interface {
	Remediate(*models.RemediateResource) error
	GetRemediations() (*models.Remediations, error)
}

//Invoker is responsible for invoking Remediation Lambda
type Invoker struct {
	lambdaClient lambdaiface.LambdaAPI
	awsSession   *session.Session
}

//NewInvoker method returns a new instance of Invoker
func NewInvoker(sess *session.Session) *Invoker {
	return &Invoker{
		lambdaClient: lambda.New(sess),
		awsSession:   sess,
	}
}
