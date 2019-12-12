package api

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
	tableName                               = os.Getenv("TABLE_NAME")
)

// API provides receiver methods for each route handler.
type API struct{}
