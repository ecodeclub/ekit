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
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcurrentLinkedBlockingQueue_Enqueue(t *testing.T) {
	testCases := []struct {
		name      string
		q         func() *ConcurrentLinkedBlockingQueue[int]
		val       int
		timeout   time.Duration
		wantErr   error
		wantSlice []int
		wantLen   int
	}{
		{
			name: "empty and enqueued",
			q: func() *ConcurrentLinkedBlockingQueue[int] {
				return NewConcurrentLinkedBlockingQueue[int](3)
			},
			val:       123,
			timeout:   time.Second,
			wantSlice: []int{123},
			wantLen:   1,
		},
		{
			name: "invalid context",
			q: func() *ConcurrentLinkedBlockingQueue[int] {
				return NewConcurrentLinkedBlockingQueue[int](3)
			},
			val:       123,
			timeout:   -time.Second,
			wantSlice: []int{},
			wantErr:   context.DeadlineExceeded,
		},
		{
			// 入队之后就满了，恰好放在切片的最后一个位置
			name: "enqueued full last index",
			q: func() *ConcurrentLinkedBlockingQueue[int] {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q := NewConcurrentLinkedBlockingQueue[int](3)
				err := q.Enqueue(ctx, 123)
				require.NoError(t, err)
				err = q.Enqueue(ctx, 234)
				require.NoError(t, err)
				return q
			},
			val:       345,
			timeout:   time.Second,
			wantSlice: []int{123, 234, 345},
			wantLen:   3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()
			q := tc.q()
			err := q.Enqueue(ctx, tc.val)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantSlice, q.AsSlice())
			assert.Equal(t, tc.wantLen, q.Len())
		})
	}

	t.Run("enqueue timeout", func(t *testing.T) {
		q := NewConcurrentLinkedBlockingQueue[int](3)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := q.Enqueue(ctx, 123)
		require.NoError(t, err)
		err = q.Enqueue(ctx, 234)
		require.NoError(t, err)
		err = q.Enqueue(ctx, 345)
		require.NoError(t, err)
		err = q.Enqueue(ctx, 456)
		require.Equal(t, context.DeadlineExceeded, err)
	})

	// 入队阻塞，而后出队，于是入队成功
	t.Run("enqueue blocking and dequeue", func(t *testing.T) {
		q := NewConcurrentLinkedBlockingQueue[int](3)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		go func() {
			time.Sleep(time.Millisecond * 100)
			val, err := q.Dequeue(ctx)
			require.NoError(t, err)
			require.Equal(t, 123, val)
		}()
		err := q.Enqueue(ctx, 123)
		require.NoError(t, err)
		err = q.Enqueue(ctx, 234)
		require.NoError(t, err)
		err = q.Enqueue(ctx, 345)
		require.NoError(t, err)
		err = q.Enqueue(ctx, 456)
		require.NoError(t, err)
	})

	// 无界的情况下，可以无限添加元素，当然小心内存, 以及goroutine调度导致的超时
	// capacity <= 0 时，为无界队列
	t.Run("capacity <= 0", func(t *testing.T) {
		q := NewConcurrentLinkedBlockingQueue[int](-1)
		for i := 0; i < 10; i++ {
			go func() {
				for i := 0; i < 1000; i++ {
					ctx := context.Background()
					val := rand.Int()
					err := q.Enqueue(ctx, val)
					require.NoError(t, err)
				}

			}()
		}
	})
}

