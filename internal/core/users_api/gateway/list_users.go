package gateway

import (
	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/panther-labs/panther/pkg/genericapi"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

// ListUsersOutput is output type for ListUsers
type ListUsersOutput struct {
	Users           []*models.User
	PaginationToken *string
}

func mapCognitoUserTypeToUser(u *provider.UserType) *models.User {
	user := models.User{
		CreatedAt: aws.Int64(u.UserCreateDate.Unix()),
		ID:        u.Username,
		Status:    u.UserStatus,
	}

	for _, attribute := range u.Attributes {
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

// ListUsers calls cognito api to list users that belongs to a user pool
func (g *UsersGateway) ListUsers(limit *int64, paginationToken *string, userPoolID *string) (*ListUsersOutput, error) {
	usersOutput, err := g.userPoolClient.ListUsers(&provider.ListUsersInput{
		Limit:           limit,
		PaginationToken: paginationToken,
		UserPoolId:      userPoolID,
	})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "cognito.ListUsers", Err: err}
	}

	users := make([]*models.User, len(usersOutput.Users))
	for i, uo := range usersOutput.Users {
		users[i] = mapCognitoUserTypeToUser(uo)
	}
	return &ListUsersOutput{Users: users, PaginationToken: usersOutput.PaginationToken}, nil
}
