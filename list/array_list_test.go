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

//	func TestArrayList_Add(t *testing.T) {
//		testCases := []struct {
//			name      string
//			list      *ArrayList[int]
//			index     int
//			newVal    int
//			wantSlice []int
//			wantErr   error
//		}{
//			// 仿照这个例子，继续添加测试
//			// 你需要综合考虑下标的各种可能取值
//			// 往两边增加，往中间加
//			// 下标可能是负数，也可能超出你的长度
//			{
//				name:      "index 0",
//				list:      NewArrayListOf[int]([]int{123}),
//				newVal:    100,
//				wantSlice: []int{100, 123},
//			},
//		}
//
//		for _, tc := range testCases {
//			t.Run(tc.name, func(t *testing.T) {
//				err := tc.list.Add(tc.index, tc.newVal)
//				assert.Equal(t, tc.wantErr, err)
//				// 因为返回了 error，所以我们不用继续往下比较了
//				if err != nil {
//					return
//				}
//				assert.Equal(t, tc.wantSlice, tc.list.vals)
//			})
//		}
//	}
//
//	func TestArrayList_Append(t *testing.T) {
//		// 这个比较简单，只需要增加元素，然后判断一下 Append 之后是否符合预期
//	}

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

//	func TestArrayList_Append(t *testing.T) {
//		// 这个比较简单，只需要增加元素，然后判断一下 Append 之后是否符合预期
//	}

//	func TestArrayList_Delete(t *testing.T) {
//		testCases := []struct {
//			name      string
//			list      *ArrayList[int]
//			index     int
//			wantSlice []int
//			wantVal   int
//			wantErr   error
//		}{
//			// 仿照这个例子，继续添加测试
//			// 你需要综合考虑下标的各种可能取值
//			// 往两边增加，往中间加
//			// 下标可能是负数，也可能超出你的长度
//			{
//				name:      "index 0",
//				list:      NewArrayListOf[int]([]int{123, 100}),
//				index:     0,
//				wantSlice: []int{100},
//				wantVal:   123,
//			},
//		}
//
//		for _, tc := range testCases {
//			t.Run(tc.name, func(t *testing.T) {
//				val, err := tc.list.Delete(tc.index)
//				assert.Equal(t, tc.wantErr, err)
//				// 因为返回了 error，所以我们不用继续往下比较了
//				if err != nil {
//					return
//				}
//				assert.Equal(t, tc.wantSlice, tc.list.vals)
//				assert.Equal(t, tc.wantVal, val)
//			})
//		}
//	}
//
//	func TestArrayList_Get(t *testing.T) {
//		testCases := []struct {
//			name    string
//			list    *ArrayList[int]
//			index   int
//			wantVal int
//			wantErr error
//		}{
//			// 仿照这个例子，继续添加测试
//			// 你需要综合考虑下标的各种可能取值
//			// 往两边增加，往中间加
//			// 下标可能是负数，也可能超出你的长度
//			{
//				name:    "index 0",
//				list:    NewArrayListOf[int]([]int{123, 100}),
//				index:   0,
//				wantVal: 123,
//			},
//		}
//
//		for _, tc := range testCases {
//			t.Run(tc.name, func(t *testing.T) {
//				val, err := tc.list.Get(tc.index)
//				assert.Equal(t, tc.wantErr, err)
//				// 因为返回了 error，所以我们不用继续往下比较了
//				if err != nil {
//					return
//				}
//				assert.Equal(t, tc.wantVal, val)
//			})
//		}
//	}
func TestArrayList_Range(t *testing.T) {
	// 设计两个测试用例，用求和来作为场景
	// 一个测试用例是计算全部元素的和
	// 一个测试用例是计算元素的和，如果遇到了第一个负数，那么就中断返回
	// 测试最终断言求的和是否符合预期
	testCases := []struct {
		name    string
		list    *ArrayList[int]
		index   int
		wantVal int
		wantErr error
	}{
		{
			name: "计算全部元素的和",
			list: &ArrayList[int]{
				vals: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			wantVal: 55,
			wantErr: nil,
		},
		{
			name: "测试中断",
			list: &ArrayList[int]{
				vals: []int{1, 2, 3, 4, -5, 6, 7, 8, -9, 10},
			},
			wantVal: 41,
			wantErr: errors.New("index 4 is error"),
		},
		{
			name: "测试数组为nil",
			list: &ArrayList[int]{
				vals: nil,
			},
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

//
// func TestArrayList_Len(t *testing.T) {
//
// }
//
// func TestArrayList_Set(t *testing.T) {
//
// }
//
// func TestArrayList_AsSlice(t *testing.T) {
// 	vals := []int{1, 2, 3}
// 	a := NewArrayListOf[int](vals)
// 	slice := a.AsSlice()
// 	// 内容相同
// 	assert.Equal(t, slice, vals)
// 	aAddr := fmt.Sprintf("%p", vals)
// 	sliceAddr := fmt.Sprintf("%p", slice)
// 	// 但是地址不同，也就是意味着 slice 必须是一个新创建的
// 	assert.NotEqual(t, aAddr, sliceAddr)
// }

func TestArrayList_AsSlice(t *testing.T) {
	vals := []int{1, 2, 3}
	a := NewArrayListOf2[int](vals)
	slice := a.AsSlice()
	// 内容相同
	assert.Equal(t, slice, vals)
	aAddr := fmt.Sprintf("%p", vals)
	sliceAddr := fmt.Sprintf("%p", slice)
	// 但是地址不同，也就是意味着 slice 必须是一个新创建的
	assert.NotEqual(t, aAddr, sliceAddr)
}

//
// // 为其它所有的公开方法都加上例子
// func ExampleArrayList_Add() {
// 	list := NewArrayListOf[int]([]int{1, 2, 3})
// 	_ = list.Add(0, 9)
// 	fmt.Println(list.AsSlice())
// 	// Output:
// 	// [9 1 2 3]
// }
