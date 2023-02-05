package stack

import (
	"errors"
	"math"
)

const (
	minCapacity uint32 = 16
	maxCapacity uint32 = math.MaxInt32 - 8
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

// NewArrayDequeWithCap 构造一个指定容量的Deque
func NewArrayDequeWithCap[T any](minCap uint32) *ArrayDeque[T] {
	capp := calculateCap(minCap)
	return &ArrayDeque[T]{
		elements: make([]T, capp, capp),
		head:     0,
		tail:     0,
	}
}

type ArrayDeque[T any] struct {
	elements   []T
	head, tail uint32
	empty      T
}

func (a *ArrayDeque[T]) Enqueue(t T) error {
	return a.AddLast(t)
}

func (a *ArrayDeque[T]) Dequeue() (T, error) {
	return a.RemoveFirst()
}

func (a *ArrayDeque[T]) PeekFirst() (T, error) {
	return a.GetFirst()
}

func (a *ArrayDeque[T]) Push(t T) error {
	return a.AddLast(t)
}

func (a *ArrayDeque[T]) Pop() (T, error) {
	return a.RemoveLast()
}

func (a *ArrayDeque[T]) Peek() (T, error) {
	return a.GetLast()
}

func (a *ArrayDeque[T]) AddLast(t T) error {
	a.elements[a.tail] = t
	a.tail = (a.tail + 1) & (a.capacity() - 1)
	if a.tail == a.head {
		err := a.growDouble()
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveLast 因为泛型不支持设置为nil, 所以使用trim方法在合适的时机缩容，但是还是存在内存释放不够及时的问题，可能存在内存泄露的隐患
func (a *ArrayDeque[T]) RemoveLast() (T, error) {
	defer func() {
		a.trim()
	}()
	var result T
	if a.IsEmpty() {
		return result, ErrEmpty
	}
	t := (a.tail - 1) & (a.capacity() - 1)
	result = a.elements[t]
	a.elements[t] = a.empty
	a.tail = t
	return result, nil
}

func (a *ArrayDeque[T]) AddFirst(t T) error {
	a.elements[a.head] = t
	a.head = (a.head + 1) & (a.capacity() - 1)
	if a.head == a.tail {
		err := a.growDouble()
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveFirst 因为不支持设置nil, 所以trim方法在合适的时机缩容，但是还是存在内存释放不够及时的问题，可能存在内存泄漏的隐患
func (a *ArrayDeque[T]) RemoveFirst() (T, error) {
	defer func() {
		a.trim()
	}()
	var result T
	if a.IsEmpty() {
		return result, ErrEmpty
	}
	h := a.head
	result = a.elements[h]
	a.elements[h] = a.empty
	a.head = (h + 1) & (a.capacity() - 1)
	return result, nil
}

func (a *ArrayDeque[T]) GetFirst() (T, error) {
	var result T
	if a.IsEmpty() {
		return result, ErrEmpty
	}
	result = a.elements[a.head]
	return result, nil
}

func (a *ArrayDeque[T]) GetLast() (T, error) {
	var result T
	if a.IsEmpty() {
		return result, ErrEmpty
	}
	t := (a.tail - 1) & (a.capacity() - 1)
	result = a.elements[t]
	return result, nil
}

func (a *ArrayDeque[T]) Len() uint32 {
	return (a.tail - a.head) & (a.capacity() - 1)
}

func (a *ArrayDeque[T]) IsEmpty() bool {
	return a.Len() == 0
}

// capacity , return capacity of the deque
func (a *ArrayDeque[T]) capacity() uint32 {
	return uint32(cap(a.elements))
}

// growDouble , grow deque capacity
func (a *ArrayDeque[T]) growDouble() error {
	p := a.head
	n := a.capacity()
	r := n - p
	newCap := n << 1
	if newCap > maxCapacity {
		return ErrOutOfCapacity
	}
	newElements := make([]T, newCap, newCap)
	copyArr(a.elements, p, newElements, 0, r)
	copyArr(a.elements, 0, newElements, r, p)
	a.elements = newElements
	a.head = 0
	a.tail = n
	return nil
}

func (a *ArrayDeque[T]) trim() {
	if a.IsEmpty() {
		return
	}
	curLen := a.Len()
	if curLen <= minCapacity {
		return
	}
	expectedCapp := calculateCap(curLen)
	if expectedCapp == a.capacity() {
		return
	}
	p := a.head
	n := curLen
	r := n - p
	newElement := make([]T, expectedCapp, expectedCapp)
	copyArr(a.elements, p, newElement, 0, r)
	copyArr(a.elements, 0, newElement, r, p)
	a.elements = newElement
	a.head = 0
	a.tail = n
}

// calculateCap , calculate init capacity
func calculateCap(capacity uint32) uint32 {
	initialCapacity := minCapacity
	if capacity >= initialCapacity {
		initialCapacity = capacity
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

// copyArr , copy slice
func copyArr[T any](src []T, srcPos uint32, dest []T, destPos uint32, length uint32) {
	for i := srcPos; i < length+srcPos; i++ {
		dest[destPos] = src[i]
		destPos++
	}
}
