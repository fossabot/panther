package pollers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/resources/client/operations"
	api "github.com/panther-labs/panther/api/gateway/resources/models"
	"github.com/panther-labs/panther/pkg/lambdalogger"

	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	pollers "github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws"
)

const resourcesAPIBatchSize = 500

// loadMessage marshals the incoming SQS message into a ScanMsg.
func loadMessage(messageBody string) (*pollermodels.ScanMsg, error) {
	msg := &pollermodels.ScanMsg{}
	err := jsoniter.Unmarshal([]byte(messageBody), msg)
	if err != nil {
		return nil, err
	}

	return msg, err
}

// batchResources creates groups of 500 resources to send to the ResourcesAPI.
func batchResources(resources []*api.AddResourceEntry) (batches [][]*api.AddResourceEntry) {
	for resourcesAPIBatchSize < len(resources) {
		resources, batches = resources[resourcesAPIBatchSize:], append(
			batches,
			resources[0:resourcesAPIBatchSize:resourcesAPIBatchSize],
		)
	}
	batches = append(batches, resources)
	return
}

// Handle is the main Lambda Handler.
func Handle(ctx context.Context, event events.SQSEvent) error {
	_, logger := lambdalogger.ConfigureGlobal(ctx, nil)
	logger.Info("starting snapshot-pollers",
		zap.Int("numEvents", len(event.Records)),
	)

	for indx, message := range event.Records {
		zap.L().Debug("loading message from the queue")
		scanRequest, err := loadMessage(message.Body)
		if err != nil || scanRequest == nil {
			zap.L().Error("unable to load message from the queue",
				zap.Int("messageNumber", indx),
				zap.Error(err),
				zap.String("messageBody", message.Body),
			)
			continue
		}

		for _, entry := range scanRequest.Entries {
			_, logger := lambdalogger.ConfigureGlobal(ctx, map[string]interface{}{
				"sqsEntry": entry,
			})
			logger.Info(
				"starting poller",
				zap.Int("messageNumber", indx),
			)

			resources, err := pollers.Poll(entry)
			if err != nil {
				zap.L().Error("poller failed", zap.Error(err))
			}

			// Send data to the Resources API
			if resources != nil {
				logger.Info(
					"total resources generated",
					zap.Int("messageNumber", indx),
					zap.Int("numResources", len(resources)),
					zap.String("integrationType", "aws"),
				)

				for _, batch := range batchResources(resources) {
					params := &operations.AddResourcesParams{
						Body:       &api.AddResources{Resources: batch},
						HTTPClient: httpClient,
					}
					zap.L().Debug("adding new resources", zap.Any("params.Body", params.Body))
					if _, err := apiClient.Operations.AddResources(params); err != nil {
						return err
					}
				}
			}

			logger.Info(
				"polling complete",
				zap.Int("messageNumber", indx),
				zap.String("integrationType", "aws"),
			)
		}
	}

	return nil
}
