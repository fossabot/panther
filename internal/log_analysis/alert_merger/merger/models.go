package merger

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
