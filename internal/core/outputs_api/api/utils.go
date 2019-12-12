package api

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

// AlertOutputToItem converts an AlertOutput to an AlertOutputItem
func AlertOutputToItem(input *models.AlertOutput) (*models.AlertOutputItem, error) {
	item := &models.AlertOutputItem{
		CreatedBy:          input.CreatedBy,
		CreationTime:       input.CreationTime,
		DisplayName:        input.DisplayName,
		LastModifiedBy:     input.LastModifiedBy,
		LastModifiedTime:   input.LastModifiedTime,
		OutputID:           input.OutputID,
		OutputType:         input.OutputType,
		VerificationStatus: input.VerificationStatus,
	}

	encryptedConfig, err := encryptionKey.EncryptConfig(input.OutputConfig)

	if err != nil {
		return nil, err
	}

	item.EncryptedConfig = encryptedConfig

	return item, nil
}

// ItemToAlertOutput converts an AlertOutputItem to an AlertOutput
func ItemToAlertOutput(input *models.AlertOutputItem) (alertOutput *models.AlertOutput, err error) {
	alertOutput = &models.AlertOutput{
		CreatedBy:          input.CreatedBy,
		CreationTime:       input.CreationTime,
		DisplayName:        input.DisplayName,
		LastModifiedBy:     input.LastModifiedBy,
		LastModifiedTime:   input.LastModifiedTime,
		OutputID:           input.OutputID,
		OutputType:         input.OutputType,
		VerificationStatus: input.VerificationStatus,
	}

	alertOutput.OutputConfig = &models.OutputConfig{}
	err = encryptionKey.DecryptConfig(input.EncryptedConfig, alertOutput.OutputConfig)

	if err != nil {
		return nil, err
	}

	return alertOutput, nil
}

func getOutputType(outputConfig *models.OutputConfig) (*string, error) {
	if outputConfig.Slack != nil {
		return aws.String("slack"), nil
	}
	if outputConfig.PagerDuty != nil {
		return aws.String("pagerduty"), nil
	}
	if outputConfig.Email != nil {
		return aws.String("email"), nil
	}
	if outputConfig.Github != nil {
		return aws.String("github"), nil
	}
	if outputConfig.Jira != nil {
		return aws.String("jira"), nil
	}
	if outputConfig.Opsgenie != nil {
		return aws.String("opsgenie"), nil
	}
	if outputConfig.MsTeams != nil {
		return aws.String("msteams"), nil
	}
	if outputConfig.Sns != nil {
		return aws.String("sns"), nil
	}
	if outputConfig.Sqs != nil {
		return aws.String("sqs"), nil
	}

	return nil, errors.New("no valid output configuration specified for alert output")
}
