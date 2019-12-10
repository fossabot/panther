package osquerylogs

import (
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
)

// Snapshot contains all the data included in OsQuery differential logs
// Reference: https://osquery.readthedocs.io/en/stable/deployment/logging/
type Snapshot struct {
	Action         *string             `json:"action,omitempty" validate:"required,eq=snapshot"`
	CalendarTime   *Time               `json:"calendarTime,omitempty" validate:"required"`
	Counter        *int                `json:"counter,omitempty,string" validate:"required"`
	Decorations    map[string]string   `json:"decorations,omitempty"`
	Epoch          *int                `json:"epoch,omitempty,string" validate:"required"`
	HostIdentifier *string             `json:"hostIdentifier,omitempty" validate:"required"`
	Name           *string             `json:"name,omitempty" validate:"required"`
	Snapshot       []map[string]string `json:"snapshot,omitempty" validate:"required"`
	UnixTime       *int                `json:"unixTime,omitempty,string" validate:"required"`
}

// SnapshotParser parses OsQuery snapshot logs
type SnapshotParser struct{}

// Parse returns the parsed events or nil if parsing failed
func (p *SnapshotParser) Parse(log string) []interface{} {
	event := &Snapshot{}
	err := jsoniter.UnmarshalFromString(log, event)
	if err != nil {
		zap.L().Debug("failed to unmarshal log", zap.Error(err))
		return nil
	}

	if err := parsers.Validator.Struct(event); err != nil {
		zap.L().Debug("failed to validate log", zap.Error(err))
		return nil
	}
	return []interface{}{event}
}

// LogType returns the log type supported by this parser
func (p *SnapshotParser) LogType() string {
	return "Osquery.Snapshot"
}
