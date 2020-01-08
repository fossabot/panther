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

type mockDecryptClient struct {
	kmsiface.KMSAPI
	returnPlaintext []byte
	err             bool
}

func (m *mockDecryptClient) Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	if m.err {
		return nil, errors.New("internal error")
	}
	return &kms.DecryptOutput{Plaintext: m.returnPlaintext}, nil
}

func TestDecryptConfigServiceError(t *testing.T) {
	key := &Key{client: &mockDecryptClient{err: true}}
	err := key.DecryptConfig([]byte("ciphertext"), nil)
	assert.NotNil(t, err.(*genericapi.AWSError))
}

func TestDecryptConfigUnmarshalError(t *testing.T) {
	type detail struct {
		Name string `json:"name"`
	}
	key := &Key{client: &mockDecryptClient{returnPlaintext: []byte("access-token")}}
	err := key.DecryptConfig([]byte("ciphertext"), &detail{})
	assert.NotNil(t, err.(*genericapi.InternalError))
}

func TestDecryptConfig(t *testing.T) {
	type detail struct {
		Name string `json:"name"`
	}
	key := &Key{client: &mockDecryptClient{returnPlaintext: []byte("{\"name\": \"panther\"}")}}
	var output detail
	assert.Nil(t, key.DecryptConfig([]byte("ciphertext"), &output))
	assert.Equal(t, detail{Name: "panther"}, output)
}
