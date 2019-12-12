package models

import "gopkg.in/go-playground/validator.v9"

// Validator builds a custom struct validator.
func Validator() *validator.Validate {
	result := validator.New()
	result.RegisterStructValidation(atLeastOneUpdate, &UpdateUserInput{})
	return result
}

func atLeastOneUpdate(sl validator.StructLevel) {
	in := sl.Current().Interface().(UpdateUserInput)
	if in.GivenName == nil && in.FamilyName == nil && in.Role == nil && in.PhoneNumber == nil {
		sl.ReportError(in, "FamilyName|GivenName|PhoneNumber|Role", "", "at_least_one_update", "")
	}
}
