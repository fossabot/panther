package utils

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"go.uber.org/zap"
)

// LogAWSError logs an AWS error to zap in a digestable format.
func LogAWSError(apiCall string, err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		zap.L().Error(
			apiCall,
			zap.String("errorCode", awsErr.Code()),
			zap.String("errorMessage", awsErr.Message()),
		)
	}
}
