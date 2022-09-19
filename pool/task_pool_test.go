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
	"testing"
	"time"

	"github.com/gotomicro/ekit/bean/option"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
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

			concurrency := 10
			pool, err := NewOnDemandBlockTaskPool(concurrency, 10)
			assert.NoError(t, err)

			assert.Equal(t, int32(concurrency), pool.initGo)
			assert.Equal(t, int32(concurrency), pool.coreGo)
			assert.Equal(t, int32(concurrency), pool.maxGo)
			assert.Equal(t, defaultMaxIdleTime, pool.maxIdleTime)

			coreGo, maxGo, maxIdleTime := int32(20), int32(30), 10*time.Second
			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime))
			assert.NoError(t, err)

			assert.Equal(t, int32(concurrency), pool.initGo)
			assert.Equal(t, coreGo, pool.coreGo)
			assert.Equal(t, maxGo, pool.maxGo)
			assert.Equal(t, maxIdleTime, pool.maxIdleTime)

			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithCoreGo(coreGo))
			assert.NoError(t, err)
			assert.Equal(t, pool.coreGo, pool.maxGo)

			concurrency, coreGo = 30, 20
			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithCoreGo(coreGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithMaxGo(maxGo))
			assert.NoError(t, err)
			assert.Equal(t, pool.maxGo, pool.coreGo)

			concurrency, maxGo = 30, 10
			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			concurrency, coreGo, maxGo = 30, 20, 10
			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			concurrency, coreGo, maxGo = 30, 10, 20
			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			concurrency, coreGo, maxGo = 20, 10, 30
			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			concurrency, coreGo, maxGo = 20, 30, 10
			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			concurrency, coreGo, maxGo = 10, 30, 20
			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithQueueBacklogRate(-0.1))
			assert.Nil(t, pool)
			assert.ErrorIs(t, err, errInvalidArgument)

			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithQueueBacklogRate(1.0))
			assert.NotNil(t, pool)
			assert.NoError(t, err)

			pool, err = NewOnDemandBlockTaskPool(concurrency, 10, WithQueueBacklogRate(1.1))
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

			concurrency, queueSize := 1, 3
			pool := testNewRunningStateTaskPool(t, concurrency, queueSize)

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

			// concurrency个tasks在运行中
			for i := 0; i < concurrency; i++ {
				<-wait
			}

			assert.Equal(t, int32(concurrency), pool.numOfGo())

			// 使运行中的tasks结束
			for i := 0; i < concurrency; i++ {
				done1 <- struct{}{}
			}

			// 积压在队列中的任务开始运行
			for i := 0; i < n-concurrency; i++ {
				<-wait
				assert.Equal(t, int32(concurrency), pool.numOfGo())
				done1 <- struct{}{}
			}

		})

		t.Run("从初始数达到核心数", func(t *testing.T) {
			t.Parallel()

			t.Run("一次性全开", func(t *testing.T) {
				t.Parallel()

				t.Run("核心数比初始数多1个", func(t *testing.T) {
					t.Parallel()

					concurrency, coreGo, maxIdleTime := int32(1), int32(2), 3*time.Millisecond
					queueSize := int(coreGo)
					testExtendNumGoFromInitGoToCoreGoAtOnce(t, concurrency, queueSize, coreGo, maxIdleTime)
				})

				t.Run("核心数比初始数多n个", func(t *testing.T) {
					t.Parallel()

					concurrency, coreGo, maxIdleTime := int32(1), int32(3), 3*time.Millisecond
					queueSize := int(coreGo)
					testExtendNumGoFromInitGoToCoreGoAtOnce(t, concurrency, queueSize, coreGo, maxIdleTime)
				})
			})

			t.Run("一次一个开", func(t *testing.T) {
				t.Parallel()

				t.Run("核心数比初始数多1个", func(t *testing.T) {
					concurrency, coreGo, maxIdleTime, queueBacklogRate := int32(2), int32(3), 3*time.Millisecond, 0.1
					queueSize := int(coreGo)
					testExtendNumGoFromInitGoToCoreGoStepByStep(t, concurrency, queueSize, coreGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
				})

				t.Run("核心数比初始数多n个", func(t *testing.T) {
					t.Parallel()

					concurrency, coreGo, maxIdleTime, queueBacklogRate := int32(2), int32(5), 3*time.Millisecond, 0.1
					queueSize := int(coreGo)
					testExtendNumGoFromInitGoToCoreGoStepByStep(t, concurrency, queueSize, coreGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))

				})
			})

			t.Run("在(初始数,核心数]区间的协程运行完任务后，在等待退出期间再次抢到任务", func(t *testing.T) {
				t.Parallel()

				concurrency, coreGo, maxIdleTime := int32(1), int32(6), 100*time.Millisecond
				queueSize := int(coreGo)

				pool := testNewRunningStateTaskPool(t, int(concurrency), queueSize, WithCoreGo(coreGo), WithMaxIdleTime(maxIdleTime))
				done := make(chan struct{}, queueSize)
				wait := make(chan struct{}, queueSize)

				for i := 0; i < queueSize; i++ {
					// i := i
					err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
						wait <- struct{}{}
						<-done
						// t.Log("task done", i)
						return nil
					}))
					assert.NoError(t, err)
				}

				for i := 0; i < queueSize; i++ {
					<-wait
				}
				assert.Equal(t, coreGo, pool.numOfGo())

				close(done)

				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					<-done
					// t.Log("task done [x]")
					return nil
				}))
				assert.NoError(t, err)

				// <-time.After(maxIdleTime * 100)
				for pool.numOfGo() > concurrency {
					// t.Log("loop", "numOfGo", pool.numOfGo(), "timeoutGroup", pool.timeoutGroup.size())
					time.Sleep(maxIdleTime)
				}
				assert.Equal(t, concurrency, pool.numOfGo())
			})
		})

		t.Run("从核心数到达最大数", func(t *testing.T) {
			t.Parallel()

			t.Run("一次性全开", func(t *testing.T) {
				t.Parallel()

				t.Run("最大数比核心数多1个", func(t *testing.T) {
					t.Parallel()

					concurrency, coreGo, maxGo, maxIdleTime := int32(2), int32(4), int32(5), 3*time.Millisecond
					queueSize := int(maxGo)
					testExtendNumGoFromInitGoToCoreGoAndMaxGoAtOnce(t, concurrency, queueSize, coreGo, maxGo, maxIdleTime)
				})

				t.Run("最大数比核心数多n个", func(t *testing.T) {
					t.Parallel()

					concurrency, coreGo, maxGo, maxIdleTime := int32(2), int32(3), int32(5), 3*time.Millisecond
					queueSize := int(maxGo)
					testExtendNumGoFromInitGoToCoreGoAndMaxGoAtOnce(t, concurrency, queueSize, coreGo, maxGo, maxIdleTime)
				})

			})

			t.Run("一次一个开", func(t *testing.T) {
				t.Parallel()

				t.Run("最大数比核心数多1个", func(t *testing.T) {
					t.Parallel()

					concurrency, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(2), int32(4), int32(5), 3*time.Millisecond, 0.1
					queueSize := int(maxGo)
					testExtendNumGoFromInitGoToCoreGoAndMaxGoStepByStep(t, concurrency, queueSize, coreGo, maxGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
				})

				t.Run("最大数比核心数多n个", func(t *testing.T) {
					t.Parallel()

					concurrency, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(1), int32(3), int32(6), 3*time.Millisecond, 0.1
					queueSize := int(maxGo)
					testExtendNumGoFromInitGoToCoreGoAndMaxGoStepByStep(t, concurrency, queueSize, coreGo, maxGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
				})
			})
		})
	})

}

