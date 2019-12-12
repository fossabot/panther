package models

// AlertOutputItem is the output configuration stored in DynamoDB.
type AlertOutputItem struct {

	// The user ID of the user that created the alert output
	CreatedBy *string `json:"createdBy"`

	// The time in epoch seconds when the alert output was created
	CreationTime *string `json:"creationTime"`

	// DisplayName is the user-provided name, e.g. "alert-channel".
	DisplayName *string `json:"displayName"`

	// EncryptedConfig is the encrypted JSON of the specific output details.
	EncryptedConfig []byte `json:"encryptedConfig"`

	// The user ID of the user that last modified the alert output last
	LastModifiedBy *string `json:"lastModifiedBy"`

	// The time in epoch seconds when the alert output was last modified
	LastModifiedTime *string `json:"lastModifiedTime"`

	// Identifies uniquely an alert output (table sort key)
	OutputID *string `json:"outputId"`

	// OutputType is the output class, e.g. "slack", "sns".
	// ("type" is a reserved Dynamo keyword, so we use "OutputType" instead)
	OutputType *string `json:"outputType"`

	// VerificationStatus is the current state of the output destination.
	// When an AlertOutput is not in 'VERIFIED' state it cannot be used to send notifications
	VerificationStatus *string `json:"verificationStatus"`
}

// DefaultOutputsItem is the default output configuration stored in DynamoDB.
type DefaultOutputsItem struct {

	//The severity of the (table sort key)
	Severity *string `json:"severity"`

	// Identifies uniquely an alert output
	OutputIDs []*string `json:"outputIds" dynamodbav:"outputIds,stringset"`
}
