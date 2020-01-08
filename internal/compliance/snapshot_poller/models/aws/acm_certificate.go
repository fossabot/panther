package aws

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
	"time"

	"github.com/aws/aws-sdk-go/service/acm"
)

const (
	AcmCertificateSchema = "AWS.ACM.Certificate"
)

// AcmCertificate contains all the information about an ACM certificate
type AcmCertificate struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from acm.CertificateDetail
	CertificateAuthorityArn *string
	DomainName              *string
	DomainValidationOptions []*acm.DomainValidation
	ExtendedKeyUsages       []*acm.ExtendedKeyUsage
	FailureReason           *string
	InUseBy                 []*string
	IssuedAt                *time.Time
	Issuer                  *string
	KeyAlgorithm            *string
	KeyUsages               []*acm.KeyUsage
	NotAfter                *time.Time
	NotBefore               *time.Time
	Options                 *acm.CertificateOptions
	RenewalEligibility      *string
	RenewalSummary          *acm.RenewalSummary
	RevocationReason        *string
	RevokedAt               *time.Time
	Serial                  *string
	SignatureAlgorithm      *string
	Status                  *string
	Subject                 *string
	SubjectAlternativeNames []*string
	Type                    *string
}
