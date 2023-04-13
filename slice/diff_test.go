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

func TestDiffSet(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			want: []int{7},
			src:  []int{1, 3, 5, 7},
			dst:  []int{1, 3, 5},
			name: "diff 1",
		},
		{
			src:  []int{1, 3, 5},
			dst:  []int{1, 3, 5, 7},
			want: []int{},
			name: "src less than dst",
		},
		{
			src:  []int{1, 3, 5, 7, 7},
			dst:  []int{1, 3, 5},
			want: []int{7},
			name: "diff deduplicate",
		},
		{
			src:  []int{1, 1, 3, 5, 7},
			dst:  []int{1, 3, 5, 5},
			want: []int{7},
			name: "dst duplicate ele",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := DiffSet[int](tt.src, tt.dst)
			assert.ElementsMatch(t, tt.want, res)
		})
	}
}

func TestDiffSetFunc(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			want: []int{7},
			src:  []int{1, 3, 5, 7},
			dst:  []int{1, 3, 5},
			name: "diff 1",
		},
		{
			src:  []int{1, 3, 5},
			dst:  []int{1, 3, 5, 7},
			want: []int{},
			name: "src less than dst",
		},
		{
			src:  []int{1, 3, 5, 7, 7},
			dst:  []int{1, 3, 5},
			want: []int{7},
			name: "diff deduplicate",
		},
		{
			src:  []int{1, 1, 3, 5, 7},
			dst:  []int{1, 3, 5, 5},
			want: []int{7},
			name: "dst duplicate ele",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := DiffSetFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.ElementsMatch(t, tt.want, res)
		})
	}
}

func ExampleDiffSet() {
	res := DiffSet[int]([]int{1, 3, 2, 2, 4}, []int{3, 4, 5, 6})
	sort.Ints(res)
	fmt.Println(res)
	// Output:
	// [1 2]
}

func ExampleDiffSetFunc() {
	res := DiffSetFunc[int]([]int{1, 3, 2, 2, 4}, []int{3, 4, 5, 6}, func(src, dst int) bool {
		return src == dst
	})
	fmt.Println(res)
	// Output:
	// [1 2]
}
