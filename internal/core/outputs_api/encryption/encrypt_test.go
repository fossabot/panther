package encryption

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
