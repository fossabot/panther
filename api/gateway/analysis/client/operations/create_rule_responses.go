// Code generated by go-swagger; DO NOT EDIT.

package operations

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
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	strfmt "github.com/go-openapi/strfmt"

	models "github.com/panther-labs/panther/api/gateway/analysis/models"
)

// CreateRuleReader is a Reader for the CreateRule structure.
type CreateRuleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateRuleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateRuleCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewCreateRuleBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewCreateRuleConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewCreateRuleInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewCreateRuleCreated creates a CreateRuleCreated with default headers values
func NewCreateRuleCreated() *CreateRuleCreated {
	return &CreateRuleCreated{}
}

/*CreateRuleCreated handles this case with default header values.

Rule created successfully
*/
type CreateRuleCreated struct {
	Payload *models.Rule
}

func (o *CreateRuleCreated) Error() string {
	return fmt.Sprintf("[POST /rule][%d] createRuleCreated  %+v", 201, o.Payload)
}

func (o *CreateRuleCreated) GetPayload() *models.Rule {
	return o.Payload
}

func (o *CreateRuleCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Rule)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateRuleBadRequest creates a CreateRuleBadRequest with default headers values
func NewCreateRuleBadRequest() *CreateRuleBadRequest {
	return &CreateRuleBadRequest{}
}

/*CreateRuleBadRequest handles this case with default header values.

Bad request
*/
type CreateRuleBadRequest struct {
	Payload *models.Error
}

func (o *CreateRuleBadRequest) Error() string {
	return fmt.Sprintf("[POST /rule][%d] createRuleBadRequest  %+v", 400, o.Payload)
}

func (o *CreateRuleBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *CreateRuleBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateRuleConflict creates a CreateRuleConflict with default headers values
func NewCreateRuleConflict() *CreateRuleConflict {
	return &CreateRuleConflict{}
}

/*CreateRuleConflict handles this case with default header values.

Rule or policy with the given ID already exists
*/
type CreateRuleConflict struct {
}

func (o *CreateRuleConflict) Error() string {
	return fmt.Sprintf("[POST /rule][%d] createRuleConflict ", 409)
}

func (o *CreateRuleConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateRuleInternalServerError creates a CreateRuleInternalServerError with default headers values
func NewCreateRuleInternalServerError() *CreateRuleInternalServerError {
	return &CreateRuleInternalServerError{}
}

/*CreateRuleInternalServerError handles this case with default header values.

Internal server error
*/
type CreateRuleInternalServerError struct {
}

func (o *CreateRuleInternalServerError) Error() string {
	return fmt.Sprintf("[POST /rule][%d] createRuleInternalServerError ", 500)
}

func (o *CreateRuleInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
