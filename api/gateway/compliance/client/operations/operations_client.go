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

	"github.com/go-openapi/runtime"
	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new operations API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for operations API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
DeleteStatus deletes the status associated with one or more policies or resources
*/
func (a *Client) DeleteStatus(params *DeleteStatusParams) (*DeleteStatusOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteStatusParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "DeleteStatus",
		Method:             "POST",
		PathPattern:        "/delete",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DeleteStatusReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteStatusOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DeleteStatus: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
DescribeOrg lists pass fail status for every policy and resource in the org
*/
func (a *Client) DescribeOrg(params *DescribeOrgParams) (*DescribeOrgOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDescribeOrgParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "DescribeOrg",
		Method:             "GET",
		PathPattern:        "/describe-org",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DescribeOrgReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DescribeOrgOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DescribeOrg: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
DescribePolicy pages through resources affected by a specific policy
*/
func (a *Client) DescribePolicy(params *DescribePolicyParams) (*DescribePolicyOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDescribePolicyParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "DescribePolicy",
		Method:             "GET",
		PathPattern:        "/describe-policy",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DescribePolicyReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DescribePolicyOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DescribePolicy: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
DescribeResource pages through policies which affect this resource
*/
func (a *Client) DescribeResource(params *DescribeResourceParams) (*DescribeResourceOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDescribeResourceParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "DescribeResource",
		Method:             "GET",
		PathPattern:        "/describe-resource",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DescribeResourceReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DescribeResourceOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DescribeResource: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetOrgOverview gets account totals and top failing policies resources
*/
func (a *Client) GetOrgOverview(params *GetOrgOverviewParams) (*GetOrgOverviewOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetOrgOverviewParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetOrgOverview",
		Method:             "GET",
		PathPattern:        "/org-overview",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetOrgOverviewReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetOrgOverviewOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetOrgOverview: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetStatus gets compliance status for a single policy resource pair
*/
func (a *Client) GetStatus(params *GetStatusParams) (*GetStatusOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetStatusParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetStatus",
		Method:             "GET",
		PathPattern:        "/status",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetStatusReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetStatusOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetStatus: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SetStatus sets the compliance status for a batch of resource policy pairs
*/
func (a *Client) SetStatus(params *SetStatusParams) (*SetStatusCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetStatusParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "SetStatus",
		Method:             "POST",
		PathPattern:        "/status",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &SetStatusReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*SetStatusCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for SetStatus: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
UpdateMetadata updates
*/
func (a *Client) UpdateMetadata(params *UpdateMetadataParams) (*UpdateMetadataOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateMetadataParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "UpdateMetadata",
		Method:             "POST",
		PathPattern:        "/update",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &UpdateMetadataReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*UpdateMetadataOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for UpdateMetadata: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
