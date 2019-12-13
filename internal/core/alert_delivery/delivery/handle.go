package delivery

import (
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

func mustParseInt(text string) int {
	val, err := strconv.Atoi(text)
	if err != nil {
		panic(err)
	}
	return val
}

func getMaxRetryDuration() time.Duration {
	return time.Duration(mustParseInt(os.Getenv("ALERT_RETRY_DURATION_MINS"))) * time.Minute
}

// HandleAlerts sends each alert to its outputs and puts failed alerts back on the queue to retry.
func HandleAlerts(alerts []*models.Alert) {
	var failedAlerts []*models.Alert

	zap.L().Info("starting processing alerts", zap.Int("alerts", len(alerts)))

	for _, alert := range alerts {
		if !dispatch(alert) {
			if time.Since(*alert.CreatedAt) > getMaxRetryDuration() {
				zap.L().Error(
					"alert delivery permanently failed, exceeded max retry duration",
					zap.Strings("failedOutputs", aws.StringValueSlice(alert.OutputIDs)),
					zap.Time("alertCreatedAt", *alert.CreatedAt),
					zap.String("policyId", *alert.PolicyID),
					zap.String("severity", *alert.Severity),
				)
			} else {
				zap.L().Warn("will retry delivery of alert",
					zap.String("policyId", *alert.PolicyID),
					zap.String("severity", *alert.Severity),
				)
				failedAlerts = append(failedAlerts, alert)
			}
		}
	}

	if len(failedAlerts) > 0 {
		retry(failedAlerts)
	}
}