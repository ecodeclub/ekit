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

package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ecodeclub/ekit/bean/option"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestOnDemandBlockTaskPool_States(t *testing.T) {
	t.Parallel()
	t.Run("ctx canceled", func(t *testing.T) {
		p1, err := NewOnDemandBlockTaskPool(2, 5)
		assert.Nil(t, err)
		testTaskPoolStatesCtxCanceled(t, p1, context.Canceled)

		p2, err := NewOnDemandBlockTaskPool(2, 5)
		assert.Nil(t, err)
		testTaskPoolStatesCtxRunningCanceled(t, p2,
			State{PoolState: stateRunning, GoCnt: 2,
				WaitingTasksCnt: 3, QueueSize: 5, RunningTasksCnt: 2})
	})

	t.Run("pool not running", func(t *testing.T) {
		p, err := NewOnDemandBlockTaskPool(2, 5)
		assert.Nil(t, err)
		testTaskPoolStatesPoolNotRunning(t, p,
			State{PoolState: stateCreated, GoCnt: 0, WaitingTasksCnt: 5, QueueSize: 5, RunningTasksCnt: 0})
	})

	t.Run("pool Shutdown", func(t *testing.T) {
		p, err := NewOnDemandBlockTaskPool(2, 5)
		assert.Nil(t, err)
		testTaskPoolStatesPoolShutdown(t, p,
			State{PoolState: stateClosing, GoCnt: 2, WaitingTasksCnt: 3, QueueSize: 5, RunningTasksCnt: 2},
			State{PoolState: stateStopped, GoCnt: 0, WaitingTasksCnt: 0, QueueSize: 5, RunningTasksCnt: 0})
	})

	t.Run("pool Shutdown Now", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			p, err := NewOnDemandBlockTaskPool(1, 2)
			assert.Nil(t, err)
			testTaskPoolStatesPoolShutdownNow(t, p,
				State{PoolState: stateStopped, GoCnt: 1, WaitingTasksCnt: 0, QueueSize: 2, RunningTasksCnt: 0})
		}
	})
}

func testTaskPoolStatesCtxCanceled(t *testing.T, pool *OnDemandBlockTaskPool, wantErr error) {
	done := make(chan struct{})
	err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		<-done
		return nil
	}))
	assert.NoError(t, err)

	err = pool.Start()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()
	_, err = pool.States(ctx, time.Millisecond)
	assert.Equal(t, wantErr, err)
	close(done)
}

func testTaskPoolStatesCtxRunningCanceled(t *testing.T, pool *OnDemandBlockTaskPool, wantState State) {
	err := pool.Start()
	assert.NoError(t, err)

	done := make(chan struct{})
	n := cap(pool.queue)
	eg := new(errgroup.Group)

	for i := 0; i < n; i++ {
		eg.Go(func() error {
			return pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				<-done
				return nil
			}))
		})
	}
	_ = eg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ch, err := pool.States(ctx, time.Millisecond)
	assert.NoError(t, err)
	state1 := <-ch
	assert.Equal(t, wantState.PoolState, state1.PoolState)
	assert.Equal(t, wantState.QueueSize, state1.QueueSize)
	assert.Equal(t, wantState.GoCnt, state1.GoCnt)
	assert.Equal(t, wantState.WaitingTasksCnt, state1.WaitingTasksCnt)
	assert.Equal(t, wantState.RunningTasksCnt, state1.RunningTasksCnt)

	cancel()
	for {
		state2, ok := <-ch
		if !ok {
			break
		}
		assert.Equal(t, wantState.PoolState, state2.PoolState)
		assert.Equal(t, wantState.QueueSize, state2.QueueSize)
		assert.Equal(t, wantState.GoCnt, state2.GoCnt)
		assert.Equal(t, wantState.WaitingTasksCnt, state2.WaitingTasksCnt)
		assert.Equal(t, wantState.RunningTasksCnt, state2.RunningTasksCnt)
	}
	close(done)
}

