package models

import (
	"github.com/aws/aws-sdk-go/aws/arn"
	"gopkg.in/go-playground/validator.v9"
)

// Validator builds a custom struct validator.
func Validator() (*validator.Validate, error) {
	result := validator.New()
	if err := result.RegisterValidation("roleArn", validateRoleArn); err != nil {
		return nil, err
	}
	return result, nil
}

func validateRoleArn(fl validator.FieldLevel) bool {
	fieldArn, err := arn.Parse(fl.Field().String())
	return err == nil && fieldArn.Service == "iam"
}
