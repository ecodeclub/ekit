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

func TestUnion(t *testing.T) {
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
			want: []int{},
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
			name: "src and dst empty",
			src:  []int{},
			dst:  []int{},
			want: []int{},
		},
		{
			name: "only src empty",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "only dst empty",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "src and dst not empty",
			src:  []int{1, 2, 3},
			dst:  []int{2, 3, 4},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "src and dst repeat",
			src:  []int{1, 2, 2, 3},
			dst:  []int{2, 3, 3, 4, 4},
			want: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Union[int](tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestUnionAny(t *testing.T) {
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
			want: []int{},
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
			name: "src and dst empty",
			src:  []int{},
			dst:  []int{},
			want: []int{},
		},
		{
			name: "only src empty",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "only dst empty",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "src and dst not empty",
			src:  []int{1, 2, 3},
			dst:  []int{2, 3, 4},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "src and dst repeat",
			src:  []int{1, 2, 2, 3},
			dst:  []int{2, 3, 3, 4},
			want: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := UnionByFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}

func ExampleUnion() {
	src := []int{1, 2, 3}
	dst := []int{2, 3, 4}
	union := Union(src, dst)
	fmt.Println(union)
	//Output: [1 2 3 4]
}

func ExampleUnionByFunc() {
	src := []int{1, 2, 3}
	dst := []int{2, 3, 4}
	equalFunc := func(src int, dst int) bool {
		return src == dst
	}
	union := UnionByFunc(src, dst, equalFunc)
	fmt.Println(union)
	//Output: [1 2 3 4]
}
