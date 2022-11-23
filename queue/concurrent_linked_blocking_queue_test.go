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
		q         func() *ConcurrentLinkBlockingQueue[int]
		val       int
		timeout   time.Duration
		wantErr   error
		wantSlice []int
		wantLen   int
		wantHead  *Node[int]
		wantTail  *Node[int]
	}{
		{
			name: "empty and enqueued",
			q: func() *ConcurrentLinkBlockingQueue[int] {
				return NewConcurrentLinkBlockingQueue[int](3)
			},
			val:       123,
			timeout:   time.Second,
			wantSlice: []int{123},
			wantLen:   1,
			wantTail:  &Node[int]{data: 123, next: nil},
			wantHead:  &Node[int]{data: 0, next: &Node[int]{data: 123, next: nil}},
		},
		{
			name: "invalid context",
			q: func() *ConcurrentLinkBlockingQueue[int] {
				return NewConcurrentLinkBlockingQueue[int](3)
			},
			val:       123,
			timeout:   -time.Second,
			wantSlice: []int{},
			wantErr:   context.DeadlineExceeded,
			wantTail:  &Node[int]{data: 0, next: nil},
			wantHead:  &Node[int]{data: 0, next: nil},
		},
		{
			// 入队之后就满了，恰好放在切片的最后一个位置
			name: "enqueued full last index",
			q: func() *ConcurrentLinkBlockingQueue[int] {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q := NewConcurrentLinkBlockingQueue[int](3)
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
			wantTail:  &Node[int]{data: 345, next: nil},
			wantHead:  &Node[int]{data: 0, next: &Node[int]{data: 123, next: &Node[int]{data: 234, next: &Node[int]{data: 345, next: nil}}}},
		},
		{
			// 入队之后就满了，恰好放在中间
			name: "enqueued full first index",
			q: func() *ConcurrentLinkBlockingQueue[int] {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q := NewConcurrentLinkBlockingQueue[int](3)
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
				err = q.Enqueue(ctx, 456)
				require.NoError(t, err)
				return q
			},
			val:       567,
			timeout:   time.Second,
			wantSlice: []int{345, 456, 567},
			wantLen:   3,
			wantTail:  &Node[int]{data: 567, next: nil},
			wantHead:  &Node[int]{data: 0, next: &Node[int]{data: 345, next: &Node[int]{data: 456, next: &Node[int]{data: 567, next: nil}}}},
		},
		{
			// 元素本身就是零值
			name: "all zero value ",
			q: func() *ConcurrentLinkBlockingQueue[int] {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q := NewConcurrentLinkBlockingQueue[int](3)
				err := q.Enqueue(ctx, 0)
				require.NoError(t, err)
				err = q.Enqueue(ctx, 0)
				require.NoError(t, err)
				return q
			},
			val:       0,
			timeout:   time.Second,
			wantSlice: []int{0, 0, 0},
			wantLen:   3,
			wantTail:  &Node[int]{data: 0, next: nil},
			wantHead:  &Node[int]{data: 0, next: &Node[int]{data: 0, next: &Node[int]{data: 0, next: &Node[int]{data: 0, next: nil}}}},
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
			assert.Equal(t, tc.wantHead, q.head)
			assert.Equal(t, tc.wantTail, q.tail)
		})
	}

	t.Run("enqueue timeout", func(t *testing.T) {
		q := NewConcurrentLinkBlockingQueue[int](3)
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
		q := NewConcurrentLinkBlockingQueue[int](3)
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
}

func TestConcurrentLinkedBlockingQueue_Dequeue(t *testing.T) {
	testCases := []struct {
		name      string
		q         func() *ConcurrentLinkBlockingQueue[int]
		val       int
		timeout   time.Duration
		wantErr   error
		wantVal   int
		wantSlice []int
		wantLen   int
		wantHead  *Node[int]
		wantTail  *Node[int]
	}{
		{
			name: "dequeued",
			q: func() *ConcurrentLinkBlockingQueue[int] {
				q := NewConcurrentLinkBlockingQueue[int](3)
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
			wantTail:  &Node[int]{data: 234, next: nil},
			wantHead:  &Node[int]{data: 0, next: &Node[int]{data: 234, next: nil}},
		},
		{
			name: "invalid context",
			q: func() *ConcurrentLinkBlockingQueue[int] {
				q := NewConcurrentLinkBlockingQueue[int](3)
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
			wantTail:  &Node[int]{data: 234, next: nil},
			wantHead:  &Node[int]{data: 0, next: &Node[int]{data: 123, next: &Node[int]{data: 234, next: nil}}},
		},
		{
			name: "dequeue and empty first",
			q: func() *ConcurrentLinkBlockingQueue[int] {
				q := NewConcurrentLinkBlockingQueue[int](3)
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
			wantTail:  &Node[int]{data: 0, next: nil},
			wantHead:  &Node[int]{data: 0, next: nil},
		},
		{
			name: "dequeue and empty middle",
			q: func() *ConcurrentLinkBlockingQueue[int] {
				q := NewConcurrentLinkBlockingQueue[int](3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := q.Enqueue(ctx, 123)
				require.NoError(t, err)
				err = q.Enqueue(ctx, 234)
				require.NoError(t, err)
				val, err := q.Dequeue(ctx)
				require.NoError(t, err)
				require.Equal(t, 123, val)
				return q
			},
			wantVal:   234,
			timeout:   time.Second,
			wantSlice: []int{},
			wantLen:   0,
			wantTail:  &Node[int]{data: 0, next: nil},
			wantHead:  &Node[int]{data: 0, next: nil},
		},
		{
			name: "dequeue and empty last",
			q: func() *ConcurrentLinkBlockingQueue[int] {
				q := NewConcurrentLinkBlockingQueue[int](3)
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
			wantTail:  &Node[int]{data: 0, next: nil},
			wantHead:  &Node[int]{data: 0, next: nil},
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
			assert.Equal(t, tc.wantHead, q.head)
			assert.Equal(t, tc.wantTail, q.tail)
		})
	}

	t.Run("dequeue timeout", func(t *testing.T) {
		q := NewConcurrentLinkBlockingQueue[int](3)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		val, err := q.Dequeue(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
		require.Equal(t, 0, val)
	})

	// 出队阻塞，然后入队，然后出队成功
	t.Run("dequeue blocking and enqueue", func(t *testing.T) {
		q := NewConcurrentLinkBlockingQueue[int](3)
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
	q := NewConcurrentLinkBlockingQueue[int](100)
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
