package osquerylogs

import (
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
)

// Differential contains all the data included in OsQuery differential logs
// Reference: https://osquery.readthedocs.io/en/stable/deployment/logging/
type Differential struct {
	Action               *string           `json:"action,omitempty" validate:"required"`
	CalendarTime         *Time             `json:"calendartime,omitempty" validate:"required"`
	Columns              map[string]string `json:"columns,omitempty" validate:"required"`
	Counter              *int              `json:"counter,omitempty,string"`
	Decorations          map[string]string `json:"decorations,omitempty"`
	Epoch                *int              `json:"epoch,omitempty,string" validate:"required"`
	HostIdentifier       *string           `json:"hostIdentifier,omitempty" validate:"required"`
	LogType              *string           `json:"logType,omitempty" validate:"required,eq=result"`
	LogUnderscoreType    *string           `json:"log_type,omitempty"`
	Name                 *string           `json:"name,omitempty" validate:"required"`
	UnixTime             *int              `json:"unixTime,omitempty,string" validate:"required"`
	LogNumericsAsNumbers *bool             `json:"logNumericsAsNumbers,omitempty,string"`
}

// DifferentialParser parses OsQuery Differential logs
type DifferentialParser struct{}

// Parse returns the parsed events or nil if parsing failed
func (p *DifferentialParser) Parse(log string) []interface{} {
	event := &Differential{}
	err := jsoniter.UnmarshalFromString(log, event)
	if err != nil {
		zap.L().Debug("failed to unmarshal log", zap.Error(err))
		return nil
	}

	// Populating LogType with LogTypeInput value
	// This is needed because we want the JSON field with key `log_type` to be marshalled
	// with key `logtype`
	event.LogType = event.LogUnderscoreType
	event.LogUnderscoreType = nil

	if err := parsers.Validator.Struct(event); err != nil {
		zap.L().Debug("failed to validate log", zap.Error(err))
		return nil
	}

	return []interface{}{event}
}

// LogType returns the log type supported by this parser
func (p *DifferentialParser) LogType() string {
	return "Osquery.Differential"
}