func testTaskPoolStatesPoolNotRunning(t *testing.T, pool *OnDemandBlockTaskPool, wantState State) {
	done := make(chan struct{})
	n := cap(pool.queue)
	eg := new(errgroup.Group)

	for i := 0; i < n; i++ {
		eg.Go(func() error {
			return pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				<-done
				return nil
			}))
		})
	}
	_ = eg.Wait()

	ch, err := pool.States(context.Background(), time.Second)
	assert.NoError(t, err)
	state1 := <-ch
	assert.Equal(t, wantState.PoolState, state1.PoolState)
	assert.Equal(t, wantState.QueueSize, state1.QueueSize)
	assert.Equal(t, wantState.GoCnt, state1.GoCnt)
	assert.Equal(t, wantState.WaitingTasksCnt, state1.WaitingTasksCnt)
	assert.Equal(t, wantState.RunningTasksCnt, state1.RunningTasksCnt)
	close(done)
}

func testTaskPoolStatesPoolShutdown(t *testing.T, pool *OnDemandBlockTaskPool, closingState, stoppedState State) {
	done := make(chan struct{})
	n := cap(pool.queue)
	eg := new(errgroup.Group)

	for i := 0; i < n; i++ {
		eg.Go(func() error {
			return pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				<-done
				return nil
			}))
		})
	}
	_ = eg.Wait()

	err := pool.Start()
	assert.Nil(t, err)
	// 调用 Shutdown 后， 即使 pool 状态关闭也不会调用 shutdownNowCancel
	shutdownDone, err := pool.Shutdown()
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ch, err := pool.States(ctx, time.Millisecond)
	assert.NoError(t, err)
	defer cancel()
	state1 := <-ch
	assert.Equal(t, closingState.PoolState, state1.PoolState)
	assert.Equal(t, closingState.QueueSize, state1.QueueSize)
	assert.Equal(t, closingState.GoCnt, state1.GoCnt)
	assert.Equal(t, closingState.WaitingTasksCnt, state1.WaitingTasksCnt)
	assert.Equal(t, closingState.RunningTasksCnt, state1.RunningTasksCnt)

	close(done)
	<-shutdownDone
	for {
		state2, ok := <-ch
		if !ok {
			break
		}
		assert.Equal(t, stoppedState.PoolState, state2.PoolState)
		assert.Equal(t, stoppedState.QueueSize, state2.QueueSize)
		assert.Equal(t, stoppedState.GoCnt, state2.GoCnt)
		assert.Equal(t, stoppedState.WaitingTasksCnt, state2.WaitingTasksCnt)
		assert.Equal(t, stoppedState.RunningTasksCnt, state2.RunningTasksCnt)
	}
}

func testTaskPoolStatesPoolShutdownNow(t *testing.T, pool *OnDemandBlockTaskPool, stoppedState State) {
	done := make(chan struct{})

	err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		<-done
		return nil
	}))
	assert.NoError(t, err)

	err = pool.Start()
	assert.NoError(t, err)

	ch, err := pool.States(context.Background(), time.Second)
	assert.NoError(t, err)
	done <- struct{}{}
	_, err = pool.ShutdownNow()
	assert.NoError(t, err)

	for {
		state, ok := <-ch
		if !ok {
			break
		}
		assert.Equal(t, stoppedState.PoolState, state.PoolState)
	}

	close(done)
}

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

func TestOnDemandBlockTaskPool_In_Created_State(t *testing.T) {
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

		t.Run("With Options", func(t *testing.T) {
			t.Parallel()

			initGo := 10
			pool, err := NewOnDemandBlockTaskPool(initGo, 10)
			assert.NoError(t, err)

			assert.Equal(t, int32(initGo), pool.initGo)
			assert.Equal(t, int32(initGo), pool.coreGo)
			assert.Equal(t, int32(initGo), pool.maxGo)
			assert.Equal(t, defaultMaxIdleTime, pool.maxIdleTime)

			coreGo, maxGo, maxIdleTime := int32(20), int32(30), 10*time.Second
			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime))
			assert.NoError(t, err)

			assert.Equal(t, int32(initGo), pool.initGo)
			assert.Equal(t, coreGo, pool.coreGo)
			assert.Equal(t, maxGo, pool.maxGo)
			assert.Equal(t, maxIdleTime, pool.maxIdleTime)

			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo))
			assert.NoError(t, err)
			assert.Equal(t, pool.coreGo, pool.maxGo)

			initGo, coreGo = 30, 20
			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithMaxGo(maxGo))
			assert.NoError(t, err)
			assert.Equal(t, pool.maxGo, pool.coreGo)

			initGo, maxGo = 30, 10
			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			initGo, coreGo, maxGo = 30, 20, 10
			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			initGo, coreGo, maxGo = 30, 10, 20
			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			initGo, coreGo, maxGo = 20, 10, 30
			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			initGo, coreGo, maxGo = 20, 30, 10
			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			initGo, coreGo, maxGo = 10, 30, 20
			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithQueueBacklogRate(-0.1))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithQueueBacklogRate(1.0))
			assert.NotNil(t, pool)
			assert.NoError(t, err)

			pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithQueueBacklogRate(1.1))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

		})
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

