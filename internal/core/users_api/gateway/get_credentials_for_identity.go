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
