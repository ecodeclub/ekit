// Copyright 2021 ecodeclub
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

	"github.com/ecodeclub/ekit/internal/errs"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name      string
		slice     []int
		addVal    int
		index     int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "index 0",
			slice:     []int{123, 100},
			addVal:    233,
			index:     0,
			wantSlice: []int{233, 123, 100},
		},
		{
			name:      "index middle",
			slice:     []int{123, 124, 125},
			addVal:    233,
			index:     1,
			wantSlice: []int{123, 233, 124, 125},
		},
		{
			name:    "index out of range",
			slice:   []int{123, 100},
			index:   12,
			wantErr: errs.NewErrIndexOutOfRange(2, 12),
		},
		{
			name:    "index less than 0",
			slice:   []int{123, 100},
			index:   -1,
			wantErr: errs.NewErrIndexOutOfRange(2, -1),
		},
		{
			name:      "index last",
			slice:     []int{123, 100, 101, 102, 102, 102},
			addVal:    233,
			index:     5,
			wantSlice: []int{123, 100, 101, 102, 102, 233, 102},
		},
		{
			name:      "append on last",
			slice:     []int{123, 100, 101, 102, 102, 102},
			addVal:    233,
			index:     6,
			wantSlice: []int{123, 100, 101, 102, 102, 102, 233},
		},
		{
			name:    "index out of range",
			slice:   []int{123, 100, 101, 102, 102, 102},
			addVal:  233,
			index:   7,
			wantErr: errs.NewErrIndexOutOfRange(6, 7),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Add(tc.slice, tc.addVal, tc.index)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSlice, res)
		})
	}
}

func OldAdd[T any](src []T, element T, index int) ([]T, error) {
	length := len(src)
	if index < 0 || index > length {
		return nil, errs.NewErrIndexOutOfRange(length, index)
	}

	//先将src扩展一个元素
	var zeroValue T
	src = append(src, zeroValue)
	for i := len(src) - 1; i > index; i-- {
		if i-1 >= 0 {
			src[i] = src[i-1]
		}
	}
	src[index] = element
	return src, nil
}

func Benchmark_Add(b *testing.B) {

	vals := []int{1, 5, 6, 3, 7}

	b.Run("Newadd", func(b *testing.B) {
		b.ReportAllocs() // 开启内存分配报告
		for i := 0; i < b.N; i++ {
			Add(vals, 785, 3)
		}
	})
	b.Run("Oldadd", func(b *testing.B) {
		b.ReportAllocs() // 开启内存分配报告
		for i := 0; i < b.N; i++ {
			OldAdd(vals, 785, 3)
		}
	})
}