func TestOnDemandBlockTaskPool_In_Running_State(t *testing.T) {
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

	t.Run("Start —— 在TaskPool启动前队列中已有任务，启动后不再Submit", func(t *testing.T) {

		t.Run("WithCoreGo,WithMaxIdleTime，所需要协程数 <= 允许创建的协程数", func(t *testing.T) {

			initGo, coreGo, maxIdleTime := 1, 3, 3*time.Millisecond
			queueSize := coreGo

			needGo, allowGo := queueSize-initGo, coreGo-initGo
			assert.LessOrEqual(t, needGo, allowGo)

			pool, err := NewOnDemandBlockTaskPool(initGo, queueSize, WithCoreGo(int32(coreGo)), WithMaxIdleTime(maxIdleTime))
			assert.NoError(t, err)

			assert.Equal(t, int32(0), pool.numOfGo())

			done := make(chan struct{}, coreGo)
			wait := make(chan struct{}, coreGo)

			for i := 0; i < coreGo; i++ {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					wait <- struct{}{}
					<-done
					return nil
				}))
				assert.NoError(t, err)
			}

			assert.Equal(t, int32(0), pool.numOfGo())

			assert.NoError(t, pool.Start())

			for i := 0; i < coreGo; i++ {
				<-wait
			}
			assert.Equal(t, int32(coreGo), pool.numOfGo())
		})

		t.Run("WithMaxGo, 所需要协程数 > 允许创建的协程数", func(t *testing.T) {
			initGo, maxGo := 3, 5
			queueSize := maxGo + 1

			needGo, allowGo := queueSize-initGo, maxGo-initGo
			assert.Greater(t, needGo, allowGo)

			pool, err := NewOnDemandBlockTaskPool(initGo, queueSize, WithMaxGo(int32(maxGo)))
			assert.NoError(t, err)

			assert.Equal(t, int32(0), pool.numOfGo())

			done := make(chan struct{}, queueSize)
			wait := make(chan struct{}, queueSize)

			for i := 0; i < queueSize; i++ {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					wait <- struct{}{}
					<-done
					return nil
				}))
				assert.NoError(t, err)
			}

			assert.Equal(t, int32(0), pool.numOfGo())

			assert.NoError(t, pool.Start())

			for i := 0; i < maxGo; i++ {
				<-wait
			}
			assert.Equal(t, int32(maxGo), pool.numOfGo())
		})
	})

	t.Run("Start —— 与Submit并发调用,WithCoreGo,WithMaxIdleTime,WithMaxGo，所需要协程数 < 允许创建的协程数", func(t *testing.T) {

		initGo, coreGo, maxGo, maxIdleTime := 2, 4, 6, 3*time.Millisecond
		queueSize := coreGo

		needGo, allowGo := queueSize-initGo, maxGo-initGo
		assert.Less(t, needGo, allowGo)

		pool, err := NewOnDemandBlockTaskPool(initGo, queueSize, WithCoreGo(int32(coreGo)), WithMaxGo(int32(maxGo)), WithMaxIdleTime(maxIdleTime))
		assert.NoError(t, err)

		assert.Equal(t, int32(0), pool.numOfGo())

		done := make(chan struct{}, queueSize)
		wait := make(chan struct{}, queueSize)

		// 与下方阻塞提交并发调用
		errChan := make(chan error)
		go func() {
			time.Sleep(10 * time.Millisecond)
			errChan <- pool.Start()
		}()

		// 模拟阻塞提交
		for i := 0; i < maxGo; i++ {
			err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				wait <- struct{}{}
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}

		assert.NoError(t, <-errChan)

		for i := 0; i < maxGo; i++ {
			<-wait
		}

		assert.Equal(t, int32(maxGo), pool.numOfGo())
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

	t.Run("工作协程", func(t *testing.T) {
		t.Parallel()

		t.Run("保持在初始数不变", func(t *testing.T) {
			t.Parallel()

			initGo, queueSize := 1, 3
			pool := testNewRunningStateTaskPool(t, initGo, queueSize)

			n := queueSize
			done1 := make(chan struct{}, n)
			wait := make(chan struct{}, n)

			// 队列中有积压任务
			for i := 0; i < n; i++ {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					wait <- struct{}{}
					<-done1
					return nil
				}))
				assert.NoError(t, err)
			}

			// initGo个tasks在运行中
			for i := 0; i < initGo; i++ {
				<-wait
			}

			assert.Equal(t, int32(initGo), pool.numOfGo())

			// 使运行中的tasks结束
			for i := 0; i < initGo; i++ {
				done1 <- struct{}{}
			}

			// 积压在队列中的任务开始运行
			for i := 0; i < n-initGo; i++ {
				<-wait
				assert.Equal(t, int32(initGo), pool.numOfGo())
				done1 <- struct{}{}
			}

		})

		t.Run("从初始数达到核心数", func(t *testing.T) {
			t.Parallel()

			t.Run("核心数比初始数多1个", func(t *testing.T) {
				t.Parallel()

				initGo, coreGo, maxIdleTime, queueBacklogRate := int32(2), int32(3), 3*time.Millisecond, 0.1
				queueSize := int(coreGo)
				testExtendGoFromInitGoToCoreGo(t, initGo, queueSize, coreGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
			})

			t.Run("核心数比初始数多n个", func(t *testing.T) {
				t.Parallel()

				initGo, coreGo, maxIdleTime, queueBacklogRate := int32(2), int32(5), 3*time.Millisecond, 0.1
				queueSize := int(coreGo)
				testExtendGoFromInitGoToCoreGo(t, initGo, queueSize, coreGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
			})

			t.Run("在(初始数,核心数]区间的协程运行完任务后，在等待退出期间再次抢到任务", func(t *testing.T) {
				t.Parallel()

				initGo, coreGo, maxIdleTime := int32(1), int32(6), 100*time.Millisecond
				queueSize := int(coreGo)

				pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, WithCoreGo(coreGo), WithMaxIdleTime(maxIdleTime))

				assert.Equal(t, initGo, pool.numOfGo())
				t.Log("1")
				done := make(chan struct{}, queueSize)
				wait := make(chan struct{}, queueSize)

				for i := 0; i < queueSize; i++ {
					i := i
					err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
						wait <- struct{}{}
						<-done
						t.Log("task done", i)
						return nil
					}))
					assert.NoError(t, err)
				}
				t.Log("2")
				for i := 0; i < queueSize; i++ {
					t.Log("wait ", i)
					<-wait
				}
				assert.Equal(t, coreGo, pool.numOfGo())

				close(done)
				t.Log("3")
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					<-done
					t.Log("task done [x]")
					return nil
				}))
				assert.NoError(t, err)
				t.Log("4")
				// <-time.After(maxIdleTime * 100)
				for pool.numOfGo() > initGo {
					t.Log("loop", "numOfGo", pool.numOfGo(), "timeoutGroup", pool.timeoutGroup.size())
					time.Sleep(maxIdleTime)
				}
				assert.Equal(t, initGo, pool.numOfGo())
			})
		})

		t.Run("从核心数到达最大数", func(t *testing.T) {
			t.Parallel()

			t.Run("最大数比核心数多1个", func(t *testing.T) {
				t.Parallel()

				initGo, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(2), int32(4), int32(5), 3*time.Millisecond, 0.1
				queueSize := int(maxGo)
				testExtendGoFromInitGoToCoreGoAndMaxGo(t, initGo, queueSize, coreGo, maxGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
			})

			t.Run("最大数比核心数多n个", func(t *testing.T) {
				t.Parallel()

				initGo, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(1), int32(3), int32(6), 3*time.Millisecond, 0.1
				queueSize := int(maxGo)
				testExtendGoFromInitGoToCoreGoAndMaxGo(t, initGo, queueSize, coreGo, maxGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
			})
		})
	})

}

