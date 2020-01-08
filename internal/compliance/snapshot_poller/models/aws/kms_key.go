package aws

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
	"time"
)

const (
	KmsKeySchema = "AWS.KMS.Key"
)

// KmsKey contains all information about a kms key
type KmsKey struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from kms.KeyMetaData
	CloudHsmClusterId *string
	CustomKeyStoreId  *string
	DeletionDate      *time.Time
	Description       *string
	Enabled           *bool
	ExpirationModel   *string
	KeyManager        *string
	KeyState          *string
	KeyUsage          *string
	Origin            *string
	ValidTo           *time.Time

	// Additional fields
	KeyRotationEnabled *bool
	Policy             *string
}
