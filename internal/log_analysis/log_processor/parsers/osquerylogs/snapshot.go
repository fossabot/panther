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

var SnapshotDesc = `Snapshot contains all the data included in OsQuery differential logs
Reference: https://osquery.readthedocs.io/en/stable/deployment/logging/`

type Snapshot struct {
	Action         *string                `json:"action,omitempty" validate:"required,eq=snapshot"`
	CalendarTime   *timestamp.ANSICwithTZ `json:"calendarTime,omitempty" validate:"required"`
	Counter        *int                   `json:"counter,omitempty,string" validate:"required"`
	Decorations    map[string]string      `json:"decorations,omitempty"`
	Epoch          *int                   `json:"epoch,omitempty,string" validate:"required"`
	HostIdentifier *string                `json:"hostIdentifier,omitempty" validate:"required"`
	Name           *string                `json:"name,omitempty" validate:"required"`
	Snapshot       []map[string]string    `json:"snapshot,omitempty" validate:"required"`
	UnixTime       *int                   `json:"unixTime,omitempty,string" validate:"required"`
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
