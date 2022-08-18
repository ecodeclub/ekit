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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentList_Add(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ConcurrentList[int]
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
			name:      "add num to index left",
			list:      NewConcurrentListOfArrayList[int]([]int{1, 2, 3}),
			newVal:    100,
			index:     0,
			wantSlice: []int{100, 1, 2, 3},
		},
		{
			name:      "add num to index right",
			list:      NewConcurrentListOfArrayList[int]([]int{1, 2, 3}),
			newVal:    100,
			index:     3,
			wantSlice: []int{1, 2, 3, 100},
		},
		{
			name:      "add num to index mid",
			list:      NewConcurrentListOfArrayList[int]([]int{1, 2, 3}),
			newVal:    100,
			index:     1,
			wantSlice: []int{1, 100, 2, 3},
		},
		{
			name:    "add num to index -1",
			list:    NewConcurrentListOfArrayList[int]([]int{1, 2, 3}),
			newVal:  100,
			index:   -1,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 3, -1),
		},
		{
			name:    "add num to index OutOfRange",
			list:    NewConcurrentListOfArrayList[int]([]int{1, 2, 3}),
			newVal:  100,
			index:   4,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 3, 4),
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
			assert.Equal(t, tc.wantSlice, tc.list.AsSlice())
		})
	}
}

//
//	func TestConcurrentList_Append(t *testing.T) {
//		// 这个比较简单，只需要增加元素，然后判断一下 Append 之后是否符合预期
//	}

func TestConcurrentList_Cap(t *testing.T) {
	testCases := []struct {
		name      string
		expectCap int
		list      *ConcurrentList[int]
	}{
		{
			name:      "与实际容量相等",
			expectCap: 3,
			list:      NewConcurrentListOfArrayList[int]([]int{1, 2, 3}),
		},
		{
			name:      "用户传入nil",
			expectCap: 0,
			list:      NewConcurrentListOfArrayList[int]([]int{}),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.list.Cap()
			assert.Equal(t, testCase.expectCap, actual)
		})
	}
}

func TestConcurrentList_Append(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ConcurrentList[int]
		newVal    int
		wantSlice []int
	}{
		{
			name:      "append 234",
			list:      NewConcurrentListOfArrayList[int]([]int{123}),
			newVal:    234,
			wantSlice: []int{123, 234},
		},
		{
			name:      "nil append 123",
			list:      NewConcurrentListOfArrayList[int](nil),
			newVal:    123,
			wantSlice: []int{123},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Append(tc.newVal)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantSlice, tc.list.AsSlice())
		})
	}
}

func TestConcurrentList_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ConcurrentList[int]
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
			list:      NewConcurrentListOfArrayList([]int{123, 100}),
			index:     0,
			wantSlice: []int{100},
			wantVal:   123,
		},
		{
			name:      "index middle",
			list:      NewConcurrentListOfArrayList([]int{123, 124, 125}),
			index:     1,
			wantSlice: []int{123, 125},
			wantVal:   124,
		},
		{
			name:    "index out of range",
			list:    NewConcurrentListOfArrayList([]int{123, 100}),
			index:   12,
			wantErr: newErrIndexOutOfRange(2, 12),
		},
		{
			name:    "index less than 0",
			list:    NewConcurrentListOfArrayList([]int{123, 100}),
			index:   -1,
			wantErr: newErrIndexOutOfRange(2, -1),
		},
		{
			name:      "index last",
			list:      NewConcurrentListOfArrayList([]int{123, 100, 101, 102, 102, 102}),
			index:     5,
			wantSlice: []int{123, 100, 101, 102, 102},
			wantVal:   102,
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
			assert.Equal(t, tc.wantSlice, tc.list.AsSlice())
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestConcurrentList_Len(t *testing.T) {
	testCases := []struct {
		name      string
		expectLen int
		list      *ConcurrentList[int]
	}{
		{
			name:      "与实际元素数相等",
			expectLen: 5,
			list:      NewConcurrentListOfArrayList([]int{1, 2, 3, 4, 5}),
		},
		{
			name:      "用户传入nil",
			expectLen: 0,
			list:      NewConcurrentListOfArrayList([]int{}),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.list.Cap()
			assert.Equal(t, testCase.expectLen, actual)
		})
	}
}

func TestConcurrentList_Get(t *testing.T) {
	testCases := []struct {
		name    string
		list    *ConcurrentList[int]
		index   int
		wantVal int
		wantErr error
	}{
		{
			name:    "index 0",
			list:    NewConcurrentListOfArrayList[int]([]int{123, 100}),
			index:   0,
			wantVal: 123,
		},
		{
			name:    "index 2",
			list:    NewConcurrentListOfArrayList[int]([]int{123, 100}),
			index:   2,
			wantVal: 0,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 2, 2),
		},
		{
			name:    "index -1",
			list:    NewConcurrentListOfArrayList[int]([]int{123, 100}),
			index:   -1,
			wantVal: 0,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 2, -1),
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
func TestConcurrentList_Range(t *testing.T) {
	// 设计两个测试用例，用求和来作为场景
	// 一个测试用例是计算全部元素的和
	// 一个测试用例是计算元素的和，如果遇到了第一个负数，那么就中断返回
	// 测试最终断言求的和是否符合预期
	testCases := []struct {
		name    string
		list    *ConcurrentList[int]
		index   int
		wantVal int
		wantErr error
	}{
		{
			name:    "计算全部元素的和",
			list:    NewConcurrentListOfArrayList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			wantVal: 55,
			wantErr: nil,
		},
		{
			name:    "测试中断",
			list:    NewConcurrentListOfArrayList([]int{1, 2, 3, 4, -5, 6, 7, 8, -9, 10}),
			wantVal: 41,
			wantErr: errors.New("index 4 is error"),
		},
		{
			name:    "测试数组为nil",
			list:    NewConcurrentListOfArrayList([]int{}),
			wantVal: 0,
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := 0
			err := tc.list.Range(func(index int, num int) error {
				if num < 0 {
					return fmt.Errorf("index %d is error", index)
				}
				result += num
				return nil
			})

			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, result)
		})
	}
}

func TestConcurrentList_AsSlice(t *testing.T) {
	vals := []int{1, 2, 3}
	a := NewConcurrentListOfArrayList[int](vals)
	slice := a.AsSlice()
	// 内容相同
	assert.Equal(t, slice, vals)
	aAddr := fmt.Sprintf("%p", vals)
	sliceAddr := fmt.Sprintf("%p", slice)
	// 但是地址不同，也就是意味着 slice 必须是一个新创建的
	assert.NotEqual(t, aAddr, sliceAddr)
}

func TestConcurrentList_Set(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ConcurrentList[int]
		index     int
		newVal    int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "set 5 by index  1",
			list:      NewConcurrentListOfArrayList([]int{0, 1, 2, 3, 4}),
			index:     1,
			newVal:    5,
			wantSlice: []int{0, 5, 2, 3, 4},
			wantErr:   nil,
		},
		{
			name:      "index  -1",
			list:      NewConcurrentListOfArrayList([]int{0, 1, 2, 3, 4}),
			index:     -1,
			newVal:    5,
			wantSlice: []int{},
			wantErr:   fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 5, -1),
		},
		{
			name:      "index  100",
			list:      NewConcurrentListOfArrayList([]int{0, 1, 2, 3, 4}),
			index:     100,
			newVal:    5,
			wantSlice: []int{},
			wantErr:   fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 5, 100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Set(tc.index, tc.newVal)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.Equal(t, tc.wantSlice, tc.list.AsSlice())
		})
	}

}
