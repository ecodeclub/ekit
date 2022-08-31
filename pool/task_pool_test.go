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

package pool

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
TaskPool有限状态机
                                                       Start/Submit/ShutdownNow() Error
                                                                \     /
                                               Shutdown()  --> CLOSING  ---等待所有任务结束
         Submit()nil--执行中状态迁移--Submit()      /    \----------/ \----------/
           \    /                    \   /      /
New() --> CREATED -- Start() --->  RUNNING -- --
           \   /                    \   /       \           Start/Submit/Shutdown() Error
  Shutdown/ShutdownNow()Error      Start()       \                \    /
                                               ShutdownNow() ---> STOPPED  -- ShutdownNow() --> STOPPED
*/

func TestTaskPool_In_Created_State(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		pool, err := NewOnDemandBlockTaskPool(1, -1)
		assert.ErrorIs(t, err, errInvalidArgument)
		assert.Nil(t, pool)

		pool, err = NewOnDemandBlockTaskPool(1, 0)
		assert.NoError(t, err)
		assert.NotNil(t, pool)

		pool, err = NewOnDemandBlockTaskPool(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, pool)

		pool, err = NewOnDemandBlockTaskPool(-1, 1)
		assert.ErrorIs(t, err, errInvalidArgument)
		assert.Nil(t, pool)

		pool, err = NewOnDemandBlockTaskPool(0, 1)
		assert.ErrorIs(t, err, errInvalidArgument)
		assert.Nil(t, pool)

		pool, err = NewOnDemandBlockTaskPool(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, pool)
	})

	// Start()导致TaskPool状态迁移，测试见TestTaskPool_In_Running_State/Start

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		t.Run("提交非法Task", func(t *testing.T) {
			t.Parallel()

			pool, _ := NewOnDemandBlockTaskPool(1, 1)
			assert.Equal(t, stateCreated, pool.internalState())
			assert.ErrorIs(t, pool.Submit(context.Background(), nil), errTaskIsInvalid)
			assert.Equal(t, stateCreated, pool.internalState())
		})

		t.Run("正常提交Task", func(t *testing.T) {
			t.Parallel()

			pool, _ := NewOnDemandBlockTaskPool(1, 3)
			assert.Equal(t, stateCreated, pool.internalState())
			testSubmitValidTask(t, pool)
			assert.Equal(t, stateCreated, pool.internalState())
		})

		t.Run("阻塞提交并导致超时", func(t *testing.T) {
			t.Parallel()

			pool, _ := NewOnDemandBlockTaskPool(1, 1)
			assert.Equal(t, stateCreated, pool.internalState())
			testSubmitBlockingAndTimeout(t, pool)
			assert.Equal(t, stateCreated, pool.internalState())
		})
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Parallel()

		pool, err := NewOnDemandBlockTaskPool(1, 1)
		assert.NoError(t, err)
		assert.Equal(t, stateCreated, pool.internalState())

		done, err := pool.Shutdown()
		assert.Nil(t, done)
		assert.ErrorIs(t, err, errTaskPoolIsNotRunning)
		assert.Equal(t, stateCreated, pool.internalState())
	})

	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()

		pool, err := NewOnDemandBlockTaskPool(1, 1)
		assert.NoError(t, err)
		assert.Equal(t, stateCreated, pool.internalState())

		tasks, err := pool.ShutdownNow()
		assert.Nil(t, tasks)
		assert.ErrorIs(t, err, errTaskPoolIsNotRunning)
		assert.Equal(t, stateCreated, pool.internalState())
	})
}

