// Package encryption handles all KMS operations.
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
