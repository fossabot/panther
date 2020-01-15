// Code generated by go-swagger; DO NOT EDIT.

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

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// TypeSet List of resource/log types to which this policy applies
// swagger:model TypeSet
type TypeSet []string

// Validate validates this type set
func (m TypeSet) Validate(formats strfmt.Registry) error {
	var res []error

	iTypeSetSize := int64(len(m))

	if err := validate.MaxItems("", "body", iTypeSetSize, 500); err != nil {
		return err
	}

	if err := validate.UniqueItems("", "body", m); err != nil {
		return err
	}

	for i := 0; i < len(m); i++ {

		if err := validate.MinLength(strconv.Itoa(i), "body", string(m[i]), 1); err != nil {
			return err
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
