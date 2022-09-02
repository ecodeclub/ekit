// Copyright 2021 gotomicro
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

func TestDiff(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "src and dst nil",
			src:  nil,
			dst:  nil,
			want: nil,
		},
		{
			name: "only src nil",
			src:  nil,
			dst:  []int{1, 2, 3},
			want: nil,
		},
		{
			name: "only dst nil",
			src:  []int{1, 2, 3},
			dst:  nil,
			want: []int{1, 2, 3},
		},
		{
			name: "src and dst empty slice",
			src:  []int{},
			dst:  []int{},
			want: []int{},
		},
		{
			name: "only src empty slice",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: []int{},
		},
		{
			name: "only dst empty slice",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "src have diff element",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{1, 2, 3, 5},
			want: []int{4},
		},
		{
			name: "src not have diff element",
			src:  []int{1, 2, 3, 5},
			dst:  []int{1, 2, 3, 4, 5},
			want: []int{},
		},
		{
			name: "src and dst have diff element",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{3, 4, 5, 6, 7},
			want: []int{1, 2},
		},
		{
			name: "not sorted array",
			src:  []int{3, 4, 5, 1, 2},
			dst:  []int{6, 7, 3, 4, 5},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Diff[int](tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestDiffAny(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "src and dst nil",
			src:  nil,
			dst:  nil,
			want: nil,
		},
		{
			name: "only src nil",
			src:  nil,
			dst:  []int{1, 2, 3},
			want: nil,
		},
		{
			name: "only dst nil",
			src:  []int{1, 2, 3},
			dst:  nil,
			want: []int{1, 2, 3},
		},
		{
			name: "src and dst empty slice",
			src:  []int{},
			dst:  []int{},
			want: []int{},
		},
		{
			name: "only src empty slice",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: []int{},
		},
		{
			name: "only dst empty slice",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "src have diff element",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{1, 2, 3, 5},
			want: []int{4},
		},
		{
			name: "src not have diff element",
			src:  []int{1, 2, 3, 5},
			dst:  []int{1, 2, 3, 4, 5},
			want: []int{},
		},
		{
			name: "src and dst have diff element",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{3, 4, 5, 6, 7},
			want: []int{1, 2},
		},
		{
			name: "not sorted array",
			src:  []int{3, 4, 5, 1, 2},
			dst:  []int{6, 7, 3, 4, 5},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := DiffFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSymmetricDiff(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "src and dst nil",
			src:  nil,
			dst:  nil,
			want: nil,
		},
		{
			name: "only src nil",
			src:  nil,
			dst:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "only dst nil",
			src:  []int{1, 2, 3},
			dst:  nil,
			want: []int{1, 2, 3},
		},
		{
			name: "src and dst empty slice",
			src:  []int{},
			dst:  []int{},
			want: []int{},
		},
		{
			name: "only src empty slice",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "only dst empty slice",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "src have diff element",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{1, 2, 3, 5},
			want: []int{4},
		},
		{
			name: "dst have diff element",
			src:  []int{1, 2, 3, 5},
			dst:  []int{1, 2, 3, 4, 5},
			want: []int{4},
		},
		{
			name: "src and dst have diff element",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{3, 4, 5, 6, 7},
			want: []int{1, 2, 6, 7},
		},
		{
			name: "not sorted array",
			src:  []int{3, 4, 5, 1, 2},
			dst:  []int{6, 7, 3, 4, 5},
			want: []int{1, 2, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := SymmetricDiff[int](tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSymmetricDiffFunc(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "src and dst nil",
			src:  nil,
			dst:  nil,
			want: nil,
		},
		{
			name: "only src nil",
			src:  nil,
			dst:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "only dst nil",
			src:  []int{1, 2, 3},
			dst:  nil,
			want: []int{1, 2, 3},
		},
		{
			name: "src and dst empty slice",
			src:  []int{},
			dst:  []int{},
			want: []int{},
		},
		{
			name: "only src empty slice",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "only dst empty slice",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "src have diff element",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{1, 2, 3, 5},
			want: []int{4},
		},
		{
			name: "dst have diff element",
			src:  []int{1, 2, 3, 5},
			dst:  []int{1, 2, 3, 4, 5},
			want: []int{4},
		},
		{
			name: "src and dst have diff element",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{3, 4, 5, 6, 7},
			want: []int{1, 2, 6, 7},
		},
		{
			name: "not sorted array",
			src:  []int{3, 4, 5, 1, 2},
			dst:  []int{6, 7, 3, 4, 5},
			want: []int{1, 2, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := SymmetricDiffFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}

func ExampleDiff() {
	src := []int{1, 2, 3, 4, 5}
	des := []int{3, 4, 5, 6, 7}
	diff := Diff(src, des)
	fmt.Println(diff)
	// Output: [1 2]
}

func ExampleDiffFunc() {
	src := []int{1, 2, 3, 4, 5}
	des := []int{3, 4, 5, 6, 7}
	equalFunc := func(src int, dst int) bool {
		return src == dst
	}
	diff := DiffFunc(src, des, equalFunc)
	fmt.Println(diff)
	// Output: [1 2]
}

func ExampleSymmetricDiff() {
	src := []int{1, 2, 3, 4, 5}
	des := []int{3, 4, 5, 6, 7}
	diff := SymmetricDiff(src, des)
	fmt.Println(diff)
	// Output: [1 2 6 7]
}

func ExampleSymmetricDiffFunc() {
	src := []int{1, 2, 3, 4, 5}
	des := []int{3, 4, 5, 6, 7}
	equalFunc := func(src int, dst int) bool {
		return src == dst
	}
	diff := SymmetricDiffFunc(src, des, equalFunc)
	fmt.Println(diff)
	// Output: [1 2 6 7]
}
