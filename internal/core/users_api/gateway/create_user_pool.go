package gateway

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
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	identityProvider "github.com/aws/aws-sdk-go/service/cognitoidentity"
	userPoolProvider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/iam"
	jsoniter "github.com/json-iterator/go"
	"github.com/matcornic/hermes"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/core/users_api/email"
)

const (
	// IdentityPoolAuthenticatedAdminsRole is the role the identity pool assumed when the admin users is authenticated
	IdentityPoolAuthenticatedAdminsRole = "IdentityPoolAuthenticatedAdminsRole"
)

// UserPool contains user pool metadata
type UserPool struct {
	UserPoolID     *string
	AppClientID    *string
	IdentityPoolID *string
}

// CreateUserPool creates a new user pool with app client and MFA enabled
func (g *UsersGateway) CreateUserPool(displayName *string) (*UserPool, error) {
	userPoolOutput, err := createUserPool(g, displayName)
	if err != nil {
		zap.L().Error("error creating pool", zap.Error(err))
		return nil, err
	}
	userPoolID := userPoolOutput.UserPool.Id

	appClient, err := createAppClient(g, userPoolID)
	if err != nil {
		zap.L().Error("error creating user pool client", zap.Error(err))
		return nil, err
	}
	_, err = setMfaConfig(g, userPoolID)
	if err != nil {
		zap.L().Error("error setting MFA config", zap.Error(err))
		return nil, err
	}
	identityPool, err := createIdentityPool(g, userPoolID, appClient.UserPoolClient.ClientId, displayName)
	if err != nil {
		zap.L().Error("error creating identity pool", zap.Error(err))
		return nil, err
	}
	return &UserPool{
		UserPoolID:     userPoolID,
		AppClientID:    appClient.UserPoolClient.ClientId,
		IdentityPoolID: identityPool.IdentityPoolId,
	}, nil
}

