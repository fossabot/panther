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
	"github.com/aws/aws-sdk-go/service/kms"
	jsoniter "github.com/json-iterator/go"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// DecryptConfig uses KMS to decrypt an output configuration.
func (key *Key) DecryptConfig(ciphertext []byte, config interface{}) error {
	response, err := key.client.Decrypt(&kms.DecryptInput{CiphertextBlob: ciphertext})
	if err != nil {
		return &genericapi.AWSError{Method: "kms.Decrypt", Err: err}
	}

	if err = jsoniter.Unmarshal(response.Plaintext, config); err != nil {
		return &genericapi.InternalError{
			Message: "failed to unmarshal config to json " + err.Error()}
	}
	return nil
}
