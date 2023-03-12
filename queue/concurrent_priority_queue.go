// Copyright 2021 ecodeclub
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

package queue

import (
	"sync"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/internal/queue"
)

type ConcurrentPriorityQueue[T any] struct {
	pq queue.PriorityQueue[T]
	m  sync.RWMutex
}

func (c *ConcurrentPriorityQueue[T]) Len() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.pq.Len()
}

// Cap 无界队列返回0，有界队列返回创建队列时设置的值
func (c *ConcurrentPriorityQueue[T]) Cap() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.pq.Cap()
}

func (c *ConcurrentPriorityQueue[T]) Peek() (T, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.pq.Peek()
}

func (c *ConcurrentPriorityQueue[T]) Enqueue(t T) error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.pq.Enqueue(t)
}

func (c *ConcurrentPriorityQueue[T]) Dequeue() (T, error) {
	c.m.Lock()
	defer c.m.Unlock()
	return c.pq.Dequeue()
}

// NewConcurrentPriorityQueue 创建优先队列 capacity <= 0 时，为无界队列
func NewConcurrentPriorityQueue[T any](capacity int, compare ekit.Comparator[T]) *ConcurrentPriorityQueue[T] {
	return &ConcurrentPriorityQueue[T]{
		pq: *queue.NewPriorityQueue[T](capacity, compare),
	}
}
