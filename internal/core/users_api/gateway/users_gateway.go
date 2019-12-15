package gateway

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	cfnIface "github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	fedIdentityProvider "github.com/aws/aws-sdk-go/service/cognitoidentity"
	fedIdentityProviderI "github.com/aws/aws-sdk-go/service/cognitoidentity/cognitoidentityiface"
	userPoolProvider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	userPoolProviderI "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

// API defines the interface for the user gateway which can be used for mocking.
type API interface {
	AddUserToGroup(id *string, groupName *string, userPoolID *string) error
	CreateUser(input *CreateUserInput) (*string, error)
	CreateUserPool(displayName *string) (*UserPool, error)
	CreateUserPoolGroups(userPoolID *string) error
	DeleteUser(id *string, userPoolID *string) error
	GetUser(id *string, userPoolID *string) (*models.User, error)
	ListGroups(userPoolID *string) ([]*models.Group, error)
	ListGroupsForUser(id *string, userPoolID *string) ([]*models.Group, error)
	ListUsers(limit *int64, paginationToken *string, userPoolID *string) (*ListUsersOutput, error)
	RemoveUserFromGroup(id *string, groupName *string, userPoolID *string) error
	ResetUserPassword(id *string, userPoolID *string) error
	ValidateToken(identityID *string, token *string) (map[string]interface{}, error)
	UpdateUser(input *UpdateUserInput) error
}

// UsersGateway encapsulates a service to Cognito Client.
type UsersGateway struct {
	userPoolClient        userPoolProviderI.CognitoIdentityProviderAPI
	fedIdentityClient     fedIdentityProviderI.CognitoIdentityAPI
	iamService            IAMService
	cloudFormationService cfnIface.CloudFormationAPI
}

// The UsersGateway must satisfy the API interface.
var _ API = (*UsersGateway)(nil)

// AppDomainURL is used to set up users email domain
var AppDomainURL = os.Getenv("APP_DOMAIN_URL")

// SESSourceEmailArn is used when Amazon Cognito emails the users with this address by calling Amazon SES on your behalf
var SESSourceEmailArn = os.Getenv("SES_SOURCE_EMAIL_ARN")

// AwsRegion is the region where the Lambda is deployed
var AwsRegion = os.Getenv("AWS_REGION")

// CustomMessageLambdaArn is used to handle user's custom message events such as forget password
var CustomMessageLambdaArn = os.Getenv("CUSTOM_MESSAGES_TRIGGER_HANDLER")

// IAMService is an interface for unit testing.  It must be satisfied by UsersGateway.iamService.
type IAMService interface {
	GetRole(*iam.GetRoleInput) (*iam.GetRoleOutput, error)
	UpdateAssumeRolePolicy(input *iam.UpdateAssumeRolePolicyInput) (*iam.UpdateAssumeRolePolicyOutput, error)
}

// New creates a new CognitoIdentityProvider client which talks to the given user pool.
func New(sess *session.Session) *UsersGateway {
	return &UsersGateway{
		userPoolClient:        userPoolProvider.New(sess),
		fedIdentityClient:     fedIdentityProvider.New(sess),
		iamService:            iam.New(sess),
		cloudFormationService: cloudformation.New(sess),
	}
}
