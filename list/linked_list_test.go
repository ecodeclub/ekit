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

	"github.com/gotomicro/ekit/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"math/rand"
	"testing"
)

// 使用无界链表，不考虑容量限制，测试Add行为本身是否符合预期
func TestLinkedList_Add(t *testing.T) {
	testCases := []struct {
		name           string
		list           *LinkedList[int]
		index          int
		newVal         int
		wantLinkedList *LinkedList[int]
		wantErr        error
	}{
		{
			name:           "add num to index left",
			list:           NewLinkedListOf[int]([]int{1, 2, 3}),
			newVal:         100,
			index:          0,
			wantLinkedList: NewLinkedListOf[int]([]int{100, 1, 2, 3}),
		},
		{
			name:           "add num to index left1",
			list:           NewLinkedListOf([]int{1, 2, 3, 44, 55, 66, 77}),
			newVal:         100,
			index:          3,
			wantLinkedList: NewLinkedListOf([]int{1, 2, 3, 100, 44, 55, 66, 77}),
		},
		{
			name:           "add num to index right",
			list:           NewLinkedListOf([]int{1, 2, 3}),
			newVal:         100,
			index:          3,
			wantLinkedList: NewLinkedListOf([]int{1, 2, 3, 100}),
		},
		{
			name:           "add num to index right1",
			list:           NewLinkedListOf([]int{1, 2, 3, 44, 55, 66, 77}),
			newVal:         100,
			index:          5,
			wantLinkedList: NewLinkedListOf([]int{1, 2, 3, 44, 55, 100, 66, 77}),
		},
		{
			name:           "add num to index mid",
			list:           NewLinkedListOf[int]([]int{1, 2, 3}),
			newVal:         100,
			index:          1,
			wantLinkedList: NewLinkedListOf([]int{1, 100, 2, 3}),
		},
		{
			name:    "add num to index -1",
			list:    NewLinkedListOf[int]([]int{1, 2, 3}),
			newVal:  100,
			index:   -1,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 3, -1),
		},
		{
			name:    "add num to index OutOfRange",
			list:    NewLinkedListOf[int]([]int{1, 2, 3}),
			newVal:  100,
			index:   4,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 3, 4),
		},
		{
			name:           "add num to index 0",
			list:           NewLinkedListOf[int]([]int{}),
			newVal:         100,
			index:          0,
			wantErr:        nil,
			wantLinkedList: NewLinkedListOf([]int{100}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Add(tc.index, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLinkedList.AsSlice(), tc.list.AsSlice())
		})
	}
}

// 测试有界链表与容量边界相关的Add行为是否符合预期
func TestBoundedLinkedList_Add(t *testing.T) {
	testCases := []struct {
		name           string
		list           *LinkedList[int]
		index          int
		newVal         int
		wantLinkedList *LinkedList[int]
		wantErr        error
	}{
		{
			name:           "under capacity, add to head",
			list:           NewLinkedListOf[int]([]int{1, 2, 3}, WithCapacityOption[int](10)),
			newVal:         100,
			index:          0,
			wantLinkedList: NewLinkedListOf[int]([]int{100, 1, 2, 3}),
		},
		{
			name:           "under capacity, add in middle",
			list:           NewLinkedListOf[int]([]int{1, 2, 3}, WithCapacityOption[int](10)),
			newVal:         100,
			index:          1,
			wantLinkedList: NewLinkedListOf[int]([]int{1, 100, 2, 3}),
		},
		{
			name:           "under capacity, add to tail",
			list:           NewLinkedListOf[int]([]int{1, 2, 3}, WithCapacityOption[int](10)),
			newVal:         100,
			index:          3,
			wantLinkedList: NewLinkedListOf[int]([]int{1, 2, 3, 100}),
		},
		{
			name:    "out of capacity, add to head",
			list:    NewLinkedListOf([]int{1, 2, 3}, WithCapacityOption[int](0)),
			newVal:  100,
			index:   0,
			wantErr: errs.ErrOutOfCapacity,
		},
		{
			name:    "out of capacity, add in middle",
			list:    NewLinkedListOf([]int{1, 2, 3}, WithCapacityOption[int](0)),
			newVal:  100,
			index:   1,
			wantErr: errs.ErrOutOfCapacity,
		},
		{
			name:    "out of capacity, add to tail",
			list:    NewLinkedListOf([]int{1, 2, 3}, WithCapacityOption[int](0)),
			newVal:  100,
			index:   2,
			wantErr: errs.ErrOutOfCapacity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Add(tc.index, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLinkedList.AsSlice(), tc.list.AsSlice())
		})
	}
}

