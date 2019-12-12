package utils

import "strings"

type ParsedResourceID struct {
	AccountID string
	Region    string
	Schema    string
}

// GenerateResourceID returns a formatted custom Resource ID.
func GenerateResourceID(awsAccountID string, region string, schema string) string {
	return strings.Join([]string{awsAccountID, region, schema}, ":")
}

func ParseResourceID(resourceID string) *ParsedResourceID {
	parsedResourceID := strings.Split(resourceID, ":")
	if len(parsedResourceID) != 3 {
		return nil
	}
	return &ParsedResourceID{
		AccountID: parsedResourceID[0],
		Region:    parsedResourceID[1],
		Schema:    parsedResourceID[2],
	}
}
