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

func TestIndex(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: 0,
			name: "first one",
		},
		{
			src:  []int{},
			dst:  1,
			want: -1,
			name: "the length of src is 0",
		},
		{
			dst:  1,
			want: -1,
			name: "src nil",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: -1,
			name: "dst not exist",
		},
		{
			src:  []int{1, 3, 4, 2, 0},
			dst:  0,
			want: 4,
			name: "last one",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, Index[int](test.src, test.dst))
		})
	}
}

func TestIndexFunc(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: 0,
			name: "first one",
		},
		{
			src:  []int{},
			dst:  1,
			want: -1,
			name: "the length of src is 0",
		},
		{
			dst:  1,
			want: -1,
			name: "src nil",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: -1,
			name: "dst not exist",
		},
		{
			src:  []int{1, 3, 4, 2, 0},
			dst:  0,
			want: 4,
			name: "last one",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, IndexFunc[int](test.src, test.dst, func(src, dst int) bool {
				return src == dst
			}))
		})
	}
}

func TestLastIndex(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: 1,
			name: "first one",
		},
		{
			src:  []int{},
			dst:  1,
			want: -1,
			name: "the length of src is 0",
		},
		{
			dst:  1,
			want: -1,
			name: "src nil",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: -1,
			name: "dst not exist",
		},
		{
			src:  []int{0, 1, 3, 4, 2, 0},
			dst:  0,
			want: 5,
			name: "last one",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, LastIndex[int](test.src, test.dst))
	}
}

func TestLastIndexFunc(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: 1,
			name: "first one",
		},
		{
			src:  []int{},
			dst:  1,
			want: -1,
			name: "the length of src is 0",
		},
		{
			dst:  1,
			want: -1,
			name: "src nil",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: -1,
			name: "dst not exist",
		},
		{
			src:  []int{0, 1, 3, 4, 2, 0},
			dst:  0,
			want: 5,
			name: "last one",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, LastIndexFunc[int](test.src, test.dst, func(src, dst int) bool {
			return src == dst
		}))
	}
}

func TestIndexAll(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want []int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: []int{0, 1},
			name: "normal test",
		},
		{
			src:  []int{},
			dst:  1,
			want: []int{},
			name: "the length of src is 0",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: []int{},
			name: "dst not exist",
		},
		{
			src:  []int{0, 1, 3, 4, 2, 0},
			dst:  0,
			want: []int{0, 5},
			name: "normal test",
		},
	}
	for _, test := range tests {
		res := IndexAll[int](test.src, test.dst)
		assert.ElementsMatch(t, test.want, res)
	}
}

func TestIndexAllFunc(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want []int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: []int{0, 1},
			name: "normal test",
		},
		{
			src:  []int{},
			dst:  1,
			want: []int{},
			name: "the length of src is 0",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: []int{},
			name: "dst not exist",
		},
		{
			src:  []int{0, 1, 3, 4, 2, 0},
			dst:  0,
			want: []int{0, 5},
			name: "normal test",
		},
	}
	for _, test := range tests {
		res := IndexAllFunc[int](test.src, test.dst, func(src, dst int) bool {
			return src == dst
		})
		assert.ElementsMatch(t, test.want, res)
	}
}

func ExampleIndex() {
	res := Index[int]([]int{1, 2, 3}, 1)
	fmt.Println(res)
	res = Index[int]([]int{1, 2, 3}, 4)
	fmt.Println(res)
	// Output:
	// 0
	// -1
}

func ExampleIndexFunc() {
	res := IndexFunc[int]([]int{1, 2, 3}, 1, func(src, dst int) bool {
		return src == dst
	})
	fmt.Println(res)
	res = IndexFunc[int]([]int{1, 2, 3}, 4, func(src, dst int) bool {
		return src == dst
	})
	fmt.Println(res)
	// Output:
	// 0
	// -1
}

func ExampleIndexAll() {
	res := IndexAll[int]([]int{1, 2, 3, 4, 5, 3, 9}, 3)
	fmt.Println(res)
	res = IndexAll[int]([]int{1, 2, 3}, 4)
	fmt.Println(res)
	// Output:
	// [2 5]
	// []
}

func ExampleIndexAllFunc() {
	res := IndexAllFunc[int]([]int{1, 2, 3, 4, 5, 3, 9}, 3, func(src, dst int) bool {
		return src == dst
	})
	fmt.Println(res)
	res = IndexAllFunc[int]([]int{1, 2, 3}, 4, func(src, dst int) bool {
		return src == dst
	})
	fmt.Println(res)
	// Output:
	// [2 5]
	// []
}

// BenchmarkIndex 主要是为了验证即便我们在 Index 这种方法里面直接调用 IndexFunc
// 性能损失几乎没有。
func BenchmarkIndex(b *testing.B) {
	b.Run("loop directly", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			IndexByLoop[int]([]int{1, 2, 3, 4, 5, 6}, 5)
		}
	})
	b.Run("delegate to IndexFunc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Index[int]([]int{1, 2, 3, 4, 5, 6}, 5)
		}
	})
}

func IndexByLoop[T comparable](src []T, dst T) int {
	for i, val := range src {
		if val == dst {
			return i
		}
	}
	return -1
}
