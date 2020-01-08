package models

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

import "time"

// LambdaInput is the request structure for the alerts-api Lambda function.
type LambdaInput struct {
	GetAlert   *GetAlertInput   `json:"getAlert"`
	GetEvent   *GetEventInput   `json:"getEvent"`
	ListAlerts *ListAlertsInput `json:"listAlerts"`
}

// GetAlertInput retrieves details for a single alert.
//
// The response will contain by definition all of the events associated with the alert.
// If `eventPageSize` and `eventPage` are specified, it will returns only the specified events in the response.
// Example:
// {
//     "getAlert": {
// 	    "alertId": "ruleId-2"
//     }
// }
type GetAlertInput struct {
	AlertID                 *string `json:"alertId" validate:"required"`
	EventsPageSize          *int    `json:"eventsPageSize,omitempty"  validate:"omitempty,min=1,max=50"`
	EventsExclusiveStartKey *string `json:"eventsExclusiveStartKey,omitempty"`
}

// GetAlertOutput retrieves details for a single alert.
//
// Example:
// {
//     "getAlert": {
// 	    "alertId": "ruleId-2"
//     }
// }
type GetAlertOutput = Alert

// GetEventInput retrieves a single event
type GetEventInput struct {
	EventID *string `json:"eventId" validate:"required"`
}

// GetEventOutput contains the event
type GetEventOutput struct {
	Event *string `json:"event"`
}

// ListAlertsInput lists the alerts in reverse-chronological order (newest to oldest)
// If "ruleId" is not set, we return all the alerts for the organization
// If the "exclusiveStartKey" is not set, we return alerts starting from the most recent one. If it is set,
// the output will return alerts starting from the "exclusiveStartKey" exclusive.
//
//
// {
//     "listAlerts": {
//         "ruleId": "My.Rule",
//         "pageSize": 25
//     }
// }
type ListAlertsInput struct {
	RuleID            *string `json:"ruleId,omitempty"`
	PageSize          *int    `json:"pageSize,omitempty"  validate:"omitempty,min=1,max=50"`
	ExclusiveStartKey *string `json:"exclusiveStartKey,omitempty"`
}

// ListAlertsOutput is the returned alert list.
type ListAlertsOutput struct {
	// Alerts is a list of alerts sorted by timestamp descending.
	// Alerts with the same timestamp are returned in ascending order of alert ID.
	Alerts []*AlertSummary `json:"alertSummaries"`
	// LastEvaluatedKey contains the last evaluated alert Id.
	// If it is populated it means there are more alerts available
	// If it is nil, it means there are no more alerts to be returned.
	LastEvaluatedKey *string `json:"lastEvaluatedKey,omitempty"`
}

// AlertSummary contains summary information for an alert
type AlertSummary struct {
	AlertID          *string    `json:"alertId"`
	RuleID           *string    `json:"ruleId"`
	CreationTime     *time.Time `json:"creationTime"`
	LastEventMatched *time.Time `json:"lastEventMatched"`
	EventsMatched    *int       `json:"eventsMatched"`
	Severity         *string    `json:"severity"`
}

// Alert contains the details of an alert
type Alert struct {
	AlertID                *string    `json:"alertId"`
	RuleID                 *string    `json:"ruleId"`
	CreationTime           *time.Time `json:"creationTime"`
	LastEventMatched       *time.Time `json:"lastEventMatched"`
	MatchedEventNum        *int       `json:"matchedEventNum"`
	Events                 []*string  `json:"events"`
	EventsLastEvaluatedKey *string    `json:"eventsLastEvaluatedKey,omitempty"`
}
