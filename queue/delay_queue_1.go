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

//go:build demo

package queue

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gotomicro/ekit"
)

var errNoElem = errors.New("no elem")

type DelayQueue[T Delayable[T]] struct {
	q         PriorityQueue[T]
	mutex     *sync.Mutex
	available *sync.Cond
}

func (d *DelayQueue[T]) Enqueue(ctx context.Context, t T) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	err := d.q.Enqueue(ctx, t)
	if err != nil {
		return err
	}
	head, err := d.q.peek()
	if err != nil {
		return err
	}
	if t.CompareTo(head) == 0 {
		d.available.Signal()
	}
	return nil
}

func (d *DelayQueue[T]) Dequeue(ctx context.Context) (T, error) {
	ticker := time.NewTicker(0)
	ticker.Stop()
	wakeup := make(chan struct{}, 1)
	for {
		head, err := d.q.peek()
		go func() {
			d.available.Wait()
			// 可能 panic，要考虑检测有没有 close 掉 wakeup
			// 并发安全难以做到
			// 如果 wakeup 已经被关闭了，那么意味着这个调用者已经拿到值了
			// 所以它被唤醒，其实是错误的
			// 需要在 wakeup 之后，唤醒别的调用者
			wakeup <- struct{}{}
		}()
		if err == errNoElem {
			select {
			case <-ctx.Done():
				var t T
				return t, ctx.Err()
			case <-wakeup:

			}
		} else {
			ticker.Reset(head.Delay())
			select {
			case <-ctx.Done():
				var t T
				return t, ctx.Err()
			case <-ticker.C:
				return d.q.Dequeue(ctx)
			case <-wakeup:
			}
		}

	}
}

func NewDelayQueue() *DelayQueue[user] {
	mutex := &sync.Mutex{}
	return &DelayQueue[user]{
		mutex:     mutex,
		available: sync.NewCond(mutex),
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
