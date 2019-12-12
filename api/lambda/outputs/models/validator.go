package models

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"gopkg.in/go-playground/validator.v9"
)

// Validator builds a custom struct validator.
func Validator() (*validator.Validate, error) {
	result := validator.New()
	result.RegisterStructValidation(ensureOneOutput, &OutputConfig{})
	if err := result.RegisterValidation("snsArn", validateAwsArn); err != nil {
		return nil, err
	}
	return result, nil
}

var outputTypes = []string{"Slack", "Sns", "Email", "PagerDuty", "Github", "Jira", "Opsgenie", "MsTeams", "Sqs"}

func ensureOneOutput(sl validator.StructLevel) {
	input := sl.Current()

	count := 0
	for _, outputType := range outputTypes {
		if !input.FieldByName(outputType).IsNil() {
			count++
		}
	}

	if count != 1 {
		sl.ReportError(input, strings.Join(outputTypes, "|"), "", "exactly_one_output", "")
	}
}

func validateAwsArn(fl validator.FieldLevel) bool {
	fieldArn, err := arn.Parse(fl.Field().String())
	return err == nil && fieldArn.Service == "sns"
}
