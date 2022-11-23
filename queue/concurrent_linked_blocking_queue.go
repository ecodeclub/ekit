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
	"sync"
)

// 实现的单链表结构
type Node[T any] struct {
	data T
	next *Node[T]
}

func NewNode[T any](data T) *Node[T] {
	return &Node[T]{
		data: data,
		next: nil,
	}
}

// ConcurrentLinkBlockingQueue 有界并发阻塞队列
type ConcurrentLinkBlockingQueue[T any] struct {
	mutex *sync.RWMutex

	//最大容量
	maxSize int

	// 队头元素下标
	head *Node[T]
	// 队尾元素下标
	tail *Node[T]
	// 包含多少个元素
	count int

	notEmpty *cond
	notFull  *cond
}

// NewConcurrentLinkBlockingQueue 创建一个有界链式阻塞队列
// 容量会在最开始的时候就初始化好
// capacity 必须为正数
func NewConcurrentLinkBlockingQueue[T any](capacity int) *ConcurrentLinkBlockingQueue[T] {
	mutex := &sync.RWMutex{}
	var t T
	res := &ConcurrentLinkBlockingQueue[T]{
		maxSize:  capacity,
		mutex:    mutex,
		notEmpty: newCond(mutex),
		notFull:  newCond(mutex),
		head:     &Node[T]{next: nil, data: t},
		count:    0,
	}
	res.tail = res.head
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
	for c.count == c.maxSize {
		signal := c.notFull.signalCh()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-signal:
			// 收到信号要重新加锁
			c.mutex.Lock()
		}
	}

	newNode := NewNode(t)
	c.tail.next = newNode
	c.tail = newNode
	c.count++

	// 这里会释放锁
	c.notEmpty.broadcast()
	return nil
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
	for c.count == 0 {
		signal := c.notEmpty.signalCh()
		select {
		case <-ctx.Done():
			var t T
			return t, ctx.Err()
		case <-signal:
			c.mutex.Lock()
		}
	}

	tmpNode := c.head.next
	c.head.next = c.head.next.next
	//只有一个元素时，需要更新tail
	if c.tail == tmpNode {
		c.tail = c.head
	}

	val := tmpNode.data
	c.count--

	c.notFull.broadcast()
	return val, nil
}

func (c *ConcurrentLinkBlockingQueue[T]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.count
}

func (c *ConcurrentLinkBlockingQueue[T]) AsSlice() []T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	res := make([]T, 0, c.count)
	tmpNode := c.head.next
	for tmpNode != nil {
		res = append(res, tmpNode.data)
		tmpNode = tmpNode.next
	}
	return res
}
