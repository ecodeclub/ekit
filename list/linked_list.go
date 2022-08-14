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

// LinkedList 双向链表
type LinkedList[T any] struct {
	head *node[T]
	tail *node[T]
	// length 有多少个元素
	length int
}

func (l *LinkedList[T]) Get(index int) (T, error) {
	// TODO implement me
	panic("implement me")
}

func (l *LinkedList[T]) Append(t T) error {
	// TODO implement me
	panic("implement me")
}

func (l *LinkedList[T]) Add(index int, t T) error {
	// TODO implement me
	panic("implement me")
}

func (l *LinkedList[T]) Set(index int, t T) error {
	// TODO implement me
	panic("implement me")
}

func (l *LinkedList[T]) Delete(index int) (T, error) {
	// TODO implement me
	panic("implement me")
}

func (l *LinkedList[T]) Len() int {
	// TODO implement me
	panic("implement me")
}

func (l *LinkedList[T]) Cap() int {
	// TODO implement me
	panic("implement me")
}

func (l *LinkedList[T]) Range(fn func(index int, t T) error) error {
	// TODO implement me
	panic("implement me")
}

func NewLinkedListOf[T any](ts []T) *LinkedList[T] {
	if ts == nil || len(ts) < 1 {
		return nil
	}

	dummyHead := &node[T]{}
	prev := dummyHead

	for i := 0; i < len(ts); i++ {
		cur := &node[T]{
			prev:  prev,
			value: &ts[i],
		}
		if dummyHead.next == nil {
			dummyHead.next = cur
		}
		prev.next = cur
		prev = cur
	}

	return &LinkedList[T]{
		head:   dummyHead.next,
		tail:   prev,
		length: len(ts),
	}
}

func (l *LinkedList[T]) AsSlice() []T {
	//slice := make([]T, 0)
	//head := l.head
	//for head != nil {
	//	slice = append(slice, *head.value)
	//	head = head.next
	//}
	slice := make([]T, l.length)
	head := l.head
	for i := 0; i < l.length; i++ {
		slice[i] = *head.value
		head = head.next
	}
	return slice
}

type node[T any] struct {
	next  *node[T]
	prev  *node[T]
	value *T
}
