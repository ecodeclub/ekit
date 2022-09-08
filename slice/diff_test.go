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

func TestDiff(t *testing.T) {
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
			name: "src and dst nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Diff[int](tt.src, tt.dst)
			assert.True(t, equal[int](res, tt.want), true)
		})
	}
}

func TestDiffAny(t *testing.T) {
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
			name: "src and dst nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := DiffFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.True(t, equal[int](res, tt.want), true)
		})
	}
}
