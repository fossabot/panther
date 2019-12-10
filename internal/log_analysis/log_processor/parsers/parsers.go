package parsers

import "gopkg.in/go-playground/validator.v9"

// LogParser represents a parser for a supported log type
type LogParser interface {
	// LogType returns the log type supported by this parser
	LogType() string

	// Parse attempts to parse the provided log line
	// If the provided log is not of the supported type the method returns nil
	Parse(log string) []interface{}
}

// Validator can be used to validate schemas of log fields
var Validator = validator.New()
