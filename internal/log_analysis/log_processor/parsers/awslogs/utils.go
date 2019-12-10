package awslogs

import (
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
)

func csvStringToPointer(value string) *string {
	if value == "-" {
		return nil
	}
	return aws.String(value)
}

func csvStringToIntPointer(value string) *int {
	if value == "-" {
		return nil
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return aws.Int(result)
}

func csvStringToFloat64Pointer(value string) *float64 {
	if value == "-" {
		return nil
	}
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return aws.Float64(result)
}

func csvStringToArray(value string) []string {
	if value == "-" {
		return []string{}
	}

	return strings.Split(value, ",")
}
