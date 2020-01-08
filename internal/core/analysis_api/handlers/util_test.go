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
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/gateway/analysis/models"
)

func TestLowerSet(t *testing.T) {
	result := lowerSet([]string{"AWS", "aws", "CIS", "cis", "Panther"})
	sort.Strings(result)
	assert.Equal(t, []string{"aws", "cis", "panther"}, result)
}

func TestIntMin(t *testing.T) {
	assert.Equal(t, -2, intMin(-2, 0))
	assert.Equal(t, 0, intMin(0, 0))
	assert.Equal(t, 5, intMin(10, 5))
}

func TestSetDifference(t *testing.T) {
	assert.Empty(t, setDifference([]string{}, []string{}))
	assert.Empty(t, setDifference([]string{"a", "b", "c"}, []string{"c", "a", "b"}))
	assert.Equal(t, []string{"a", "b"}, setDifference([]string{"a", "b"}, nil))
	assert.Empty(t, setDifference(nil, []string{"a", "b"}))
	assert.Equal(t, []string{"panther", "labs"},
		setDifference([]string{"panther", "labs", "inc"}, []string{"inc", "runpanther.io"}))
}

func TestSetEquality(t *testing.T) {
	assert.True(t, setEquality(nil, []string{}))
	assert.True(t, setEquality([]string{"panther", "labs", "inc"}, []string{"inc", "labs", "panther"}))
	assert.False(t, setEquality([]string{"panther"}, []string{"panther", "labs"}))
	assert.False(t, setEquality([]string{"panther", "labs"}, []string{"panther", "inc"}))
}

func TestPoliciesEqual(t *testing.T) {
	first := &tableItem{
		Body:          "def policy(resource): return True",
		ID:            "My:Favorite:Policy",
		ResourceTypes: nil,
		Tags:          []string{"AWS", "pci", "S3"},
		Tests: []*models.UnitTest{
			{
				ExpectedResult: true,
				Resource:       `{"first": "a", "second": "b", "third": "c"}`,
			},
		},

		// These fields shouldn't matter for equality testing
		CreatedAt: models.ModifyTime(time.Now()),
		CreatedBy: "austin",
		VersionID: "version-the-first",
		LowerID:   "my:favorite:policy",
		LowerTags: []string{"aws", "pci", "s3"},
	}
	second := &tableItem{
		Body:          first.Body,
		ID:            first.ID,
		ResourceTypes: []string{},
		Tags:          []string{"pci", "S3", "AWS"}, // same tags, different order
		Tests: []*models.UnitTest{
			{
				ExpectedResult: true,
				Resource:       `{"second":"b","first":"a","third":"c"}`, // same resource, different order
			},
		},
	}
	equal, err := policiesEqual(first, second)
	assert.True(t, equal)
	assert.NoError(t, err)

	second.Body = "def policy(resource): return False"
	equal, err = policiesEqual(first, second)
	assert.False(t, equal)
	assert.NoError(t, err)
}
