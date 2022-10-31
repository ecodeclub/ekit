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
	"time"

	"github.com/gotomicro/ekit/list"

	"github.com/gotomicro/ekit/internal/queue"
)

type DelayQueue[T Delayable] struct {
	q     queue.PriorityQueue[T]
	mutex sync.RWMutex

	enqueueReqs *list.LinkedList[delayQueueReq]
	dequeueReqs *list.LinkedList[delayQueueReq]
}

type delayQueueReq struct {
	ch chan struct{}
}

func NewDelayQueue[T Delayable](c int) *DelayQueue[T] {
	return &DelayQueue[T]{
		q: *queue.NewPriorityQueue[T](c, func(src T, dst T) int {
			srcDelay := src.Delay()
			dstDelay := dst.Delay()
			if srcDelay > dstDelay {
				return 1
			}
			if srcDelay == dstDelay {
				return 0
			}
			return -1
		}),
		enqueueReqs: list.NewLinkedList[delayQueueReq](),
		dequeueReqs: list.NewLinkedList[delayQueueReq](),
	}
}

func (d *DelayQueue[T]) Enqueue(ctx context.Context, t T) error {
	// 确保 ctx 没有过期
	if ctx.Err() != nil {
		return ctx.Err()
	}
	for {
		d.mutex.Lock()
		err := d.q.Enqueue(t)
		if err == queue.ErrOutOfCapacity {
			ch := make(chan struct{}, 1)
			_ = d.enqueueReqs.Append(delayQueueReq{ch: ch})
			d.mutex.Unlock()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ch:
			}
			continue
		}
		if err == nil {
			// 这里使用写锁，是为了在 Dequeue 那边
			// 当一开始的 Peek 返回 queue.ErrEmptyQueue 的时候不会错过这个入队信号
			if d.dequeueReqs.Len() == 0 {
				// 没人等。
				d.mutex.Unlock()
				return nil
			}
			req, err := d.dequeueReqs.Delete(0)
			if err == nil {
				// 唤醒出队的
				req.ch <- struct{}{}
			}
		}
		d.mutex.Unlock()
		return err
	}
}

func (d *DelayQueue[T]) Dequeue(ctx context.Context) (T, error) {
	// 确保 ctx 没有过期
	if ctx.Err() != nil {
		var t T
		return t, ctx.Err()
	}
	ticker := time.NewTicker(time.Second)
	ticker.Stop()
	defer func() {
		ticker.Stop()
	}()
	for {
		d.mutex.Lock()
		head, err := d.q.Peek()
		if err != nil && err != queue.ErrEmptyQueue {
			var t T
			return t, err
		}
		if err == queue.ErrEmptyQueue {
			ch := make(chan struct{}, 1)
			_ = d.dequeueReqs.Append(delayQueueReq{ch: ch})
			d.mutex.Unlock()
			select {
			case <-ctx.Done():
				var t T
				return t, ctx.Err()
			case <-ch:
			}
			continue
		}

		delay := head.Delay()
		// 已经到期了
		if delay <= 0 {
			// 拿着锁，所以不然不可能返回 error
			t, _ := d.q.Dequeue()
			d.wakeEnqueue()
			d.mutex.Unlock()
			return t, nil
		}

		// 在进入 select 之前必须要释放锁
		d.mutex.Unlock()
		ticker.Reset(delay)
		select {
		case <-ctx.Done():
			var t T
			return t, ctx.Err()
		case <-ticker.C:
			var t T
			d.mutex.Lock()
			t, err = d.q.Dequeue()
			// 被人抢走了，理论上是不会出现这个可能的
			if err != nil {
				d.mutex.Unlock()
				continue
			}
			d.wakeEnqueue()
			d.mutex.Unlock()
			return t, nil
		}
	}
}

func (d *DelayQueue[T]) wakeEnqueue() {
	req, err := d.enqueueReqs.Delete(0)
	if err == nil {
		// 唤醒等待入队的
		req.ch <- struct{}{}
	}
}

type Delayable interface {
	Delay() time.Duration
}