func TestLinkedList_Delete(t *testing.T) {
	testCases := []struct {
		name           string
		list           *LinkedList[int]
		wantLinkedList *LinkedList[int]
		delVal         int
		index          int
		wantErr        error
	}{
		{
			name:    "delete num to index -1",
			list:    NewLinkedListOf[int]([]int{1, 2, 3}),
			index:   -1,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 3, -1),
		},
		{
			name:    "delete beyond length index 99",
			list:    NewLinkedListOf[int]([]int{1, 2, 3}),
			index:   99,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 3, 99),
		},
		{
			name:    "delete beyond length index 3",
			list:    NewLinkedListOf[int]([]int{1, 2, 3}),
			index:   3,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 3, 3),
		},
		{
			name:    "delete empty node",
			list:    NewLinkedListOf[int]([]int{}),
			index:   3,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 0, 3),
		},
		{
			name:           "delete num to index 0",
			list:           NewLinkedListOf[int]([]int{1, 2, 3}),
			index:          0,
			delVal:         1,
			wantLinkedList: NewLinkedListOf([]int{2, 3}),
		},
		{
			name:           "delete num to index by tail",
			list:           NewLinkedListOf[int]([]int{1, 2, 3, 4, 5}),
			index:          4,
			delVal:         5,
			wantLinkedList: NewLinkedListOf([]int{1, 2, 3, 4}),
		},
		{
			name:           "delete num to index 1",
			list:           NewLinkedListOf[int]([]int{11, 22, 33, 44, 55}),
			index:          1,
			delVal:         22,
			wantLinkedList: NewLinkedListOf([]int{11, 33, 44, 55}),
		},
		{
			name:           "deleting an element with only one",
			list:           NewLinkedListOf[int]([]int{888}),
			index:          0,
			delVal:         888,
			wantLinkedList: NewLinkedListOf([]int{}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			delVal, err := tc.list.Delete(tc.index)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
			} else {
				assert.Equal(t, tc.delVal, delVal)
				assert.Equal(t, tc.wantLinkedList.AsSlice(), tc.list.AsSlice())
			}
		})
	}
}

