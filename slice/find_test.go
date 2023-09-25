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

package slice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	testCases := []struct {
		name  string
		input []Number
		match matchFunc[Number]

		wantVal Number
		found   bool
	}{
		{
			name: "找到了",
			input: []Number{
				{val: 123},
				{val: 234},
			},
			match: func(src Number) bool {
				return src.val == 123
			},
			wantVal: Number{val: 123},
			found:   true,
		},
		{
			name: "没找到",
			input: []Number{
				{val: 123},
				{val: 234},
			},
			match: func(src Number) bool {
				return src.val == 456
			},
		},
		{
			name: "nil",
			match: func(src Number) bool {
				return src.val == 123
			},
		},
		{
			name:  "没有元素",
			input: []Number{},
			match: func(src Number) bool {
				return src.val == 123
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, found := Find[Number](tc.input, tc.match)
			assert.Equal(t, tc.found, found)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestFindAll(t *testing.T) {
	testCases := []struct {
		name  string
		input []Number
		match matchFunc[Number]

		wantVals []Number
	}{
		{
			name:  "没有符合条件的",
			input: []Number{{val: 2}, {val: 4}},
			match: func(src Number) bool {
				return src.val%2 == 1
			},
			wantVals: []Number{},
		},
		{
			name:  "找到了",
			input: []Number{{val: 2}, {val: 3}, {val: 4}},
			match: func(src Number) bool {
				return src.val%2 == 1
			},
			wantVals: []Number{{val: 3}},
		},
		{
			name: "nil",
			match: func(src Number) bool {
				return src.val%2 == 1
			},
			wantVals: []Number{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vals := FindAll[Number](tc.input, tc.match)
			assert.Equal(t, tc.wantVals, vals)
		})
	}
}

func ExampleFind() {
	val, ok := Find[int]([]int{1, 2, 3}, func(src int) bool {
		return src == 2
	})
	fmt.Println(val, ok)
	val, ok = Find[int]([]int{1, 2, 3}, func(src int) bool {
		return src == 4
	})
	fmt.Println(val, ok)
	// Output:
	// 2 true
	// 0 false
}

func ExampleFindAll() {
	vals := FindAll[int]([]int{2, 3, 4}, func(src int) bool {
		return src%2 == 1
	})
	fmt.Println(vals)
	vals = FindAll[int]([]int{2, 3, 4}, func(src int) bool {
		return src > 5
	})
	fmt.Println(vals)
	// Output:
	// [3]
	// []
}

type Number struct {
	val int
}
