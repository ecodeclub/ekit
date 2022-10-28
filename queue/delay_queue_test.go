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
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestDelayQueue_New(t *testing.T) {
	t.Parallel()

	_, err := testNewDelayQueue[*Int](0)
	assert.ErrorIs(t, err, errInvalidArgument)

	_, err = testNewDelayQueue[*Int](-1)
	assert.ErrorIs(t, err, errInvalidArgument)

	_, err = NewDelayQueue[*Int](1, nil)
	assert.ErrorIs(t, err, errInvalidArgument)

	q, err := testNewDelayQueue[*Int](1)
	assert.NoError(t, err)
	assert.Equal(t, 0, q.Len())
}

func TestDelayQueue_Enqueue(t *testing.T) {
	t.Parallel()

	t.Run("超时控制", func(t *testing.T) {
		t.Parallel()

		t.Run("直接超时", func(t *testing.T) {
			t.Parallel()

			testPassCanceledCtxToEnqueueOrDequeueOperation(t, func(q *DelayQueue[*Int], i int, ctx context.Context) error {
				return q.Enqueue(ctx, newInt(i, time.Second))
			})
		})

		t.Run("等待超时", func(t *testing.T) {
			t.Parallel()

			testContextTimeoutInEnqueueOrDequeueOperation(t, func(q *DelayQueue[*Int], i int, ctx context.Context) error {
				return q.Enqueue(ctx, newInt(i, 100*time.Millisecond))
			})

		})
	})

	t.Run("Enqueue与Dequeue串行", func(t *testing.T) {
		t.Parallel()

		t.Run("无Dequeue，Enqueue之间串行", func(t *testing.T) {
			t.Parallel()

			t.Run("队列未满", func(t *testing.T) {
				t.Parallel()

				capacity, numOfEnqueue, expectedError := 3, 2, (error)(nil)
				testCallEnqueueSequentially(t, capacity, numOfEnqueue, expectedError, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueue int) {
					assert.Equal(t, numOfEnqueue, q.Len())
				})
			})

			t.Run("队列刚满", func(t *testing.T) {
				t.Parallel()

				capacity, numOfEnqueue, expectedError := 3, 3, (error)(nil)
				testCallEnqueueSequentially(t, capacity, numOfEnqueue, expectedError, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueue int) {
					assert.Equal(t, numOfEnqueue, q.Len())
				})
			})

			t.Run("队列超满", func(t *testing.T) {
				t.Parallel()

				capacity, numOfEnqueue, expectedError := 3, 4, errOutOfCapacity
				testCallEnqueueSequentially(t, capacity, numOfEnqueue, expectedError, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueue int) {
					assert.Equal(t, capacity, q.Len())
				})
			})
		})

		t.Run("无Dequeue, Enqueue之间并发", func(t *testing.T) {
			t.Parallel()

			t.Run("队列未满", func(t *testing.T) {
				t.Parallel()

				capacity, numOfEnqueueGo, expectedErr := 3, 2, (error)(nil)
				testCallEnqueueConcurrently(t, capacity, numOfEnqueueGo, expectedErr, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int) {
					assert.Equal(t, numOfEnqueueGo, q.Len())
				})
			})

			t.Run("队列刚满", func(t *testing.T) {
				t.Parallel()

				capacity, numOfEnqueueGo, expectedErr := 3, 3, (error)(nil)
				testCallEnqueueConcurrently(t, capacity, numOfEnqueueGo, expectedErr, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int) {
					assert.Equal(t, numOfEnqueueGo, q.Len())
				})
			})

			t.Run("队列超满", func(t *testing.T) {
				t.Parallel()

				capacity, numOfEnqueueGo, expectedErr := 3, 5, errOutOfCapacity
				testCallEnqueueConcurrently(t, capacity, numOfEnqueueGo, expectedErr, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int) {
					assert.Equal(t, capacity, q.Len())
				})
			})
		})

		// 有Dequeue，Enqueue之间串行，Dequeue之间串行/并发（无所谓）详见下方 TestDelayQueue_Dequeue/Enqueue与Dequeue串行

		// 有Dequeue，Enqueue之间并发，Dequeue之间串行/并行（无所谓）详见下方 TestDelayQueue_Dequeue/Enqueue与Dequeue串行

	})

	// Enqueue与Dequeue之间并发，详见下发 TestDelayQueue_Enqueue_Dequeue
}

