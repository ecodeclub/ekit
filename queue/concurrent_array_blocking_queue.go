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

// ConcurrentArrayBlockingQueue 有界并发阻塞队列
type ConcurrentArrayBlockingQueue[T any] struct {
	data  []T
	mutex *sync.RWMutex

	// 队头元素下标
	head int
	// 队尾元素下标
	tail int
	// 包含多少个元素
	count int

	notEmpty *cond
	notFull  *cond

	// zero 不能作为返回值返回，防止用户篡改
	zero T
}

// NewConcurrentBlockingQueue 创建一个有界阻塞队列
// 容量会在最开始的时候就初始化好
// capacity 必须为正数
func NewConcurrentBlockingQueue[T any](capacity int) *ConcurrentArrayBlockingQueue[T] {
	mutex := &sync.RWMutex{}
	res := &ConcurrentArrayBlockingQueue[T]{
		data:     make([]T, capacity),
		mutex:    mutex,
		notEmpty: newCond(mutex),
		notFull:  newCond(mutex),
	}
	return res
}

// Enqueue 入队
// 注意：目前我们还没实现超时控制，即我们只能部分利用 ctx 里面的超时或者取消机制
// 核心在于当 goroutine 被阻塞之后，再无法监听超时或者取消
// 只有在被唤醒之后我们才会再次检测是否已经超时或者取消
func (c *ConcurrentArrayBlockingQueue[T]) Enqueue(ctx context.Context, t T) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mutex.Lock()
	for c.count == len(c.data) {
		signal := c.notFull.signalCh()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-signal:
			// 收到信号要重新加锁
			c.mutex.Lock()
		}
	}
	c.data[c.tail] = t
	c.tail++
	c.count++
	// c.tail 已经是最后一个了，重置下标
	if c.tail == cap(c.data) {
		c.tail = 0
	}
	// 这里会释放锁
	c.notEmpty.broadcast()
	return nil
}

// Dequeue 出队
// 注意：目前我们还没实现超时控制，即我们只能部分利用 ctx 里面的超时或者取消机制
// 核心在于当 goroutine 被阻塞之后，再无法监听超时或者取消
// 只有在被唤醒之后我们才会再次检测是否已经超时或者取消
func (c *ConcurrentArrayBlockingQueue[T]) Dequeue(ctx context.Context) (T, error) {
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
	val := c.data[c.head]
	// 为了释放内存，GC
	c.data[c.head] = c.zero
	c.count--
	c.head++
	// 重置下标
	if c.head == cap(c.data) {
		c.head = 0
	}
	c.notFull.broadcast()
	return val, nil
}

func (c *ConcurrentArrayBlockingQueue[T]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.count
}

func (c *ConcurrentArrayBlockingQueue[T]) AsSlice() []T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	res := make([]T, 0, c.count)
	cnt := 0
	capacity := cap(c.data)
	for cnt < c.count {
		index := (c.head + cnt) % capacity
		res = append(res, c.data[index])
		cnt++
	}
	return res
}
