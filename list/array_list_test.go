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
		// 仿照这个例子，继续添加测试
		// 你需要综合考虑下标的各种可能取值
		// 往两边增加，往中间加
		// 下标可能是负数，也可能超出你的长度
		{
			name:      "index 0",
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    100,
			wantSlice: []int{100, 123},
		},
		{
			name:    "index -1",
			list:    NewArrayListOf[int]([]int{123}),
			index:   -1,
			newVal:  12,
			wantErr: newErrIndexOutOfRange(1, -1),
		},
		{
			name:    "index 10",
			list:    NewArrayListOf[int]([]int{123}),
			index:   10,
			newVal:  12,
			wantErr: newErrIndexOutOfRange(1, 10),
		},
		{
			name:      "index 1",
			list:      NewArrayListOf[int]([]int{123}),
			index:     1,
			newVal:    12,
			wantSlice: []int{123, 12},
		},
		{
			name:      "index 2",
			list:      NewArrayListOf[int]([]int{1, 2, 3, 4, 5, 6}),
			index:     2,
			newVal:    12,
			wantSlice: []int{1, 2, 12, 3, 4, 5, 6},
		},
		{
			name:      "index == len",
			list:      NewArrayListOf[int]([]int{1, 2, 3, 4, 5, 6}),
			index:     6,
			newVal:    12,
			wantSlice: []int{1, 2, 3, 4, 5, 6, 12},
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
	// 这个比较简单，只需要增加元素，然后判断一下 Append 之后是否符合预期
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		newVal    int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "append 2",
			list:      NewArrayListOf[int]([]int{1, 2, 3, 4}),
			newVal:    2,
			wantSlice: []int{1, 2, 3, 4, 2},
			wantErr:   nil,
		},
		{
			name:      "append 3",
			list:      NewArrayListOf[int]([]int{1, 2, 3, 4}),
			newVal:    3,
			wantSlice: []int{1, 2, 3, 4, 3},
			wantErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Append(tc.newVal)
			//这个可以省略 没有err 都是nil start
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			//这个可以省略 没有err 都是nil end
			assert.Equal(t, tc.wantSlice, tc.list.vals)
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
				vals: make([]int, 5),
			},
		},
		{
			name:      "用户传入nil",
			expectCap: 0,
			list: &ArrayList[int]{
				vals: nil,
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
		vals: make([]int, 0),
	}

	b.Run("Cap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			list.Cap()
		}
	})

	b.Run("Runtime cap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cap(list.vals)
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
			name:    "index -1",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   -1,
			wantErr: newErrIndexOutOfRange(2, -1),
		},
		{
			name:    "index 10",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   10,
			wantErr: newErrIndexOutOfRange(2, 10),
		},
		{
			name:      "index 1",
			list:      NewArrayListOf[int]([]int{123, 100, 150}),
			index:     1,
			wantSlice: []int{123, 150},
			wantVal:   100,
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
			name:    "index 2",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   2,
			wantErr: newErrIndexOutOfRange(2, 2),
		},
		{
			name:    "index -1",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   -1,
			wantErr: newErrIndexOutOfRange(2, -1),
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

//
func TestArrayList_Range(t *testing.T) {
	// 设计两个测试用例，用求和来作为场景
	// 一个测试用例是计算全部元素的和
	// 一个测试用例是计算元素的和，如果遇到了第一个负数，那么就中断返回
	// 测试最终断言求的和是否符合预期
	pluralErr := errors.New("value cannot be complex")
	testCases := []struct {
		name    string
		list    *ArrayList[int]
		wantVal int
		wantErr error
	}{
		{
			name:    "Range sum",
			list:    NewArrayListOf[int]([]int{1, 2, 3}),
			wantVal: 6,
		},
		{
			name:    "Range -1",
			list:    NewArrayListOf[int]([]int{-1, 2, 3}),
			wantErr: pluralErr,
		},
		{
			name: "Range nil",
			list: &ArrayList[int]{
				vals: nil,
			},
			wantVal: 0,
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var sum int
			err := tc.list.Range(func(index int, t int) error {
				if t < 0 {
					return pluralErr
				}
				sum += t
				return nil
			})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, sum)
		})
	}
}

func TestArrayList_Len(t *testing.T) {
	testCases := []struct {
		name string
		len  int
		list *ArrayList[int]
	}{
		{
			name: "len 1",
			len:  1,
			list: &ArrayList[int]{
				vals: make([]int, 1),
			},
		},
		{
			name: "len 0",
			len:  0,
			list: &ArrayList[int]{
				vals: make([]int, 0),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.list.Len()
			assert.Equal(t, testCase.len, actual)
		})
	}
}

//
func TestArrayList_Set(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		index     int
		wantSlice []int
		newVal    int
		wantErr   error
	}{
		{
			name:      "index 0",
			list:      NewArrayListOf[int]([]int{123, 100}),
			index:     0,
			wantSlice: []int{1, 100},
			newVal:    1,
		},
		{
			name:    "index -1",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   -1,
			wantErr: newErrIndexOutOfRange(2, -1),
		},
		{
			name:    "index 10",
			list:    NewArrayListOf[int]([]int{123, 100}),
			index:   10,
			wantErr: newErrIndexOutOfRange(2, 10),
		},
		{
			name:      "index 1",
			list:      NewArrayListOf[int]([]int{123, 100, 150}),
			index:     1,
			wantSlice: []int{123, 3, 150},
			newVal:    3,
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

//
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

func ExampleArrayList_Append() {
	list := NewArrayListOf[int]([]int{1, 2, 3})
	_ = list.Append(9)
	fmt.Println(list.AsSlice())
	// Output:
	// [1 2 3 9]
}

func ExampleArrayList_Get() {
	list := NewArrayListOf[int]([]int{1, 2, 3})
	val, _ := list.Get(1)
	fmt.Println(val)
	// Output:
	// 2
}

func ExampleArrayList_Delete() {
	list := NewArrayListOf[int]([]int{1, 2, 3})
	val, _ := list.Delete(1)
	fmt.Println(val)
	// Output:
	// 2
}

func ExampleArrayList_Set() {
	list := NewArrayListOf[int]([]int{1, 2, 3})
	_ = list.Set(1, 3)
	fmt.Println(list.AsSlice())
	// Output:
	// [1 3 3]
}

func ExampleArrayList_Len() {
	list := NewArrayListOf[int]([]int{1, 2, 3})
	n := list.Len()
	fmt.Println(n)
	// Output:
	// 3
}

func ExampleArrayList_Cap() {
	list := NewArrayListOf[int]([]int{1, 2, 3})
	n := list.Cap()
	fmt.Println(n)
	// Output:
	// 3
}

func ExampleArrayList_Range() {
	list := NewArrayListOf[int]([]int{1, 2, 3})
	sum := 0
	_ = list.Range(func(i int, val int) error {
		sum += val
		return nil
	})
	fmt.Println(sum)
	// Output:
	// 6
}
