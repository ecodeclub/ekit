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

	"github.com/gotomicro/ekit/internal/queue"

	"github.com/gotomicro/ekit"
)

type DelayQueue[T Delayable[T]] struct {
	q     queue.PriorityQueue[T]
	mutex sync.RWMutex

	enqueueSignal chan struct{}
	dequeueSignal chan struct{}
}

func NewDelayQueue[T Delayable[T]](compare ekit.Comparator[T]) *DelayQueue[T] {
	return &DelayQueue[T]{
		q:             *queue.NewPriorityQueue[T](0, compare),
		enqueueSignal: make(chan struct{}, 1),
		dequeueSignal: make(chan struct{}, 1),
	}
}

func (d *DelayQueue[T]) Enqueue(ctx context.Context, t T) error {
	for {
		d.mutex.Lock()
		err := d.q.Enqueue(t)
		d.mutex.Unlock()
		if err == queue.ErrOutOfCapacity {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-d.dequeueSignal:
				continue
			}
		}

		if err == nil {
			// 这里使用写锁，是为了在 Dequeue 那边
			// 当一开始的 Peek 返回 queue.ErrEmptyQueue 的时候不会错过这个入队信号
			d.mutex.Lock()
			head, err := d.q.Peek()
			if err != nil {
				// 这种情况就是出现在入队成功之后，元素立刻被取走了
				// 这里 err 预期应该只有 queue.ErrEmptyQueue 一种可能
				d.mutex.Lock()
				return nil
			}
			if t.CompareTo(head) == 0 {
				select {
				case d.enqueueSignal <- struct{}{}:
				default:
				}
			}
			d.mutex.Lock()
		}
		return err
	}

}

func (d *DelayQueue[T]) Dequeue(ctx context.Context) (T, error) {
	ticker := time.NewTicker(0)
	ticker.Stop()
	for {
		d.mutex.RLock()
		head, err := d.q.Peek()
		d.mutex.RUnlock()
		if err == queue.ErrEmptyQueue {
			select {
			case <-ctx.Done():
				var t T
				return t, ctx.Err()
			case <-d.enqueueSignal:
			}
		} else {
			ticker.Reset(head.Delay())
			select {
			case <-ctx.Done():
				var t T
				return t, ctx.Err()
			case <-ticker.C:
				var t T
				d.mutex.Lock()
				t, err = d.q.Dequeue()
				d.mutex.Unlock()
				// 被人抢走了，理论上是不会出现这个可能的
				if err == queue.ErrEmptyQueue {
					continue
				}
				select {
				case d.dequeueSignal <- struct{}{}:
				default:
				}
				return t, nil
			case <-d.enqueueSignal:
			}
		}
	}
}

type Delayable[T any] interface {
	Delay() time.Duration
	ekit.Comparable[T]
}

type user struct {
}

func (u user) Delay() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (u user) CompareTo(dst user) int {
	//TODO implement me
	panic("implement me")
}
