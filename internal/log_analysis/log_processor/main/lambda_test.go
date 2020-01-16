package main

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
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
)

// Replace global logger with an in-memory observer for tests.
func mockLogger() *observer.ObservedLogs {
	core, mockLog := observer.New(zap.DebugLevel)
	zap.ReplaceGlobals(zap.New(core))
	return mockLog
}

func TestProcessOpLog(t *testing.T) {
	logs := mockLogger()
	functionName := "myfunction"
	lc := lambdacontext.LambdaContext{
		InvokedFunctionArn: functionName,
	}
	err := process(&lc, events.SQSEvent{
		Records: []events.SQSMessage{}, // empty, should do no work
	})
	require.NoError(t, err)
	message := common.OpLogNamespace + ":" + common.OpLogComponent + ":" + functionName
	require.Equal(t, 1, len(logs.FilterMessage(message).All())) // should be just one like this
	assert.Equal(t, zapcore.InfoLevel, logs.FilterMessage(message).All()[0].Level)
	assert.Equal(t, message, logs.FilterMessage(message).All()[0].Entry.Message)
	serviceDim := logs.FilterMessage(message).All()[0].ContextMap()[common.OpLogLambdaServiceDim.Key]
	assert.Equal(t, common.OpLogLambdaServiceDim.String, serviceDim)
	// deal with native int type which is how this is defined
	sqsMessageCount := logs.FilterMessage(message).All()[0].ContextMap()["sqsMessageCount"]
	switch v := sqsMessageCount.(type) {
	case int64:
		assert.Equal(t, int64(0), v)
	case int32:
		assert.Equal(t, int32(0), v)
	default:
		t.Errorf("unknown type for sqsMessageCount: %#v", sqsMessageCount)
	}
}
