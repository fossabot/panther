package aws

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
