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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnion(t *testing.T) {
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
			name: "src and dst nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Union[int](tt.src, tt.dst)
			assert.True(t, equal[int](res, tt.want), true)
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
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6, 1},
			want: []int{1, 2, 3, 4, 5, 6},
			name: "src and dst nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := UnionByFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.True(t, equal[int](res, tt.want), true)
		})
	}
}
