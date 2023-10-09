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

func TestUpdateMap(t *testing.T) {
	testCases := []struct {
		name    string
		m1      map[int]int
		m2      map[int]int
		wantRes map[int]int
	}{
		{
			name:    "update map",
			m1:      map[int]int{1: 1, 2: 2},
			m2:      map[int]int{1: 0, 2: 2, 3: 3},
			wantRes: map[int]int{1: 0, 2: 2, 3: 3},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := UpdateMap[int, int](tc.m1, tc.m2)
			if len(res) != len(tc.wantRes) {
				t.Fatal("Fail, length mismatch")
			}
			for k, v := range res {
				v1, ok := tc.wantRes[k]
				if !ok {
					t.Fatal("Fail, keys not equal")
				}
				if v != v1 {
					t.Fatal("Fail, values not equal")
				}
			}

		})
	}
}
