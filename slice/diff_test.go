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
			name: "src and dst nil",
		},
		{
			name: "测试int差级",
			src:  []int{1, 2, 3, 6},
			dst:  []int{2, 3, 8},
			want: []int{1, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Diff[int](tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
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
			name: "src and dst nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := DiffFunc[int](tt.src, tt.dst, func(src, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}
