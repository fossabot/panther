package models

// LambdaInput is the invocation event expected by the Lambda function.
//
// Exactly one action must be specified.
type LambdaInput struct {
	AddUserToOrganization     *AddUserToOrganizationInput     `json:"addUserToOrganization"`
	CreateUserInfrastructure  *CreateUserInfrastructureInput  `json:"createUserInfrastructure"`
	GetUser                   *GetUserInput                   `json:"getUser"`
	GetUserOrganizationAccess *GetUserOrganizationAccessInput `json:"getUserOrganizationAccess"`
	InviteUser                *InviteUserInput                `json:"inviteUser"`
	ListRoles                 *ListRolesInput                 `json:"listRoles"`
	ListUsers                 *ListUsersInput                 `json:"listUsers"`
	RemoveUser                *RemoveUserInput                `json:"removeUser"`
	ResetUserPassword         *ResetUserPasswordInput         `json:"resetUserPassword"`
	UpdateUser                *UpdateUserInput                `json:"updateUser"`
	ValidateCredentials       *ValidateCredentialsInput       `json:"validateCredentials"`
}

// AddUserToOrganizationInput adds a user to organization mapping
type AddUserToOrganizationInput struct {
	Email *string `json:"email" validate:"required,email"`
}

// AddUserToOrganizationOutput returns the user email and organization ID
type AddUserToOrganizationOutput struct {
	Email *string
}

// CreateUserInfrastructureInput creates Cognito infrastructure for a new user and organization.
type CreateUserInfrastructureInput struct {
	DisplayName *string `json:"displayName" validate:"required,min=1"`
	GivenName   *string `json:"givenName" validate:"required,min=1"`
	FamilyName  *string `json:"familyName" validate:"required,min=1"`
	Email       *string `json:"email" validate:"required,email"`
}

// CreateUserInfrastructureOutput returns the Panther user and user pool details.
type CreateUserInfrastructureOutput struct {
	User           *User   `json:"user"`
	UserPoolID     *string `json:"userPoolId"`
	AppClientID    *string `json:"appClientId"`
	IdentityPoolID *string `json:"identityPoolId"`
}

// GetUserInput retrieves a user's information based on id.
type GetUserInput struct {
	ID         *string `json:"id" validate:"required,uuid4"`
	UserPoolID *string `json:"userPoolId" validate:"required,min=1"`
}

// GetUserOutput returns the Panther user details.
type GetUserOutput = User

// GetUserOrganizationAccessInput retrieves a user's organization id based on email.
type GetUserOrganizationAccessInput struct {
	Email *string `json:"email" validate:"required,email"`
}

// GetUserOrganizationAccessOutput retrieves a user's organization id based on email.
type GetUserOrganizationAccessOutput struct {
	UserPoolID     *string `json:"userPoolId"`
	AppClientID    *string `json:"appClientId"`
	IdentityPoolID *string `json:"identityPoolId"`
}

// InviteUserInput creates a new user with minimal permissions and sends them an invite.
type InviteUserInput struct {
	GivenName  *string `json:"givenName" validate:"required,min=1"`
	FamilyName *string `json:"familyName" validate:"required,min=1"`
	Email      *string `json:"email" validate:"required,email"`
	UserPoolID *string `json:"userPoolId" validate:"required,min=1"`
	Role       *string `json:"role" validate:"required,min=1"`
}

// InviteUserOutput returns the randomly generated user id.
type InviteUserOutput struct {
	ID *string `json:"id"`
}

// ListRolesInput lists all available Panther groups.
type ListRolesInput struct {
	UserPoolID *string `json:"userPoolId" validate:"required,min=1"`
}

// ListRolesOutput is returned by the Lambda function.
type ListRolesOutput struct {
	Roles []*Group `json:"roles"`
}

// ListUsersInput lists all users in Panther.
type ListUsersInput struct {
	UserPoolID      *string `json:"userPoolId" validate:"required,min=1"`
	Limit           *int64  `json:"limit" validate:"omitempty,min=1"`
	PaginationToken *string `json:"paginationToken" validate:"omitempty,min=1"`
}

// ListUsersOutput returns a page of users.
type ListUsersOutput struct {
	Users           []*User `json:"users"`
	PaginationToken *string `json:"paginationToken"`
}

// RemoveUserInput deletes a user.
type RemoveUserInput struct {
	ID         *string `json:"id" validate:"required,uuid4"`
	UserPoolID *string `json:"userPoolId" validate:"required,min=1"`
}

// ResetUserPasswordInput resets the password for a user.
type ResetUserPasswordInput struct {
	ID         *string `json:"id" validate:"required,uuid4"`
	UserPoolID *string `json:"userPoolId" validate:"required,min=1"`
}

// UpdateUserInput updates user details.
type UpdateUserInput struct {
	ID         *string `json:"id" validate:"required,uuid4"`
	UserPoolID *string `json:"userPoolId" validate:"required,min=1"`

	// At least one of the following must be specified:
	GivenName   *string `json:"givenName" validate:"omitempty,min=1"`
	FamilyName  *string `json:"familyName" validate:"omitempty,min=1"`
	Email       *string `json:"email" validate:"omitempty,min=1"`
	PhoneNumber *string `json:"phoneNumber" validate:"omitempty,min=1"`
	Role        *string `json:"role" validate:"omitempty,min=1"`
}

// ValidateCredentialsInput validates the identities token of a user
type ValidateCredentialsInput struct {
	IdentityPoolID *string `json:"identityPoolId" validate:"required"`
	IdentityID     *string `json:"identityId" validate:"required"`
	JWT            *string `json:"jwt" validate:"required"`
}

// ValidateCredentialsOutput is returned by the lambda function.
// The identity value is the access token jwt value the user received after authentication
// nolint: lll
// example https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html#user-pool-access-token-payload
type ValidateCredentialsOutput struct {
	Identity map[string]interface{} `json:"identity"`
}
