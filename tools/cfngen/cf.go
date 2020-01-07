// Package cfngen generates CloudFormation from Go objects.
package cfngen

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
