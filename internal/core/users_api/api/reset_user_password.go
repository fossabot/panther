package api

import (
	"github.com/panther-labs/panther/api/lambda/users/models"
)

// ResetUserPassword resets a user password.
func (API) ResetUserPassword(input *models.ResetUserPasswordInput) error {
	return userGateway.ResetUserPassword(input.ID, input.UserPoolID)
}