func TestDelayQueue_Dequeue(t *testing.T) {

	t.Run("超时控制", func(t *testing.T) {
		t.Parallel()

		t.Run("直接超时", func(t *testing.T) {
			t.Parallel()

			testPassCanceledCtxToEnqueueOrDequeueOperation(t, func(q *DelayQueue[*Int], i int, ctx context.Context) error {
				_, err := q.Dequeue(ctx)
				return err
			})
		})

		t.Run("等待超时", func(t *testing.T) {
			t.Parallel()

			testContextTimeoutInEnqueueOrDequeueOperation(t, func(q *DelayQueue[*Int], i int, ctx context.Context) error {
				_, err := q.Dequeue(ctx)
				return err
			})
		})
	})

	t.Run("Enqueue与Dequeue串行", func(t *testing.T) {
		t.Parallel()

		t.Run("无Enqueue，Dequeue之间串行", func(t *testing.T) {
			t.Parallel()

			t.Run("队列为空，Dequeue协程超时退出", func(t *testing.T) {
				t.Parallel()

				capacity, numOfDequeue := 1, 3
				q, err := testNewDelayQueue[*Int](capacity)
				assert.NoError(t, err)
				assert.Equal(t, 0, q.Len())

				for i := 0; i < numOfDequeue; i++ {
					func() {
						ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond)
						defer cancelFunc()
						_, err := q.Dequeue(ctx)
						assert.Equal(t, context.DeadlineExceeded, err)
					}()
				}

				assert.Equal(t, 0, q.Len())
			})
		})

		t.Run("无Enqueue，Dequeue之间并发", func(t *testing.T) {
			t.Parallel()

			t.Run("队列为空，并发Dequeue协程超时退出", func(t *testing.T) {
				t.Parallel()

				capacity, numOfDequeueGo := 1, 3
				q, err := testNewDelayQueue[*Int](capacity)
				assert.NoError(t, err)
				assert.Equal(t, 0, q.Len())

				errChan := make(chan error, numOfDequeueGo)
				for i := 0; i < numOfDequeueGo; i++ {
					go func() {
						ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond)
						defer cancelFunc()
						_, err := q.Dequeue(ctx)
						errChan <- err
					}()
				}

				for i := 0; i < numOfDequeueGo; i++ {
					assert.Equal(t, context.DeadlineExceeded, <-errChan)
				}

				assert.Equal(t, 0, q.Len())
			})

		})

		t.Run("有Enqueue, Enqueue之间并发，Dequeue之间串行", func(t *testing.T) {
			t.Parallel()

			t.Run("队列从满到空", func(t *testing.T) {
				t.Parallel()

				n := 3
				q, err := testNewDelayQueue[*Int](n)
				assert.NoError(t, err)

				// 并发Enqueue
				var eg errgroup.Group
				for i := 0; i < n; i++ {
					i := i
					eg.Go(func() error {
						return q.Enqueue(context.Background(), newInt(i, time.Millisecond))
					})
				}

				assert.NoError(t, eg.Wait())
				// 队列已满
				assert.Equal(t, n, q.Len())

				// 与上方Enqueue串行，与后续Dequeue也串行
				for i := 0; i < n; i++ {
					d, err := q.Dequeue(context.Background())
					assert.NoError(t, err)
					assert.True(t, d.isExpired(time.Now()))
				}

				assert.Equal(t, 0, q.Len())
			})
		})

		t.Run("有Enqueue，Enqueue之间串行，Dequeue之间并发", func(t *testing.T) {
			t.Parallel()

			t.Run("队列从满到空", func(t *testing.T) {
				t.Parallel()

				n := 3
				q, err := testNewDelayQueue[*Int](n)
				assert.NoError(t, err)

				// 串行Enqueue
				for i := 0; i < n; i++ {
					duration := 10*time.Millisecond + time.Duration(i)*10*time.Millisecond
					assert.NoError(t, q.Enqueue(context.Background(), newInt(i, duration)))
				}

				// 队列已满
				assert.Equal(t, n, q.Len())

				// 未调用Dequeue方法， dequeueProxy 协程一定不存在
				assert.Equal(t, int64(0), atomic.LoadInt64(&q.numOfDequeueProxyGo))

				var eg errgroup.Group
				expiredElements := make(chan *Int, n)
				for i := 0; i < n; i++ {
					eg.Go(func() error {
						e, err := q.Dequeue(context.Background())
						expiredElements <- e
						return err
					})
				}

				assert.NoError(t, eg.Wait())
				assert.Equal(t, 0, q.Len())

				// 最后一个Dequeue返回后，dequeueProxy 协程必须退出
				assert.Equal(t, int64(0), atomic.LoadInt64(&q.numOfDequeueProxyGo))

				now := time.Now()

				// 取出的元素必须过期，间接验证 dequeueProxy 协程被创建
				for i := 0; i < n; i++ {
					elem := <-expiredElements
					assert.True(t, elem.isExpired(now))
				}
			})
		})

	})

	// Enqueue与Dequeue之间并发，详见下发 TestDelayQueue_Enqueue_Dequeue
}

