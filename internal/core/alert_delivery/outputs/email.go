package outputs

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

const emailTemplate = "<h2>Message</h2>%s<br>" +
	"<h2>Severity</h2>%s<br>" +
	"<h2>Runbook</h2>%s<br>" +
	"<h2>Description</h2>%s"

var sesConfigurationSet = os.Getenv("SES_CONFIGURATION_SET")

func generateEmailContent(alert *alertmodels.Alert) *string {
	messageField := fmt.Sprintf("<a href='%s'>%s</a>",
		generateURL(alert),
		aws.StringValue(generateAlertMessage(alert)))
	return aws.String(fmt.Sprintf(
		emailTemplate,
		messageField,
		aws.StringValue(alert.Severity),
		aws.StringValue(alert.Runbook),
		aws.StringValue(alert.PolicyDescription),
	))
}

// Email sends email to destination
func (client *OutputClient) Email(alert *alertmodels.Alert, config *outputmodels.EmailConfig) *AlertDeliveryError {
	emailInput := &ses.SendEmailInput{
		ConfigurationSetName: aws.String(sesConfigurationSet),
		Source:               client.mailFrom,
		Destination:          &ses.Destination{ToAddresses: []*string{config.DestinationAddress}},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    generateAlertTitle(alert),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    generateEmailContent(alert),
				},
			},
		},
	}

	if _, err := client.sesClient.SendEmail(emailInput); err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == ses.ErrCodeMessageRejected {
			return &AlertDeliveryError{Message: "request failed " + err.Error(), Permanent: false}
		}
		return &AlertDeliveryError{Message: "request failed " + err.Error(), Permanent: true}
	}

	return nil
}
