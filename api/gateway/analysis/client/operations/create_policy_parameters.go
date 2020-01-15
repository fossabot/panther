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
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	strfmt "github.com/go-openapi/strfmt"

	models "github.com/panther-labs/panther/api/gateway/analysis/models"
)

// NewCreatePolicyParams creates a new CreatePolicyParams object
// with the default values initialized.
func NewCreatePolicyParams() *CreatePolicyParams {
	var ()
	return &CreatePolicyParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreatePolicyParamsWithTimeout creates a new CreatePolicyParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreatePolicyParamsWithTimeout(timeout time.Duration) *CreatePolicyParams {
	var ()
	return &CreatePolicyParams{

		timeout: timeout,
	}
}

// NewCreatePolicyParamsWithContext creates a new CreatePolicyParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreatePolicyParamsWithContext(ctx context.Context) *CreatePolicyParams {
	var ()
	return &CreatePolicyParams{

		Context: ctx,
	}
}

// NewCreatePolicyParamsWithHTTPClient creates a new CreatePolicyParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreatePolicyParamsWithHTTPClient(client *http.Client) *CreatePolicyParams {
	var ()
	return &CreatePolicyParams{
		HTTPClient: client,
	}
}

/*CreatePolicyParams contains all the parameters to send to the API endpoint
for the create policy operation typically these are written to a http.Request
*/
type CreatePolicyParams struct {

	/*Body*/
	Body *models.UpdatePolicy

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the create policy params
func (o *CreatePolicyParams) WithTimeout(timeout time.Duration) *CreatePolicyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create policy params
func (o *CreatePolicyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create policy params
func (o *CreatePolicyParams) WithContext(ctx context.Context) *CreatePolicyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create policy params
func (o *CreatePolicyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create policy params
func (o *CreatePolicyParams) WithHTTPClient(client *http.Client) *CreatePolicyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create policy params
func (o *CreatePolicyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create policy params
func (o *CreatePolicyParams) WithBody(body *models.UpdatePolicy) *CreatePolicyParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create policy params
func (o *CreatePolicyParams) SetBody(body *models.UpdatePolicy) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreatePolicyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