func TestTaskPool_In_Running_State(t *testing.T) {
	t.Parallel()

	t.Run("Start —— 使TaskPool状态由Created变为Running", func(t *testing.T) {
		t.Parallel()

		pool, _ := NewOnDemandBlockTaskPool(1, 1)

		// 与下方 testSubmitBlockingAndTimeout 并发执行
		errChan := make(chan error)
		go func() {
			time.Sleep(1 * time.Millisecond)
			errChan <- pool.Start()
		}()

		assert.Equal(t, stateCreated, pool.internalState())

		testSubmitBlockingAndTimeout(t, pool)

		assert.NoError(t, <-errChan)
		assert.Equal(t, stateRunning, pool.internalState())

		// 重复调用
		assert.ErrorIs(t, pool.Start(), errTaskPoolIsStarted)
		assert.Equal(t, stateRunning, pool.internalState())
	})

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		t.Run("提交非法Task", func(t *testing.T) {
			t.Parallel()

			pool := testNewRunningStateTaskPool(t, 1, 1)
			assert.ErrorIs(t, pool.Submit(context.Background(), nil), errTaskIsInvalid)
			assert.Equal(t, stateRunning, pool.internalState())
		})

		t.Run("正常提交Task", func(t *testing.T) {
			t.Parallel()

			pool := testNewRunningStateTaskPool(t, 1, 3)
			testSubmitValidTask(t, pool)
			assert.Equal(t, stateRunning, pool.internalState())
		})

		t.Run("阻塞提交并导致超时", func(t *testing.T) {
			t.Parallel()

			pool := testNewRunningStateTaskPool(t, 1, 1)
			testSubmitBlockingAndTimeout(t, pool)
			assert.Equal(t, stateRunning, pool.internalState())
		})
	})

	// Shutdown()导致TaskPool状态迁移，TestTaskPool_In_Closing_State/Shutdown

	// ShutdownNow()导致TaskPool状态迁移，TestTestPool_In_Stopped_State/ShutdownNow
}

func TestTaskPool_In_Closing_State(t *testing.T) {
	t.Parallel()

	t.Run("Shutdown —— 使TaskPool状态由Running变为Closing", func(t *testing.T) {
		t.Parallel()

		queueSize := 2
		pool := testNewRunningStateTaskPool(t, 1, queueSize)

		// 模拟阻塞提交
		n := queueSize * 5
		firstSubmitErrChan := make(chan error, 1)
		for i := 0; i < n; i++ {
			go func() {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					time.Sleep(10 * time.Millisecond)
					return nil
				}))
				if err != nil {
					firstSubmitErrChan <- err
				}
			}()
		}

		// 调用Shutdown使TaskPool状态发生迁移
		type ShutdownResult struct {
			done <-chan struct{}
			err  error
		}
		resultChan := make(chan ShutdownResult)
		go func() {
			time.Sleep(time.Millisecond)
			done, err := pool.Shutdown()
			resultChan <- ShutdownResult{done: done, err: err}
		}()
		r := <-resultChan

		// Closing过程中Submit会报错间接证明TaskPool处于StateClosing状态
		assert.ErrorIs(t, <-firstSubmitErrChan, errTaskPoolIsClosing)

		// Shutdown调用成功
		assert.NoError(t, r.err)
		select {
		case <-r.done:
			break
		default:
			// 第二次调用
			done2, err2 := pool.Shutdown()
			assert.Nil(t, done2)
			assert.ErrorIs(t, err2, errTaskPoolIsClosing)
			assert.Equal(t, stateClosing, pool.internalState())
		}

		<-r.done
		assert.Equal(t, stateStopped, pool.internalState())

		// 第一个Shutdown将状态迁移至StateStopped
		// 第三次调用
		done, err := pool.Shutdown()
		assert.Nil(t, done)
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
	})

	t.Run("Start", func(t *testing.T) {
		t.Parallel()

		queueSize := 10
		pool := testNewRunningStateTaskPool(t, 2, queueSize)

		// 提交任务
		for i := 0; i < queueSize; i++ {
			go func() {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					time.Sleep(10 * time.Millisecond)
					return nil
				}))
				if err != nil {
					return
				}
			}()
		}

		done, err := pool.Shutdown()
		assert.NoError(t, err)

		select {
		case <-done:
		default:
			assert.ErrorIs(t, pool.Start(), errTaskPoolIsClosing)
		}

		<-done
		assert.Equal(t, stateStopped, pool.internalState())
	})

	// Submit()在状态中会报错，因为Closing是一个中间状态，故在上面的Shutdown间接测到了

	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()

		concurrency := 2
		pool := testNewRunningStateTaskPool(t, concurrency, 0)

		// 提交任务
		num := concurrency * 5
		for i := 0; i < num; i++ {
			go func() {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					time.Sleep(10 * time.Millisecond)
					return nil
				}))
				if err != nil {
					return
				}
			}()
		}

		done, err := pool.Shutdown()
		assert.NoError(t, err)

		select {
		case <-done:
		default:
			tasks, err := pool.ShutdownNow()
			assert.ErrorIs(t, err, errTaskPoolIsClosing)
			assert.Nil(t, tasks)
		}

		<-done
		assert.Equal(t, stateStopped, pool.internalState())
	})
}

