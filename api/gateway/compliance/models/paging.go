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
	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Paging paging
// swagger:model Paging
type Paging struct {

	// this page
	// Required: true
	// Minimum: 1
	ThisPage *int64 `json:"thisPage"`

	// total items
	// Required: true
	// Minimum: 1
	TotalItems *int64 `json:"totalItems"`

	// total pages
	// Required: true
	// Minimum: 1
	TotalPages *int64 `json:"totalPages"`
}

// Validate validates this paging
func (m *Paging) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateThisPage(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTotalItems(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTotalPages(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Paging) validateThisPage(formats strfmt.Registry) error {

	if err := validate.Required("thisPage", "body", m.ThisPage); err != nil {
		return err
	}

	if err := validate.MinimumInt("thisPage", "body", int64(*m.ThisPage), 1, false); err != nil {
		return err
	}

	return nil
}

func (m *Paging) validateTotalItems(formats strfmt.Registry) error {

	if err := validate.Required("totalItems", "body", m.TotalItems); err != nil {
		return err
	}

	if err := validate.MinimumInt("totalItems", "body", int64(*m.TotalItems), 1, false); err != nil {
		return err
	}

	return nil
}

func (m *Paging) validateTotalPages(formats strfmt.Registry) error {

	if err := validate.Required("totalPages", "body", m.TotalPages); err != nil {
		return err
	}

	if err := validate.MinimumInt("totalPages", "body", int64(*m.TotalPages), 1, false); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Paging) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Paging) UnmarshalBinary(b []byte) error {
	var res Paging
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
