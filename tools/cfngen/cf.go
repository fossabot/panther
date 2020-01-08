// Package cfngen generates CloudFormation from Go objects.
package cfngen

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
	"io"

	jsoniter "github.com/json-iterator/go"
)

// FIXME: consider replacing this with CDK when a Go version is available.

// enable compatibility with encoding/json
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Represents a CF reference
type Ref struct {
	Ref string
}

// Represents a simple CF Fn::Sub using template params
type Sub struct {
	Sub string `json:"Fn::Sub"`
}

// Represents CF Parameter
type Parameter struct {
	Type          string
	Default       interface{} `json:",omitempty"`
	Description   string
	AllowedValues []interface{} `json:",omitempty"`
	MinValue      interface{}   `json:",omitempty"`
	MaxValue      interface{}   `json:",omitempty"`
}

type Output struct {
	Description string
	Value       interface{}
}

// Represents a CF template
type Template struct {
	AWSTemplateFormatVersion string
	Description              string                 `json:",omitempty"`
	Parameters               map[string]interface{} `json:",omitempty"`
	Resources                map[string]interface{} `json:",omitempty"`
	Outputs                  map[string]interface{} `json:",omitempty"`
}

// Emit CF as JSON
func (t *Template) WriteCloudFormation(w io.Writer) (err error) {
	jsonBytes, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return
	}
	_, err = w.Write(jsonBytes)
	return
}

// Create a CF template , use WriteCloudFormation() to emit.
func NewTemplate(description string, parameters map[string]interface{}, resources map[string]interface{},
	outputs map[string]interface{}) (t *Template) {

	t = &Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Description:              description,
		Parameters:               parameters,
		Resources:                resources,
		Outputs:                  outputs,
	}
	return
}
