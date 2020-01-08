package outputs

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

// AlertDeliveryError indicates whether a failed alert should be retried.
type AlertDeliveryError struct {
	// Message is the description of the problem: what went wrong.
	Message string

	// Permanent indicates whether the alert output should be retried.
	// For example, outputs which don't exist or errors creating the request are permanent failures.
	// But any error talking to the output itself can be retried by the Lambda function later.
	Permanent bool
}

func (e *AlertDeliveryError) Error() string { return e.Message }
