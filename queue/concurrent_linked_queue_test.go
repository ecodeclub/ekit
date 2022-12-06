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
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcurrentQueue_Enqueue(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		q    func() *ConcurrentLinkedQueue[int]
		val  int

		timeout  time.Duration
		wantData []int
		wantErr  error
	}{
		{
			name: "empty",
			q: func() *ConcurrentLinkedQueue[int] {
				return NewConcurrentLinkedQueue[int]()
			},
			timeout:  time.Second,
			val:      123,
			wantData: []int{123},
		},
		{
			name: "multiple",
			q: func() *ConcurrentLinkedQueue[int] {
				q := NewConcurrentLinkedQueue[int]()
				ctx, cancle := context.WithTimeout(context.Background(), time.Second)
				defer cancle()
				err := q.Enqueue(ctx, 123)
				require.NoError(t, err)
				return q
			},
			timeout:  time.Second,
			val:      234,
			wantData: []int{123, 234},
		},
		{
			// 已经超时了的 context 设置
			name: "invalid contex",
			q: func() *ConcurrentLinkedQueue[int] {
				return NewConcurrentLinkedQueue[int]()
			},
			timeout: -time.Second,
			val:     123,
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := tc.q()
			ctx, cancle := context.WithTimeout(context.Background(), tc.timeout)
			defer cancle()
			err := q.Enqueue(ctx, tc.val)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantData, q.asSlice())
		})
	}
}

func TestConcurrentQueue_Dequeue(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		q        func() *ConcurrentLinkedQueue[int]
		timeout  time.Duration
		wantVal  int
		wantData []int
		wantErr  error
	}{
		{
			name: "empty",
			q: func() *ConcurrentLinkedQueue[int] {
				q := NewConcurrentLinkedQueue[int]()
				return q
			},
			timeout: time.Second,
			wantErr: errEmptyQueue,
		},
		{
			name: "single",
			q: func() *ConcurrentLinkedQueue[int] {
				q := NewConcurrentLinkedQueue[int]()
				ctx, cancle := context.WithTimeout(context.Background(), time.Second)
				defer cancle()
				err := q.Enqueue(ctx, 123)
				assert.NoError(t, err)
				return q
			},
			timeout: time.Second,
			wantVal: 123,
		},
		{
			name: "multiple",
			q: func() *ConcurrentLinkedQueue[int] {
				q := NewConcurrentLinkedQueue[int]()
				ctx, cancle := context.WithTimeout(context.Background(), time.Second)
				defer cancle()
				err := q.Enqueue(ctx, 123)
				assert.NoError(t, err)
				err = q.Enqueue(ctx, 234)
				assert.NoError(t, err)
				return q
			},
			timeout:  time.Second,
			wantVal:  123,
			wantData: []int{234},
		},
		{
			name: "enqueue and dequeue",
			q: func() *ConcurrentLinkedQueue[int] {
				q := NewConcurrentLinkedQueue[int]()
				ctx := context.Background()
				err := q.Enqueue(ctx, 123)
				assert.NoError(t, err)
				err = q.Enqueue(ctx, 234)
				assert.NoError(t, err)
				val, err := q.Dequeue(ctx)
				assert.Equal(t, 123, val)
				assert.NoError(t, err)
				err = q.Enqueue(ctx, 345)
				assert.NoError(t, err)
				return q
			},
			timeout:  time.Second,
			wantVal:  234,
			wantData: []int{345},
		},
		{
			// context 本身已经过期了
			name: "invalid context",
			q: func() *ConcurrentLinkedQueue[int] {
				q := NewConcurrentLinkedQueue[int]()
				ctx, cancle := context.WithTimeout(context.Background(), time.Second)
				defer cancle()
				err := q.Enqueue(ctx, 123)
				assert.NoError(t, err)
				return q
			},
			timeout: -time.Second,
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := tc.q()
			ctx, cancle := context.WithTimeout(context.Background(), tc.timeout)
			defer cancle()
			val, err := q.Dequeue(ctx)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
			assert.Equal(t, tc.wantData, q.asSlice())
		})
	}
}

func TestConcurrentLinkedQueue(t *testing.T) {
	t.Parallel()
	// 仅仅是为了测试在入队出队期间不会出现 panic 或者死循环之类的问题
	// FIFO 特性参考其余测试
	q := NewConcurrentLinkedQueue[int]()
	var wg sync.WaitGroup
	wg.Add(10000)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				ctx := context.Background()
				val := rand.Int()
				_ = q.Enqueue(ctx, val)
			}
		}()
	}
	var cnt int32 = 0
	for i := 0; i < 10; i++ {
		go func() {
			for {
				if atomic.LoadInt32(&cnt) >= 10000 {
					return
				}
				ctx := context.Background()
				_, err := q.Dequeue(ctx)
				if err == nil {
					atomic.AddInt32(&cnt, 1)
					wg.Done()
				}
			}
		}()
	}
	wg.Wait()
}

func (c *ConcurrentLinkedQueue[T]) asSlice() []T {
	var res []T
	cur := (*node[T])((*node[T])(c.head).next)
	for cur != nil {
		res = append(res, cur.val)
		cur = (*node[T])(cur.next)
	}
	return res
}
