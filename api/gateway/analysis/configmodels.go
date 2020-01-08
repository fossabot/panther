package analysis

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

// Config defines the file format when parsing a bulk upload.
//
// YAML tags required because the YAML unmarshaller needs them
// JSON tags not present because the JSON unmarshaller is easy
type Config struct {
	AnalysisType              string            `yaml:"AnalysisType"`
	AutoRemediationID         string            `yaml:"AutoRemediationID"`
	AutoRemediationParameters map[string]string `yaml:"AutoRemediationParameters"`
	Description               string            `yaml:"Description"`
	DisplayName               string            `yaml:"DisplayName"`
	Enabled                   bool              `yaml:"Enabled"`
	Filename                  string            `yaml:"Filename"`
	PolicyID                  string            `yaml:"PolicyID"`
	ResourceTypes             []string          `yaml:"ResourceTypes"`
	Reference                 string            `yaml:"Reference"`
	Runbook                   string            `yaml:"Runbook"`
	Severity                  string            `yaml:"Severity"`
	Suppressions              []string          `yaml:"Suppressions"`
	Tags                      []string          `yaml:"Tags"`
	Tests                     []Test            `yaml:"Tests"`
}

// Test is a unit test definition when parsing policies in a bulk upload.
type Test struct {
	ExpectedResult bool        `yaml:"ExpectedResult"`
	Name           string      `yaml:"Name"`
	Resource       interface{} `yaml:"Resource"`
	ResourceType   string      `yaml:"ResourceType"`
}
