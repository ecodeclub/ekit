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

package list

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayList_Add(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		index     int
		newVal    int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "index 0",
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    100,
			wantSlice: []int{100, 123},
		},
		{
			name:      "insert end",
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    100,
			index:     1,
			wantSlice: []int{123, 100},
		},
		{
			name:      "index -1",
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    100,
			index:     -1,
			wantSlice: []int{123},
			wantErr:   newErrIndexOutOfRange(1, -1),
		},
		{
			name:      "out of length index",
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    100,
			index:     100,
			wantSlice: []int{123},
			wantErr:   newErrIndexOutOfRange(1, 100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Add(tc.index, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSlice, tc.list.data)
		})
	}
}

func TestArrayList_Append(t *testing.T) {
	// 这个比较简单，只需要增加元素，然后判断一下 Append 之后是否符合预期
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		newVal    int
		wantSlice []int
		size      int
	}{
		{
			name:      "newVal 100",
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    100,
			wantSlice: []int{123, 100},
			size:      2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Append(tc.newVal)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSlice, tc.list.data)
			assert.Equal(t, tc.size, tc.list.size)
		})
	}
}

func TestArrayList_Cap(t *testing.T) {
	testCases := []struct {
		name      string
		expectCap int
		list      *ArrayList[int]
	}{
		{
			name:      "与实际容量相等",
			expectCap: 5,
			list: &ArrayList[int]{
				data: make([]int, 5),
			},
		},
		{
			name:      "用户传入nil",
			expectCap: 0,
			list: &ArrayList[int]{
				data: nil,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.list.Cap()
			assert.Equal(t, testCase.expectCap, actual)
		})
	}
}

func BenchmarkArrayList_Cap(b *testing.B) {
	list := &ArrayList[int]{
		data: make([]int, 0),
	}

	b.Run("Cap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			list.Cap()
		}
	})

	b.Run("Runtime cap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cap(list.data)
		}
	})
}

func TestArrayList_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		index     int
		wantSlice []int
		wantVal   int
		wantErr   error
	}{
		// 仿照这个例子，继续添加测试
		// 你需要综合考虑下标的各种可能取值
		// 往两边增加，往中间加
		// 下标可能是负数，也可能超出你的长度
		{
			name:      "index 0",
			list:      NewArrayListOf[int]([]int{123, 100}),
			index:     0,
			wantSlice: []int{100},
			wantVal:   123,
		},
		{
			name:      "index 1",
			list:      NewArrayListOf[int]([]int{123, 100}),
			index:     1,
			wantSlice: []int{123},
			wantVal:   100,
		},
		{
			name:      "index -1",
			list:      NewArrayListOf[int]([]int{123, 100}),
			index:     -1,
			wantSlice: []int{123, 100},
			wantVal:   0,
			wantErr:   newErrIndexOutOfRange(2, -1),
		},
		{
			name:      "index 100",
			list:      NewArrayListOf[int]([]int{123, 100}),
			index:     100,
			wantSlice: []int{123, 100},
			wantVal:   0,
			wantErr:   newErrIndexOutOfRange(2, 100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.list.Delete(tc.index)
			assert.Equal(t, tc.wantErr, err)
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSlice, tc.list.data)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestArrayList_Get(t *testing.T) {
	testCases := []struct {
		name    string
		list    *ArrayList[int]
		index   int
		wantVal int
		wantErr error
	}{
		// 仿照这个例子，继续添加测试
		// 你需要综合考虑下标的各种可能取值
		// 往两边增加，往中间加
		// 下标可能是负数，也可能超出你的长度
		{
			name:    "index 0",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   0,
			wantVal: 123,
		},
		{
			name:    "index -1",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   -1,
			wantVal: 0,
			wantErr: newErrIndexOutOfRange(2, -1),
		},
		{
			name:    "index 100",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   100,
			wantVal: 0,
			wantErr: newErrIndexOutOfRange(2, 100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.list.Get(tc.index)
			assert.Equal(t, tc.wantErr, err)
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestArrayList_Range(t *testing.T) {
	// 设计两个测试用例，用求和来作为场景
	// 一个测试用例是计算全部元素的和
	// 一个测试用例是计算元素的和，如果遇到了第一个负数，那么就中断返回
	// 测试最终断言求的和是否符合预期
	var testCases = []struct {
		name    string
		list    *ArrayList[int]
		wantVal int
	}{
		{
			name:    "sum",
			list:    NewArrayListOf([]int{123, 100}),
			wantVal: 223,
		},
		{
			name:    "reduce",
			list:    NewArrayListOf([]int{123, 100, -1, 20}),
			wantVal: 223,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var sum int
			err := tc.list.Range(func(index int, t int) error {
				if t < 0 {
					return errors.New("negative number")
				}
				sum += t
				return nil
			})
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, sum)
		})
	}

}

//func TestArrayList_Len(t *testing.T) {
//
//}
//
//func TestArrayList_Set(t *testing.T) {
//
//}

func TestArrayList_AsSlice(t *testing.T) {
	data := []int{1, 2, 3}
	a := NewArrayListOf[int](data)
	slice := a.AsSlice()
	// 内容相同
	assert.Equal(t, slice, data)
	aAddr := fmt.Sprintf("%p", data)
	sliceAddr := fmt.Sprintf("%p", slice)
	// 但是地址不同，也就是意味着 slice 必须是一个新创建的
	assert.NotEqual(t, aAddr, sliceAddr)
}
