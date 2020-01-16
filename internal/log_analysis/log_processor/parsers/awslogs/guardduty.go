package awslogs

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

var GuardDutyDesc = `Amazon GuardDuty is a threat detection service that continuously monitors for malicious activity 
and unauthorized behavior inside AWS Accounts. 
See also GuardDuty Finding Format : https://docs.aws.amazon.com/guardduty/latest/ug/guardduty_finding-format.html`

type GuardDuty struct {
	SchemaVersion *string            `json:"schemaVersion" validate:"required"`
	AccountID     *string            `json:"accountId" validate:"len=12,numeric"`
	Region        *string            `json:"region" validate:"required"`
	Partition     *string            `json:"partition" validate:"required"`
	ID            *string            `json:"id,omitempty" validate:"required"`
	Arn           *string            `json:"arn" validate:"required"`
	Type          *string            `json:"type" validate:"required"`
	Resource      interface{}        `json:"resource" validate:"required"`
	Severity      *int               `json:"severity" validate:"required,min=0"`
	CreatedAt     *timestamp.RFC3339 `json:"createdAt" validate:"required,min=0"`
	UpdatedAt     *timestamp.RFC3339 `json:"updatedAt" validate:"required,min=0"`
	Title         *string            `json:"title" validate:"required"`
	Description   *string            `json:"description" validate:"required"`
	Service       *GuardDutyService  `json:"service" validate:"required"`
}

type GuardDutyService struct {
	AdditionalInfo interface{}        `json:"additionalInfo"`
	Action         interface{}        `json:"action"`
	ServiceName    *string            `json:"serviceName" validate:"required"`
	DetectorID     *string            `json:"detectorId" validate:"required"`
	ResourceRole   *string            `json:"resourceRole"`
	EventFirstSeen *timestamp.RFC3339 `json:"eventFirstSeen"`
	EventLastSeen  *timestamp.RFC3339 `json:"eventLastSeen"`
	Archived       *bool              `json:"archived"`
	Count          *int               `json:"count"`
}

// VPCFlowParser parses AWS VPC Flow Parser logs
type GuardDutyParser struct{}

// Parse returns the parsed events or nil if parsing failed
func (p *GuardDutyParser) Parse(log string) []interface{} {
	event := &GuardDuty{}
	err := jsoniter.UnmarshalFromString(log, event)
	if err != nil {
		zap.L().Debug("failed to parse log", zap.Error(err))
		return nil
	}

	if err := parsers.Validator.Struct(event); err != nil {
		zap.L().Debug("failed to validate log", zap.Error(err))
		return nil
	}
	return []interface{}{event}
}

// LogType returns the log type supported by this parser
func (p *GuardDutyParser) LogType() string {
	return "AWS.GuardDuty"
}
