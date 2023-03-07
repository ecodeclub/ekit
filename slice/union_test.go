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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnionSet(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6, 1},
			want: []int{1, 2, 3, 4, 5, 6},
			name: "not empty",
		},
		{
			src:  []int{},
			dst:  []int{1, 3},
			want: []int{1, 3},
			name: "src is empty",
		},
		{

			src:  []int{1, 3},
			dst:  []int{},
			want: []int{1, 3},
			name: "dst is empty",
		},
		{
			src:  []int{},
			dst:  []int{},
			want: []int{},
			name: "src and dst are empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := UnionSet[int](tt.src, tt.dst)
			assert.ElementsMatch(t, tt.want, res)
		})
	}
}

func TestUnionSetFunc(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6, 1},
			want: []int{1, 2, 3, 4, 5, 6},
			name: "not empty",
		},
		{
			src:  []int{},
			dst:  []int{1, 3},
			want: []int{1, 3},
			name: "src is empty",
		},
		{

			src:  []int{1, 3},
			dst:  []int{},
			want: []int{1, 3},
			name: "dst is empty",
		},
		{
			src:  []int{},
			dst:  []int{},
			want: []int{},
			name: "src and dst are empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := UnionSetFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.ElementsMatch(t, tt.want, res)
		})
	}
}

func ExampleUnionSet() {
	res := UnionSet[int]([]int{1, 3, 4, 5}, []int{1, 4, 7})
	sort.Ints(res)
	fmt.Println(res)
	// Output:
	// [1 3 4 5 7]
}

func ExampleUnionSetFunc() {
	res := UnionSetFunc[int]([]int{1, 3, 4, 5}, []int{1, 4, 7}, func(src, dst int) bool {
		return src == dst
	})
	sort.Ints(res)
	fmt.Println(res)
	// Output:
	// [1 3 4 5 7]
}