func testExtendGoFromInitGoToCoreGo(t *testing.T, initGo int32, queueSize int, coreGo int32, maxIdleTime time.Duration, opts ...option.Option[OnDemandBlockTaskPool]) {

	opts = append(opts, WithCoreGo(coreGo), WithMaxIdleTime(maxIdleTime))
	pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, opts...)

	assert.Equal(t, initGo, pool.numOfGo())

	assert.LessOrEqual(t, initGo, coreGo)

	done := make(chan struct{})
	wait := make(chan struct{}, coreGo)

	// 稳定在initGo
	t.Log("XX")
	for i := int32(0); i < initGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		t.Log("submit ", i)
	}

	t.Log("YY")
	for i := int32(0); i < initGo; i++ {
		<-wait
	}

	// 至少initGo个协程
	assert.GreaterOrEqual(t, pool.numOfGo(), initGo)

	t.Log("ZZ")

	// 逐步添加任务
	for i := int32(1); i <= coreGo-initGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		<-wait
		t.Log("after wait coreGo", coreGo, i, pool.numOfGo())
	}

	t.Log("UU")

	assert.Equal(t, pool.numOfGo(), coreGo)
	close(done)

	// 等待最大空闲时间后稳定在initGo
	for pool.numOfGo() > initGo {
	}

	assert.Equal(t, initGo, pool.numOfGo())
}

