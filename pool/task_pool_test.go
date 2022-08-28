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

func TestTaskPool_New(t *testing.T) {
	t.Parallel()

	pool, err := NewBlockQueueTaskPool(1, -1)
	assert.ErrorIs(t, err, errInvalidArgument)
	assert.Nil(t, pool)

	pool, err = NewBlockQueueTaskPool(1, 0)
	assert.NoError(t, err)
	assert.NotNil(t, pool)

	pool, err = NewBlockQueueTaskPool(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, pool)

	pool, err = NewBlockQueueTaskPool(-1, 1)
	assert.ErrorIs(t, err, errInvalidArgument)
	assert.Nil(t, pool)

	pool, err = NewBlockQueueTaskPool(0, 1)
	assert.ErrorIs(t, err, errInvalidArgument)
	assert.Nil(t, pool)

	pool, err = NewBlockQueueTaskPool(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, pool)

}

func TestTaskPool_Start(t *testing.T) {
	/*
					todo: Start
		              1. happy - [x] change state from CREATED to RUNNING  - done
		                       -  非阻塞 task 调度器开始工作，开始执行工作
		                       - [x] Start多次，保证只运行一次,或者报错——TaskPool已经启动
		              2. sad   -
					         CLOSING state -> start error,多次运行结果一致
			                 STOPPED state -> start error多次运行结果一致
	*/

	t.Parallel()

	pool, _ := NewBlockQueueTaskPool(1, 1)
	assert.Equal(t, stateCreated, pool.internalState())

	n := 5
	errChan := make(chan error, n)

	// 第一次调用
	assert.NoError(t, pool.Start())
	assert.Equal(t, stateRunning, pool.internalState())

	// 多次调用
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			errChan <- pool.Start()
			wg.Done()
		}()
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		assert.NoError(t, err)
	}
	assert.Equal(t, stateRunning, pool.internalState())
}

func TestTaskPool_Submit(t *testing.T) {
	t.Parallel()
	// todo: Submit
	//     在TaskPool所有状态中都可以提交，有的成功/阻塞，有的立即失败。
	//          [x]ctx时间内提交成功
	//			[x]ctx超时，提交失败给出错误信息
	//    [x] 监听状态变化，从running->closing/stopped
	//    [x] Shutdown后,状态变迁，需要检查出并报错 errTaskPoolIsClosing
	//   [x] ShutdownNow后状态变迁，需要检查出并报错,errTaskPoolIsStopped

	t.Run("提交Task阻塞", func(t *testing.T) {
		t.Parallel()

		t.Run("TaskPool状态由Created变为Running", func(t *testing.T) {
			t.Parallel()

			pool, _ := NewBlockQueueTaskPool(1, 1)

			// 与下方 testSubmitBlockingAndTimeout 并发执行
			errChan := make(chan error)
			go func() {
				<-time.After(1 * time.Millisecond)
				errChan <- pool.Start()
			}()

			assert.Equal(t, stateCreated, pool.internalState())

			testSubmitBlockingAndTimeout(t, pool)

			assert.NoError(t, <-errChan)
			assert.Equal(t, stateRunning, pool.internalState())
		})

		t.Run("TaskPool状态由Running变为Closing", func(t *testing.T) {
			t.Parallel()

			pool := testNewRunningStateTaskPool(t, 1, 2)

			// 模拟阻塞提交
			n := 10
			firstSubmitErrChan := make(chan error, 1)
			for i := 0; i < n; i++ {
				go func() {
					err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
						<-time.After(10 * time.Millisecond)
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
				<-time.After(time.Millisecond)
				done, err := pool.Shutdown()
				resultChan <- ShutdownResult{done: done, err: err}
			}()

			r := <-resultChan

			// 阻塞中的任务报错，证明处于TaskPool处于StateClosing状态
			assert.ErrorIs(t, <-firstSubmitErrChan, errTaskPoolIsClosing)

			// Shutdown调用成功
			assert.NoError(t, r.err)
			<-r.done
			// 等待状态迁移完成，并最终进入StateStopped状态
			<-time.After(100 * time.Millisecond)
			assert.Equal(t, stateStopped, pool.internalState())
		})

		t.Run("TaskPool状态由Running变为Stopped", func(t *testing.T) {
			t.Parallel()

			pool := testNewRunningStateTaskPool(t, 1, 2)

			// 模拟阻塞提交
			n := 5
			firstSubmitErrChan := make(chan error, 1)
			for i := 0; i < n; i++ {
				go func() {
					err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
						<-time.After(10 * time.Millisecond)
						return nil
					}))
					if err != nil {
						firstSubmitErrChan <- err
					}
				}()
			}

			// 并发调用ShutdownNow

			result := make(chan ShutdownNowResult, 1)
			go func() {
				<-time.After(time.Millisecond)
				tasks, err := pool.ShutdownNow()
				result <- ShutdownNowResult{tasks: tasks, err: err}
			}()

			r := <-result
			assert.NoError(t, r.err)
			assert.NotEmpty(t, r.tasks)

			assert.ErrorIs(t, <-firstSubmitErrChan, errTaskPoolIsStopped)
			assert.Equal(t, stateStopped, pool.internalState())
		})

	})
}

