package handlers

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/analysis/models"
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
