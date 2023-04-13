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

	"github.com/ecodeclub/ekit"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	testCases := []struct {
		name  string
		input []Integer
		want  Integer
	}{
		{
			name:  "value",
			input: []Integer{1},
			want:  1,
		},
		{
			name:  "values",
			input: []Integer{2, 3, 1},
			want:  3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Max[Integer](tc.input)
			assert.Equal(t, tc.want, res)
		})
	}

	assert.Panics(t, func() {
		Max[int](nil)
	})
	assert.Panics(t, func() {
		Max[int]([]int{})
	})

	testMaxTypes[uint](t)
	testMaxTypes[uint8](t)
	testMaxTypes[uint16](t)
	testMaxTypes[uint32](t)
	testMaxTypes[uint64](t)
	testMaxTypes[int](t)
	testMaxTypes[int8](t)
	testMaxTypes[int16](t)
	testMaxTypes[int32](t)
	testMaxTypes[int64](t)
	testMaxTypes[float32](t)
	testMaxTypes[float64](t)
}

func TestMin(t *testing.T) {
	testCases := []struct {
		name  string
		input []Integer
		want  Integer
	}{
		{
			name:  "value",
			input: []Integer{3},
			want:  3,
		},
		{
			name:  "values",
			input: []Integer{3, 1, 2},
			want:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Min[Integer](tc.input)
			assert.Equal(t, tc.want, res)
		})
	}

	assert.Panics(t, func() {
		Min[int](nil)
	})
	assert.Panics(t, func() {
		Min[int]([]int{})
	})

	testMinTypes[uint](t)
	testMinTypes[uint8](t)
	testMinTypes[uint16](t)
	testMinTypes[uint32](t)
	testMinTypes[uint64](t)
	testMinTypes[int](t)
	testMinTypes[int8](t)
	testMinTypes[int16](t)
	testMinTypes[int32](t)
	testMinTypes[int64](t)
	testMinTypes[float32](t)
	testMinTypes[float64](t)
}

func TestSum(t *testing.T) {
	testCases := []struct {
		name  string
		input []Integer
		want  Integer
	}{
		{
			name: "nil",
		},
		{
			name:  "empty",
			input: []Integer{},
		},
		{
			name:  "value",
			input: []Integer{1},
			want:  1,
		},
		{
			name:  "values",
			input: []Integer{1, 2, 3},
			want:  6,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Sum[Integer](tc.input)
			assert.Equal(t, tc.want, res)
		})
	}

	testSumTypes[uint](t)
	testSumTypes[uint8](t)
	testSumTypes[uint16](t)
	testSumTypes[uint32](t)
	testSumTypes[uint64](t)
	testSumTypes[int](t)
	testSumTypes[int8](t)
	testSumTypes[int16](t)
	testSumTypes[int32](t)
	testSumTypes[int64](t)
	testSumTypes[float32](t)
	testSumTypes[float64](t)
}

// testMaxTypes 只是用来测试一下满足 Max 方法约束的所有类型
func testMaxTypes[T ekit.RealNumber](t *testing.T) {
	res := Max[T]([]T{1, 2, 3})
	assert.Equal(t, T(3), res)
}

// testMinTypes 只是用来测试一下满足 Min 方法约束的所有类型
func testMinTypes[T ekit.RealNumber](t *testing.T) {
	res := Min[T]([]T{1, 2, 3})
	assert.Equal(t, T(1), res)
}

// testSumTypes 只是用来测试一下满足 Sum 方法约束的所有类型
func testSumTypes[T ekit.RealNumber](t *testing.T) {
	res := Sum[T]([]T{1, 2, 3})
	assert.Equal(t, T(6), res)
}

type Integer int

func ExampleSum() {
	res := Sum[int]([]int{1, 2, 3})
	fmt.Println(res)
	res = Sum[int](nil)
	fmt.Println(res)
	// Output:
	// 6
	// 0
}

func ExampleMin() {
	res := Min[int]([]int{1, 2, 3})
	fmt.Println(res)
	// Output:
	// 1
}

func ExampleMax() {
	res := Max[int]([]int{1, 2, 3})
	fmt.Println(res)
	// Output:
	// 3
}
