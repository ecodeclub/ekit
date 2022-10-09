package queue

import (
	"errors"
)

type Less[T any] func(a, b T) bool

var (
	ErrOutOfCapacity = errors.New("超出最大容量限制")
	ErrEmptyQueue    = errors.New("队列为空")
)

type PriorityArrayQueue[T any] struct {
	// 用于比较前一个元素是否小于后一个元素
	less Less[T]
	// 队列容量
	capacity int
	// 队列中的元素，为便于计算父子节点的index，0位置留空，根节点从1开始
	data []T
}

func (p *PriorityArrayQueue[T]) Len() int {
	return len(p.data) - 1
}

func (p *PriorityArrayQueue[T]) Cap() int {
	return p.capacity
}

func (p *PriorityArrayQueue[T]) Enqueue(t T) error {
	if p.capacity > 0 && len(p.data)-1 == p.capacity {
		return ErrOutOfCapacity
	}

	p.data = append(p.data, t)
	node, parent := len(p.data)-1, (len(p.data)-1)/2
	for parent > 0 && p.less(p.data[node], p.data[parent]) {
		p.data[parent], p.data[node] = p.data[node], p.data[parent]
		node = parent
		parent = parent / 2
	}

	return nil
}

func (p *PriorityArrayQueue[T]) Dequeue() (T, error) {
	if len(p.data) < 2 {
		t := new(T)
		return *t, ErrEmptyQueue
	}
	pop := p.data[1]
	p.data[1] = p.data[len(p.data)-1]
	p.data = p.data[:len(p.data)-1]
	p.heapify(p.data, len(p.data)-1, 1)
	return pop, nil
}

func (p *PriorityArrayQueue[T]) heapify(data []T, n, i int) {
	minPos := i
	for true {
		if left := i * 2; left <= n && p.less(data[left], data[minPos]) {
			minPos = left
		}
		if right := i*2 + 1; right <= n && p.less(data[right], data[minPos]) {
			minPos = right
		}
		if minPos == i {
			break
		}
		data[i], data[minPos] = data[minPos], data[i]
		i = minPos
	}
}

func (p *PriorityArrayQueue[T]) buildHeap() {
	last := len(p.data) - 1
	for i := last / 2; i > 0; i-- {
		p.heapify(p.data, len(p.data)-1, i)
	}
}

func NewBoundlessPriorityArrayQueue[T any](less Less[T]) *PriorityArrayQueue[T] {
	return &PriorityArrayQueue[T]{
		capacity: 0,
		data:     make([]T, 1, 64),
		less:     less,
	}
}

func NewPriorityArrayQueue[T any](capacity int, less Less[T]) *PriorityArrayQueue[T] {
	return &PriorityArrayQueue[T]{
		capacity: capacity,
		data:     make([]T, 1, capacity+1),
		less:     less,
	}
}

func NewPriorityArrayQueueFromArray[T any](data []T, less Less[T], opts ...PriorityArrayQueueOption[T]) *PriorityArrayQueue[T] {
	p := &PriorityArrayQueue[T]{
		capacity: len(data),
		data:     make([]T, 1, len(data)+1),
		less:     less,
	}
	p.data = append(p.data, data...)
	p.buildHeap()
	for _, opt := range opts {
		opt(p)
	}
	return p
}

type PriorityArrayQueueOption[T any] func(p *PriorityArrayQueue[T])

func WithNewCapacity[T any](capacity int) PriorityArrayQueueOption[T] {
	return func(p *PriorityArrayQueue[T]) {
		if capacity <= p.capacity {
			return
		}
		p.capacity = capacity
		old := p.data[1:]
		p.data = make([]T, 1, capacity+1)
		p.data = append(p.data, old...)
	}
}
