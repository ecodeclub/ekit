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
	"context"
	"sync"

	"github.com/ecodeclub/ekit/list"
)

// ConcurrentLinkedBlockingQueue 基于链表的并发阻塞队列
// 如果 maxSize 是正数。那么就是有界并发阻塞队列
// 如果不是，就是无界并发阻塞队列, 在这种情况下，入队永远能够成功
type ConcurrentLinkedBlockingQueue[T any] struct {
	mutex *sync.RWMutex

	// 最大容量
	maxSize int
	// 链表
	linkedlist *list.LinkedList[T]

	notEmpty *cond
	notFull  *cond
}

// NewConcurrentLinkedBlockingQueue 创建链式阻塞队列 capacity <= 0 时，为无界队列
func NewConcurrentLinkedBlockingQueue[T any](capacity int) *ConcurrentLinkedBlockingQueue[T] {
	mutex := &sync.RWMutex{}
	res := &ConcurrentLinkedBlockingQueue[T]{
		maxSize:    capacity,
		mutex:      mutex,
		notEmpty:   newCond(mutex),
		notFull:    newCond(mutex),
		linkedlist: list.NewLinkedList[T](),
	}
	return res
}

// Enqueue 入队
// 注意：目前我们已经通过broadcast实现了超时控制
func (c *ConcurrentLinkedBlockingQueue[T]) Enqueue(ctx context.Context, t T) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mutex.Lock()
	for c.maxSize > 0 && c.linkedlist.Len() == c.maxSize {
		signal := c.notFull.signalCh()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-signal:
			// 收到信号要重新加锁
			c.mutex.Lock()
		}
	}

	err := c.linkedlist.Append(t)

	// 这里会释放锁
	c.notEmpty.broadcast()
	return err
}

// Dequeue 出队
// 注意：目前我们已经通过broadcast实现了超时控制
func (c *ConcurrentLinkedBlockingQueue[T]) Dequeue(ctx context.Context) (T, error) {
	if ctx.Err() != nil {
		var t T
		return t, ctx.Err()
	}
	c.mutex.Lock()
	for c.linkedlist.Len() == 0 {
		signal := c.notEmpty.signalCh()
		select {
		case <-ctx.Done():
			var t T
			return t, ctx.Err()
		case <-signal:
			c.mutex.Lock()
		}
	}

	val, err := c.linkedlist.Delete(0)
	c.notFull.broadcast()
	return val, err
}

func (c *ConcurrentLinkedBlockingQueue[T]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.linkedlist.Len()
}

func (c *ConcurrentLinkedBlockingQueue[T]) AsSlice() []T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	res := c.linkedlist.AsSlice()
	return res
}
