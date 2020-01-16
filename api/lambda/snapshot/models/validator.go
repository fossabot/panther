package models

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
