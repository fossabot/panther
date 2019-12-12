package encryption

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/pkg/genericapi"
)

type mockEncryptClient struct {
	kmsiface.KMSAPI
	err bool
}

func (m *mockEncryptClient) Encrypt(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	if m.err {
		return nil, errors.New("internal error")
	}
	return &kms.EncryptOutput{CiphertextBlob: []byte("super secret")}, nil
}

func TestEncryptConfigMarshalError(t *testing.T) {
	key := &Key{client: &mockEncryptClient{}}
	result, err := key.EncryptConfig(key.EncryptConfig)
	assert.Nil(t, result)
	assert.NotNil(t, err.(*genericapi.InternalError))
}

func TestEncryptConfigServiceError(t *testing.T) {
	key := &Key{client: &mockEncryptClient{err: true}}
	result, err := key.EncryptConfig("access-token")
	assert.Nil(t, result)
	assert.NotNil(t, err.(*genericapi.AWSError))
}

func TestEncrypt(t *testing.T) {
	key := &Key{client: &mockEncryptClient{}}
	result, err := key.EncryptConfig("access-token")
	assert.NotNil(t, result)
	assert.Nil(t, err)
}
