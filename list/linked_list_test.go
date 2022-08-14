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

func TestLinkedList_Add(t *testing.T) {
	fmt.Println("仿照 ArrayList 的测试写代码")
}

func TestLinkedList_Append(t *testing.T) {
	fmt.Println("仿照 ArrayList 的测试写代码")
}

func TestNewLinkedListOf(t *testing.T) {
	ts := []int{1, 2, 3, 4, 5}
	linkedList := NewLinkedListOf(ts)
	head, tail := linkedList.head, linkedList.tail
	fmt.Println("forward:")
	for head != nil {
		fmt.Printf("%d ", *head.value)
		head = head.next
	}

	fmt.Printf("\nbackward\n")
	for tail != nil {
		fmt.Printf("%d ", *tail.value)
		tail = tail.prev
	}
}

func TestLinkedList_AsSlice(t *testing.T) {
	fmt.Println("仿照 ArrayList 的测试写代码")
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
