package api

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

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/panther-labs/panther/internal/compliance/snapshot_api/ddb"
)

var (
	db                                      = ddb.New(tableName)
	sess                                    = session.Must(session.NewSession())
	SQSClient               sqsiface.SQSAPI = sqs.New(sess)
	maxElapsedTime                          = 5 * time.Second
	snapshotPollersQueueURL                 = os.Getenv("SNAPSHOT_POLLERS_QUEUE_URL")
	logAnalysisQueueURL                     = os.Getenv("LOG_ANALYSIS_QUEUE_URL")
	tableName                               = os.Getenv("TABLE_NAME")
)

// API provides receiver methods for each route handler.
type API struct{}
