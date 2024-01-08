// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mapx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	testCases := []struct {
		name    string
		input   map[int]int
		wantRes []int
	}{
		{
			name:    "nil",
			input:   nil,
			wantRes: []int{},
		},
		{
			name:    "empty",
			input:   map[int]int{},
			wantRes: []int{},
		},
		{
			name: "single",
			input: map[int]int{
				1: 11,
			},
			wantRes: []int{1},
		},
		{
			name: "multiple",
			input: map[int]int{
				1: 11,
				2: 12,
			},
			wantRes: []int{1, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Keys[int, int](tc.input)
			assert.ElementsMatch(t, tc.wantRes, res)
		})
	}
}
func TestValues(t *testing.T) {
	testCases := []struct {
		name    string
		input   map[int]int
		wantRes []int
	}{
		{
			name:    "nil",
			input:   nil,
			wantRes: []int{},
		},
		{
			name:    "empty",
			input:   map[int]int{},
			wantRes: []int{},
		},
		{
			name: "single",
			input: map[int]int{
				1: 11,
			},
			wantRes: []int{11},
		},
		{
			name: "multiple",
			input: map[int]int{
				1: 11,
				2: 12,
			},
			wantRes: []int{11, 12},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Values[int, int](tc.input)
			assert.ElementsMatch(t, tc.wantRes, res)
		})
	}
}

func TestKeysValues(t *testing.T) {
	testCases := []struct {
		name       string
		input      map[int]int
		wantKeys   []int
		wantValues []int
	}{
		{
			name:       "nil",
			input:      nil,
			wantKeys:   []int{},
			wantValues: []int{},
		},
		{
			name:       "empty",
			input:      map[int]int{},
			wantKeys:   []int{},
			wantValues: []int{},
		},
		{
			name: "single",
			input: map[int]int{
				1: 11,
			},
			wantKeys:   []int{1},
			wantValues: []int{11},
		},
		{
			name: "multiple",
			input: map[int]int{
				1: 11,
				2: 12,
			},
			wantKeys:   []int{1, 2},
			wantValues: []int{11, 12},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keys, values := KeysValues[int, int](tc.input)
			assert.ElementsMatch(t, tc.wantKeys, keys)
			assert.ElementsMatch(t, tc.wantValues, values)
		})
	}
}

func TestToMap(t *testing.T) {
	type caseType struct {
		keys   []int
		values []string

		result map[int]string
		err    error
	}
	for _, c := range []caseType{
		{
			keys:   []int{1, 2, 3},
			values: []string{"1", "2", "3"},
			result: map[int]string{
				1: "1",
				2: "2",
				3: "3",
			},
			err: nil,
		},
		{
			keys:   []int{1, 2, 3},
			values: []string{"1", "2"},
			result: nil,
			err:    fmt.Errorf("keys与values的长度不同, len(keys)=3, len(values)=2"),
		},
		{
			keys:   []int{1, 2, 3},
			values: nil,
			result: nil,
			err:    fmt.Errorf("keys与values均不可为nil"),
		},
		{
			keys:   nil,
			values: []string{"1", "2"},
			result: nil,
			err:    fmt.Errorf("keys与values均不可为nil"),
		},
		{
			keys:   nil,
			values: nil,
			result: nil,
			err:    fmt.Errorf("keys与values均不可为nil"),
		},
		{
			keys:   []int{1, 2, 3, 1, 1},
			values: []string{"1", "2", "3", "10", "100"},
			result: map[int]string{
				1: "100",
				2: "2",
				3: "3",
			},
			err: nil,
		},
	} {
		result, err := ToMap(c.keys, c.values)
		assert.Equal(t, c.err, err)
		assert.Equal(t, c.result, result)
	}
}
