package osquerylogs

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
)

// Batch contains all the data included in OsQuery batch logs
// Reference : https://osquery.readthedocs.io/en/stable/deployment/logging/
type Batch struct {
	CalendarTime *Time             `json:"calendarTime,omitempty" validate:"required"`
	Counter      *int              `json:"counter,omitempty,string"  validate:"required"`
	Decorations  map[string]string `json:"decorations,omitempty"`
	DiffResults  *BatchDiffResults `json:"diffResults,omitempty" validate:"required"`
	Epoch        *int              `json:"epoch,omitempty,string"  validate:"required"`
	Hostname     *string           `json:"hostname,omitempty"  validate:"required"`
	Name         *string           `json:"name,omitempty"  validate:"required"`
	UnixTime     *int              `json:"unixTime,omitempty,string"  validate:"required"`
}

// OsqueryBatchDiffResults contains diff data for OsQuery batch results
type BatchDiffResults struct {
	Added   []map[string]string `json:"added,omitempty"`
	Removed []map[string]string `json:"removed,omitempty"`
}

// BatchParser parses OsQuery Batch logs
type BatchParser struct{}

// Parse returns the parsed events or nil if parsing failed
func (p *BatchParser) Parse(log string) []interface{} {
	event := &Batch{}
	err := jsoniter.UnmarshalFromString(log, event)
	if err != nil {
		zap.L().Debug("failed to unmarshal log", zap.Error(err))
		return nil
	}

	tsa, _ := jsoniter.MarshalToString(event)
	fmt.Println(tsa)

	if err := parsers.Validator.Struct(event); err != nil {
		zap.L().Debug("failed to validate log", zap.Error(err))
		return nil
	}
	return []interface{}{event}
}

// LogType returns the log type supported by this parser
func (p *BatchParser) LogType() string {
	return "Osquery.Batch"
}
