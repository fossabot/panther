package verification

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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

// VerifyOutput performs verification on a AlertOutput
// Note that in case the output is not an email, no action is performed.
// In case it is an email, we use SES's email verification mechanism.
func (verification *OutputVerification) VerifyOutput(input *models.AlertOutput) (*models.AlertOutput, error) {
	if *input.OutputType != "email" {
		return input, nil
	}
	request := &ses.SendCustomVerificationEmailInput{
		EmailAddress:         input.OutputConfig.Email.DestinationAddress,
		ConfigurationSetName: aws.String(sesConfigurationSet),
		TemplateName:         aws.String(emailVerificationTemplate),
	}
	response, err := verification.sesClient.SendCustomVerificationEmail(request)

	if err != nil {
		return nil, err
	}

	zap.L().Info("sent a verification email", zap.String("messageId", response.String()))
	input.VerificationStatus = aws.String(models.VerificationStatusPending)
	return input, nil
}
