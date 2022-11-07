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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelayQueue_Dequeue(t *testing.T) {
	t.Parallel()
	now := time.Now()
	testCases := []struct {
		name    string
		q       *DelayQueue[delayElem]
		timeout time.Duration
		wantVal int
		wantErr error
	}{
		{
			name: "dequeued",
			q: newDelayQueue(t, delayElem{
				deadline: now.Add(time.Millisecond * 10),
				val:      11,
			}),
			timeout: time.Second,
			wantVal: 11,
		},
		{
			// 元素本身就已经过期了
			name: "already deadline",
			q: newDelayQueue(t, delayElem{
				deadline: now.Add(-time.Millisecond * 10),
				val:      11,
			}),
			timeout: time.Second,
			wantVal: 11,
		},
		{
			// 已经超时了的 context 设置
			name: "invalid context",
			q: newDelayQueue(t, delayElem{
				deadline: now.Add(time.Millisecond * 10),
				val:      11,
			}),
			timeout: -time.Second,
			wantErr: context.DeadlineExceeded,
		},
		{
			name:    "empty and timeout",
			q:       NewDelayQueue[delayElem](10),
			timeout: time.Second,
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "not empty but timeout",
			q: newDelayQueue(t, delayElem{
				deadline: now.Add(time.Second * 10),
				val:      11,
			}),
			timeout: time.Second,
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tt := range testCases {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()
			ele, err := tc.q.Dequeue(ctx)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, ele.val)
		})
	}

	// 最开始没有元素，然后进去了一个元素
	t.Run("dequeue while enqueue", func(t *testing.T) {
		q := NewDelayQueue[delayElem](3)
		go func() {
			time.Sleep(time.Millisecond * 500)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err := q.Enqueue(ctx, delayElem{
				val:      123,
				deadline: time.Now().Add(time.Millisecond * 100),
			})
			require.NoError(t, err)
		}()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		ele, err := q.Dequeue(ctx)
		require.NoError(t, err)
		require.Equal(t, 123, ele.val)
	})

	// 进去了一个更加短超时时间的元素
	// 于是后面两个都会拿出来，但是时间短的会先拿出来
	t.Run("enqueue short ele", func(t *testing.T) {
		q := NewDelayQueue[delayElem](3)
		// 长时间过期的元素
		err := q.Enqueue(context.Background(), delayElem{
			val:      234,
			deadline: time.Now().Add(time.Second),
		})
		require.NoError(t, err)

		go func() {
			time.Sleep(time.Millisecond * 200)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err := q.Enqueue(ctx, delayElem{
				val:      123,
				deadline: time.Now().Add(time.Millisecond * 300),
			})
			require.NoError(t, err)
		}()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		// 先拿出短时间的
		ele, err := q.Dequeue(ctx)
		require.NoError(t, err)
		require.Equal(t, 123, ele.val)
		require.True(t, ele.deadline.Before(time.Now()))

		// 再拿出长时间的
		ele, err = q.Dequeue(ctx)
		require.NoError(t, err)
		require.Equal(t, 234, ele.val)
		require.True(t, ele.deadline.Before(time.Now()))

		// 没有元素了，会超时
		_, err = q.Dequeue(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
	})
}

func TestDelayQueue_Enqueue(t *testing.T) {
	t.Parallel()
	now := time.Now()
	testCases := []struct {
		name    string
		q       *DelayQueue[delayElem]
		timeout time.Duration
		val     delayElem
		wantErr error
	}{
		{
			name:    "enqueued",
			q:       NewDelayQueue[delayElem](3),
			timeout: time.Second,
			val:     delayElem{val: 123, deadline: now.Add(time.Minute)},
		},
		{
			// context 本身已经过期了
			name:    "invalid context",
			q:       NewDelayQueue[delayElem](3),
			timeout: -time.Second,
			val:     delayElem{val: 123, deadline: now.Add(time.Minute)},
			wantErr: context.DeadlineExceeded,
		},
		{
			// enqueue 的时候阻塞住了，直到超时
			name:    "enqueue timeout",
			q:       newDelayQueue(t, delayElem{val: 123, deadline: now.Add(time.Minute)}),
			timeout: time.Millisecond * 100,
			val:     delayElem{val: 234, deadline: now.Add(time.Minute)},
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()
			err := tc.q.Enqueue(ctx, tc.val)
			assert.Equal(t, tc.wantErr, err)
		})
	}

	// 队列满了，这时候入队。
	// 在等待一段时间之后，队列元素被取走一个
	t.Run("enqueue while dequeue", func(t *testing.T) {
		t.Parallel()
		q := newDelayQueue(t, delayElem{val: 123, deadline: time.Now().Add(time.Second)})
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()
			ele, err := q.Dequeue(ctx)
			require.NoError(t, err)
			require.Equal(t, 123, ele.val)
		}()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		err := q.Enqueue(ctx, delayElem{val: 345, deadline: time.Now().Add(time.Millisecond * 1500)})
		require.NoError(t, err)
	})

	// 入队相同过期时间的元素
	// 但是因为我们在入队的时候是分别计算 Delay 的
	// 那么就会导致虽然过期时间是相同的，但是因为调用 Delay 有先后之分
	// 所以会造成 dstDelay 就是要比 srcDelay 小一点
	t.Run("enqueue with same deadline", func(t *testing.T) {
		t.Parallel()
		q := NewDelayQueue[delayElem](3)
		deadline := time.Now().Add(time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		err := q.Enqueue(ctx, delayElem{val: 123, deadline: deadline})
		require.NoError(t, err)
		err = q.Enqueue(ctx, delayElem{val: 456, deadline: deadline})
		require.NoError(t, err)
		err = q.Enqueue(ctx, delayElem{val: 789, deadline: deadline})
		require.NoError(t, err)

		ele, err := q.Dequeue(ctx)
		require.NoError(t, err)
		require.Equal(t, 123, ele.val)

		ele, err = q.Dequeue(ctx)
		require.NoError(t, err)
		require.Equal(t, 789, ele.val)

		ele, err = q.Dequeue(ctx)
		require.NoError(t, err)
		require.Equal(t, 456, ele.val)
	})
}

func newDelayQueue(t *testing.T, eles ...delayElem) *DelayQueue[delayElem] {
	q := NewDelayQueue[delayElem](len(eles))
	for _, ele := range eles {
		err := q.Enqueue(context.Background(), ele)
		require.NoError(t, err)
	}
	return q
}

type delayElem struct {
	deadline time.Time
	val      int
}

func (d delayElem) Delay() time.Duration {
	return time.Until(d.deadline)
}
