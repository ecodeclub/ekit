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

	"math/rand"
	"testing"
)

func TestLinkedList_Add(t *testing.T) {

	testCases := []struct {
		name           string
		list           *LinkedList[int]
		index          int
		newVal         int
		wantLinkedList *LinkedList[int]
		wantErr        error
	}{
		// 仿照这个例子，继续添加测试
		// 你需要综合考虑下标的各种可能取值
		// 往两边增加，往中间加
		// 下标可能是负数，也可能超出你的长度
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Add(tc.index, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
			assert.True(t, linkedListEqual(tc.list, tc.wantLinkedList))

		})
	}
}

func TestLinkedList_Append(t *testing.T) {
	fmt.Println("仿照 ArrayList 的测试写代码")
}

func TestNewLinkedListOf(t *testing.T) {
	ts := []int{1, 2, 3, 4, 5}
	linkedList := NewLinkedListOf(ts)
	head, tail := linkedList.head, linkedList.tail
	t.Log(head.val, tail.val)
	fmt.Println("forward")
	for head != nil {
		fmt.Printf("%d ", head.val)
		head = head.next
	}
	fmt.Println("\nbackward")
	for tail != nil && &tail.val != nil {
		fmt.Printf("%d ", tail.val)
		tail = tail.prev
	}
	t.Log("Ok")
}

func TestLinkedList_AsSlice(t *testing.T) {
	vals := []int{1, 2, 3}
	a := NewLinkedListOf[int](vals)
	slice := a.AsSlice()
	fmt.Println(vals, slice)
	// 内容相同
	assert.Equal(t, slice, vals)
	aAddr := fmt.Sprintf("%p", vals)
	sliceAddr := fmt.Sprintf("%p", slice)
	// 但是地址不同，也就是意味着 slice 必须是一个新创建的
	assert.NotEqual(t, aAddr, sliceAddr)
}

func TestLinkedList_Cap(t *testing.T) {
	fmt.Println("仿照 ArrayList 的测试写代码")
}

func TestLinkedList_Delete(t *testing.T) {
	fmt.Println("仿照 ArrayList 的测试写代码")
}

func TestLinkedList_Get(t *testing.T) {
	fmt.Println("仿照 ArrayList 的测试写代码")
}

func TestLinkedList_Len(t *testing.T) {
	fmt.Println("仿照 ArrayList 的测试写代码")
}

func TestLinkedList_Range(t *testing.T) {
	fmt.Println("仿照 ArrayList 的测试写代码")
}

func TestLinkedList_Set(t *testing.T) {
	fmt.Println("仿照 ArrayList 的测试写代码")
}

func linkedListEqual[T comparable](l1 *LinkedList[T], l2 *LinkedList[T]) bool {
	if l1.length != l2.length {
		fmt.Println(l1.length, l2.length)
		return false
	}

	if l1.length == 0 {
		return true
	}

	l1Pos := l1.head
	l2Pos := l2.head
	for l1Pos != nil && l2Pos != nil {
		if l1Pos.val != l2Pos.val {
			fmt.Println(l1Pos.val, l2Pos.val)
			return false
		}
		l1Pos = l1Pos.next
		l2Pos = l2Pos.next
	}
	return l1Pos == l2Pos
}

func BenchmarkLinkedList_Add(b *testing.B) {
	l := NewLinkedListOf[int]([]int{1, 2, 3})
	testCase := make([]int, 0, b.N)
	for i := 1; i <= b.N; i++ {
		testCase = append(testCase, rand.Intn(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Add(testCase[i], testCase[i])
	}
}
