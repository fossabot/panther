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

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/magefile/mage/mg"

	"github.com/panther-labs/panther/api/gateway/analysis/client"
	"github.com/panther-labs/panther/api/gateway/analysis/client/operations"
	"github.com/panther-labs/panther/api/gateway/analysis/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// Install Python rules/policies if they don't already exist.
func initializeAnalysisSets(awsSession *session.Session, endpoint string, config *PantherConfig) error {
	httpClient := gatewayapi.GatewayClient(awsSession)
	apiClient := client.NewHTTPClientWithConfig(nil, client.DefaultTransportConfig().
		WithBasePath("/v1").WithHost(endpoint))

	policies, err := apiClient.Operations.ListPolicies(&operations.ListPoliciesParams{
		PageSize:   aws.Int64(1),
		HTTPClient: httpClient,
	})
	if err != nil {
		return err
	}

	rules, err := apiClient.Operations.ListRules(&operations.ListRulesParams{
		PageSize:   aws.Int64(1),
		HTTPClient: httpClient,
	})
	if err != nil {
		return err
	}

	if len(policies.Payload.Policies) > 0 || len(rules.Payload.Rules) > 0 {
		if mg.Verbose() {
			fmt.Println("deploy: initial analysis set ignored: policies and/or rules already exist")
		}

		return nil
	}

	var newRules, newPolicies int64
	for _, path := range config.InitialAnalysisSets {
		fmt.Println("deploy: uploading " + path)
		var contents []byte
		if strings.HasPrefix(path, "file://") {
			contents, err = ioutil.ReadFile(strings.TrimPrefix(path, "file://"))
			if err != nil {
				return err
			}
		} else {
			contents, err = download(path)
			if err != nil {
				return err
			}
		}

		// BulkUpload to panther-analysis-api
		encoded := base64.StdEncoding.EncodeToString(contents)
		response, err := apiClient.Operations.BulkUpload(&operations.BulkUploadParams{
			Body: &models.BulkUpload{
				Data:   models.Base64zipfile(encoded),
				UserID: "00000000-0000-0000-0000-000000000000",
			},
			HTTPClient: httpClient,
		})
		if err != nil {
			return err
		}

		newRules += *response.Payload.NewRules
		newPolicies += *response.Payload.NewPolicies
	}

	fmt.Printf("deploy: added %d policies and %d rules\n", newPolicies, newRules)
	return nil
}
