package custommessage

import (
	"github.com/aws/aws-lambda-go/events"
)

// HandleEvent routes Custom Message event based on the triggerSource
func HandleEvent(event *events.CognitoEventUserPoolsCustomMessage) (*events.CognitoEventUserPoolsCustomMessage, error) {
	switch ts := event.TriggerSource; ts {
	case "CustomMessage_ForgotPassword":
		return handleForgotPassword(event)
	default:
		return event, nil
	}
}
