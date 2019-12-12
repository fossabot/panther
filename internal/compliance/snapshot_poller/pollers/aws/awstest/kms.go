package awstest

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/stretchr/testify/mock"
)

// Example KMS API return values
var (
	ExampleKeyId = aws.String("188c57ed-b28a-4c0e-9821-f4940d15cb0a")

	ExampleListKeysOutput = &kms.ListKeysOutput{
		Keys: []*kms.KeyListEntry{
			{
				KeyArn: aws.String("arn:aws:kms:us-west-2:857418155548:key/188c57ed-b28a-4c0e-9821-f4940d15cb0a"),
				KeyId:  aws.String("188c57ed-b28a-4c0e-9821-f4940d15cb0a"),
			},
			{
				KeyArn: aws.String("arn:aws:kms:us-west-2:857418155548:key/d15a1e37-3ef7-4882-9be5-ef3a024114db"),
				KeyId:  aws.String("d15a1e37-3ef7-4882-9be5-ef3a024114db"),
			},
		},
	}

	ExampleGetKeyRotationStatusOutput = &kms.GetKeyRotationStatusOutput{
		KeyRotationEnabled: aws.Bool(true),
	}

	ExampleDescribeKeyOutput = &kms.DescribeKeyOutput{
		KeyMetadata: &kms.KeyMetadata{
			AWSAccountId: aws.String("857418155548"),
			Arn:          aws.String("arn:aws:kms:us-west-2:857418155548:key/188c57ed-b28a-4c0e-9821-f4940d15cb0a"),
			CreationDate: ExampleDate,
			Description:  aws.String("Encryption key for panther-snapshot-queue data"),
			Enabled:      aws.Bool(true),
			KeyId:        aws.String("188c57ed-b28a-4c0e-9821-f4940d15cb0a"),
			KeyManager:   aws.String("CUSTOMER"),
			KeyState:     aws.String("Enabled"),
			KeyUsage:     aws.String("ENCRYPT_DECRYPT"),
			Origin:       aws.String("AWS_KMS"),
		},
	}

	ExampleDescribeKeyOutputAWSManaged = &kms.DescribeKeyOutput{
		KeyMetadata: &kms.KeyMetadata{
			AWSAccountId: aws.String("123456789012"),
			Arn:          aws.String("arn:aws:kms:us-west-2:857418155548:key/188c57ed-b28a-4c0e-9821-f4940d15cb0a"),
			CreationDate: ExampleDate,
			Description:  aws.String("Default master key"),
			Enabled:      aws.Bool(true),
			KeyId:        aws.String("188c57ed-b28a-4c0e-9821-f4940d15cb0a"),
			KeyManager:   aws.String("AWS"),
			KeyState:     aws.String("Enabled"),
			KeyUsage:     aws.String("ENCRYPT_DECRYPT"),
			Origin:       aws.String("AWS_KMS"),
		},
	}

	ExampleGetKeyPolicyOutput = &kms.GetKeyPolicyOutput{
		Policy: aws.String("{\n  \"Version\" : \"2012-10-17\",\n  \"Id\" : \"auto-awslambda\",\n  \"Statement\" : [ {\n    \"Sid\" : \"Allow access through AWS Lambda for all principals in the account that are authorized to use AWS Lambda\",\n    \"Effect\" : \"Allow\",\n    \"Principal\" : {\n      \"AWS\" : \"*\"\n    },\n    \"Action\" : [ \"kms:Encrypt\", \"kms:Decrypt\", \"kms:ReEncrypt*\", \"kms:GenerateDataKey*\", \"kms:CreateGrant\", \"kms:DescribeKey\" ],\n    \"Resource\" : \"*\",\n    \"Condition\" : {\n      \"StringEquals\" : {\n        \"kms:CallerAccount\" : \"857418155548\",\n        \"kms:ViaService\" : \"lambda.us-west-2.amazonaws.com\"\n      }\n    }\n  }, {\n    \"Sid\" : \"Allow direct access to key metadata to the account\",\n    \"Effect\" : \"Allow\",\n    \"Principal\" : {\n      \"AWS\" : \"arn:aws:iam::857418155548:root\"\n    },\n    \"Action\" : \"kms:Describe*\",\n    \"Resource\" : \"*\"\n  } ]\n}"),
	}

	ExampleListResourceTags = &kms.ListResourceTagsOutput{
		Tags: []*kms.Tag{
			{
				TagKey:   aws.String("Key1"),
				TagValue: aws.String("Value1"),
			},
		},
	}

	svcKmsSetupCalls = map[string]func(*MockKms){
		"ListKeys": func(svc *MockKms) {
			svc.On("ListKeys", mock.Anything).
				Return(ExampleListKeysOutput, nil)
		},
		"GetKeyRotationStatus": func(svc *MockKms) {
			svc.On("GetKeyRotationStatus", mock.Anything).
				Return(ExampleGetKeyRotationStatusOutput, nil)
		},
		"DescribeKey": func(svc *MockKms) {
			svc.On("DescribeKey", mock.Anything).
				Return(ExampleDescribeKeyOutput, nil)
		},
		"GetKeyPolicy": func(svc *MockKms) {
			svc.On("GetKeyPolicy", mock.Anything).
				Return(ExampleGetKeyPolicyOutput, nil)
		},
		"ListResourceTags": func(svc *MockKms) {
			svc.On("ListResourceTags", mock.Anything).
				Return(ExampleListResourceTags, nil)
		},
	}

	svcKmsSetupCallsError = map[string]func(*MockKms){
		"ListKeys": func(svc *MockKms) {
			svc.On("ListKeys", mock.Anything).
				Return(&kms.ListKeysOutput{},
					errors.New("KMS.ListKeys error"),
				)
		},
		"GetKeyRotationStatus": func(svc *MockKms) {
			svc.On("GetKeyRotationStatus", mock.Anything).
				Return(&kms.GetKeyRotationStatusOutput{},
					errors.New("KMS.GetKeyRotationStatus error"),
				)
		},
		"DescribeKey": func(svc *MockKms) {
			svc.On("DescribeKey", mock.Anything).
				Return(&kms.DescribeKeyOutput{},
					errors.New("KMS.DescribeKey error"),
				)
		},
		"GetKeyPolicy": func(svc *MockKms) {
			svc.On("GetKeyPolicy", mock.Anything).
				Return(&kms.GetKeyPolicyOutput{},
					errors.New("KMS.GetKeyPolicy error"),
				)
		},
		"ListResourceTags": func(svc *MockKms) {
			svc.On("ListResourceTags", mock.Anything).
				Return(&kms.ListResourceTagsOutput{},
					errors.New("KMS.ListResourceTags error"),
				)
		},
	}

	MockKmsForSetup = &MockKms{}
)

