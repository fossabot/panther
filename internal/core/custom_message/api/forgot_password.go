package custommessage

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

	"github.com/aws/aws-lambda-go/events"
	"github.com/matcornic/hermes"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/core/users_api/email"
)

func handleForgotPassword(event *events.CognitoEventUserPoolsCustomMessage) (*events.CognitoEventUserPoolsCustomMessage, error) {
	zap.L().Info("generate forget password email for:" + event.UserName)

	user, err := userGateway.GetUser(&event.UserName, &event.UserPoolID)
	if err != nil {
		zap.L().Error("failed to generate forget password html email for:"+event.UserName, zap.Error(err))
		return nil, err
	}

	emailParams := hermes.Email{
		Body: hermes.Body{
			Name: *user.GivenName + " " + *user.FamilyName,
			Intros: []string{
				`A password reset has been requested for this email address.
If you did not request a password reset, you can ignore this email.`,
			},
			Actions: []hermes.Action{
				{
					Instructions: "To set a new password for your Panther account, please click here:",
					Button: hermes.Button{
						TextColor: "#FFFFFF",
						Color:     "#6967F4", // Optional action button color
						Text:      "Reset my password",
						Link:      "https://" + appDomainURL + "/password-reset?token=" + event.Request.CodeParameter + "&email=" + *user.Email,
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := email.PantherEmailTemplate.GeneratePlainText(emailParams)

	// We have to do this because most email clients are not friendly with basic new line markup
	// replacing \n with a <br /> is the easiest way to mitigate this issue
	emailBody = strings.Replace(emailBody, "\n", "<br />", -1)
	if err != nil {
		zap.L().Error("failed to generate forget password html email for:"+event.UserName, zap.Error(err))
		return nil, err
	}
	event.Response.EmailMessage = emailBody
	event.Response.EmailSubject = "Panther Password Reset"
	return event, nil
}
