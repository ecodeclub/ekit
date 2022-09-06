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
	"strconv"
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

			want: []int{},
		},
		{
			name: "only src nil",
			dst:  []int{1, 2, 3},

			want: []int{},
		},
		{
			name: "only dst nil",
			src:  []int{1, 2, 3},

			want: []int{},
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

			want: []int{},
		},
		{
			name: "only dst empty",
			src:  []int{1, 2, 3},
			dst:  []int{},

			want: []int{},
		},
		{
			name: "src contains all dst",
			src:  []int{1, 2, 3, 4, 5, 6},
			dst:  []int{1, 2, 3},

			want: []int{1, 2, 3},
		},
		{
			name: "src equal to dst",
			src:  []int{1, 2, 3},
			dst:  []int{1, 2, 3},

			want: []int{1, 2, 3},
		},
		{
			name: "src contains few dst",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{4, 5, 6, 7, 8},

			want: []int{4, 5},
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},

			want: []int{},
		},
		{
			name: "duplicates in the intersection",
			src:  []int{1, 2, 3, 3, 4},
			dst:  []int{1, 3, 3},

			want: []int{1, 3, 3},
		},
		{
			name: "duplicates in src but not in intersection",
			src:  []int{1, 2, 3, 3, 4},
			dst:  []int{1, 3},

			want: []int{1, 3},
		},
		{
			name: "duplicates in dst but not in intersection",
			src:  []int{1, 2, 3, 4},
			dst:  []int{1, 3, 3},

			want: []int{1, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Intersect[int](tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestIntersectByFunc(t *testing.T) {
	tests := []struct {
		name  string
		src   []any
		dst   []any
		equal EqualFunc[any]
		want  []any
	}{
		{
			name: "src and dst nil",

			want: []any{},
		},
		{
			name: "only src nil",
			dst:  []any{1, 2, 3},

			want: []any{},
		},
		{
			name: "only dst nil",
			src:  []any{1, 2, 3},

			want: []any{},
		},
		{
			name: "src and dst empty",
			src:  []any{},
			dst:  []any{},

			want: []any{},
		},
		{
			name: "only src empty",
			src:  []any{},
			dst:  []any{1, 2, 3},

			want: []any{},
		},
		{
			name: "only dst empty",
			src:  []any{1, 2, 3},
			dst:  []any{},

			want: []any{},
		},
		{
			name: "src contains all dst",
			src:  []any{1, 2, 3, 4, 5, 6},
			dst:  []any{1, 2, 3},

			want: []any{1, 2, 3},
		},
		{
			name: "src equal to dst",
			src:  []any{1, 2, 3},
			dst:  []any{1, 2, 3},

			want: []any{1, 2, 3},
		},
		{
			name: "src contains few dst",
			src:  []any{1, 2, 3, 4, 5},
			dst:  []any{4, 5, 6, 7, 8},

			want: []any{4, 5},
		},
		{
			name: "src not contains dst",
			src:  []any{1, 2, 3},
			dst:  []any{4, 5, 6},

			want: []any{},
		},
		{
			name: "duplicates in the intersection",
			src:  []any{1, 2, 3, 3, 4},
			dst:  []any{1, 3, 3},

			want: []any{1, 3, 3},
		},
		{
			name: "duplicates in src but not in intersection",
			src:  []any{1, 2, 3, 3, 4},
			dst:  []any{1, 3},

			want: []any{1, 3},
		},
		{
			name: "duplicates in dst but not in intersection",
			src:  []any{1, 2, 3, 4},
			dst:  []any{1, 3, 3},

			want: []any{1, 3},
		},
		{
			name: "src int and dst string",
			src:  []any{1, 2, 3},
			dst:  []any{"1"},
			equal: func(x, y any) bool {
				xVal := strconv.Itoa(x.(int))
				return xVal == y
			},

			want: []any{1},
		},
		{
			name: "equal panic",
			src:  []any{1, 2, 3},
			dst:  []any{4, 5, 6},
			equal: func(x, y any) bool {
				panic("panic test")
			},

			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func(src, dst any) bool {
				return src == dst
			}
			if tt.equal != nil {
				f = tt.equal
			}
			res := IntersectByFunc[any](tt.src, tt.dst, f)
			assert.Equal(t, tt.want, res)
		})
	}
}

func ExampleIntersect() {
	src := []int{1, 2, 3, 4, 5}
	dst := []int{3, 4, 5, 6, 7}
	result := Intersect(src, dst)
	fmt.Println(result)
	// Output: [3 4 5]
}

func ExampleIntersectByFunc() {
	src := []int{1, 2, 3, 4, 5}
	dst := []int{3, 4, 5, 6, 7}
	result := IntersectByFunc(src, dst, func(x, y int) bool {
		return x == y
	})
	fmt.Println(result)
	// Output: [3 4 5]
}