func testExtendGoFromInitGoToCoreGoAndMaxGo(t *testing.T, initGo int32, queueSize int, coreGo, maxGo int32, maxIdleTime time.Duration, opts ...option.Option[OnDemandBlockTaskPool]) {

	opts = append(opts, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime))
	pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, opts...)

	assert.Equal(t, initGo, pool.numOfGo())

	assert.LessOrEqual(t, initGo, coreGo)
	assert.LessOrEqual(t, coreGo, maxGo)

	done := make(chan struct{})
	wait := make(chan struct{}, maxGo)

	// 稳定在initGo
	t.Log("00")
	for i := int32(0); i < initGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		t.Log("submit ", i)
	}
	t.Log("AA")
	for i := int32(0); i < initGo; i++ {
		<-wait
	}

	assert.GreaterOrEqual(t, pool.numOfGo(), initGo)

	t.Log("BB")

	// 逐步添加任务
	for i := int32(1); i <= coreGo-initGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		<-wait
		t.Log("after wait coreGo", coreGo, i, pool.numOfGo())
	}

	t.Log("CC")

	assert.GreaterOrEqual(t, pool.numOfGo(), coreGo)

	for i := int32(1); i <= maxGo-coreGo; i++ {

		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		<-wait
		t.Log("after wait maxGo", maxGo, i, pool.numOfGo())
	}

	t.Log("DD")

	assert.Equal(t, pool.numOfGo(), maxGo)
	close(done)

	// 等待最大空闲时间后稳定在initGo
	for pool.numOfGo() > initGo {
	}
	assert.Equal(t, initGo, pool.numOfGo())
}