func TestDelayQueue_Enqueue_Dequeue(t *testing.T) {
	t.Parallel()

	increasingOrderExpireFunc := func(i int) time.Duration {
		return 1000*time.Millisecond + time.Duration(i)*200*time.Millisecond
	}

	getDecreasingOrderExpireFunc := func(numOfEnqueue int) func(i int) time.Duration {
		return func(i int) time.Duration {
			return 20*time.Duration(numOfEnqueue)*time.Millisecond - time.Duration(i)*20*time.Millisecond
		}
	}

	fixedExpireFunc := func(i int) time.Duration {
		return time.Millisecond
	}

	backgroundCtxFunc := func() (context.Context, context.CancelFunc) {
		return context.Background(), func() {}
	}

	dequeueCtxFunc := func() (context.Context, context.CancelFunc) {
		return context.WithTimeout(context.Background(), 1000*time.Millisecond)
	}

	// 注意：并发测试的结果有是很难预测，尽最大努力构造重复的测试场景

	t.Run("Enqueue与Dequeue并发", func(t *testing.T) {
		t.Parallel()

		t.Run("1:N", func(t *testing.T) {
			t.Parallel()

			t.Run("初始队列为空，1个Dequeue阻塞，N个Enqueue并发，验证数据过期且队列长度为N-1", func(t *testing.T) {
				t.Parallel()

				capacity, numOfEnqueueGo, numOfDequeueGo := 4, 4, 1
				testCallEnqueueAndDequeueConcurrently(t, capacity, numOfEnqueueGo, numOfDequeueGo, increasingOrderExpireFunc, backgroundCtxFunc, backgroundCtxFunc, nil, nil, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int, numOfDequeueGo int) {
					assert.Equal(t, numOfEnqueueGo-numOfDequeueGo, q.Len())
				})
			})

			t.Run("初始队列为空，1个Enqueue阻塞，N个Dequeue并发，验证数据过期且队列为空", func(t *testing.T) {
				t.Parallel()

				capacity, numOfEnqueueGo, numOfDequeueGo := 6, 1, 6

				testCallEnqueueAndDequeueConcurrently(t, capacity, numOfEnqueueGo, numOfDequeueGo, fixedExpireFunc, backgroundCtxFunc, dequeueCtxFunc, nil, context.DeadlineExceeded, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int, numOfDequeueGo int) {
					assert.Equal(t, 0, q.Len())
				})
			})

		})

		t.Run("N:N", func(t *testing.T) {
			t.Parallel()

			t.Run("初始队列为空，1个Dequeue阻塞，1个Enqueue并发，验证数据过期且队列为空", func(t *testing.T) {
				t.Parallel()

				capacity, numOfEnqueueGo, numOfDequeueGo := 1, 1, 1
				testCallEnqueueAndDequeueConcurrently(t, capacity, numOfEnqueueGo, numOfDequeueGo, fixedExpireFunc, backgroundCtxFunc, backgroundCtxFunc, nil, nil, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int, numOfDequeueGo int) {
					assert.Equal(t, 0, q.Len())
				})
			})

			t.Run("初始队列为空，N个Dequeue阻塞，N个Enqueue并发，验证数据过期且队列为空", func(t *testing.T) {
				t.Parallel()

				capacity, numOfEnqueueGo, numOfDequeueGo := 10, 10, 10
				testCallEnqueueAndDequeueConcurrently(t, capacity, numOfEnqueueGo, numOfDequeueGo, getDecreasingOrderExpireFunc(numOfEnqueueGo), backgroundCtxFunc, backgroundCtxFunc, nil, nil, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int, numOfDequeueGo int) {
					assert.Equal(t, 0, q.Len())
				})
			})
		})

		t.Run("N:M", func(t *testing.T) {
			t.Parallel()

			t.Run("Enqueue多，Dequeue少", func(t *testing.T) {
				t.Parallel()

				t.Run("新元素入队，队头改变，返回新队头", func(t *testing.T) {
					t.Parallel()

					capacity, numOfEnqueueGo, numOfDequeueGo := 100, 80, 50
					testCallEnqueueAndDequeueConcurrently(t, capacity, numOfEnqueueGo, numOfDequeueGo, getDecreasingOrderExpireFunc(numOfEnqueueGo), backgroundCtxFunc, backgroundCtxFunc, nil, nil, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int, numOfDequeueGo int) {
						assert.Equal(t, numOfEnqueueGo-numOfDequeueGo, q.Len())
					})
				})
			})

			t.Run("Dequeue多，Enqueue少", func(t *testing.T) {
				t.Parallel()

				t.Run("新元素入队，队头改变，返回新队头", func(t *testing.T) {
					t.Parallel()
					capacity, numOfEnqueueGo, numOfDequeueGo := 10, 3, 9
					testCallEnqueueAndDequeueConcurrently(t, capacity, numOfEnqueueGo, numOfDequeueGo, getDecreasingOrderExpireFunc(numOfEnqueueGo), backgroundCtxFunc, dequeueCtxFunc, nil, context.DeadlineExceeded, func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int, numOfDequeueGo int) {
						assert.Equal(t, 0, q.Len())
					})
				})
			})
		})
	})
}

