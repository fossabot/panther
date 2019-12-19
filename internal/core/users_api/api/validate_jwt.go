package api

import (
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// ValidateCredentials validate the JWT token and returns the claims
func (API) ValidateCredentials(input *models.ValidateCredentialsInput) (*models.ValidateCredentialsOutput, error) {
	jwtc, err := userGateway.ValidateToken(input.IdentityID, input.JWT)
	if err != nil {
		zap.L().Error("error parsing JWT token", zap.Error(err))
		return nil, err
	}
	// Check if the credentials from the token matches with what we got in the database
	o, err := GetOrganizations()
	// if of these checks below error, we ought to get paged because someone is messing with us
	if o == nil {
		zap.L().Error("organization not found", zap.Error(err))
		return nil, &genericapi.InvalidInputError{Message: "organization not found"}
	}
	if *o.IdentityPoolID != *input.IdentityPoolID {
		zap.L().Error("identity Pool id to validate: " +
			*input.IdentityPoolID + " and identity pool id stored: " +
			*o.IdentityPoolID + " not matching")
		return nil, &genericapi.InvalidInputError{Message: "identity Pool Invalid"}
	}
	return &models.ValidateCredentialsOutput{
		Identity: jwtc,
	}, nil
}
