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

	"golang.org/x/sync/semaphore"
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

	enqueueCap *semaphore.Weighted
	dequeueCap *semaphore.Weighted

	// zero 不能作为返回值返回，防止用户篡改
	zero T
}

// NewConcurrentArrayBlockingQueue 创建一个有界阻塞队列
// 容量会在最开始的时候就初始化好
// capacity 必须为正数
func NewConcurrentArrayBlockingQueue[T any](capacity int) *ConcurrentArrayBlockingQueue[T] {
	mutex := &sync.RWMutex{}

	semaForEnqueue := semaphore.NewWeighted(int64(capacity))
	semaForDequeue := semaphore.NewWeighted(int64(capacity))

	// error暂时不处理，因为目前没办法处理，只能考虑panic掉
	// 相当于将信号量置空
	_ = semaForDequeue.Acquire(context.TODO(), int64(capacity))

	res := &ConcurrentArrayBlockingQueue[T]{
		data:       make([]T, capacity),
		mutex:      mutex,
		enqueueCap: semaForEnqueue,
		dequeueCap: semaForDequeue,
	}
	return res
}

// Enqueue 入队
// 通过sema来控制容量、超时、阻塞问题
func (c *ConcurrentArrayBlockingQueue[T]) Enqueue(ctx context.Context, t T) error {

	// 能拿到，说明队列还有空位，可以入队，拿不到则阻塞
	err := c.enqueueCap.Acquire(ctx, 1)

	if err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 拿到锁，先判断是否超时，防止在抢锁时已经超时
	if ctx.Err() != nil {

		// 超时应该主动归还信号量，避免容量泄露
		c.enqueueCap.Release(1)

		return ctx.Err()
	}

	c.data[c.tail] = t
	c.tail++
	c.count++

	// c.tail 已经是最后一个了，重置下标
	if c.tail == cap(c.data) {
		c.tail = 0
	}

	// 往出队的sema放入一个元素，出队的goroutine可以拿到并出队
	c.dequeueCap.Release(1)

	return nil

}

// Dequeue 出队
// 通过sema来控制容量、超时、阻塞问题
func (c *ConcurrentArrayBlockingQueue[T]) Dequeue(ctx context.Context) (T, error) {

	// 能拿到，说明队列有元素可以取，可以出队，拿不到则阻塞
	err := c.dequeueCap.Acquire(ctx, 1)

	var res T

	if err != nil {
		return res, err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 拿到锁，先判断是否超时，防止在抢锁时已经超时
	if ctx.Err() != nil {

		// 超时应该主动归还信号量，有元素消费不到
		c.dequeueCap.Release(1)

		return res, ctx.Err()
	}

	res = c.data[c.head]
	// 为了释放内存，GC
	c.data[c.head] = c.zero

	c.head++
	c.count--
	if c.head == cap(c.data) {
		c.head = 0
	}

	// 往入队的sema放入一个元素，入队的goroutine可以拿到并入队
	c.enqueueCap.Release(1)

	return res, nil

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
