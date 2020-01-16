package gateway

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