func TestTaskPool_Shutdown(t *testing.T) {
	t.Parallel()

	pool := testNewRunningStateTaskPool(t, 1, 1)

	// 第一次调用
	done, err := pool.Shutdown()
	assert.NoError(t, err)

	select {
	case <-done:
		break
	default:
		// 第二次调用
		done2, err2 := pool.Shutdown()
		assert.Equal(t, done2, done)
		assert.Equal(t, err2, err)
		assert.Equal(t, stateClosing, pool.internalState())
	}

	<-done
	assert.Equal(t, stateStopped, pool.internalState())

	// 第一个Shutdown将状态迁移至StateStopped
	// 第三次调用
	done, err = pool.Shutdown()
	assert.Nil(t, done)
	assert.ErrorIs(t, err, errTaskPoolIsStopped)
}

func TestTestPool_ShutdownNow(t *testing.T) {

	t.Parallel()

	pool := testNewRunningStateTaskPool(t, 1, 1)

	n := 3
	c := make(chan ShutdownNowResult, n)

	for i := 0; i < n; i++ {
		go func() {
			tasks, er := pool.ShutdownNow()
			c <- ShutdownNowResult{tasks: tasks, err: er}
		}()
	}

	for i := 0; i < n; i++ {
		r := <-c
		assert.Nil(t, r.tasks)
		assert.NoError(t, r.err)
		assert.Equal(t, stateStopped, pool.internalState())
	}
}

func TestTaskPool__Created_(t *testing.T) {
	t.Parallel()

	pool, err := NewBlockQueueTaskPool(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, pool)
	assert.Equal(t, stateCreated, pool.internalState())

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		t.Run("提交非法Task", func(t *testing.T) {
			t.Parallel()

			pool, _ := NewBlockQueueTaskPool(1, 1)
			assert.Equal(t, stateCreated, pool.internalState())
			testSubmitInvalidTask(t, pool)
			assert.Equal(t, stateCreated, pool.internalState())
		})

		t.Run("正常提交Task", func(t *testing.T) {
			t.Parallel()

			pool, _ := NewBlockQueueTaskPool(1, 3)
			assert.Equal(t, stateCreated, pool.internalState())
			testSubmitValidTask(t, pool)
			assert.Equal(t, stateCreated, pool.internalState())
		})

		t.Run("阻塞提交并导致超时", func(t *testing.T) {
			t.Parallel()

			pool, _ := NewBlockQueueTaskPool(1, 1)
			assert.Equal(t, stateCreated, pool.internalState())
			testSubmitBlockingAndTimeout(t, pool)
			assert.Equal(t, stateCreated, pool.internalState())
		})
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Parallel()

		done, err := pool.Shutdown()
		assert.Nil(t, done)
		assert.ErrorIs(t, err, errTaskPoolIsNotRunning)
		assert.Equal(t, stateCreated, pool.internalState())
	})

	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()

		tasks, err := pool.ShutdownNow()
		assert.Nil(t, tasks)
		assert.ErrorIs(t, err, errTaskPoolIsNotRunning)
		assert.Equal(t, stateCreated, pool.internalState())
	})

}

