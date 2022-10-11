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
	// 已经到期的
	ch chan T
}

func (d *DelayQueue[T]) Enqueue(ctx context.Context, t T) error {
	d.mutex.Lock()
	err := d.q.Enqueue(ctx, t)
	if err != nil {
		d.mutex.Unlock()
		return err
	}
	// 这里释放锁并不会有任何的问题
	d.mutex.Unlock()
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
	select {
	case <-ctx.Done():
		var t T
		return t, ctx.Err()
	case t := <-d.ch:
		return t, nil
	}
}

func NewDelayQueue[T Delayable[T]]() *DelayQueue[T] {
	mutex := &sync.Mutex{}
	res := &DelayQueue[T]{
		mutex:     mutex,
		available: sync.NewCond(mutex),
	}

	go func() {
		for {
			t, err := res.q.peek()
			if err != nil && err != errNoElem {
				return
			}
			if err == errNoElem {
				res.available.Wait()
			} else {
				delay := t.Delay()
				if delay <= 0 {
					res.mutex.Lock()
					t, err = res.q.peek()
					if err == nil && t.Delay() <= 0 {
						// 这里应该能够立刻获得一个元素
						t, err = res.q.Dequeue(context.Background())
						res.mutex.Unlock()
						res.ch <- t
					}
					continue
				}
				go func() {
					time.Sleep(delay)
					res.available.Signal()
				}()
				res.available.Wait()
			}
		}
	}()

	return res
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
