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

package queue

import (
	"context"
	"github.com/gotomicro/ekit/list"
	"sync"
)

// ConcurrentLinkBlockingQueue 有界并发阻塞队列
type ConcurrentLinkBlockingQueue[T any] struct {
	mutex *sync.RWMutex

	// 最大容量
	maxSize int
	// 链表
	linkedlist *list.LinkedList[T]

	notEmpty *cond
	notFull  *cond
}

// NewConcurrentLinkBlockingQueue 创建一个有界链式阻塞队列
// 容量会在最开始的时候就初始化好
// capacity 必须为正数
func NewConcurrentLinkBlockingQueue[T any](capacity int) *ConcurrentLinkBlockingQueue[T] {
	mutex := &sync.RWMutex{}
	res := &ConcurrentLinkBlockingQueue[T]{
		maxSize:    capacity,
		mutex:      mutex,
		notEmpty:   newCond(mutex),
		notFull:    newCond(mutex),
		linkedlist: list.NewLinkedListOf[T]([]T{}),
	}
	return res
}

// Enqueue 入队
// 注意：目前我们还没实现超时控制，即我们只能部分利用 ctx 里面的超时或者取消机制
// 核心在于当 goroutine 被阻塞之后，再无法监听超时或者取消
// 只有在被唤醒之后我们才会再次检测是否已经超时或者取消
func (c *ConcurrentLinkBlockingQueue[T]) Enqueue(ctx context.Context, t T) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mutex.Lock()
	for c.linkedlist.Len() == c.maxSize {
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
// 注意：目前我们还没实现超时控制，即我们只能部分利用 ctx 里面的超时或者取消机制
// 核心在于当 goroutine 被阻塞之后，再无法监听超时或者取消
// 只有在被唤醒之后我们才会再次检测是否已经超时或者取消
func (c *ConcurrentLinkBlockingQueue[T]) Dequeue(ctx context.Context) (T, error) {
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

func (c *ConcurrentLinkBlockingQueue[T]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.linkedlist.Len()
}

func (c *ConcurrentLinkBlockingQueue[T]) AsSlice() []T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	res := c.linkedlist.AsSlice()
	return res
}