type extendStrategyCheckFunc func(t *testing.T, i int32, pool *OnDemandBlockTaskPool)

func testExtendNumGoFromInitGoToCoreGoAtOnce(t *testing.T, concurrency int32, queueSize int, coreGo int32, maxIdleTime time.Duration, opts ...option.Option[OnDemandBlockTaskPool]) {
	extendToCoreGoAtOnce := func(t *testing.T, i int32, pool *OnDemandBlockTaskPool) {
		assert.Equal(t, coreGo, pool.numOfGo())
	}
	opts = append(opts, WithCoreGo(coreGo), WithMaxIdleTime(maxIdleTime))
	testExtendNumGoFromInitGoToCoreGoAndMaxGo(t, concurrency, queueSize, coreGo, coreGo, extendToCoreGoAtOnce, nil, opts...)
}

func testExtendNumGoFromInitGoToCoreGoStepByStep(t *testing.T, concurrency int32, queueSize int, coreGo int32, maxIdleTime time.Duration, opts ...option.Option[OnDemandBlockTaskPool]) {
	extendToCoreGoAtOnce := func(t *testing.T, i int32, pool *OnDemandBlockTaskPool) {
		assert.Equal(t, i, pool.numOfGo())
	}
	opts = append(opts, WithCoreGo(coreGo), WithMaxIdleTime(maxIdleTime))
	testExtendNumGoFromInitGoToCoreGoAndMaxGo(t, concurrency, queueSize, coreGo, coreGo, extendToCoreGoAtOnce, nil, opts...)
}