// 使用无界链表，不考虑容量限制，测试Append行为本身是否符合预期
func TestLinkedList_Append(t *testing.T) {
	testCases := []struct {
		name           string
		list           *LinkedList[int]
		newVal         []int
		wantLinkedList *LinkedList[int]
		wantErr        error
	}{
		{
			name:           "append non-empty values to non-empty list",
			list:           NewLinkedListOf[int]([]int{123}),
			newVal:         []int{234, 456},
			wantLinkedList: NewLinkedListOf[int]([]int{123, 234, 456}),
		},
		{
			name:           "append empty values to non-empty list",
			list:           NewLinkedListOf[int]([]int{123}),
			newVal:         []int{},
			wantLinkedList: NewLinkedListOf[int]([]int{123}),
		},
		{
			name:           "append nil to non-empty list",
			list:           NewLinkedListOf[int]([]int{123}),
			newVal:         nil,
			wantLinkedList: NewLinkedListOf[int]([]int{123}),
		},
		{
			name:           "append non-empty values to empty list",
			list:           NewLinkedListOf[int]([]int{}),
			newVal:         []int{234, 456},
			wantLinkedList: NewLinkedListOf[int]([]int{234, 456}),
		},
		{
			name:           "append empty values to empty list",
			list:           NewLinkedListOf[int]([]int{}),
			newVal:         []int{},
			wantLinkedList: NewLinkedListOf[int]([]int{}),
		},
		{
			name:           "append nil to empty list",
			list:           NewLinkedListOf[int]([]int{}),
			newVal:         nil,
			wantLinkedList: NewLinkedListOf[int]([]int{}),
		},
		{
			name:           "append non-empty values to nil list",
			list:           NewLinkedListOf[int](nil),
			newVal:         []int{234, 456},
			wantLinkedList: NewLinkedListOf[int]([]int{234, 456}),
		},
		{
			name:           "append empty values to nil list",
			list:           NewLinkedListOf[int](nil),
			newVal:         []int{},
			wantLinkedList: NewLinkedListOf[int]([]int{}),
		},
		{
			name:           "append nil to nil list",
			list:           NewLinkedListOf[int](nil),
			newVal:         nil,
			wantLinkedList: NewLinkedListOf[int]([]int{}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Append(tc.newVal...)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantLinkedList.AsSlice(), tc.list.AsSlice())
		})
	}
}

// 测试有界链表与容量边界相关的Append行为是否符合预期
func TestBoundedLinkedList_Append(t *testing.T) {
	testCases := []struct {
		name           string
		list           *LinkedList[int]
		newVal         []int
		wantLinkedList *LinkedList[int]
		wantErr        error
	}{
		{
			name:           "under capacity",
			list:           NewLinkedListOf[int]([]int{1, 2, 3}, WithCapacityOption[int](10)),
			newVal:         []int{4, 5},
			wantLinkedList: NewLinkedListOf[int]([]int{1, 2, 3, 4, 5}, WithCapacityOption[int](10)),
		},
		{
			name:           "out of capacity",
			list:           NewLinkedListOf[int]([]int{1, 2, 3}, WithCapacityOption[int](3)),
			newVal:         []int{4, 5},
			wantLinkedList: NewLinkedListOf[int]([]int{1, 2, 3}, WithCapacityOption[int](3)),
			wantErr:        errs.ErrOutOfCapacity,
		},
		{
			name:           "out of capacity partially",
			list:           NewLinkedListOf[int]([]int{1, 2, 3}, WithCapacityOption[int](4)),
			newVal:         []int{4, 5, 6},
			wantLinkedList: NewLinkedListOf[int]([]int{1, 2, 3}, WithCapacityOption[int](4)),
			wantErr:        errs.ErrOutOfCapacity,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Append(tc.newVal...)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantLinkedList.AsSlice(), tc.list.AsSlice())
		})
	}
}

func TestNewLinkedList_Default(t *testing.T) {
	l := NewLinkedList[int]()
	assert.Equal(t, true, l.IsBoundless())
	assert.Equal(t, 0, l.Cap())
	assert.Equal(t, 0, l.Len())
	// 首尾相连的双向链表
	assert.Equal(t, l.tail, l.head.next)
	assert.Equal(t, l.tail, l.head.prev)
	assert.Equal(t, l.head, l.tail.prev)
	assert.Equal(t, l.head, l.tail.next)
}

func TestNewLinkedList_WithCapacityOption(t *testing.T) {
	testCases := []struct {
		name        string
		capacity    int
		isBoundless bool
	}{
		{
			name:        "boundless",
			capacity:    0,
			isBoundless: true,
		},
		{
			name:        "boundless",
			capacity:    10,
			isBoundless: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLinkedList[int](WithCapacityOption[int](tc.capacity))
			assert.Equal(t, tc.isBoundless, l.IsBoundless())
			assert.Equal(t, tc.capacity, l.Cap())
			assert.Equal(t, 0, l.Len())
		})
	}
}