func TestTaskPool__Running_(t *testing.T) {
	t.Parallel()

	t.Run("Start", func(t *testing.T) {
		t.Parallel()

		pool := testNewRunningStateTaskPool(t, 1, 3)
		assert.NoError(t, pool.Start())
		assert.Equal(t, stateRunning, pool.internalState())
	})

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		t.Run("提交非法Task", func(t *testing.T) {
			t.Parallel()
			pool := testNewRunningStateTaskPool(t, 1, 1)
			testSubmitInvalidTask(t, pool)
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
}

func TestTaskPool__Closing_(t *testing.T) {
	t.Parallel()

	t.Run("Start", func(t *testing.T) {
		t.Parallel()

		num := 10
		pool := testNewRunningStateTaskPool(t, 2, num)

		// 提交任务
		for i := 0; i < num; i++ {
			go func() {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					<-time.After(10 * time.Millisecond)
					return nil
				}))
				t.Log(err)
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
		<-time.After(10 * time.Millisecond)
		assert.Equal(t, stateStopped, pool.internalState())
	})

	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()

		pool := testNewRunningStateTaskPool(t, 1, 0)

		// 提交任务
		num := 10
		for i := 0; i < num; i++ {
			go func() {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					<-time.After(10 * time.Millisecond)
					return nil
				}))
				t.Log(err)
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
		<-time.After(50 * time.Millisecond)
		assert.Equal(t, stateStopped, pool.internalState())
	})

}

func TestTestPool__Stopped_(t *testing.T) {
	t.Parallel()

	concurrency, n := 2, 20
	pool := testNewRunningStateTaskPool(t, concurrency, n)

	// 模拟阻塞提交
	for i := 0; i < n; i++ {
		go func() {
			t.Log(pool.NumGo())
			err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				<-time.After(time.Second)
				return nil
			}))
			t.Log(err)
		}()
	}

	<-time.After(100 * time.Millisecond)
	assert.Equal(t, int32(concurrency), pool.NumGo())

	tasks, err := pool.ShutdownNow()
	assert.NoError(t, err)
	assert.NotEmpty(t, tasks)
	assert.Equal(t, stateStopped, pool.internalState())

	t.Run("Start", func(t *testing.T) {
		t.Parallel()

		assert.ErrorIs(t, pool.Start(), errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error { return nil }))
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Parallel()

		done, err := pool.Shutdown()
		assert.Nil(t, done)
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})

	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()

		// 多次调用返回相同结果
		tasks2, err := pool.ShutdownNow()
		assert.NoError(t, err)
		assert.NotEmpty(t, tasks2)
		assert.Equal(t, tasks2, tasks)
		assert.Equal(t, stateStopped, pool.internalState())
	})
}

func testSubmitBlockingAndTimeout(t *testing.T, pool *BlockQueueTaskPool) {

	err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		<-time.After(2 * time.Millisecond)
		return nil
	}))
	assert.NoError(t, err)

	n := 2
	errChan := make(chan error, n)
	for i := 0; i < n; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()
			err := pool.Submit(ctx, TaskFunc(func(ctx context.Context) error {
				<-time.After(2 * time.Millisecond)
				return nil
			}))
			if err != nil {
				errChan <- err
			}
		}()
	}

	assert.ErrorIs(t, <-errChan, context.DeadlineExceeded)
}

func testSubmitValidTask(t *testing.T, pool *BlockQueueTaskPool) {

	err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error { return nil }))
	assert.NoError(t, err)

	err = pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error { panic("task panic") }))
	assert.NoError(t, err)

	err = pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error { return errors.New("fake error") }))
	assert.NoError(t, err)
}

func testSubmitInvalidTask(t *testing.T, pool *BlockQueueTaskPool) {

	invalidTasks := map[string]Task{"*FakeTask": (*FakeTask)(nil), "nil": nil, "TaskFunc(nil)": TaskFunc(nil)}

	for name, task := range invalidTasks {
		t.Run(name, func(t *testing.T) {
			err := pool.Submit(context.Background(), task)
			assert.ErrorIs(t, err, errTaskIsInvalid)
		})
	}
}

type ShutdownNowResult struct {
	tasks []Task
	err   error
}

func testNewRunningStateTaskPool(t *testing.T, concurrency int, queueSize int) *BlockQueueTaskPool {
	pool, _ := NewBlockQueueTaskPool(concurrency, queueSize)
	assert.Equal(t, stateCreated, pool.internalState())
	assert.NoError(t, pool.Start())
	assert.Equal(t, stateRunning, pool.internalState())
	return pool
}

type FakeTask struct{}

func (f *FakeTask) Run(ctx context.Context) error { return nil }