func testExtendNumGoFromInitGoToCoreGoAndMaxGoAtOnce(t *testing.T, concurrency int32, queueSize int, coreGo int32, maxGo int32, maxIdleTime time.Duration, opts ...option.Option[OnDemandBlockTaskPool]) {
	extendToCoreGoAtOnce := func(t *testing.T, i int32, pool *OnDemandBlockTaskPool) {
		assert.Equal(t, coreGo, pool.numOfGo())
	}
	extendToMaxGoAtOnce := func(t *testing.T, i int32, pool *OnDemandBlockTaskPool) {
		assert.Equal(t, maxGo, pool.numOfGo())
	}
	opts = append(opts, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime))
	testExtendNumGoFromInitGoToCoreGoAndMaxGo(t, concurrency, queueSize, coreGo, maxGo, extendToCoreGoAtOnce, extendToMaxGoAtOnce, opts...)
}

func testExtendNumGoFromInitGoToCoreGoAndMaxGoStepByStep(t *testing.T, concurrency int32, queueSize int, coreGo, maxGo int32, maxIdleTime time.Duration, opts ...option.Option[OnDemandBlockTaskPool]) {
	extendStepByStep := func(t *testing.T, i int32, pool *OnDemandBlockTaskPool) {
		assert.Equal(t, i, pool.numOfGo())
	}
	opts = append(opts, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime))
	testExtendNumGoFromInitGoToCoreGoAndMaxGo(t, concurrency, queueSize, coreGo, maxGo, extendStepByStep, extendStepByStep, opts...)
}

func testExtendNumGoFromInitGoToCoreGoAndMaxGo(t *testing.T, initGo int32, queueSize int, coreGo, maxGo int32, duringExtendToCoreGo extendStrategyCheckFunc, duringExtendToMaxGo extendStrategyCheckFunc, opts ...option.Option[OnDemandBlockTaskPool]) {

	pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, opts...)
	// waitTime := (maxIdleTime + 1) * 330

	assert.LessOrEqual(t, initGo, coreGo)
	assert.LessOrEqual(t, coreGo, maxGo)

	done := make(chan struct{})
	wait := make(chan struct{}, maxGo)

	// 稳定在concurrency
	for i := int32(0); i < initGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
	}
	// t.Log("AA")
	for i := int32(0); i < initGo; i++ {
		<-wait
	}
	assert.Equal(t, initGo, pool.numOfGo())
	// t.Log("BB")

	// 逐步添加任务
	for i := int32(1); i <= coreGo-initGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		// t.Log("before wait", "m", m, "i", i, "m-n", m-n, "len(wait)", len(wait), "len(queue)", len(pool.queue), "numGO", pool.numOfGo(), "nextnumGO", pool.expectedNumGo)
		<-wait
		// t.Log("after wait coreGo", coreGo, i, pool.numOfGo())

		duringExtendToCoreGo(t, i+initGo, pool)
		// assert.Equal(t, i+n, pool.numOfGo())
	}

	// t.Log("CC")

	assert.Equal(t, coreGo, pool.numOfGo())

	for i := int32(1); i <= maxGo-coreGo; i++ {

		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		// t.Log("before wait", "m", m, "i", i, "m-n", m-n, "len(wait)", len(wait), "len(queue)", len(pool.queue), "numGO", pool.numOfGo(), "nextnumGO", pool.expectedNumGo)
		<-wait
		// t.Log("after wait maxGo", maxGo, i, pool.numOfGo())

		duringExtendToMaxGo(t, i+coreGo, pool)
	}

	// t.Log("DD")

	assert.Equal(t, maxGo, pool.numOfGo())
	close(done)

	// 等待最大空闲时间后，稳定在n
	// <-time.After(waitTime)
	for pool.numOfGo() > initGo {
	}
	assert.Equal(t, initGo, pool.numOfGo())
}

