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
		},
		{
			name: "only src nil",
			dst:  []int{1, 2, 3},

			want: []int{1, 2, 3},
		},
		{
			name: "src and dst empty",
			src:  []int{},
			dst:  []int{},

			want: nil,
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
			name: "src contains all dst",
			src:  []int{1, 2, 3, 4, 5, 6},
			dst:  []int{1, 2, 3},

			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "src contains few dst",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{4, 5, 6, 7, 8},

			want: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},

			want: []int{1, 2, 3, 4, 5, 6},
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
		name  string
		src   []any
		dst   []any
		equal EqualFunc[any]
		want  []any
	}{
		{
			name: "src and dst nil",
		},
		{
			name: "only src nil",
			dst:  []any{1, 2, 3},

			want: []any{1, 2, 3},
		},
		{
			name: "src and dst empty",
			src:  []any{},
			dst:  []any{},

			want: nil,
		},
		{
			name: "only src empty",
			src:  []any{},
			dst:  []any{1, 2, 3},

			want: []any{1, 2, 3},
		},
		{
			name: "only dst empty",
			src:  []any{1, 2, 3},
			dst:  []any{},

			want: []any{1, 2, 3},
		},
		{
			name: "src contains all dst",
			src:  []any{1, 2, 3, 4, 5, 6},
			dst:  []any{1, 2, 3},

			want: []any{1, 2, 3, 4, 5, 6},
		},
		{
			name: "src contains few dst",
			src:  []any{1, 2, 3, 4, 5},
			dst:  []any{4, 5, 6, 7, 8},

			want: []any{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name: "src not contains dst",
			src:  []any{1, 2, 3},
			dst:  []any{4, 5, 6},

			want: []any{1, 2, 3, 4, 5, 6},
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
			res := UnionByFunc[any](tt.src, tt.dst, f)
			assert.Equal(t, tt.want, res)
		})
	}
}

func ExampleUnion() {
	src := []int{1, 2, 3, 4, 5}
	dst := []int{3, 4, 5, 6, 7}
	result := Union(src, dst)
	fmt.Println(result)
	//Output: [1 2 3 4 5 6 7]
}

func ExampleUnionByFunc() {
	src := []int{1, 2, 3, 4, 5}
	dst := []int{3, 4, 5, 6, 7}
	result := UnionByFunc(src, dst, func(x, y int) bool {
		return x == y
	})
	fmt.Println(result)
	//Output: [1 2 3 4 5 6 7]
}
