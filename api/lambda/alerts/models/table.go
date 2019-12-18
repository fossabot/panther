package models

import "time"

// AlertItem is a DDB representation of an Alert
type AlertItem struct {
	AlertID          *string    `json:"alertId"`
	RuleID           *string    `json:"ruleId"`
	CreationTime     *time.Time `json:"creationTime"`
	LastEventMatched *time.Time `json:"lastEventMatched"`
	EventHashes      [][]byte   `json:"eventHashes"`
}
