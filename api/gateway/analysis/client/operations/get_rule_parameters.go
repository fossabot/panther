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
)

// NewGetRuleParams creates a new GetRuleParams object
// with the default values initialized.
func NewGetRuleParams() *GetRuleParams {
	var ()
	return &GetRuleParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetRuleParamsWithTimeout creates a new GetRuleParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetRuleParamsWithTimeout(timeout time.Duration) *GetRuleParams {
	var ()
	return &GetRuleParams{

		timeout: timeout,
	}
}

// NewGetRuleParamsWithContext creates a new GetRuleParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetRuleParamsWithContext(ctx context.Context) *GetRuleParams {
	var ()
	return &GetRuleParams{

		Context: ctx,
	}
}

// NewGetRuleParamsWithHTTPClient creates a new GetRuleParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetRuleParamsWithHTTPClient(client *http.Client) *GetRuleParams {
	var ()
	return &GetRuleParams{
		HTTPClient: client,
	}
}

/*GetRuleParams contains all the parameters to send to the API endpoint
for the get rule operation typically these are written to a http.Request
*/
type GetRuleParams struct {

	/*RuleID
	  Unique ASCII rule identifier

	*/
	RuleID string
	/*VersionID
	  Optional version ID to retrieve (for older versions)

	*/
	VersionID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get rule params
func (o *GetRuleParams) WithTimeout(timeout time.Duration) *GetRuleParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get rule params
func (o *GetRuleParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get rule params
func (o *GetRuleParams) WithContext(ctx context.Context) *GetRuleParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get rule params
func (o *GetRuleParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get rule params
func (o *GetRuleParams) WithHTTPClient(client *http.Client) *GetRuleParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get rule params
func (o *GetRuleParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRuleID adds the ruleID to the get rule params
func (o *GetRuleParams) WithRuleID(ruleID string) *GetRuleParams {
	o.SetRuleID(ruleID)
	return o
}

// SetRuleID adds the ruleId to the get rule params
func (o *GetRuleParams) SetRuleID(ruleID string) {
	o.RuleID = ruleID
}

// WithVersionID adds the versionID to the get rule params
func (o *GetRuleParams) WithVersionID(versionID *string) *GetRuleParams {
	o.SetVersionID(versionID)
	return o
}

// SetVersionID adds the versionId to the get rule params
func (o *GetRuleParams) SetVersionID(versionID *string) {
	o.VersionID = versionID
}

// WriteToRequest writes these params to a swagger request
func (o *GetRuleParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param ruleId
	qrRuleID := o.RuleID
	qRuleID := qrRuleID
	if qRuleID != "" {
		if err := r.SetQueryParam("ruleId", qRuleID); err != nil {
			return err
		}
	}

	if o.VersionID != nil {

		// query param versionId
		var qrVersionID string
		if o.VersionID != nil {
			qrVersionID = *o.VersionID
		}
		qVersionID := qrVersionID
		if qVersionID != "" {
			if err := r.SetQueryParam("versionId", qVersionID); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
