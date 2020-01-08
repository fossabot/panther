package models

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
