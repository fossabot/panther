package handlers

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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/gateway/analysis/models"
)

func TestPagePoliciesPageSize1(t *testing.T) {
	policies := []*models.PolicySummary{{ID: "a"}, {ID: "b"}, {ID: "c"}, {ID: "d"}}
	result := pagePolicies(policies, 1, 1)
	expected := &models.PolicyList{
		Paging: &models.Paging{
			ThisPage:   aws.Int64(1),
			TotalItems: aws.Int64(4),
			TotalPages: aws.Int64(4),
		},
		Policies: []*models.PolicySummary{{ID: "a"}},
	}
	assert.Equal(t, expected, result)

	result = pagePolicies(policies, 1, 2)
	expected.Paging.ThisPage = aws.Int64(2)
	expected.Policies = []*models.PolicySummary{{ID: "b"}}
	assert.Equal(t, expected, result)

	result = pagePolicies(policies, 1, 3)
	expected.Paging.ThisPage = aws.Int64(3)
	expected.Policies = []*models.PolicySummary{{ID: "c"}}
	assert.Equal(t, expected, result)

	result = pagePolicies(policies, 1, 4)
	expected.Paging.ThisPage = aws.Int64(4)
	expected.Policies = []*models.PolicySummary{{ID: "d"}}
	assert.Equal(t, expected, result)
}

func TestPagePoliciesSinglePage(t *testing.T) {
	policies := []*models.PolicySummary{{ID: "a"}, {ID: "b"}, {ID: "c"}, {ID: "d"}}
	result := pagePolicies(policies, 25, 1)
	expected := &models.PolicyList{
		Paging: &models.Paging{
			ThisPage:   aws.Int64(1),
			TotalItems: aws.Int64(4),
			TotalPages: aws.Int64(1),
		},
		Policies: policies,
	}
	assert.Equal(t, expected, result)
}

func TestPagePoliciesPageOutOfBounds(t *testing.T) {
	policies := []*models.PolicySummary{{ID: "a"}, {ID: "b"}, {ID: "c"}, {ID: "d"}}
	result := pagePolicies(policies, 1, 10)
	expected := &models.PolicyList{
		Paging: &models.Paging{
			ThisPage:   aws.Int64(10),
			TotalItems: aws.Int64(4),
			TotalPages: aws.Int64(4),
		},
		Policies: []*models.PolicySummary{}, // empty list - page out of bounds
	}
	assert.Equal(t, expected, result)
}