func testCallEnqueueSequentially(t *testing.T, capacity int, numOfEnqueue int, expectedError error, lenAssertFunc func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueue int)) {

	q, err := testNewDelayQueue[*Int](capacity)
	assert.NoError(t, err)

	for i := 0; i < numOfEnqueue-1; i++ {
		assert.NoError(t, q.Enqueue(context.Background(), newInt(i, time.Microsecond)))
		// 顺序调用，Enqueue调用者即是第一个（启动 enqueueProxy）又是最后一个（关闭 enqueueProxy）
		assert.Equal(t, int64(0), atomic.LoadInt64(&q.numOfEnqueueProxyGo))
	}

	// 超时控制
	ctx, cancelFunc := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancelFunc()
	assert.Equal(t, expectedError, q.Enqueue(ctx, newInt(numOfEnqueue, time.Microsecond)))

	lenAssertFunc(t, q, capacity, numOfEnqueue)
}

func testCallEnqueueConcurrently(t *testing.T, capacity int, numOfEnqueueGo int, expectedErr error, lenAssertFunc func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int)) {

	q, err := testNewDelayQueue[*Int](capacity)
	assert.NoError(t, err)

	// Enqueue之前，enqueueProxy 协程必须未启动
	assert.Equal(t, int64(0), atomic.LoadInt64(&q.numOfEnqueueProxyGo))

	var eg errgroup.Group
	for i := 0; i < numOfEnqueueGo; i++ {
		i := i
		eg.Go(func() error {
			err := q.Enqueue(context.Background(), newInt(i, time.Microsecond))
			// syncChan <- struct{}{}
			return err
		})
	}

	assert.Equal(t, expectedErr, eg.Wait())

	// 最后一个Enqueue方法返回后，enqueueProxy 协程必须退出
	assert.Equal(t, int64(0), atomic.LoadInt64(&q.numOfEnqueueProxyGo))

	// 通过队列的长度来间接验证 enqueueProxy 协程正常工作
	lenAssertFunc(t, q, capacity, numOfEnqueueGo)
}

