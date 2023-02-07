package stack

import (
	"errors"
)

const (
	minCapacity int = 1 << 4
	maxCapacity int = 1 << 30
)

var (
	ErrOutOfCapacity = errors.New("ekit: 超出最大容量限制")
	ErrEmpty         = errors.New("ekit: deque 为空")
)

var (
	_ Stack[any] = &ArrayDeque[any]{}
	_ Queue[any] = &ArrayDeque[any]{}
)

// NewArrayDeque 构造一个最小容量的Deque
func NewArrayDeque[T any]() *ArrayDeque[T] {
	return NewArrayDequeWithCap[T](minCapacity)
}

// NewArrayDequeWithCap 构造一个期望容量的Deque
func NewArrayDequeWithCap[T any](expectedCapacity int) *ArrayDeque[T] {
	capacity := calculateCapacity(expectedCapacity)
	return &ArrayDeque[T]{
		elements: make([]T, capacity, capacity),
		head:     0,
		tail:     0,
	}
}

func NewArrayDequeOf[T any](elements ...T) *ArrayDeque[T] {
	length := len(elements)
	deque := NewArrayDequeWithCap[T](length)
	for _, val := range elements {
		_ = deque.AddLast(val)
	}
	return deque
}

// ArrayDeque 基于切片实现的可扩容循环队列。
//  1.相较于SDK中的list.List, ArrayDeque无需维护复杂的链表关系
//  2.相较于SDK中的list.List, 基于切片的随机访问，可以快速获取数据
//  3.相较于SDK中的list.List, 基于切片的连续性，可以降低内存的碎片化
// 	4.此循环队列不是无限制添加的，容量限定范围为[1 << 4, 1 << 30]
// 	5.此循环队列不是线程安全的，它不支持多线程的并发访问
// 	6.实现了Queue和Stack接口，你可以将他作为栈或者队列使用
//  7.ArrayDeque没有对空值作限制，但是不推荐向其中写入nil
//  8.ArrayDeque具有扩容机制，并且容量始终为2的幂，这造成了内存空间的浪费，在使用前可以预估容量，减少扩容导致的内存浪费
type ArrayDeque[T any] struct {
	elements   []T
	head, tail int
	empty      T
}

func (a *ArrayDeque[T]) Front() (T, error) {
	return a.GetFirst()
}

func (a *ArrayDeque[T]) Enqueue(t T) error {
	return a.AddLast(t)
}

func (a *ArrayDeque[T]) Dequeue() (T, error) {
	return a.RemoveFirst()
}

func (a *ArrayDeque[T]) Push(t T) error {
	return a.AddLast(t)
}

func (a *ArrayDeque[T]) Pop() (T, error) {
	return a.RemoveLast()
}

func (a *ArrayDeque[T]) Top() (T, error) {
	return a.GetLast()
}

// AddFirst 从队首添加元素，如果失败，那么返回错误
func (a *ArrayDeque[T]) AddFirst(t T) error {
	a.head = (a.head - 1) & (a.Cap() - 1)
	a.elements[a.head] = t
	if a.head == a.tail {
		err := a.doubleCapacity()
		if err != nil {
			return err
		}
	}
	return nil
}

// AddLast 从尾部添加元素，如果失败，那么返回错误
func (a *ArrayDeque[T]) AddLast(t T) error {
	a.elements[a.tail] = t
	a.tail = (a.tail + 1) & (a.Cap() - 1)
	if a.head == a.tail {
		err := a.doubleCapacity()
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveFirst 删除第一个元素，如果为空，那么返回错误
func (a *ArrayDeque[T]) RemoveFirst() (T, error) {
	var result T
	if a.IsEmpty() {
		return result, ErrEmpty
	}
	result = a.elements[a.head]
	a.elements[a.head] = a.empty
	a.head = (a.head + 1) & (a.Cap() - 1)
	return result, nil
}

// RemoveLast 删除最后一个元素，如果为空，那么返回错误
func (a *ArrayDeque[T]) RemoveLast() (T, error) {
	var result T
	if a.IsEmpty() {
		return result, ErrEmpty
	}
	a.tail = (a.tail - 1) & (a.Cap() - 1)
	result = a.elements[a.tail]
	a.elements[a.tail] = a.empty
	return result, nil
}

// GetFirst 获取arrayDeque的队首元素，如果没有，那么返回错误
func (a *ArrayDeque[T]) GetFirst() (T, error) {
	var result T
	if a.IsEmpty() {
		return result, ErrEmpty
	}
	result = a.elements[a.head]
	return result, nil
}

// GetLast 获取arrayDeque末尾的元素，如果没有，那么返回ErrEmpty
func (a *ArrayDeque[T]) GetLast() (T, error) {
	var result T
	if a.IsEmpty() {
		return result, ErrEmpty
	}
	t := (a.tail - 1) & (a.Cap() - 1)
	result = a.elements[t]
	return result, nil
}

// Cap 返回arrayDeque此时的容量
func (a *ArrayDeque[T]) Cap() int {
	return cap(a.elements)
}

// Len 返回arrayDeque中的元素个数
func (a *ArrayDeque[T]) Len() int {
	return (a.tail - a.head) & (a.Cap() - 1)
}

// IsEmpty 返回arrayDeque是否为空
func (a *ArrayDeque[T]) IsEmpty() bool {
	return a.Len() == 0
}

// doubleCapacity 给arrayDeque扩容，如果容量超过了最大容量，那么返回错误
func (a *ArrayDeque[T]) doubleCapacity() error {
	n := a.Cap()
	newCapacity := n << 1
	if newCapacity < 0 || newCapacity > maxCapacity {
		return ErrOutOfCapacity
	}
	newElements := make([]T, newCapacity, newCapacity)
	r := n - a.head
	copy(newElements[0:r], a.elements[a.head:r+a.head])
	copy(newElements[r:], a.elements[:a.tail+1])
	a.elements = newElements
	a.head = 0
	a.tail = n
	return nil
}

// calculateCapacity 从给定的预期容量中计算一个合适的预期容量
func calculateCapacity(expected int) int {
	initialCapacity := minCapacity
	if expected > initialCapacity {
		initialCapacity = expected
		initialCapacity |= initialCapacity >> 1
		initialCapacity |= initialCapacity >> 2
		initialCapacity |= initialCapacity >> 4
		initialCapacity |= initialCapacity >> 8
		initialCapacity |= initialCapacity >> 16
		initialCapacity++
		if initialCapacity > maxCapacity {
			initialCapacity = maxCapacity
		}
	}
	return initialCapacity
}
