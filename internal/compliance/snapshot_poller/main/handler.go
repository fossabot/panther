package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers"
)

func main() {
	lambda.Start(pollers.Handle)
}
