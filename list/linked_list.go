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
	"github.com/gotomicro/ekit/bean/option"
	"github.com/gotomicro/ekit/internal/errs"
)

var (
	_ List[any] = &LinkedList[any]{}
)

// node 双向循环链表结点
type node[T any] struct {
	prev *node[T]
	next *node[T]
	val  T
}

// LinkedList 双向循环链表
type LinkedList[T any] struct {
	head     *node[T]
	tail     *node[T]
	length   int
	capacity int
}

// NewLinkedList 创建一个双向循环链表
func NewLinkedList[T any](opts ...option.Option[LinkedList[T]]) *LinkedList[T] {
	head := &node[T]{}
	tail := &node[T]{next: head, prev: head}
	head.next, head.prev = tail, tail
	res := &LinkedList[T]{
		head: head,
		tail: tail,
	}
	option.Apply(res, opts...)
	return res
}

// NewLinkedListOf 将切片转换为双向循环链表, 直接使用了切片元素的值，而没有进行复制
func NewLinkedListOf[T any](ts []T, opts ...option.Option[LinkedList[T]]) *LinkedList[T] {
	list := NewLinkedList[T]()
	if err := list.Append(ts...); err != nil {
		panic(err)
	}
	option.Apply(list, opts...)
	return list
}

// WithCapacityOption 用于为LinkedList设置capacity
// 配合NewLinkedListOf()使用时，若这里设置的capacity值小于切片的长度，则用切片长度作为它的capacity
func WithCapacityOption[T any](capacity int) option.Option[LinkedList[T]] {
	return func(l *LinkedList[T]) {
		if capacity < l.Len() {
			capacity = l.Len()
		}
		l.capacity = capacity
	}
}

func (l *LinkedList[T]) IsBoundless() bool {
	return l.capacity <= 0
}

func (l *LinkedList[T]) findNode(index int) *node[T] {
	var cur *node[T]
	if index <= l.Len()/2 {
		cur = l.head
		for i := -1; i < index; i++ {
			cur = cur.next
		}
	} else {
		cur = l.tail
		for i := l.Len(); i > index; i-- {
			cur = cur.prev
		}
	}

	return cur
}

func (l *LinkedList[T]) Get(index int) (T, error) {
	if !l.checkIndex(index) {
		var zeroValue T
		return zeroValue, errs.NewErrIndexOutOfRange(l.Len(), index)
	}
	n := l.findNode(index)
	return n.val, nil
}

func (l *LinkedList[T]) checkIndex(index int) bool {
	return 0 <= index && index < l.Len()
}

// Append 往链表最后添加元素
func (l *LinkedList[T]) Append(ts ...T) error {
	if !l.IsBoundless() && l.length+len(ts) > l.capacity {
		return errs.ErrOutOfCapacity
	}
	for _, t := range ts {
		node := &node[T]{prev: l.tail.prev, next: l.tail, val: t}
		node.prev.next, node.next.prev = node, node
		l.length++
	}
	return nil
}

// Add 在 LinkedList 下标为 index 的位置插入一个元素
// 当 index 等于 LinkedList 长度等同于 Append
func (l *LinkedList[T]) Add(index int, t T) error {
	if !l.IsBoundless() && l.length == l.capacity {
		return errs.ErrOutOfCapacity
	}
	if index < 0 || index > l.length {
		return errs.NewErrIndexOutOfRange(l.length, index)
	}
	if index == l.length {
		return l.Append(t)
	}
	next := l.findNode(index)
	node := &node[T]{prev: next.prev, next: next, val: t}
	node.prev.next, node.next.prev = node, node
	l.length++
	return nil
}

// Set 设置链表中index索引处的值为t
func (l *LinkedList[T]) Set(index int, t T) error {
	if !l.checkIndex(index) {
		return errs.NewErrIndexOutOfRange(l.Len(), index)
	}
	node := l.findNode(index)
	node.val = t
	return nil
}

// Delete 删除指定位置的元素
func (l *LinkedList[T]) Delete(index int) (T, error) {
	if !l.checkIndex(index) {
		var zeroValue T
		return zeroValue, errs.NewErrIndexOutOfRange(l.Len(), index)
	}
	node := l.findNode(index)
	node.prev.next = node.next
	node.next.prev = node.prev
	node.prev, node.next = nil, nil
	l.length--
	return node.val, nil
}

// Len 返回LinkedList当前长度
func (l *LinkedList[T]) Len() int {
	return l.length
}

// Cap 返回LinkedList的容量
func (l *LinkedList[T]) Cap() int {
	return l.capacity
}

func (l *LinkedList[T]) Range(fn func(index int, t T) error) error {
	for cur, i := l.head.next, 0; i < l.length; i++ {
		err := fn(i, cur.val)
		if err != nil {
			return err
		}
		cur = cur.next
	}
	return nil
}

func (l *LinkedList[T]) AsSlice() []T {
	slice := make([]T, l.length)
	for cur, i := l.head.next, 0; i < l.length; i++ {
		slice[i] = cur.val
		cur = cur.next
	}
	return slice
}
