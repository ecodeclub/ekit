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

func TestIntersect(t *testing.T) {
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
			name: "src nil",
			src:  nil,
			want: []int{},
		},
		{
			name: "dst nil",
			dst:  nil,
			want: []int{},
		},
		{
			name: "src and dst empty",
			src:  []int{},
			dst:  []int{},
			want: []int{},
		},
		{
			name: "src empty",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: []int{},
		},
		{
			name: "dst empty",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: []int{},
		},
		{
			name: "src and dst have intersect element",
			src:  []int{1, 2, 3},
			dst:  []int{2, 3, 4},
			want: []int{2, 3},
		},
		{
			name: "src and dst have repeat intersect element",
			src:  []int{1, 2, 2, 3},
			dst:  []int{2, 2, 3, 4},
			want: []int{2, 3},
		},
		{
			name: "src and dst not have intersect element",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Intersect[int](tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestIntersectAny(t *testing.T) {
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
			name: "src nil",
			src:  nil,
			want: []int{},
		},
		{
			name: "dst nil",
			dst:  nil,
			want: []int{},
		},
		{
			name: "src and dst empty",
			src:  []int{},
			dst:  []int{},
			want: []int{},
		},
		{
			name: "src empty",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: []int{},
		},
		{
			name: "dst empty",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: []int{},
		},
		{
			name: "src and dst have intersect element",
			src:  []int{1, 2, 3},
			dst:  []int{2, 3, 4},
			want: []int{2, 3},
		},
		{
			name: "src and dst have repeat intersect element",
			src:  []int{1, 2, 2, 3},
			dst:  []int{2, 2, 3, 4},
			want: []int{2, 3},
		},
		{
			name: "src and dst not have intersect element",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IntersectByFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}

func ExampleIntersect() {
	src := []int{1, 2, 3}
	dst := []int{2, 3, 4}
	intersect := Intersect(src, dst)
	fmt.Println(intersect)
	//Output: [2 3]
}

func ExampleIntersectByFunc() {
	src := []int{1, 2, 3}
	dst := []int{2, 3, 4}
	equalFunc := func(src int, dst int) bool {
		return src == dst
	}
	intersect := IntersectByFunc(src, dst, equalFunc)
	fmt.Println(intersect)
	//Output: [2 3]
}
