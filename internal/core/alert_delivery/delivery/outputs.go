package delivery

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
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

type outputCacheKey struct {
	OutputID string
}

type cachedOutput struct {
	Output    *outputmodels.AlertOutput
	Timestamp time.Time
}

type cachedOutputIDs struct {
	//Map from Severity -> List of output Ids
	Outputs   map[string][]*string
	Timestamp time.Time
}

func getRefreshInterval() time.Duration {
	intervalMins := os.Getenv("OUTPUTS_REFRESH_INTERVAL_MIN")
	if intervalMins == "" {
		intervalMins = "5"
	}
	return time.Duration(mustParseInt(intervalMins)) * time.Minute
}

var (
	alertOutputCache      = make(map[outputCacheKey]cachedOutput) // Map outputID to its credentials
	defaultOutputIDsCache *cachedOutputIDs                        // Map of organizationId to default output ids
	outputsAPI            = os.Getenv("OUTPUTS_API")
	refreshInterval       = getRefreshInterval()
)

// Get output ids for an alert
func getAlertOutputIds(alert *alertmodels.Alert) ([]*string, error) {
	if len(alert.OutputIDs) > 0 {
		return alert.OutputIDs, nil
	}

	if defaultOutputIDsCache != nil && time.Since(defaultOutputIDsCache.Timestamp) < refreshInterval {
		zap.L().Info("using cached output Ids")
		return defaultOutputIDsCache.Outputs[*alert.Severity], nil
	}

	zap.L().Info("getting default outputs")
	input := outputmodels.LambdaInput{GetDefaultOutputs: &outputmodels.GetDefaultOutputsInput{}}
	var defaultOutputs outputmodels.GetDefaultOutputsOutput
	if err := genericapi.Invoke(lambdaClient, outputsAPI, &input, &defaultOutputs); err != nil {
		return nil, err
	}

	defaultOutputIDsCache = &cachedOutputIDs{
		Timestamp: time.Now(),
		Outputs:   make(map[string][]*string, len(defaultOutputs.Defaults)),
	}

	for _, output := range defaultOutputs.Defaults {
		defaultOutputIDsCache.Outputs[*output.Severity] = output.OutputIDs
	}

	zap.L().Debug("default output ids cache", zap.Any("cache", defaultOutputIDsCache))
	return defaultOutputIDsCache.Outputs[*alert.Severity], nil
}

// Get output details, either from in-memory cache or the outputs-api
func getOutput(outputID string) (*outputmodels.GetOutputOutput, error) {
	key := outputCacheKey{OutputID: outputID}

	if cached, ok := alertOutputCache[key]; ok && time.Since(cached.Timestamp) < refreshInterval {
		zap.L().Info("using cached outputs",
			zap.String("outputID", outputID))
		return cached.Output, nil
	}

	zap.L().Info("getting outputs from outputs-api",
		zap.String("outputID", outputID))

	input := outputmodels.LambdaInput{GetOutput: &outputmodels.GetOutputInput{OutputID: aws.String(outputID)}}
	var result outputmodels.GetOutputOutput
	if err := genericapi.Invoke(lambdaClient, outputsAPI, &input, &result); err != nil {
		return nil, err
	}

	alertOutputCache[key] = cachedOutput{Output: &result, Timestamp: time.Now()}
	return &result, nil
}
