package common

import (
	"go.uber.org/zap"

	"github.com/panther-labs/panther/pkg/oplog"
)

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

// labels for oplog
const (
	OpLogNamespace  = "Panther" // FIXME: move "up" in the stack
	OpLogComponent  = "LogProcessor"
	OpLogServiceDim = "Service"
)

var (
	OpLogManager = oplog.NewManager(OpLogNamespace, OpLogComponent)

	// cross cutting dimensions

	OpLogLambdaServiceDim    = zap.String(OpLogServiceDim, "lambda")
	OpLogS3ServiceDim        = zap.String(OpLogServiceDim, "s3")
	OpLogSNSServiceDim       = zap.String(OpLogServiceDim, "sns")
	OpLogProcessorServiceDim = zap.String(OpLogServiceDim, "processor")
	OpLogGlueServiceDim      = zap.String(OpLogServiceDim, "glue")

	/*
			  Example CloudWatch Insight queries this structure enables:

			  -- show latest activity
			  filter namespace="Panther" and component="LogProcessor"
				| fields @timestamp, operation, stats.LogType, stats.LogLineCount, stats.BytesProcessedCount, stats.EventCount,
		                   stats.SuccessfullyClassifiedCount, stats.ClassificationFailureCount, error
				| sort @timestamp desc
			    | limit 200

			  -- show latest errors
			  filter namespace="Panther" and component="LogProcessor"
			  | filter level='error'
			  | fields @timestamp, operation, stats.LogType, stats.LogLineCount, stats.BytesProcessedCount, stats.EventCount,
		                   stats.SuccessfullyClassifiedCount, stats.ClassificationFailureCount, error
			  | sort @timestamp desc
			  | limit 200

			  -- show all sns activity
			  filter namespace="Panther" and component="LogProcessor"
			  | filter Service='sns'
			  | fields @timestamp, topicArn
			  | sort @timestamp desc
			  | limit 200

			   -- show all s3 activity
			   filter namespace="Panther" and component="LogProcessor"
			   | filter Service='s3'
			   | fields @timestamp, bucket, key
			   | sort @timestamp desc
			   | limit 200

	*/

)
