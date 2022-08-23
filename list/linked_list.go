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

func (l *LinkedList[T]) getNode(index int) *node[T] {
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
	return cur
}

func (l *LinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= l.length {
		var t T
		return t, newErrIndexOutOfRange(l.length, index)
	}
	n := l.getNode(index)
	return n.val, nil
}

// Append 往链表最后添加元素
func (l *LinkedList[T]) Append(t T) error {
	newLastNode := &node[T]{val: t}
	if l.length == 0 {
		l.head = newLastNode
		l.tail = newLastNode
	} else {
		l.tail.next = newLastNode
		newLastNode.prev = l.tail
		l.tail = newLastNode
	}
	l.length += 1
	return nil
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

// Set 设置链表中index索引处的值为t
func (l *LinkedList[T]) Set(index int, t T) error {
	if index < 0 || index >= l.length {
		return newErrIndexOutOfRange(l.length, index)
	}
	rv := l.getNode(index)
	rv.val = t
	return nil
}

// Delete 删除指定位置的元素
func (l *LinkedList[T]) Delete(index int) (T, error) {
	nLen := l.length
	var delVal T // 需要删除的节点val
	if index < 0 || index >= nLen {
		return delVal, newErrIndexOutOfRange(nLen, index)
	}
	defer func() {
		l.length -= 1
	}()

	// 删除head
	if index == 0 {
		delVal = l.head.val
		if nLen > 1 {
			l.head.next.prev = nil
			l.head = l.head.next
		} else {
			l.head = nil
			l.tail = nil
		}
		return delVal, nil
	}
	// 删除tail
	if index == nLen-1 {
		delVal = l.tail.val
		l.tail = l.tail.prev
		l.tail.next = nil
		return delVal, nil
	}

	n := l.head
	for i := 0; i < index-1; i++ {
		n = n.next
	}
	delVal = n.next.val
	n.next = n.next.next
	n.next.prev = n

	return delVal, nil
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
	slice := make([]T, l.length)
	head := l.head
	for i := 0; i < l.length; i++ {
		slice[i] = head.val
		head = head.next
	}
	return slice
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
