package mage

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

const configFile = "deployments/panther_config.yml"

type bucketsParameters struct {
	AccessLogsBucketName string `yaml:"AccessLogsBucketName"`
}

type appParameters struct {
	CloudWatchLogRetentionDays   int    `yaml:"CloudWatchLogRetentionDays"`
	Debug                        bool   `yaml:"Debug"`
	LayerVersionArns             string `yaml:"LayerVersionArns"`
	PythonLayerVersionArn        string `yaml:"PythonLayerVersionArn"`
	WebApplicationCertificateArn string `yaml:"WebApplicationCertificateArn"`
	TracingMode                  string `yaml:"TracingMode"`
}

// PantherConfig describes the panther_config.yml file.
type PantherConfig struct {
	BucketsParameterValues bucketsParameters `yaml:"BucketsParameterValues"`
	AppParameterValues     appParameters     `yaml:"AppParameterValues"`
	PipLayer               []string          `yaml:"PipLayer"`
	InitialAnalysisSets    []string          `yaml:"InitialAnalysisSets"`
}