func testCallEnqueueAndDequeueConcurrently(t *testing.T, capacity int, numOfEnqueueGo int, numOfDequeueGo int, expireFunc func(i int) time.Duration,
	enqueueCtxFunc func() (context.Context, context.CancelFunc), dequeueCtxFunc func() (context.Context, context.CancelFunc),
	enqueueError error, dequeueError error, qLenAssertFunc func(t *testing.T, q *DelayQueue[*Int], capacity int, numOfEnqueueGo int, numOfDequeueGo int)) {

	q, err := testNewDelayQueue[*Int](capacity)
	assert.NoError(t, err)
	assert.Equal(t, 0, q.Len())

	// 未调用Dequeue方法， dequeueProxy 协程一定不存在
	assert.Equal(t, int64(0), atomic.LoadInt64(&q.numOfDequeueProxyGo))
	// 未调用Enqueue方法，enqueueProxy 协程一定不存在
	assert.Equal(t, int64(0), atomic.LoadInt64(&q.numOfEnqueueProxyGo))

	var dequeueErrGroup errgroup.Group
	expiredElements := make(chan *Int, numOfDequeueGo)

	for i := 0; i < numOfDequeueGo; i++ {
		dequeueErrGroup.Go(func() error {
			ctx, cancelFunc := dequeueCtxFunc()
			defer cancelFunc()
			e, err := q.Dequeue(ctx)
			if err == nil {
				expiredElements <- e
			}
			return err
		})
	}

	// 并发入队
	var enqueueErrGroup errgroup.Group
	for i := 0; i < numOfEnqueueGo; i++ {
		i := i
		enqueueErrGroup.Go(func() error {
			ctx, cancelFunc := enqueueCtxFunc()
			defer cancelFunc()
			return q.Enqueue(ctx, newInt(i, expireFunc(i)))
		})
	}

	assert.Equal(t, dequeueError, dequeueErrGroup.Wait())
	now := time.Now()
	assert.Equal(t, enqueueError, enqueueErrGroup.Wait())

	// 最后一个Dequeue返回后，dequeueProxy 协程必须退出
	assert.Equal(t, int64(0), atomic.LoadInt64(&q.numOfDequeueProxyGo))
	// 最后一个Enqueue返回后，enqueueProxy 协程必须退出
	assert.Equal(t, int64(0), atomic.LoadInt64(&q.numOfEnqueueProxyGo))

	// 取出的元素必须过期
	// 间接证明 dequeueProxy 协程与 enqueueProxy 协程正常工作
	n := len(expiredElements)
	for i := 0; i < n; i++ {
		elem := <-expiredElements
		assert.True(t, elem.isExpired(now))
	}

	qLenAssertFunc(t, q, capacity, numOfEnqueueGo, numOfDequeueGo)
}

func testPassCanceledCtxToEnqueueOrDequeueOperation(t *testing.T, op func(q *DelayQueue[*Int], i int, ctx context.Context) error) {
	q, err := testNewDelayQueue[*Int](1)
	assert.NoError(t, err)

	createContextFns := []func() (context.Context, context.CancelFunc){
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(context.Background(), time.Nanosecond)
		},
		func() (context.Context, context.CancelFunc) {
			return context.WithDeadline(context.Background(), time.Now().Add(1))
		},
	}

	errChan := make(chan error, len(createContextFns))
	for i, fn := range createContextFns {
		i, fn := i, fn
		go func() {
			ctx, cancelFunc := fn()
			cancelFunc()
			errChan <- op(q, i, ctx)
		}()
	}

	for i := 0; i < len(createContextFns); i++ {
		assert.Error(t, <-errChan)
	}
}

func testContextTimeoutInEnqueueOrDequeueOperation(t *testing.T, op func(q *DelayQueue[*Int], i int, ctx context.Context) error) {
	var eg errgroup.Group

	q, err := testNewDelayQueue[*Int](10)
	assert.NoError(t, err)

	n := 10
	waitChan := make(chan int, n)
	for i := 0; i < n; i++ {
		i := i
		eg.Go(func() error {
			waitChan <- i
			timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Microsecond)
			defer cancelFunc()
			return op(q, i, timeoutCtx)
		})
	}
	for i := 0; i < n; i++ {
		<-waitChan
	}
	assert.Equal(t, context.DeadlineExceeded, eg.Wait())
}

type Int struct {
	id       int
	deadline time.Time
}

func newInt(id int, expire time.Duration) *Int {
	return &Int{id: id, deadline: time.Now().Add(expire)}
}

func (i *Int) Deadline() time.Time {
	return i.deadline
}

func (i *Int) isExpired(now time.Time) bool {
	return now.Unix()-i.deadline.Unix() >= 0
}

func testNewDelayQueue[T Delayable[T]](capacity int) (*DelayQueue[T], error) {
	return NewDelayQueue[T](capacity, func(t1 T, t2 T) int {
		t1Unix := t1.Deadline().Unix()
		t2Unix := t2.Deadline().Unix()
		if t1Unix < t2Unix {
			return -1
		} else if t1Unix == t2Unix {
			return 0
		} else {
			return 1
		}
	})
}