func TestOnDemandBlockTaskPool_In_Closing_State(t *testing.T) {
	t.Parallel()

	t.Run("Shutdown —— 使TaskPool状态由Running变为Closing", func(t *testing.T) {
		t.Parallel()

		initGo, queueSize := 2, 4
		pool := testNewRunningStateTaskPool(t, initGo, queueSize)

		// 模拟阻塞提交
		n := initGo + queueSize*2
		eg := new(errgroup.Group)
		waitChan := make(chan struct{}, n)
		taskDone := make(chan struct{})
		for i := 0; i < n; i++ {
			eg.Go(func() error {
				return pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					waitChan <- struct{}{}
					<-taskDone
					return nil
				}))
			})
		}
		for i := 0; i < initGo; i++ {
			<-waitChan
		}
		done, err := pool.Shutdown()
		assert.NoError(t, err)
		// Closing过程中Submit会报错间接证明TaskPool处于StateClosing状态
		assert.ErrorIs(t, eg.Wait(), errTaskPoolIsClosing)

		// 第二次调用
		done2, err2 := pool.Shutdown()
		assert.Nil(t, done2)
		assert.ErrorIs(t, err2, errTaskPoolIsClosing)
		assert.Equal(t, stateClosing, pool.internalState())

		assert.Equal(t, int32(initGo), pool.numOfGo())

		close(taskDone)
		<-done
		assert.Equal(t, stateStopped, pool.internalState())

		// 第一个Shutdown将状态迁移至StateStopped
		// 第三次调用
		done3, err := pool.Shutdown()
		assert.Nil(t, done3)
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
	})

	t.Run("Shutdown —— 协程数按需扩展至maxGo，调用Shutdown成功后，所有协程运行完任务后可以自动退出", func(t *testing.T) {
		t.Parallel()

		initGo, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(1), int32(3), int32(5), 10*time.Millisecond, 0.1
		queueSize := int(maxGo)
		pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime), WithQueueBacklogRate(queueBacklogRate))

		assert.LessOrEqual(t, initGo, coreGo)
		assert.LessOrEqual(t, coreGo, maxGo)

		taskDone := make(chan struct{})
		wait := make(chan struct{})

		for i := int32(0); i < maxGo; i++ {
			err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				wait <- struct{}{}
				<-taskDone
				return nil
			}))
			assert.NoError(t, err)
		}

		// 提交任务后立即Shutdown
		shutdownDone, err := pool.Shutdown()
		assert.NoError(t, err)

		// 已提交的任务应该正常运行并能扩展至maxGo
		for i := int32(0); i < maxGo; i++ {
			<-wait
		}
		assert.Equal(t, maxGo, pool.numOfGo())

		// 让所有任务结束
		close(taskDone)
		<-shutdownDone

		// 用循环取代time.After/time.Sleep
		for pool.numOfGo() != 0 {

		}

		// 最终全部退出
		assert.Equal(t, int32(0), pool.numOfGo())
	})

	t.Run("Start", func(t *testing.T) {
		t.Parallel()

		pool, wait := testNewRunningStateTaskPoolWithQueueFullFilled(t, 2, 4)

		done, err := pool.Shutdown()
		assert.NoError(t, err)

		select {
		case <-done:
		default:
			assert.ErrorIs(t, pool.Start(), errTaskPoolIsClosing)
		}

		close(wait)
		<-done
		assert.Equal(t, stateStopped, pool.internalState())
	})

	// Submit()在状态中会报错，因为Closing是一个中间状态，故在上面的Shutdown间接测到了

	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()

		pool, wait := testNewRunningStateTaskPoolWithQueueFullFilled(t, 2, 4)

		done, err := pool.Shutdown()
		assert.NoError(t, err)

		select {
		case <-done:
		default:
			tasks, err := pool.ShutdownNow()
			assert.ErrorIs(t, err, errTaskPoolIsClosing)
			assert.Nil(t, tasks)
		}

		close(wait)
		<-done
		assert.Equal(t, stateStopped, pool.internalState())
	})
}

func TestOnDemandBlockTaskPool_In_Stopped_State(t *testing.T) {
	t.Parallel()

	t.Run("ShutdownNow —— 使TaskPool状态由Running变为Stopped", func(t *testing.T) {
		t.Parallel()

		initGo, queueSize := 2, 4
		pool, wait := testNewRunningStateTaskPoolWithQueueFullFilled(t, initGo, queueSize)

		// 模拟阻塞提交
		eg := new(errgroup.Group)
		for i := 0; i < queueSize; i++ {
			eg.Go(func() error {
				return pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					<-wait
					return nil
				}))
			})
		}

		// 并发调用ShutdownNow
		result := make(chan ShutdownNowResult, 1)
		go func() {
			tasks, err := pool.ShutdownNow()
			result <- ShutdownNowResult{tasks: tasks, err: err}
			close(wait)
		}()

		// 阻塞的Submit在ShutdownNow后会报错间接证明TaskPool处于StateStopped状态
		assert.ErrorIs(t, eg.Wait(), errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())

		r := <-result
		assert.NoError(t, r.err)
		assert.Equal(t, queueSize, len(r.tasks))

		// 重复调用
		tasks, err := pool.ShutdownNow()
		assert.Nil(t, tasks)
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})

	t.Run("ShutdownNow —— 工作协程数扩展至maxGo后，调用ShutdownNow成功，所有协程能够接收到信号", func(t *testing.T) {
		t.Parallel()

		initGo, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(1), int32(3), int32(5), 10*time.Millisecond, 0.1
		queueSize := int(maxGo)
		pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime), WithQueueBacklogRate(queueBacklogRate))

		assert.LessOrEqual(t, initGo, coreGo)
		assert.LessOrEqual(t, coreGo, maxGo)

		taskDone := make(chan struct{})
		wait := make(chan struct{}, queueSize)

		for i := 0; i < queueSize; i++ {
			err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				wait <- struct{}{}
				<-taskDone
				return nil
			}))
			assert.NoError(t, err)
		}

		tasks, err := pool.ShutdownNow()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(tasks), 0)

		// 让所有任务结束
		close(taskDone)

		// 用循环取代time.After/time.Sleep
		for pool.numOfGo() != 0 {
		}

		assert.Equal(t, int32(0), pool.numOfGo())
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
	done := make(chan struct{})
	err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		<-done
		return nil
	}))
	assert.NoError(t, err)

	n := cap(pool.queue) + 2
	eg := new(errgroup.Group)

	for i := 0; i < n; i++ {
		eg.Go(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()
			return pool.Submit(ctx, TaskFunc(func(ctx context.Context) error {
				<-done
				return nil
			}))
		})
	}
	assert.ErrorIs(t, eg.Wait(), context.DeadlineExceeded)
	close(done)
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

