package utils

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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	key1 = "key1"
	key2 = "key2"
	val1 = "value1"
	val2 = "value2"
)

func TestParseTagSlice(t *testing.T) {
	tagSlice := []*struct {
		Key   *string
		Value *string
	}{
		{
			Key:   &key1,
			Value: &val1,
		},
		{
			Key:   &key2,
			Value: &val2,
		},
	}

	tagMap := ParseTagSlice(tagSlice)
	require.NotEmpty(t, tagMap)
	assert.Equal(t, val1, *tagMap[key1])
	assert.Equal(t, val2, *tagMap[key2])
}

func TestParseTagSliceEmpty(t *testing.T) {
	tagSlice := make([]*struct {
		Key   *string
		Value *string
	}, 0)

	tagMap := ParseTagSlice(tagSlice)
	require.Empty(t, tagMap)
}

func TestParseTagSliceBadTag(t *testing.T) {
	tagSlice := []*struct {
		Key      *string
		NotValue *string
	}{
		{
			Key:      &key1,
			NotValue: &val1,
		},
		{
			Key:      &key2,
			NotValue: &val2,
		},
	}

	tagMap := ParseTagSlice(tagSlice)
	require.Nil(t, tagMap)
}

func TestParseTagSliceNilSlice(t *testing.T) {
	var tagSlice []*struct {
		Key   *string
		Value *string
	}

	tagMap := ParseTagSlice(tagSlice)
	require.Empty(t, tagMap)
}
