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

func (l *LinkedList[T]) getNode(index int) (n *node[T]) {
	if l.length == 0 {
		return
	}
	cur := l.head
	curIndex := 0

	if l.fromTailToHead(index) {
		cur = l.tail
		index = l.length - index - 1
		for curIndex < index {
			curIndex += 1
			cur = cur.prev
		}
	} else {
		for curIndex < index {
			curIndex += 1
			cur = cur.next
		}
	}

	n = cur
	return

}

func (l *LinkedList[T]) Get(index int) (t T, err error) {
	if index < 0 || index >= l.length {
		err = newErrIndexOutOfRange(l.length, index)
		return
	}
	var node *node[T]
	node = l.getNode(index)
	return node.val, nil
}

func (l *LinkedList[T]) Append(t T) error {
	// TODO implement me
	panic("implement me")
}

// Add 在 LinkedList 下标为 index 的位置插入一个元素
// 当 index 等于 LinkedList 长度等同于 Append
func (l *LinkedList[T]) Add(index int, t T) error {
	if index < 0 || index > l.length {
		return newErrIndexOutOfRange(l.length, index)
	}
	defer func() {
		l.length += 1
	}()

	newNode := &node[T]{
		val: t,
	}

	if l.length == 0 {
		l.head = newNode
		l.tail = newNode
		return nil
	}
	if index == 0 {
		newNode.insertAfter(l.head)
		l.head = newNode
		return nil
	}
	if index == l.length {
		l.tail.insertAfter(newNode)
		l.tail = newNode
		return nil
	}

	cur := l.getNode(index)
	prev := cur.prev
	prev.insertAfter(newNode)
	newNode.insertAfter(cur)
	return nil
}

func (l *LinkedList[T]) fromTailToHead(index int) bool {
	return index > (l.length / 2)
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

func (l *LinkedList[T]) AsSlice() []T {
	// TODO implement me
	panic("implement me")
}

type node[T any] struct {
	next *node[T]
	prev *node[T]
	val  T
}

func (n *node[T]) insertAfter(newNode *node[T]) {
	n.next = newNode
	newNode.prev = n
}

// NewLinkedListOf 将切片转换为链表, 数组的值是浅拷贝.
func NewLinkedListOf[T any](ts []T) *LinkedList[T] {
	var head *node[T] = nil
	var tail *node[T] = nil

	for _, ele := range ts {
		newNode := &node[T]{
			val: ele,
		}
		if head == nil {
			head = newNode
		} else {
			tail.insertAfter(newNode)
		}
		tail = newNode
	}

	return &LinkedList[T]{
		head:   head,
		tail:   tail,
		length: len(ts),
	}
}