func TestConcurrentLinkedBlockingQueue_Dequeue(t *testing.T) {
	testCases := []struct {
		name      string
		q         func() *ConcurrentLinkedBlockingQueue[int]
		val       int
		timeout   time.Duration
		wantErr   error
		wantVal   int
		wantSlice []int
		wantLen   int
	}{
		{
			name: "dequeued",
			q: func() *ConcurrentLinkedBlockingQueue[int] {
				q := NewConcurrentLinkedBlockingQueue[int](3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := q.Enqueue(ctx, 123)
				require.NoError(t, err)
				err = q.Enqueue(ctx, 234)
				require.NoError(t, err)
				return q
			},
			wantVal:   123,
			timeout:   time.Second,
			wantSlice: []int{234},
			wantLen:   1,
		},
		{
			name: "invalid context",
			q: func() *ConcurrentLinkedBlockingQueue[int] {
				q := NewConcurrentLinkedBlockingQueue[int](3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := q.Enqueue(ctx, 123)
				require.NoError(t, err)
				err = q.Enqueue(ctx, 234)
				require.NoError(t, err)
				return q
			},
			wantErr:   context.DeadlineExceeded,
			timeout:   -time.Second,
			wantSlice: []int{123, 234},
			wantLen:   2,
		},
		{
			name: "dequeue and empty first",
			q: func() *ConcurrentLinkedBlockingQueue[int] {
				q := NewConcurrentLinkedBlockingQueue[int](3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := q.Enqueue(ctx, 123)
				require.NoError(t, err)
				return q
			},
			wantVal:   123,
			timeout:   time.Second,
			wantSlice: []int{},
			wantLen:   0,
		},
		{
			name: "dequeue and empty last",
			q: func() *ConcurrentLinkedBlockingQueue[int] {
				q := NewConcurrentLinkedBlockingQueue[int](3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := q.Enqueue(ctx, 123)
				require.NoError(t, err)
				err = q.Enqueue(ctx, 234)
				require.NoError(t, err)
				err = q.Enqueue(ctx, 345)
				require.NoError(t, err)
				val, err := q.Dequeue(ctx)
				require.NoError(t, err)
				require.Equal(t, 123, val)
				val, err = q.Dequeue(ctx)
				require.NoError(t, err)
				require.Equal(t, 234, val)
				return q
			},
			wantVal:   345,
			timeout:   time.Second,
			wantSlice: []int{},
			wantLen:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()
			q := tc.q()
			val, err := q.Dequeue(ctx)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantVal, val)
			assert.Equal(t, tc.wantSlice, q.AsSlice())
			assert.Equal(t, tc.wantLen, q.Len())
		})
	}

	t.Run("dequeue timeout", func(t *testing.T) {
		q := NewConcurrentLinkedBlockingQueue[int](3)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		val, err := q.Dequeue(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
		require.Equal(t, 0, val)
	})

	// 出队阻塞，然后入队，然后出队成功
	t.Run("dequeue blocking and enqueue", func(t *testing.T) {
		q := NewConcurrentLinkedBlockingQueue[int](3)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		go func() {
			time.Sleep(time.Millisecond * 100)
			err := q.Enqueue(ctx, 123)
			require.NoError(t, err)
		}()
		val, err := q.Dequeue(ctx)
		require.NoError(t, err)
		require.Equal(t, 123, val)
	})
}

func TestConcurrentLinkedBlockingQueue(t *testing.T) {
	// 并发测试，只是测试有没有死锁之类的问题
	// 先进先出这个特性依赖于其它单元测试
	// 也依赖于代码审查
	q := NewConcurrentLinkedBlockingQueue[int](100)
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			val := rand.Int()
			err := q.Enqueue(ctx, val)
			cancel()
			require.NoError(t, err)
		}()
	}
	go func() {
		for i := 0; i < 1000; i++ {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				_, err := q.Dequeue(ctx)
				cancel()
				require.NoError(t, err)
				wg.Done()
			}()
		}
	}()
	wg.Wait()
}

func ExampleNewConcurrentLinkedBlockingQueue() {
	// 创建一个容量为 10 的有界并发阻塞队列，如果传入 0 或者负数，那么创建的是无界并发阻塞队列
	q := NewConcurrentLinkedBlockingQueue[int](10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = q.Enqueue(ctx, 22)
	val, err := q.Dequeue(ctx)
	// 这是例子，实际中你不需要写得那么复杂
	switch err {
	case context.Canceled:
		// 有人主动取消了，即调用了 cancel 方法。在这个例子里不会出现这个情况
	case context.DeadlineExceeded:
		// 超时了
	case nil:
		fmt.Println(val)
	default:
		// 其它乱七八糟的
	}
	// Output:
	// 22
}