func constructWelcomeEmail() (*string, error) {
	hemail := hermes.Email{
		Body: hermes.Body{
			Title: "Welcome to Panther",
			Intros: []string{
				"We're very excited to have you on board",
			},
			Dictionary: []hermes.Entry{
				{Key: "Username", Value: "{username}"},
				{Key: "Password", Value: "{####}"},
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started, please click here to visit our Panther dashboard and use your username and password above:",
					Button: hermes.Button{
						TextColor: "#FFFFFF",
						Color:     "#6967F4", // Optional action button color
						Text:      "Sign Me In",
						Link:      "https://" + AppDomainURL + "/sign-in",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just email us at support@runpanther.io, we'd love to help.",
			},
		},
	}
	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := email.PantherEmailTemplate.GeneratePlainText(hemail)
	if err != nil {
		zap.L().Error("failed to generate welcome email", zap.Error(err))
		return nil, err
	}
	// We have to do this because most email clients are not friendly with basic new line markup
	// replacing \n with a <br /> is the easiest way to mitigate this issue
	emailBody = strings.Replace(emailBody, "\n", "<br />", -1)
	return &emailBody, nil
}

func createUserPool(g *UsersGateway, displayName *string) (*userPoolProvider.CreateUserPoolOutput, error) {
	// Get IAM Role allowing emails to be sent via SNS
	role, err := g.iamService.GetRole(&iam.GetRoleInput{
		RoleName: aws.String("CognitoSNSRole"),
	})
	if err != nil {
		zap.L().Error("error getting role:", zap.Error(err))
		return nil, err
	}
	roleArn := role.Role.Arn

	emailMessageBody, err := constructWelcomeEmail()

	// If SESSourceEmailArn is passed in as a cloudformation parameter, use it to send cognito user email
	var emailConfig = &userPoolProvider.EmailConfigurationType{}
	if len(SESSourceEmailArn) != 0 {
		emailConfig = &userPoolProvider.EmailConfigurationType{
			EmailSendingAccount: aws.String(userPoolProvider.EmailSendingAccountTypeDeveloper),
			SourceArn:           aws.String(SESSourceEmailArn),
		}
	}

	if err != nil {
		zap.L().Error("failed to generate welcome email", zap.Error(err))
		return nil, err
	}

	userPoolInput := &userPoolProvider.CreateUserPoolInput{
		AdminCreateUserConfig: &userPoolProvider.AdminCreateUserConfigType{
			AllowAdminCreateUserOnly: aws.Bool(true),
			InviteMessageTemplate: &userPoolProvider.MessageTemplateType{
				EmailSubject: aws.String("Welcome to Panther!"),
				EmailMessage: emailMessageBody,
			},
		},
		EmailConfiguration:     emailConfig,
		AutoVerifiedAttributes: []*string{aws.String("email")},
		MfaConfiguration:       aws.String("ON"),
		LambdaConfig: &userPoolProvider.LambdaConfigType{
			CustomMessage: &CustomMessageLambdaArn,
		},
		Policies: &userPoolProvider.UserPoolPolicyType{
			PasswordPolicy: &userPoolProvider.PasswordPolicyType{
				MinimumLength:    aws.Int64(16),
				RequireLowercase: aws.Bool(true),
				RequireNumbers:   aws.Bool(true),
				RequireSymbols:   aws.Bool(true),
				RequireUppercase: aws.Bool(true),
			},
		},
		PoolName: aws.String(*displayName + " User Pool"),
		Schema: []*userPoolProvider.SchemaAttributeType{
			{
				AttributeDataType:      aws.String("String"),
				DeveloperOnlyAttribute: aws.Bool(false),
				Mutable:                aws.Bool(true),
				Name:                   aws.String("email"),
				Required:               aws.Bool(true),
			},
			{
				AttributeDataType:      aws.String("String"),
				DeveloperOnlyAttribute: aws.Bool(false),
				Mutable:                aws.Bool(true),
				Name:                   aws.String("given_name"),
				Required:               aws.Bool(true),
			},
			{
				AttributeDataType:      aws.String("String"),
				DeveloperOnlyAttribute: aws.Bool(false),
				Mutable:                aws.Bool(true),
				Name:                   aws.String("family_name"),
				Required:               aws.Bool(true),
			},
			{
				AttributeDataType:      aws.String("String"),
				DeveloperOnlyAttribute: aws.Bool(false),
				Mutable:                aws.Bool(true),
				Name:                   aws.String("email"),
				Required:               aws.Bool(true),
			},
		},
		SmsConfiguration: &userPoolProvider.SmsConfigurationType{
			SnsCallerArn: roleArn,
		},
		UsernameAttributes: []*string{aws.String("email")},
	}
	return g.userPoolClient.CreateUserPool(userPoolInput)
}

func createAppClient(g *UsersGateway, userPoolID *string) (*userPoolProvider.CreateUserPoolClientOutput, error) {
	return g.userPoolClient.CreateUserPoolClient(&userPoolProvider.CreateUserPoolClientInput{
		ClientName:           aws.String("Panther"),
		GenerateSecret:       aws.Bool(false),
		RefreshTokenValidity: aws.Int64(1),
		UserPoolId:           userPoolID,
		WriteAttributes: []*string{
			aws.String("email"),
			aws.String("given_name"),
			aws.String("family_name"),
		},
	})
}

func createIdentityPool(g *UsersGateway, userPoolID *string, appID *string, displayName *string) (*identityProvider.IdentityPool, error) {
	identityPool, err := g.fedIdentityClient.CreateIdentityPool(&identityProvider.CreateIdentityPoolInput{
		AllowUnauthenticatedIdentities: aws.Bool(false),
		CognitoIdentityProviders: []*identityProvider.Provider{
			{
				ClientId:             appID,
				ProviderName:         aws.String("cognito-idp." + AwsRegion + ".amazonaws.com/" + *userPoolID),
				ServerSideTokenCheck: aws.Bool(true),
			},
		},
		IdentityPoolName: displayName,
	})
	if err != nil {
		zap.L().Error("error creating identity pool:", zap.Error(err))
		return nil, err
	}
	firstUserPoolProvider := identityPool.CognitoIdentityProviders[0]
	poolProvider := *firstUserPoolProvider.ProviderName + ":" + *appID
	defaultRoleArn, err := getRoleArn(g, aws.String("DefaultIdentityPoolAuthenticatedRole"))
	if err != nil {
		zap.L().Error("error creating constructing role arn:", zap.Error(err))
		return nil, err
	}
	_, err = g.fedIdentityClient.SetIdentityPoolRoles(&identityProvider.SetIdentityPoolRolesInput{
		IdentityPoolId: identityPool.IdentityPoolId,
		RoleMappings: map[string]*identityProvider.RoleMapping{
			poolProvider: {
				AmbiguousRoleResolution: aws.String(identityProvider.AmbiguousRoleResolutionTypeDeny),
				Type:                    aws.String(identityProvider.RoleMappingTypeToken),
			},
		},
		Roles: map[string]*string{
			"authenticated":   defaultRoleArn,
			"unauthenticated": defaultRoleArn,
		},
	})
	if err != nil {
		zap.L().Error("error setting identity pool roles:", zap.Error(err))
		return nil, err
	}
	if err := updatePolicyCondition(g, aws.String(IdentityPoolAuthenticatedAdminsRole), identityPool.IdentityPoolId); err != nil {
		zap.L().Error("error setting trusted relationship for admin role", zap.Error(err))
		return nil, err
	}
	return identityPool, nil
}

func getRoleArn(g *UsersGateway, roleName *string) (*string, error) {
	role, err := g.iamService.GetRole(&iam.GetRoleInput{
		RoleName: roleName,
	})
	if err != nil {
		zap.L().Error("error getting role arn:", zap.Error(err))
		return nil, err
	}
	return role.Role.Arn, nil
}

// IAMStringEqualsConditions is a struct contains the Condition for aud attributes
type IAMStringEqualsConditions struct {
	AUD []*string `json:"cognito-identity.amazonaws.com:aud"`
}

// IAMStringLikeConditions is a struct contains the Condition for amr attributes
type IAMStringLikeConditions struct {
	AMR *string `json:"cognito-identity.amazonaws.com:amr"`
}

// IAMAssumeRolePolicyCondition is a struct contains the Condition for a trusted relationship
type IAMAssumeRolePolicyCondition struct {
	StringEquals *IAMStringEqualsConditions
	StringLike   *IAMStringLikeConditions `json:"ForAnyValue:StringLike"`
}

// IAMAssumeRolePolicyPrincipal is a struct contains a Federated service
type IAMAssumeRolePolicyPrincipal struct {
	Federated *string
}

// IAMAssumeRolePolicyStatement is a struct that represents the Statement for an IdentityGroup IAM Role
type IAMAssumeRolePolicyStatement struct {
	Effect    *string
	Principal *IAMAssumeRolePolicyPrincipal
	Action    *string
	Condition *IAMAssumeRolePolicyCondition
}

// IAMAssumeRolePolicyDocument is a struct that represents the Policy for an IdentityGroup IAM Role,
// mainly used for unmarshalling and editing
type IAMAssumeRolePolicyDocument struct {
	Version   *string
	Statement []*IAMAssumeRolePolicyStatement
}

// Append identityPoolId to the Condition for group IAM policy
func updatePolicyCondition(g *UsersGateway, roleName *string, identityPoolID *string) error {
	role, err := g.iamService.GetRole(&iam.GetRoleInput{
		RoleName: roleName,
	})
	if err != nil {
		zap.L().Error("error getting role arn", zap.Error(err))
		return err
	}
	pdS, err := url.QueryUnescape(*role.Role.AssumeRolePolicyDocument)
	if err != nil {
		zap.L().Error("error unescaping policy document", zap.Error(err))
		return err
	}
	var pd IAMAssumeRolePolicyDocument

	if err := jsoniter.UnmarshalFromString(pdS, &pd); err != nil {
		zap.L().Error("error unmarshalling policy document", zap.Error(err))
		return err
	}
	if pd.Statement[0].Condition == nil {
		pd.Statement[0].Condition = &IAMAssumeRolePolicyCondition{
			StringEquals: &IAMStringEqualsConditions{
				AUD: []*string{identityPoolID, identityPoolID}, // Slice gets converted to string if there's only 1 item
			},
			StringLike: &IAMStringLikeConditions{
				AMR: aws.String("authenticated"),
			},
		}
	} else {
		pd.Statement[0].Condition.StringEquals.AUD = append(pd.Statement[0].Condition.StringEquals.AUD, identityPoolID)
	}
	c, err := jsoniter.MarshalToString(&pd)
	if err != nil {
		zap.L().Error("error marshalling policy document", zap.Error(err))
		return err
	}
	if _, err = g.iamService.UpdateAssumeRolePolicy(&iam.UpdateAssumeRolePolicyInput{
		PolicyDocument: aws.String(c),
		RoleName:       roleName,
	}); err != nil {
		zap.L().Error("error updating policy document", zap.Error(err))
		return err
	}
	return nil
}

func setMfaConfig(g *UsersGateway, userPoolID *string) (*userPoolProvider.SetUserPoolMfaConfigOutput, error) {
	return g.userPoolClient.SetUserPoolMfaConfig(&userPoolProvider.SetUserPoolMfaConfigInput{
		MfaConfiguration: aws.String("ON"),
		SoftwareTokenMfaConfiguration: &userPoolProvider.SoftwareTokenMfaConfigType{
			Enabled: aws.Bool(true),
		},
		UserPoolId: userPoolID,
	})
}
