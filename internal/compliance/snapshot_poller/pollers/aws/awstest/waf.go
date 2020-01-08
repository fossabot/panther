package awstest

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
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/waf/wafiface"
	"github.com/stretchr/testify/mock"
)

// Example WAF regional API return values
var (
	PaginateListWebAcls = false

	ExampleListWebAclsOutput = &waf.ListWebACLsOutput{
		WebACLs: []*waf.WebACLSummary{
			{
				Name:     aws.String("example-web-acl-1"),
				WebACLId: aws.String("asdfasdf-f123-e123-g123-1234asdf1234"),
			},
			{
				Name:     aws.String("example-web-acl-2"),
				WebACLId: aws.String("asdfasdf-x123-y123-z123-1234asdf1234"),
			},
		},
		NextMarker: aws.String("asdfasdf-x123-y123-z123-1234asdf1234"),
	}

	ExampleGetWebAclOutput = &waf.GetWebACLOutput{
		WebACL: &waf.WebACL{
			WebACLId:   aws.String("asdfasdf-f123-e123-g123-1234asdf1234"),
			Name:       aws.String("example-web-acl-1"),
			MetricName: aws.String("examplewebacl1"),
			DefaultAction: &waf.WafAction{
				Type: aws.String("ALLOW"),
			},
			WebACLArn: aws.String("arn:aws:waf-regional:us-west-2:123456789012:webacl/asdfasdf-f123-e123-g123-1234asdf1234"),
			Rules: []*waf.ActivatedRule{
				{
					Priority: aws.Int64(1),
					RuleId:   aws.String("112233"),
					Action: &waf.WafAction{
						Type: aws.String("COUNT"),
					},
					Type: aws.String("REGULAR"),
				},
			},
		},
	}

	ExampleGetRule = &waf.GetRuleOutput{
		Rule: &waf.Rule{
			RuleId:     aws.String("112233"),
			Name:       aws.String("test-rule"),
			MetricName: aws.String("testrule"),
			Predicates: []*waf.Predicate{
				{
					Negated: aws.Bool(false),
					Type:    aws.String("XssMatch"),
					DataId:  aws.String("123abc-123def"),
				},
			},
		},
	}

	ExampleListTagsForResourceWaf = &waf.ListTagsForResourceOutput{
		TagInfoForResource: &waf.TagInfoForResource{
			TagList: []*waf.Tag{
				{
					Key:   aws.String("Key1"),
					Value: aws.String("Value1"),
				},
			},
		},
	}

	svcWafSetupCalls = map[string]func(*MockWaf){
		"ListWebACLs": func(svc *MockWaf) {
			PaginateListWebAcls = false
			svc.On("ListWebACLs", mock.Anything).
				Return(ExampleListWebAclsOutput, nil)
		},
		"GetWebACL": func(svc *MockWaf) {
			svc.On("GetWebACL", mock.Anything).
				Return(ExampleGetWebAclOutput, nil)
		},
		"ListTagsForResource": func(svc *MockWaf) {
			svc.On("ListTagsForResource", mock.Anything).
				Return(ExampleListTagsForResourceWaf, nil)
		},
		"GetRule": func(svc *MockWaf) {
			svc.On("GetRule", mock.Anything).
				Return(ExampleGetRule, nil)
		},
	}

	svcWafSetupCallsError = map[string]func(*MockWaf){
		"ListWebACLs": func(svc *MockWaf) {
			svc.On("ListWebACLs", mock.Anything).
				Return(
					&waf.ListWebACLsOutput{},
					errors.New("WAF.ListWebACLs"),
				)
		},
		"GetWebACL": func(svc *MockWaf) {
			svc.On("GetWebACL", mock.Anything).
				Return(&waf.GetWebACLOutput{},
					errors.New("WAF.GetWebACL error"),
				)
		},
		"ListTagsForResource": func(svc *MockWaf) {
			svc.On("ListTagsForResource", mock.Anything).
				Return(&waf.ListTagsForResourceOutput{},
					errors.New("WAF.ListTagsForResource error"),
				)
		},
		"GetRule": func(svc *MockWaf) {
			svc.On("GetRule", mock.Anything).
				Return(&waf.GetRuleOutput{},
					errors.New("WAF.GetRule error"),
				)
		},
	}

	MockWafForSetup = &MockWaf{}
)

// WAF mock

// SetupMockWaf is used to override the WAF Client initializer
func SetupMockWaf(sess *session.Session, cfg *aws.Config) interface{} {
	return MockWafForSetup
}

// MockWaf is a mock WAF client
type MockWaf struct {
	wafiface.WAFAPI
	mock.Mock
}

// BuildMockWafSvc builds and returns a MockWaf struct
//
// Additionally, the appropriate calls to On and Return are made based on the strings passed in
func BuildMockWafSvc(funcs []string) (mockSvc *MockWaf) {
	mockSvc = &MockWaf{}
	for _, f := range funcs {
		svcWafSetupCalls[f](mockSvc)
	}
	return
}

// BuildMockWafSvcError builds and returns a MockWaf struct with errors set
//
// Additionally, the appropriate calls to On and Return are made based on the strings passed in
func BuildMockWafSvcError(funcs []string) (mockSvc *MockWaf) {
	mockSvc = &MockWaf{}
	for _, f := range funcs {
		svcWafSetupCallsError[f](mockSvc)
	}
	return
}

// BuildWafServiceSvcAll builds and returns a MockWaf struct
//
// Additionally, the appropriate calls to On and Return are made for all possible function calls
func BuildMockWafSvcAll() (mockSvc *MockWaf) {
	mockSvc = &MockWaf{}
	for _, f := range svcWafSetupCalls {
		f(mockSvc)
	}
	return
}

// BuildMockWafSvcAllError builds and returns a MockWaf struct with errors set
//
// Additionally, the appropriate calls to On and Return are made for all possible function calls
func BuildMockWafSvcAllError() (mockSvc *MockWaf) {
	mockSvc = &MockWaf{}
	for _, f := range svcWafSetupCallsError {
		f(mockSvc)
	}
	return
}

func (m *MockWaf) ListWebACLs(in *waf.ListWebACLsInput) (*waf.ListWebACLsOutput, error) {
	PaginateListWebAcls = !PaginateListWebAcls
	args := m.Called(in)
	if PaginateListWebAcls {
		return args.Get(0).(*waf.ListWebACLsOutput), args.Error(1)
	}
	var empty []*waf.WebACLSummary
	return &waf.ListWebACLsOutput{WebACLs: empty}, args.Error(1)
}

func (m *MockWaf) GetWebACL(in *waf.GetWebACLInput) (*waf.GetWebACLOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*waf.GetWebACLOutput), args.Error(1)
}

func (m *MockWaf) ListTagsForResource(in *waf.ListTagsForResourceInput) (*waf.ListTagsForResourceOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*waf.ListTagsForResourceOutput), args.Error(1)
}

func (m *MockWaf) GetRule(in *waf.GetRuleInput) (*waf.GetRuleOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*waf.GetRuleOutput), args.Error(1)
}
