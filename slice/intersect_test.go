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
			dst:  []int{1, 2, 3, 4, 5, 6},
			name: "src and dst nil",
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
			want: make([]int, 0),
			name: "src and dst nil",
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
