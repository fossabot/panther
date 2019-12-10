package destinations

import (
	"github.com/aws/aws-sdk-go/service/firehose"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
)

// Destination defines the interface that all Destinations should follow
type Destination interface {
	SendEvents(parsedEventChannel chan *common.ParsedEvent, errorChannel chan error)
}

//CreateDestination the method returns the appropriate Destination based on configuration
func CreateDestination() Destination {
	zap.L().Info("creating destination")
	return createFirehoseDestination()
}

func createFirehoseDestination() Destination {
	client := firehose.New(common.Session)
	zap.L().Info("created Firehose destination")
	return &FirehoseDestination{
		client:         client,
		firehosePrefix: "panther",
	}
}
