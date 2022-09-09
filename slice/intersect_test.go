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
			want: []int{1, 3, 5},
			src:  []int{1, 3, 5, 7},
			dst:  []int{1, 3, 5},
			name: "normal test",
		},
		{
			src:  []int{},
			dst:  []int{1, 3, 5, 7},
			want: []int{},
			name: "length of src is 0",
		},
		{
			src:  []int{1, 3, 5, 5},
			dst:  []int{1, 3, 5},
			want: []int{1, 3, 5},
			name: "exist the same ele in src",
		},
		{
			src:  []int{1, 1, 3, 5, 7},
			dst:  []int{1, 3, 5, 5},
			want: []int{1, 3, 5},
			name: "exist the same ele in src and dst",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IntersectSet[int](tt.src, tt.dst)
			assert.Equal(t, true, equal[int](res, tt.want))
		})
	}
}

func TestIntersectFunc(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			want: []int{1, 3, 5},
			src:  []int{1, 3, 5, 7},
			dst:  []int{1, 3, 5},
			name: "normal test",
		},
		{
			src:  []int{},
			dst:  []int{1, 3, 5, 7},
			want: []int{},
			name: "length of src is 0",
		},
		{
			src:  []int{1, 3, 5, 5},
			dst:  []int{1, 3, 5},
			want: []int{1, 3, 5},
			name: "exist the same ele in src",
		},
		{
			src:  []int{1, 1, 3, 5, 7},
			dst:  []int{1, 3, 5, 5},
			want: []int{1, 3, 5},
			name: "exist the same ele in src and dst",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IntersectSetByFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.Equal(t, true, equal[int](res, tt.want))
		})
	}
}
