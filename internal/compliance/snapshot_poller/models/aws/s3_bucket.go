package aws

import "github.com/aws/aws-sdk-go/service/s3"

// S3BucketSchema is the name of the S3Bucket Schema
const S3BucketSchema = "AWS.S3.Bucket"

// S3Bucket contains all information about an S3 bucket.
type S3Bucket struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Additional fields
	EncryptionRules                []*s3.ServerSideEncryptionRule
	Grants                         []*s3.Grant
	LifecycleRules                 []*s3.LifecycleRule
	LoggingPolicy                  *s3.LoggingEnabled
	MFADelete                      *string
	ObjectLockConfiguration        *s3.ObjectLockConfiguration
	Owner                          *s3.Owner
	Policy                         *string
	PublicAccessBlockConfiguration *s3.PublicAccessBlockConfiguration
	Versioning                     *string
}
