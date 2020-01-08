package merger

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
	"crypto/sha1" // nolint: gosec
	"time"
)

// AlertNotification models a notification sent to Alert merger
type AlertNotification struct {
	RuleID        *string    `json:"ruleId"`
	RuleVersionID *string    `json:"ruleVersionId"`
	Event         *string    `json:"event"`
	Timestamp     *time.Time `json:"timestamp"`
}

// MatchedEvent represents an event matched by the Panther rule engine
type MatchedEvent struct {
	EventHash [sha1.Size]byte `json:"eventHash"`
	Timestamp *time.Time      `json:"timestamp"`
	Event     *string         `json:"event"`
}
