package handlers

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
	"archive/zip"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/panther-labs/panther/api/gateway/analysis"
	"github.com/panther-labs/panther/api/gateway/analysis/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

type writeResult struct {
	item       *tableItem
	changeType int
	err        error
}

// BulkUpload uploads multiple policies from a zipfile.
func BulkUpload(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseBulkUpload(request)
	if err != nil {
		return badRequest(err)
	}

	policies, err := extractZipFile(input)
	if err != nil {
		return badRequest(err)
	}

	// Create/modify each policy in parallel
	results := make(chan writeResult)
	for _, policy := range policies {
		go func(item *tableItem) {
			defer func() {
				// Recover from panic so we don't block forever when waiting for routines to finish.
				if r := recover(); r != nil {
					zap.L().Error("panicked while processing item",
						zap.String("id", string(item.ID)), zap.Any("panic", r))
					results <- writeResult{err: errors.New("panicked goroutine")}
				}
			}()
			changeType, err := writeItem(item, input.UserID, nil)
			results <- writeResult{item: item, changeType: changeType, err: err}
		}(policy)
	}

	counts := &models.BulkUploadResult{
		ModifiedPolicies: aws.Int64(0),
		NewPolicies:      aws.Int64(0),
		TotalPolicies:    aws.Int64(0),

		ModifiedRules: aws.Int64(0),
		NewRules:      aws.Int64(0),
		TotalRules:    aws.Int64(0),
	}

	var response *events.APIGatewayProxyResponse

	// Wait for all the goroutines to finish.
	for range policies {
		result := <-results
		if result.err != nil {
			// Set the response with an error code - 4XX first, otherwise 5XX
			if result.err == errWrongType {
				msg := fmt.Sprintf("ID %s does not have expected type %s", result.item.ID, result.item.Type)
				response = gatewayapi.MarshalResponse(&models.Error{Message: &msg}, http.StatusConflict)
			} else if response == nil {
				// errExists and errNotExists do not apply here  -
				// bulk upload automatically creates or updates depending on whether it already exists
				response = &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}

			continue
		}

		if result.item.Type == typePolicy {
			*counts.TotalPolicies++
			if result.changeType == newItem {
				*counts.NewPolicies++
			} else if result.changeType == updatedItem {
				*counts.ModifiedPolicies++
			}
		} else {
			*counts.TotalRules++
			if result.changeType == newItem {
				*counts.NewRules++
			} else if result.changeType == updatedItem {
				*counts.ModifiedRules++
			}
		}
	}

	if response != nil {
		return response
	}
	return gatewayapi.MarshalResponse(counts, http.StatusOK)
}

func parseBulkUpload(request *events.APIGatewayProxyRequest) (*models.BulkUpload, error) {
	var result models.BulkUpload
	if err := jsoniter.UnmarshalFromString(request.Body, &result); err != nil {
		return nil, err
	}

	if err := result.Validate(nil); err != nil {
		return nil, err
	}

	return &result, nil
}

func extractZipFile(input *models.BulkUpload) (map[models.ID]*tableItem, error) {
	// Base64-decode
	content, err := base64.StdEncoding.DecodeString(string(input.Data))
	if err != nil {
		return nil, fmt.Errorf("base64 decoding failed: %s", err)
	}

	// Unzip in memory (the max request size is only 6 MB, so this should easily fit)
	zipReader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return nil, fmt.Errorf("zipReader failed: %s", err)
	}

	policyBodies := make(map[string]models.Body) // map base file name to contents
	result := make(map[models.ID]*tableItem)

	// Process each file
	for _, zipFile := range zipReader.File {
		if strings.HasSuffix(zipFile.Name, "/") {
			continue // skip directories (we will see their nested files next)
		}

		unzippedBytes, err := readZipFile(zipFile)
		if err != nil {
			return nil, fmt.Errorf("file extraction failed: %s: %s", zipFile.Name, err)
		}

		if strings.Contains(zipFile.Name, "__pycache__") {
			continue
		}

		var config analysis.Config

		switch strings.ToLower(filepath.Ext(zipFile.Name)) {
		case ".py":
			// Store the Python body to be referenced later
			policyBodies[filepath.Base(zipFile.Name)] = models.Body(unzippedBytes)
			continue
		case ".json":
			err = jsoniter.Unmarshal(unzippedBytes, &config)
		case ".yml", ".yaml":
			err = yaml.Unmarshal(unzippedBytes, &config)
		default:
			zap.L().Debug("skipped unsupported file", zap.String("fileName", zipFile.Name))
		}

		if err != nil {
			return nil, err
		}

		// Map the Config struct fields over to the fields we need to store in Dynamo
		policy := tableItem{
			AutoRemediationID:         models.AutoRemediationID(config.AutoRemediationID),
			AutoRemediationParameters: models.AutoRemediationParameters(config.AutoRemediationParameters),

			// Use filename as placeholder for the body which we lookup later
			Body: models.Body(config.Filename),

			Description:   models.Description(config.Description),
			DisplayName:   models.DisplayName(config.DisplayName),
			Enabled:       models.Enabled(config.Enabled),
			ID:            models.ID(config.PolicyID),
			Reference:     models.Reference(config.Reference),
			ResourceTypes: models.TypeSet(config.ResourceTypes),
			Runbook:       models.Runbook(config.Runbook),
			Severity:      models.Severity(strings.ToUpper(config.Severity)),
			Suppressions:  models.Suppressions(config.Suppressions),
			Tags:          config.Tags,
			Tests:         make([]*models.UnitTest, len(config.Tests)),
			Type:          strings.ToUpper(config.AnalysisType),
		}

		for i, test := range config.Tests {
			resource, err := jsoniter.MarshalToString(test.Resource)
			if err != nil {
				return nil, err
			}

			policy.Tests[i] = &models.UnitTest{
				ExpectedResult: models.TestExpectedResult(test.ExpectedResult),
				Name:           models.TestName(test.Name),
				Resource:       models.TestResource(resource),
				ResourceType:   models.TestResourceType(test.ResourceType),
			}
		}

		if _, exists := result[policy.ID]; exists {
			return nil, fmt.Errorf("multiple policy specs with ID %s", policy.ID)
		}
		result[policy.ID] = &policy
	}

	// Finish each policy by adding its body and then validate it
	for _, policy := range result {
		if body, ok := policyBodies[string(policy.Body)]; ok {
			policy.Body = body
			if err := validateUploadedPolicy(policy, input.UserID); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("policy %s is missing a body", policy.ID)
		}
	}

	return result, nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			zap.L().Error("error closing zip file", zap.Error(err))
		}
	}()
	return ioutil.ReadAll(f)
}

// Ensure that the uploaded policy is valid according to the API spec for a Policy
func validateUploadedPolicy(item *tableItem, userID models.UserID) error {
	if item.Type != typePolicy && item.Type != typeRule {
		return fmt.Errorf("policy ID %s is invalid: unknown analysis type %s", item.ID, item.Type)
	}

	policy := item.Policy(models.ComplianceStatusPASS) // Convert to the external Policy model for validation
	policy.CreatedAt = models.ModifyTime(time.Now())
	policy.CreatedBy = userID
	policy.LastModified = policy.CreatedAt
	policy.LastModifiedBy = userID
	policy.VersionID = "mock.version.id.mock.version.id."

	if err := policy.Validate(nil); err != nil {
		return fmt.Errorf("policy ID %s is invalid: %s", policy.ID, err)
	}
	return nil
}