// KMS mock

// SetupMockKms is used to override the KMS Client initializer
func SetupMockKms(sess *session.Session, cfg *aws.Config) interface{} {
	return MockKmsForSetup
}

// MockKms is a mock KMS client
type MockKms struct {
	kmsiface.KMSAPI
	mock.Mock
}

// BuildMockKmsSvc builds and returns a MockKms struct
//
// Additionally, the appropriate calls to On and Return are made based on the strings passed in
func BuildMockKmsSvc(funcs []string) (mockSvc *MockKms) {
	mockSvc = &MockKms{}
	for _, f := range funcs {
		svcKmsSetupCalls[f](mockSvc)
	}
	return
}

// BuildMockKmsSvcError builds and returns a MockKms struct with errors set
//
// Additionally, the appropriate calls to On and Return are made based on the strings passed in
func BuildMockKmsSvcError(funcs []string) (mockSvc *MockKms) {
	mockSvc = &MockKms{}
	for _, f := range funcs {
		svcKmsSetupCallsError[f](mockSvc)
	}
	return
}

// BuildKmsServiceSvcAll builds and returns a MockKms struct
//
// Additionally, the appropriate calls to On and Return are made for all possible function calls
func BuildMockKmsSvcAll() (mockSvc *MockKms) {
	mockSvc = &MockKms{}
	for _, f := range svcKmsSetupCalls {
		f(mockSvc)
	}
	return
}

// BuildMockKmsSvcAllError builds and returns a MockKms struct with errors set
//
// Additionally, the appropriate calls to On and Return are made for all possible function calls
func BuildMockKmsSvcAllError() (mockSvc *MockKms) {
	mockSvc = &MockKms{}
	for _, f := range svcKmsSetupCallsError {
		f(mockSvc)
	}
	return
}

func (m *MockKms) ListKeys(in *kms.ListKeysInput) (*kms.ListKeysOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*kms.ListKeysOutput), args.Error(1)
}

func (m *MockKms) GetKeyRotationStatus(in *kms.GetKeyRotationStatusInput) (*kms.GetKeyRotationStatusOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*kms.GetKeyRotationStatusOutput), args.Error(1)
}

func (m *MockKms) DescribeKey(in *kms.DescribeKeyInput) (*kms.DescribeKeyOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*kms.DescribeKeyOutput), args.Error(1)
}

func (m *MockKms) GetKeyPolicy(in *kms.GetKeyPolicyInput) (*kms.GetKeyPolicyOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*kms.GetKeyPolicyOutput), args.Error(1)
}

func (m *MockKms) ListResourceTags(in *kms.ListResourceTagsInput) (*kms.ListResourceTagsOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*kms.ListResourceTagsOutput), args.Error(1)
}