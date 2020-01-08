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
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	usermodels "github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// GetVerificationStatus returns the verification status of an email address
func (verification *OutputVerification) GetVerificationStatus(input *models.AlertOutput) (*string, error) {
	if *input.OutputType != "email" {
		result := models.VerificationStatusSuccess
		return &result, nil
	}

	// Check if the email is a user's email
	isVerified, err := verification.isVerifiedUserEmail(input.OutputConfig.Email.DestinationAddress)
	if err != nil {
		return nil, err
	}
	if *isVerified {
		return aws.String(models.VerificationStatusSuccess), nil
	}

	// If the email is not from a verified user, check SES to see if it has been verified that way
	request := &ses.GetIdentityVerificationAttributesInput{
		Identities: []*string{input.OutputConfig.Email.DestinationAddress},
	}
	response, err := verification.sesClient.GetIdentityVerificationAttributes(request)
	if err != nil {
		return nil, err
	}
	verificationStatusAttributes := response.VerificationAttributes[aws.StringValue(input.OutputConfig.Email.DestinationAddress)]
	if verificationStatusAttributes == nil {
		return aws.String(models.VerificationStatusNotStarted), nil
	}

	switch *verificationStatusAttributes.VerificationStatus {
	case ses.VerificationStatusSuccess:
		return aws.String(models.VerificationStatusSuccess), nil
	case ses.VerificationStatusNotStarted:
		return aws.String(models.VerificationStatusNotStarted), nil
	case ses.VerificationStatusPending:
		return aws.String(models.VerificationStatusPending), nil
	default:
		return aws.String(models.VerificationStatusFailed), nil
	}
}

func (verification *OutputVerification) isVerifiedUserEmail(email *string) (*bool, error) {
	input := usermodels.LambdaInput{
		GetUserOrganizationAccess: &usermodels.GetUserOrganizationAccessInput{
			Email: email,
		},
	}
	var output usermodels.GetUserOrganizationAccessOutput
	if err := genericapi.Invoke(verification.lambdaClient, usersAPI, &input, &output); err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return aws.Bool(false), nil
		}
		return nil, err
	}
	return aws.Bool(true), nil
}
