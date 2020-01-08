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

// Group is a struct for Panther Group containing employees.
type Group struct {
	Description *string `json:"description"`
	Name        *string `json:"name"`
}

// User is a struct describing a Panther User.
type User struct {
	CreatedAt   *int64  `json:"createdAt"`
	Email       *string `json:"email"`
	FamilyName  *string `json:"familyName"`
	GivenName   *string `json:"givenName"`
	ID          *string `json:"id"`
	PhoneNumber *string `json:"phoneNumber"`
	Role        *string `json:"role"` // Roles are group name
	Status      *string `json:"status"`
}
