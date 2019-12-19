package encryption

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