func TestNewLinkedListOf(t *testing.T) {
	testCases := []struct {
		name        string
		slice       []int
		wantedSlice []int
	}{
		{
			name:        "nil",
			slice:       nil,
			wantedSlice: []int{},
		},
		{
			name:        "vacant",
			slice:       []int{},
			wantedSlice: []int{},
		},
		{
			name:        "single",
			slice:       []int{1},
			wantedSlice: []int{1},
		},
		{
			name:        "normal",
			slice:       []int{1, 2, 3},
			wantedSlice: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		list := NewLinkedListOf(tc.slice)
		// 在这里断言你的元素，可以利用 Get 方法，也可以直接用 AsSlice 来断言
		assert.Equal(t, tc.wantedSlice, list.AsSlice())
	}
}

func TestNewLinkedListOf_WithCapacity(t *testing.T) {
	testCases := []struct {
		name         string
		capacity     int
		slice        []int
		wantCapacity int
		wantedSlice  []int
	}{
		{
			name:         "normal",
			capacity:     10,
			slice:        []int{1, 2, 3},
			wantedSlice:  []int{1, 2, 3},
			wantCapacity: 10,
		},
		{
			name:         "capacity less than slice length",
			capacity:     2,
			slice:        []int{1, 2, 3},
			wantedSlice:  []int{1, 2, 3},
			wantCapacity: 3,
		},
	}

	for _, tc := range testCases {
		list := NewLinkedListOf(tc.slice, WithCapacityOption[int](tc.capacity))
		// 在这里断言你的元素，可以利用 Get 方法，也可以直接用 AsSlice 来断言
		assert.Equal(t, tc.wantedSlice, list.AsSlice())
		assert.Equal(t, tc.wantCapacity, list.Cap())
	}
}

func TestLinkedList_AsSlice(t *testing.T) {
	vals := []int{1, 2, 3}
	a := NewLinkedListOf[int](vals)
	slice := a.AsSlice()
	// 内容相同
	assert.Equal(t, slice, vals)
	aAddr := fmt.Sprintf("%p", vals)
	sliceAddr := fmt.Sprintf("%p", slice)
	// 但是地址不同，也就是意味着 slice 必须是一个新创建的
	assert.NotEqual(t, aAddr, sliceAddr)
}

func TestLinkedList_CapLen(t *testing.T) {
	testCases := []struct {
		name    string
		l       *LinkedList[int]
		wantCap int
		wantLen int
	}{
		{
			name:    "boundless",
			l:       NewLinkedListOf[int]([]int{1, 2, 3}),
			wantCap: 0,
			wantLen: 3,
		},
		{
			name:    "bounded",
			l:       NewLinkedListOf[int]([]int{1, 2, 3}, WithCapacityOption[int](5)),
			wantCap: 5,
			wantLen: 3,
		},
	}
	for _, tc := range testCases {
		t.Run(t.Name(), func(t *testing.T) {
			assert.Equal(t, tc.wantCap, tc.l.Cap())
			assert.Equal(t, tc.wantLen, tc.l.Len())
			require.NoError(t, tc.l.Append(100))
			assert.Equal(t, tc.wantCap, tc.l.Cap())
			assert.Equal(t, tc.wantLen+1, tc.l.Len())
		})
	}
}

func TestLinkedList_Get(t *testing.T) {
	tests := []struct {
		name    string
		list    *LinkedList[int]
		index   int
		wantVal int
		wantErr error
	}{
		{
			name:    "get left",
			list:    NewLinkedListOf([]int{1, 2, 3, 4, 5}),
			index:   0,
			wantVal: 1,
		},
		{
			name:    "get right",
			list:    NewLinkedListOf([]int{1, 2, 3, 4, 5}),
			index:   4,
			wantVal: 5,
		},
		{
			name:    "get middle",
			list:    NewLinkedListOf([]int{1, 2, 3, 4, 5}),
			index:   2,
			wantVal: 3,
		},
		{
			name:    "over left",
			list:    NewLinkedListOf([]int{1, 2, 3, 4, 5}),
			index:   -1,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 5, -1),
		},
		{
			name:    "over right",
			list:    NewLinkedListOf([]int{1, 2, 3, 4, 5}),
			index:   5,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 5, 5),
		},
		{
			name:    "empty list",
			list:    NewLinkedListOf([]int{}),
			index:   0,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 0, 0),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			get, err := tc.list.Get(tc.index)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantVal, get)
		})
	}
}

