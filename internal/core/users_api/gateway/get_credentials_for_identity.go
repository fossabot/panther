package gateway

import (
	"strings"

	"github.com/SermoDigital/jose/jws"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"go.uber.org/zap"
)

// ValidateToken calls cognito api and validate the identity id and the jwt token is valid
// returns the jwt claims
func (g *UsersGateway) ValidateToken(identityID *string, token *string) (map[string]interface{}, error) {
	jwt, err := jws.ParseJWT([]byte(*token))
	if err != nil {
		zap.L().Error("Error parsing JWT token", zap.Error(err))
		return nil, err
	}
	// To figure out what's available in the claims, copy the jwt token
	// in the request header on Panther app and check it out on https://jwt.io/
	jwtc := jwt.Claims()
	iss := jwtc["iss"].(string)
	iss = strings.TrimPrefix(iss, "https://")

	// Validate that token is valid
	_, err = g.fedIdentityClient.GetCredentialsForIdentity(&cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: identityID,
		Logins: map[string]*string{
			iss: token,
		},
	})
	if err != nil {
		zap.L().Error("Invalid token", zap.Error(err))
		return nil, err
	}
	return jwtc, nil
}
