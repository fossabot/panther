// Package encryption handles all KMS operations.
package encryption

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

// API defines the interface which can be used for mocking.
type API interface {
	DecryptConfig([]byte, interface{}) error
	EncryptConfig(interface{}) ([]byte, error)
}

// Key encapsulates a connection to the KMS encryption key.
type Key struct {
	ID     *string
	client kmsiface.KMSAPI
}

// The EncryptionKey must satisfy the API interface.
var _ API = (*Key)(nil)

// New creates AWS clients to interface with the encryption key.
func New(ID string, sess *session.Session) *Key {
	return &Key{ID: aws.String(ID), client: kms.New(sess)}
}
