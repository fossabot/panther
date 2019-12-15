package gateway

import (
	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/panther-labs/panther/pkg/genericapi"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

func mapGetUserOutputToPantherUser(u *provider.AdminGetUserOutput) *models.User {
	user := models.User{
		CreatedAt: aws.Int64(u.UserCreateDate.Unix()),
		ID:        u.Username,
		Status:    u.UserStatus,
	}

	for _, attribute := range u.UserAttributes {
		switch *attribute.Name {
		case "email":
			user.Email = attribute.Value
		case "phone_number":
			user.PhoneNumber = attribute.Value
		case "given_name":
			user.GivenName = attribute.Value
		case "family_name":
			user.FamilyName = attribute.Value
		}
	}

	return &user
}

// GetUser calls cognito api to get user info
func (g *UsersGateway) GetUser(id *string, userPoolID *string) (*models.User, error) {
	user, err := g.userPoolClient.AdminGetUser(&provider.AdminGetUserInput{
		Username:   id,
		UserPoolId: userPoolID,
	})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "cognito.AdminGetUser", Err: err}
	}
	return mapGetUserOutputToPantherUser(user), nil
}
