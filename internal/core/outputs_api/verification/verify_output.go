package verification

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