func TestOnDemandBlockTaskPool_In_Closing_State(t *testing.T) {
	t.Parallel()

	t.Run("Shutdown —— 使TaskPool状态由Running变为Closing", func(t *testing.T) {
		t.Parallel()

		concurrency, queueSize := 2, 4
		pool := testNewRunningStateTaskPool(t, concurrency, queueSize)

		// 模拟阻塞提交
		n := concurrency + queueSize*2
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
		for i := 0; i < concurrency; i++ {
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

		assert.Equal(t, int32(concurrency), pool.numOfGo())

		close(taskDone)
		<-done
		assert.Equal(t, stateStopped, pool.internalState())

		// 第一个Shutdown将状态迁移至StateStopped
		// 第三次调用
		done3, err := pool.Shutdown()
		assert.Nil(t, done3)
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
	})

	t.Run("Shutdown —— 协程数仍能按需扩展，调度循环也能自然退出", func(t *testing.T) {
		t.Parallel()

		concurrency, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(1), int32(3), int32(5), 10*time.Millisecond, 0.1
		queueSize := int(maxGo)
		pool := testNewRunningStateTaskPool(t, int(concurrency), queueSize, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime), WithQueueBacklogRate(queueBacklogRate))

		assert.LessOrEqual(t, concurrency, coreGo)
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

		// 等待close(b.queue)信号传递到各个协程
		time.Sleep(1 * time.Second)

		// 调度循环应该正常工作，一直按需开协程直到maxGo
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

		concurrency, queueSize := 2, 4
		pool, wait := testNewRunningStateTaskPoolWithQueueFullFilled(t, concurrency, queueSize)

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

	t.Run("ShutdownNow —— 工作协程数不再扩展，调度循环立即退出", func(t *testing.T) {
		t.Parallel()

		concurrency, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(1), int32(3), int32(5), 10*time.Millisecond, 0.1
		queueSize := int(maxGo)
		pool := testNewRunningStateTaskPool(t, int(concurrency), queueSize, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime), WithQueueBacklogRate(queueBacklogRate))

		assert.LessOrEqual(t, concurrency, coreGo)
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

		// 使调度循环进入default分支
		for i := int32(0); i < coreGo; i++ {
			<-wait
		}

		tasks, err := pool.ShutdownNow()
		assert.NoError(t, err)
		// 见下方双重检查
		assert.GreaterOrEqual(t, len(tasks)+int(pool.numOfGo()), queueSize)

		// 让所有任务结束
		close(taskDone)

		// 用循环取代time.After/time.Sleep
		// 特殊场景需要双重检查
		// 协程1工作中，调度循环处于default分支准备扩展协程（新增一个），此时调用ShutdownNow()
		// 协程1完成工作接收到ShutdownNow()信号退出，而协程2还未开启可以使pool.numOfGo()短暂为0
		// 协程2启动后直接收到ShutdownNow()信号退出
		for pool.numOfGo() != 0 {

		}
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

func testNewRunningStateTaskPool(t *testing.T, concurrency int, queueSize int, opts ...option.Option[OnDemandBlockTaskPool]) *OnDemandBlockTaskPool {
	pool, _ := NewOnDemandBlockTaskPool(concurrency, queueSize, opts...)
	assert.Equal(t, stateCreated, pool.internalState())
	assert.NoError(t, pool.Start())
	assert.Equal(t, stateRunning, pool.internalState())
	return pool
}

func testNewStoppedStateTaskPool(t *testing.T, concurrency int, queueSize int) *OnDemandBlockTaskPool {
	pool := testNewRunningStateTaskPool(t, concurrency, queueSize)
	tasks, err := pool.ShutdownNow()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tasks))
	assert.Equal(t, stateStopped, pool.internalState())
	return pool
}

func testNewRunningStateTaskPoolWithQueueFullFilled(t *testing.T, concurrency int, queueSize int) (*OnDemandBlockTaskPool, chan struct{}) {
	pool := testNewRunningStateTaskPool(t, concurrency, queueSize)
	wait := make(chan struct{})
	for i := 0; i < concurrency+queueSize; i++ {
		func() {
			err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				<-wait
				return nil
			}))
			if err != nil {
				return
			}
		}()
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
