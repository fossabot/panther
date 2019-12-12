package encryption

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/kms"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// EncryptConfig uses KMS to encrypt an output configuration.
func (key *Key) EncryptConfig(config interface{}) ([]byte, error) {
	body, err := json.Marshal(config)
	if err != nil {
		return nil, &genericapi.InternalError{
			Message: "failed to marshal config to json " + err.Error()}
	}

	response, err := key.client.Encrypt(&kms.EncryptInput{KeyId: key.ID, Plaintext: body})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "kms.Encrypt", Err: err}
	}

	return response.CiphertextBlob, nil
}
