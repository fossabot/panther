package aws

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