func TestLinkedList_Range(t *testing.T) {
	testCases := []struct {
		name    string
		list    *LinkedList[int]
		wantVal int
		wantErr error
	}{
		{
			name:    "计算全部元素的和",
			list:    NewLinkedListOf([]int{1, 2, 3, 4, 5}),
			wantVal: 15,
			wantErr: nil,
		},
		{
			name:    "测试中断",
			list:    NewLinkedListOf([]int{1, 2, 3, 4, -5, 6, 7, 8, -9, 10}),
			wantErr: errors.New("index 4 is error"),
		},
		{
			name:    "测试数组为nil",
			list:    NewLinkedListOf([]int{}),
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

func TestLinkedList_Set(t *testing.T) {
	testCases := []struct {
		name           string
		list           *LinkedList[int]
		wantLinkedList *LinkedList[int]
		index          int
		setVal         int
		wantErr        error
	}{
		{
			name:    "set num to index -1",
			list:    NewLinkedListOf[int]([]int{1, 2, 3}),
			index:   -1,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 3, -1),
		},
		{
			name:    "set beyond length index 99",
			list:    NewLinkedListOf[int]([]int{1, 2, 3}),
			index:   99,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 3, 99),
		},
		{
			name:    "set empty node",
			list:    NewLinkedListOf[int]([]int{}),
			index:   3,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 0, 3),
		},
		{
			name:           "set num to index 3",
			list:           NewLinkedListOf[int]([]int{11, 22, 33, 44, 55}),
			index:          2,
			setVal:         999,
			wantLinkedList: NewLinkedListOf([]int{11, 22, 999, 44, 55}),
		},
		{
			name:           "set num to head",
			list:           NewLinkedListOf[int]([]int{11, 22, 33, 44, 55}),
			index:          0,
			setVal:         -200,
			wantLinkedList: NewLinkedListOf([]int{-200, 22, 33, 44, 55}),
		},
		{
			name:           "set num to tail",
			list:           NewLinkedListOf[int]([]int{-11, 22, -33, 44, -55, 999, -888}),
			index:          6,
			setVal:         888,
			wantLinkedList: NewLinkedListOf([]int{-11, 22, -33, 44, -55, 999, 888}),
		},
		{
			name:    "index == len(*node)",
			list:    NewLinkedListOf[int]([]int{-11, 22, -33, 44, -55, 999, -888}),
			index:   7,
			setVal:  888,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 7, 7),
		},
		{
			name:    "len(*node) == 0",
			list:    NewLinkedListOf[int]([]int{}),
			index:   0,
			setVal:  888,
			wantErr: fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", 0, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Set(tc.index, tc.setVal)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
			} else {
				assert.Equal(t, tc.wantLinkedList, tc.list)
			}
		})
	}
}

func BenchmarkLinkedList_Add(b *testing.B) {
	l := NewLinkedListOf[int]([]int{1, 2, 3})
	testCase := make([]int, 0, b.N)
	for i := 1; i <= b.N; i++ {
		testCase = append(testCase, rand.Intn(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Add(testCase[i], testCase[i])
	}
}

func BenchmarkLinkedList_Get(b *testing.B) {
	l := NewLinkedListOf[int]([]int{1, 2, 3})
	for i := 1; i <= b.N; i++ {
		err := l.Add(i, i)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = l.Get(i)
	}
}