func TestTestPool_In_Stopped_State(t *testing.T) {
	t.Parallel()

	t.Run("ShutdownNow —— 使TaskPool状态由Running变为Stopped", func(t *testing.T) {
		t.Parallel()

		concurrency, queueSize := 2, 4
		pool := testNewRunningStateTaskPool(t, concurrency, queueSize)

		// 模拟阻塞提交
		n := queueSize + 6
		firstSubmitErrChan := make(chan error, concurrency)
		var submitWg sync.WaitGroup
		submitWg.Add(concurrency)
		for i := 0; i < n; i++ {
			go func() {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					submitWg.Done()
					time.Sleep(20 * time.Millisecond)
					return nil
				}))
				if err != nil {
					firstSubmitErrChan <- err
				}
			}()
		}

		submitWg.Wait()
		assert.Equal(t, int32(concurrency), pool.NumGo())

		// 并发调用ShutdownNow
		result := make(chan ShutdownNowResult, 1)
		go func() {
			tasks, err := pool.ShutdownNow()
			result <- ShutdownNowResult{tasks: tasks, err: err}
		}()

		r := <-result
		assert.NoError(t, r.err)
		assert.Greater(t, len(r.tasks), 0)

		// 阻塞的Submit在ShutdownNow后会报错间接证明TaskPool处于StateStopped状态
		assert.ErrorIs(t, <-firstSubmitErrChan, errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())

		// 重复调用
		tasks, err := pool.ShutdownNow()
		assert.Nil(t, tasks)
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})

	t.Run("Start", func(t *testing.T) {
		t.Parallel()

		pool := testNewStoppedStateTaskPool(t, 1, 1)
		assert.ErrorIs(t, pool.Start(), errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		pool := testNewStoppedStateTaskPool(t, 1, 1)
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error { return nil }))
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Parallel()

		pool := testNewStoppedStateTaskPool(t, 1, 1)
		done, err := pool.Shutdown()
		assert.Nil(t, done)
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})
}

func testSubmitBlockingAndTimeout(t *testing.T, pool *OnDemandBlockTaskPool) {

	err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		time.Sleep(2 * time.Millisecond)
		return nil
	}))
	assert.NoError(t, err)

	n := len(pool.queue) + 1
	errChan := make(chan error, n)
	for i := 0; i < n; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()
			err := pool.Submit(ctx, TaskFunc(func(ctx context.Context) error {
				time.Sleep(2 * time.Millisecond)
				return nil
			}))
			if err != nil {
				errChan <- err
			}
		}()
	}

	assert.ErrorIs(t, <-errChan, context.DeadlineExceeded)
}

func testSubmitValidTask(t *testing.T, pool *OnDemandBlockTaskPool) {

	err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error { return nil }))
	assert.NoError(t, err)

	err = pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error { panic("task panic") }))
	assert.NoError(t, err)

	err = pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error { return errors.New("fake error") }))
	assert.NoError(t, err)
}

type ShutdownNowResult struct {
	tasks []Task
	err   error
}

func testNewRunningStateTaskPool(t *testing.T, concurrency int, queueSize int) *OnDemandBlockTaskPool {
	pool, _ := NewOnDemandBlockTaskPool(concurrency, queueSize)
	assert.Equal(t, stateCreated, pool.internalState())
	assert.NoError(t, pool.Start())
	assert.Equal(t, stateRunning, pool.internalState())
	return pool
}

func testNewStoppedStateTaskPool(t *testing.T, concurrency int, queueSize int) *OnDemandBlockTaskPool {
	pool := testNewRunningStateTaskPool(t, concurrency, queueSize)
	_, err := pool.ShutdownNow()
	assert.NoError(t, err)
	assert.Equal(t, stateStopped, pool.internalState())
	return pool
}

type FakeTask struct{}

func (f *FakeTask) Run(_ context.Context) error { return nil }
