package api

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// AddOutput encrypts the output configuration and stores it to Dynamo.
func (API) AddOutput(input *models.AddOutputInput) (*models.AddOutputOutput, error) {
	item, err := outputsTable.GetOutputByName(input.DisplayName)
	if err != nil {
		return nil, err
	}

	if item != nil {
		return nil, &genericapi.AlreadyExistsError{
			Message: "A destination with the name" + *input.DisplayName + " already exists, please choose another display name"}
	}

	outputType, err := getOutputType(input.OutputConfig)
	if err != nil {
		return nil, &genericapi.InvalidInputError{Message: err.Error()}
	}

	alertOutput := &models.AlertOutput{
		OutputID:           aws.String(uuid.New().String()),
		DisplayName:        input.DisplayName,
		CreatedBy:          input.UserID,
		CreationTime:       aws.String(time.Now().Format(time.RFC3339)),
		LastModifiedBy:     input.UserID,
		LastModifiedTime:   aws.String(time.Now().Format(time.RFC3339)),
		OutputType:         outputType,
		OutputConfig:       input.OutputConfig,
		DefaultForSeverity: input.DefaultForSeverity,
	}

	status, err := outputVerification.GetVerificationStatus(alertOutput)
	if err != nil {
		return nil, err
	}
	alertOutput.VerificationStatus = status

	if *status != models.VerificationStatusSuccess {
		alertOutput, err = outputVerification.VerifyOutput(alertOutput)
		if err != nil {
			return nil, err
		}
	}

	alertOutputItem, err := AlertOutputToItem(alertOutput)
	if err != nil {
		return nil, err
	}

	if err = outputsTable.PutOutput(alertOutputItem); err != nil {
		return nil, err
	}

	zap.L().Info("stored new alert output",
		zap.String("outputId", *alertOutput.OutputID))

	if err = addToDefaults(input.DefaultForSeverity, alertOutput.OutputID); err != nil {
		return nil, err
	}
	return alertOutput, nil
}

func addToDefaults(severities []*string, outputID *string) error {
	for _, severity := range severities {
		defaults, err := defaultsTable.GetDefault(severity)
		if err != nil {
			return err
		}

		if defaults == nil {
			defaults = &models.DefaultOutputsItem{
				Severity:  severity,
				OutputIDs: []*string{outputID},
			}
		} else {
			defaults.OutputIDs = append(defaults.OutputIDs, outputID)
		}

		if err = defaultsTable.PutDefaults(defaults); err != nil {
			return err
		}
	}
	return nil
}
