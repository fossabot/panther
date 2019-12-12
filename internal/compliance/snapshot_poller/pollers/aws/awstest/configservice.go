package awstest

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
	"github.com/stretchr/testify/mock"
)

var (

	//
	// Example configservice variables
	//

	ExampleDescribeConfigurationRecorders = &configservice.DescribeConfigurationRecordersOutput{
		ConfigurationRecorders: []*configservice.ConfigurationRecorder{
			{
				Name: aws.String("default"),
				RecordingGroup: &configservice.RecordingGroup{
					AllSupported:               aws.Bool(true),
					IncludeGlobalResourceTypes: aws.Bool(true),
				},
				RoleARN: aws.String("arn:aws:iam::857418155548:role/aws-service-role/config.amazonaws.com/AWSServiceRoleForConfig"),
			},
		},
	}

	ExampleDescribeConfigurationRecorderStatus = &configservice.DescribeConfigurationRecorderStatusOutput{
		ConfigurationRecordersStatus: []*configservice.ConfigurationRecorderStatus{
			{
				LastStartTime:        &ExampleTimeParsed,
				LastStatus:           aws.String("SUCCESS"),
				LastStatusChangeTime: &ExampleTimeParsed,
				LastStopTime:         &ExampleTimeParsed,
				Name:                 aws.String("default"),
				Recording:            aws.Bool(true),
			},
		},
	}

	ExampleConfigName = aws.String("IAmTheConfig")

	svcConfigSetupCalls = map[string]func(*MockConfigService){
		"DescribeConfigurationRecorders": func(svc *MockConfigService) {
			svc.On("DescribeConfigurationRecorders", mock.Anything).
				Return(ExampleDescribeConfigurationRecorders, nil)
		},
		"DescribeConfigurationRecorderStatus": func(svc *MockConfigService) {
			svc.On("DescribeConfigurationRecorderStatus", mock.Anything).
				Return(ExampleDescribeConfigurationRecorderStatus, nil)
		},
	}

	svcConfigSetupCallsError = map[string]func(*MockConfigService){
		"DescribeConfigurationRecorders": func(svc *MockConfigService) {
			svc.On("DescribeConfigurationRecorders", mock.Anything).
				Return(
					&configservice.DescribeConfigurationRecordersOutput{},
					errors.New("fake ConfigService.DescribeRecorders error"),
				)
		},
		"DescribeConfigurationRecorderStatus": func(svc *MockConfigService) {
			svc.On("DescribeConfigurationRecorderStatus", mock.Anything).
				Return(
					&configservice.DescribeConfigurationRecorderStatusOutput{},
					errors.New("fake ConfigService.DescribeStatus error"),
				)
		},
	}

	MockConfigServiceForSetup = &MockConfigService{}
)

// SetupMockConfigService is used to override the ConfigService client initializer
func SetupMockConfigService(sess *session.Session, cfg *aws.Config) interface{} {
	return MockConfigServiceForSetup
}

type MockConfigService struct {
	configserviceiface.ConfigServiceAPI
	mock.Mock
}

// BuildMockConfigServiceSvc builds and returns a MockConfigService struct
//
// Additionally, the appropriate calls to On and Return are made based on the strings passed in
func BuildMockConfigServiceSvc(funcs []string) (mockSvc *MockConfigService) {
	mockSvc = &MockConfigService{}
	for _, f := range funcs {
		svcConfigSetupCalls[f](mockSvc)
	}
	return
}

// BuildMockConfigServiceSvcError builds and returns a MockConfigService struct with errors set
//
// Additionally, the appropriate calls to On and Return are made based on the strings passed in
func BuildMockConfigServiceSvcError(funcs []string) (mockSvc *MockConfigService) {
	mockSvc = &MockConfigService{}
	for _, f := range funcs {
		svcConfigSetupCallsError[f](mockSvc)
	}
	return
}

// BuildMockConfigServiceSvcAll builds and returns a MockConfigService struct
//
// Additionally, the appropriate calls to On and Return are made for all possible function calls
func BuildMockConfigServiceSvcAll() (mockSvc *MockConfigService) {
	mockSvc = &MockConfigService{}
	for _, f := range svcConfigSetupCalls {
		f(mockSvc)
	}
	return
}

// BuildMockConfigServiceSvcAllError builds and returns a MockConfigService struct with errors set
//
// Additionally, the appropriate calls to On and Return are made for all possible function calls
func BuildMockConfigServiceSvcAllError() (mockSvc *MockConfigService) {
	mockSvc = &MockConfigService{}
	for _, f := range svcConfigSetupCallsError {
		f(mockSvc)
	}
	return
}

func (m *MockConfigService) DescribeConfigurationRecorders(
	in *configservice.DescribeConfigurationRecordersInput) (*configservice.DescribeConfigurationRecordersOutput, error) {

	args := m.Called(in)
	return args.Get(0).(*configservice.DescribeConfigurationRecordersOutput), args.Error(1)
}

func (m *MockConfigService) DescribeConfigurationRecorderStatus(
	in *configservice.DescribeConfigurationRecorderStatusInput) (*configservice.DescribeConfigurationRecorderStatusOutput, error) {

	args := m.Called(in)
	return args.Get(0).(*configservice.DescribeConfigurationRecorderStatusOutput), args.Error(1)
}