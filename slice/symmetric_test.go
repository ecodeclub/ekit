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

func TestSymmetricDiff(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			src:  []int{1, 2, 3, 4},
			dst:  []int{4, 5, 6, 1},
			want: []int{2, 3, 5, 6},
			name: "normal test",
		},
		{
			src:  []int{},
			dst:  []int{1},
			want: []int{1},
			name: "src length is 0",
		},
		{
			src:  []int{1, 3},
			dst:  []int{2, 4},
			want: []int{1, 3, 2, 4},
			name: "not exist same ele",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := SymmetricDiff[int](tt.src, tt.dst)
			assert.Equal(t, true, equal[int](res, tt.want))
		})
	}
}

func TestSymmetricDiffAny(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			src:  []int{1, 2, 3, 4},
			dst:  []int{4, 5, 6, 1},
			want: []int{2, 3, 5, 6},
			name: "normal test",
		},
		{
			src:  []int{},
			dst:  []int{1},
			want: []int{1},
			name: "src length is 0",
		},
		{
			src:  []int{1, 3},
			dst:  []int{2, 4},
			want: []int{1, 3, 2, 4},
			name: "not exist same ele",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := SymmetricDiffFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.Equal(t, true, equal[int](res, tt.want))
		})
	}
}
