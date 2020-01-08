package osquerylogs

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/timestamp"
)

var StatusDesc = `Status is a diagnostic osquery log about the daemon.
Reference: https://osquery.readthedocs.io/en/stable/deployment/logging/`

type Status struct {
	CalendarTime      *timestamp.ANSICwithTZ `json:"calendarTime,omitempty" validate:"required"`
	Decorations       map[string]string      `json:"decorations,omitempty"`
	Filename          *string                `json:"filename,omitempty" validate:"required"`
	HostIdentifier    *string                `json:"hostIdentifier,omitempty" validate:"required"`
	Line              *int                   `json:"line,omitempty,string" validate:"required"`
	LogType           *string                `json:"logType,omitempty" validate:"required,eq=status"`
	LogUnderscoreType *string                `json:"log_type,omitempty"`
	Message           *string                `json:"message,omitempty"`
	Severity          *int                   `json:"severity,omitempty,string" validate:"required"`
	UnixTime          *int                   `json:"unixTime,omitempty,string" validate:"required"`
	Version           *string                `json:"version,omitempty" validate:"required"`
}

// StatusParser parses OsQuery Status logs
type StatusParser struct{}

// Parse returns the parsed events or nil if parsing failed
func (p *StatusParser) Parse(log string) []interface{} {
	event := &Status{}
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
func (p *StatusParser) LogType() string {
	return "Osquery.Status"
}
