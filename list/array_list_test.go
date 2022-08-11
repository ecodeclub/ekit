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
		// 仿照这个例子，继续添加测试
		// 你需要综合考虑下标的各种可能取值
		// 往两边增加，往中间加
		// 下标可能是负数，也可能超出你的长度
		{
			name:      "index 0",
			index:     0,
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    100,
			wantSlice: []int{100, 123},
		},
		{
			name:      "index 1",
			index:     1,
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    100,
			wantSlice: []int{123, 100},
		},
		{
			name:    "index 2",
			index:   2,
			list:    NewArrayListOf[int]([]int{123}),
			wantErr: newErrIndexOutOfRange(1, 2),
		},
		{
			name:    "index -1",
			index:   -1,
			list:    NewArrayListOf[int]([]int{123}),
			wantErr: newErrIndexOutOfRange(1, -1),
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
			assert.Equal(t, tc.wantSlice, tc.list.vals)
		})
	}
}

func TestArrayList_Append(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		newVal    int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "index 0",
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    100,
			wantSlice: []int{123, 100},
		},
		{
			name:      "index 1",
			list:      NewArrayListOf[int]([]int{123, 456}),
			newVal:    789,
			wantSlice: []int{123, 456, 789},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Append(tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSlice, tc.list.vals)
		})
	}

}

func TestArrayList_Cap(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		wantSlice int
		wantErr   error
	}{
		{
			name:      "index 0",
			list:      NewArrayList[int](10),
			wantSlice: 10,
		},
		{
			name:      "index 1",
			list:      NewArrayList[int](20),
			wantSlice: 20,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.list.Cap()

			assert.Equal(t, tc.wantSlice, res)
		})
	}
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
		{
			name:      "index 0",
			list:      NewArrayListOf[int]([]int{123, 100}),
			index:     0,
			wantSlice: []int{100},
			wantVal:   123,
		},
		{
			name:    "index -1",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   -1,
			wantErr: newErrIndexOutOfRange(2, -1),
		},
		{
			name:    "index 5",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   5,
			wantErr: newErrIndexOutOfRange(2, 5),
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
			assert.Equal(t, tc.wantSlice, tc.list.vals)
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
			wantErr: newErrIndexOutOfRange(2, -1),
		},
		{
			name:    "index 5",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   5,
			wantErr: newErrIndexOutOfRange(2, 5),
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
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		index     int
		newVal    int
		wantSlice []int
		wantVal   int
		wantErr   error
	}{
		{
			name:    "index 0",
			list:    NewArrayListOf[int]([]int{1, 2, 3, 4, 5}),
			index:   0,
			wantVal: 15,
		},
		{
			name:    "index 1",
			list:    NewArrayListOf[int]([]int{1, -1, 3, 4, 5}),
			index:   1,
			wantErr: newErrIndexOutOfRange(5, 1),
		},
	}
	var res = 0
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Range(func(index int, t int) error {
				if t < 0 {
					return newErrIndexOutOfRange(5, index)
				}
				res += t
				return nil
			})
			assert.Equal(t, tc.wantErr, err)
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, res)
		})
	}
}

func TestArrayList_Len(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		index     int
		newVal    int
		wantSlice []int
		wantVal   int
		wantErr   error
	}{
		// 仿照这个例子，继续添加测试
		// 你需要综合考虑下标的各种可能取值
		// 往两边增加，往中间加
		// 下标可能是负数，也可能超出你的长度
		{
			name:    "index 0",
			list:    NewArrayListOf[int]([]int{123, 456}),
			wantVal: 2,
		},
		{
			name:    "index 1",
			list:    NewArrayListOf[int]([]int{}),
			wantVal: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val := tc.list.Len()
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestArrayList_Set(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		index     int
		newVal    int
		wantSlice []int
		wantErr   error
	}{
		// 仿照这个例子，继续添加测试
		// 你需要综合考虑下标的各种可能取值
		// 往两边增加，往中间加
		// 下标可能是负数，也可能超出你的长度
		{
			name:      "index 0",
			index:     0,
			list:      NewArrayListOf[int]([]int{123, 456}),
			newVal:    100,
			wantSlice: []int{100, 456},
		},
		{
			name:    "index 1",
			index:   1,
			list:    NewArrayListOf[int]([]int{123}),
			newVal:  100,
			wantErr: newErrIndexOutOfRange(1, 1),
		},
		{
			name:    "index -1",
			index:   -1,
			newVal:  100,
			list:    NewArrayListOf[int]([]int{123}),
			wantErr: newErrIndexOutOfRange(1, -1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Set(tc.index, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSlice, tc.list.vals)
		})
	}

}

func TestArrayList_AsSlice(t *testing.T) {
	vals := []int{1, 2, 3}
	a := NewArrayListOf[int](vals)
	slice := a.AsSlice()
	// 内容相同
	assert.Equal(t, slice, vals)
	aAddr := fmt.Sprintf("%p", vals)
	sliceAddr := fmt.Sprintf("%p", slice)
	// 但是地址不同，也就是意味着 slice 必须是一个新创建的
	assert.NotEqual(t, aAddr, sliceAddr)
}

// 为其它所有的公开方法都加上例子
func ExampleArrayList_Add() {
	list := NewArrayListOf[int]([]int{1, 2, 3})
	_ = list.Add(0, 9)
	fmt.Println(list.AsSlice())
	// Output:
	// [9 1 2 3]
}
