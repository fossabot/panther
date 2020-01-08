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
