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

package list

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ecodeclub/ekit/internal/errs"

	"github.com/stretchr/testify/assert"
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
			name:      "add num to index left",
			list:      NewArrayListOf[int]([]int{1, 2, 3}),
			newVal:    100,
			index:     0,
			wantSlice: []int{100, 1, 2, 3},
		},
		{
			name:      "add num to index right",
			list:      NewArrayListOf[int]([]int{1, 2, 3}),
			newVal:    100,
			index:     3,
			wantSlice: []int{1, 2, 3, 100},
		},
		{
			name:      "add num to index mid",
			list:      NewArrayListOf[int]([]int{1, 2, 3}),
			newVal:    100,
			index:     1,
			wantSlice: []int{1, 100, 2, 3},
		},
		{
			name:    "add num to index -1",
			list:    NewArrayListOf[int]([]int{1, 2, 3}),
			newVal:  100,
			index:   -1,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 3, -1),
		},
		{
			name:    "add num to index OutOfRange",
			list:    NewArrayListOf[int]([]int{1, 2, 3}),
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

func TestArrayList_Append(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		newVal    []int
		wantSlice []int
	}{
		{
			name:      "append non-empty values to non-empty list",
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    []int{234, 456},
			wantSlice: []int{123, 234, 456},
		},
		{
			name:      "append empty values to non-empty list",
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    []int{},
			wantSlice: []int{123},
		},
		{
			name:      "append nil to non-empty list",
			list:      NewArrayListOf[int]([]int{123}),
			newVal:    nil,
			wantSlice: []int{123},
		},
		{
			name:      "append non-empty values to empty list",
			list:      NewArrayListOf[int]([]int{}),
			newVal:    []int{234, 456},
			wantSlice: []int{234, 456},
		},
		{
			name:      "append empty values to empty list",
			list:      NewArrayListOf[int]([]int{}),
			newVal:    []int{},
			wantSlice: []int{},
		},
		{
			name:      "append nil to empty list",
			list:      NewArrayListOf[int]([]int{}),
			newVal:    nil,
			wantSlice: []int{},
		},
		{
			name:      "append non-empty values to nil list",
			list:      NewArrayListOf[int](nil),
			newVal:    []int{234, 456},
			wantSlice: []int{234, 456},
		},
		{
			name:      "append empty values to nil list",
			list:      NewArrayListOf[int](nil),
			newVal:    []int{},
			wantSlice: []int{},
		},
		{
			name:      "append nil to nil list",
			list:      NewArrayListOf[int](nil),
			newVal:    nil,
			wantSlice: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Append(tc.newVal...)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantSlice, tc.list.AsSlice())
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
			name: "deleted",
			list: &ArrayList[int]{
				vals: []int{123, 124, 125},
			},
			index:     1,
			wantSlice: []int{123, 125},
			wantVal:   124,
		},
		{
			name: "index out of range",
			list: &ArrayList[int]{
				vals: []int{123, 100},
			},
			index:   12,
			wantErr: errs.NewErrIndexOutOfRange(2, 12),
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

// TestArrayList_Delete_Shrinkage 测试缩容
func TestArrayList_Delete_Shrink(t *testing.T) {
	testCases := []struct {
		name    string // 用例名称
		cap     int    // 原始容量
		loop    int    // 切片中元素个数
		wantCap int    // 期望缩容后的容量
	}{
		// ----- #阶段一 逻辑测试# -----
		// 只测试需求的逻辑和代码的逻辑是否一致

		// case 1: cap小于等于64，不进行缩容
		{
			name:    "case 1",
			cap:     64,
			loop:    1,
			wantCap: 64,
		},
		// case 2: cap大于2048，元素是容量的 1/2。 target:已有容量的 5/8
		{
			name:    "case 2",
			cap:     4000,
			loop:    999,
			wantCap: 2500,
		},
		// case 3: cap小于等于2048，元素是容量的四分之一。 target:缩到原本的一半
		{
			name:    "case 3",
			cap:     2048,
			loop:    300,
			wantCap: 1024,
		},
		// case 4: cap > 2048，但不满足缩容条件的例子
		{
			name:    "case 4",
			cap:     4000,
			loop:    3888,
			wantCap: 4000,
		},
		// case 5: cap <= 2048，但不满足缩容条件的例子
		{
			name:    "case 5",
			cap:     2048,
			loop:    666,
			wantCap: 2048,
		},

		// cap <= 64，但不满足缩容条件的例子
		{
			name:    "cap <= 64",
			cap:     64,
			loop:    2,
			wantCap: 64,
		},

		// ----- #阶段二 边界测试# -----
		// 测试用例边界
		// ps:测试时：
		//		（1）会默认删除一个元素，loop需要+1
		//		（2）测试结果向下取整

		// case 6: cap65
		{
			name:    "case 6",
			cap:     65,
			loop:    2,
			wantCap: 32,
		},
		// case 6-2:  cap65,loop为16
		{
			name:    "case 6-2",
			cap:     65,
			loop:    16,
			wantCap: 32,
		},
		// case 6-3:  cap130,loop为34，删除一个元素后为33，刚好不满足四分之一
		{
			name:    "case 6-3",
			cap:     130,
			loop:    34,
			wantCap: 130,
		},
		// case 7 cap2047
		{
			name:    "case 7",
			cap:     2047,
			loop:    10,
			wantCap: 1023, // 1023.5 ，向下取整
		},
		{
			name:    "case 7-1",
			cap:     2047,
			loop:    512,
			wantCap: 1023,
		},
		{
			name:    "case 7-2",
			cap:     2047,
			loop:    513, // 四分之一为511.75，513删除一个元素512刚好不满足
			wantCap: 2047,
		},
		//  case 8 : cap 2049
		{
			name:    "case 8",
			cap:     2049,
			loop:    10,
			wantCap: 1280, // 1280.625 ，向下取整
		},
		{
			name:    "case 8-1",
			cap:     2049,
			loop:    1025,
			wantCap: 1280, // 1280.625 ，向下取整
		},
		{
			name:    "case 8-2",
			cap:     2049,
			loop:    1026, // 二分之一为1024，1026删除一个元素后1025刚好不满足
			wantCap: 2049,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := NewArrayList[int](tc.cap)
			for i := 0; i < tc.loop; i++ {
				_ = list.Append(i)
			}
			_, _ = list.Delete(0)
			assert.Equal(t, tc.wantCap, list.Cap())
		})
	}
}

func TestArrayList_Len(t *testing.T) {
	testCases := []struct {
		name      string
		expectLen int
		list      *ArrayList[int]
	}{
		{
			name:      "与实际元素数相等",
			expectLen: 5,
			list: &ArrayList[int]{
				vals: make([]int, 5),
			},
		},
		{
			name:      "用户传入nil",
			expectLen: 0,
			list: &ArrayList[int]{
				vals: nil,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.list.Cap()
			assert.Equal(t, testCase.expectLen, actual)
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
			wantVal: 0,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 2, 2),
		},
		{
			name:    "index -1",
			list:    NewArrayListOf[int]([]int{123, 100}),
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
func TestArrayList_Range(t *testing.T) {
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

func TestArrayList_Set(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		index     int
		newVal    int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "set 5 by index  1",
			list:      NewArrayListOf[int]([]int{0, 1, 2, 3, 4}),
			index:     1,
			newVal:    5,
			wantSlice: []int{0, 5, 2, 3, 4},
			wantErr:   nil,
		},
		{
			name:      "index  -1",
			list:      NewArrayListOf[int]([]int{0, 1, 2, 3, 4}),
			index:     -1,
			newVal:    5,
			wantSlice: []int{},
			wantErr:   fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 5, -1),
		},
		{
			name:      "index  100",
			list:      NewArrayListOf[int]([]int{0, 1, 2, 3, 4}),
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
			assert.Equal(t, tc.wantSlice, tc.list.vals)
		})
	}
}