func testNewRunningStateTaskPool(t *testing.T, initGo int, queueSize int, opts ...option.Option[OnDemandBlockTaskPool]) *OnDemandBlockTaskPool {
	pool, _ := NewOnDemandBlockTaskPool(initGo, queueSize, opts...)
	assert.Equal(t, stateCreated, pool.internalState())
	assert.NoError(t, pool.Start())
	assert.Equal(t, stateRunning, pool.internalState())
	return pool
}

func testNewStoppedStateTaskPool(t *testing.T, initGo int, queueSize int) *OnDemandBlockTaskPool {
	pool := testNewRunningStateTaskPool(t, initGo, queueSize)
	tasks, err := pool.ShutdownNow()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tasks))
	assert.Equal(t, stateStopped, pool.internalState())
	return pool
}

func testNewRunningStateTaskPoolWithQueueFullFilled(t *testing.T, initGo int, queueSize int) (*OnDemandBlockTaskPool, chan struct{}) {
	pool := testNewRunningStateTaskPool(t, initGo, queueSize)
	wait := make(chan struct{})
	for i := 0; i < initGo+queueSize; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			<-wait
			return nil
		}))
		assert.NoError(t, err)
	}
	return pool, wait
}

func TestGroup(t *testing.T) {
	n := 10

	// g := &sliceGroup{members: make([]int, n, n)}
	g := &group{mp: make(map[int]int)}

	for i := 0; i < n; i++ {
		assert.False(t, g.isIn(i))
		g.add(i)
		assert.True(t, g.isIn(i))
		assert.Equal(t, int32(i+1), g.size())
	}

	assert.Equal(t, int32(n), g.size())

	for i := 0; i < n; i++ {
		g.delete(i)
		assert.Equal(t, int32(n-i-1), g.size())
	}

	assert.Equal(t, int32(0), g.size())

	assert.False(t, g.isIn(n+1))

	id := 100
	g.add(id)
	assert.Equal(t, int32(1), g.size())
	assert.True(t, g.isIn(id))
	g.delete(id)
	assert.Equal(t, int32(0), g.size())
}

func ExampleNewOnDemandBlockTaskPool() {
	p, _ := NewOnDemandBlockTaskPool(10, 100)
	_ = p.Start()
	// wg 只是用来确保任务执行的，你在实际使用过程中是不需要的
	var wg sync.WaitGroup
	wg.Add(1)
	_ = p.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		fmt.Println("hello, world")
		wg.Done()
		return nil
	}))
	wg.Wait()
	// Output:
	// hello, world
}

func ExampleOnDemandBlockTaskPool_States() {
	p, _ := NewOnDemandBlockTaskPool(10, 100)
	var wg sync.WaitGroup
	wg.Add(1)
	_ = p.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		wg.Done()
		return nil
	}))
	_ = p.Start()
	ch, err := p.States(context.Background(), time.Second*10)
	if err == nil {
		fmt.Println("get ch")
	}
	state := <-ch
	fmt.Println(state.PoolState)
	fmt.Println(state.RunningTasksCnt)
	fmt.Println(state.WaitingTasksCnt)
	fmt.Println(state.GoCnt)
	fmt.Println(state.QueueSize)
	wg.Wait()

	// Output:
	// get ch
	// 2
	// 0
	// 0
	// 10
	// 100
}
