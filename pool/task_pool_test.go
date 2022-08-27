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
	"fmt"
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
	assert.Equal(t, StateCreated, pool.State())
	errChan := make(chan error)
	go func() {
		// 多次运行结果一直
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				errChan <- pool.Start()
				wg.Done()
			}()
		}
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		assert.NoError(t, err)
	}
	assert.Equal(t, StateRunning, pool.State())
}

func TestTaskPool_Submit(t *testing.T) {
	t.Parallel()
	//todo: Submit
	//   TaskPoolRunning Shutdown/ShutdownNow前
	//     在TaskPool所有状态中都可以提交，有的成功/阻塞，有的立即失败。
	//            [x]ctx时间内提交成功
	//			[x]ctx超时，提交失败给出错误信息
	//    [x] 监听状态变化，从running->closing/stopped
	//    [x] Shutdown后,状态变迁，需要检查出并报错 ErrTaskPoolIsClosing
	//   [x] ShutdownNow后状态变迁，需要检查出并报错,ErrTaskPoolIsStopped

	t.Run("提交Task阻塞", func(t *testing.T) {
		t.Parallel()

		t.Run("TaskPool状态由Created变为Running", func(t *testing.T) {
			t.Parallel()
			// 为了准确模拟，内部自定义一个pool
			pool, _ := NewBlockQueueTaskPool(1, 1)
			assert.Equal(t, StateCreated, pool.State())
			errChan := make(chan error)
			go func() {
				errChan <- pool.Start()
			}()
			testSubmitBlockingAndTimeout(t, pool)
			assert.NoError(t, <-errChan)
			assert.Equal(t, StateRunning, pool.State())
		})

		t.Run("TaskPool状态由Running变为Closing", func(t *testing.T) {
			t.Parallel()
			// 为了准确模拟，内部自定义一个pool
			pool, _ := NewBlockQueueTaskPool(1, 2)
			assert.Equal(t, StateCreated, pool.State())
			err := pool.Start()
			assert.NoError(t, err)
			assert.Equal(t, StateRunning, pool.State())

			// 模拟阻塞提交
			n := 10
			submitErrChan := make(chan error, 1)
			for i := 0; i < n; i++ {
				go func() {
					err := pool.Submit(context.Background(), &SlowTask{task: TaskFunc(func(ctx context.Context) error {
						<-time.After(10 * time.Millisecond)
						return nil
					})})
					if err != nil {
						submitErrChan <- err
					}
				}()
			}

			// 调用Shutdown使TaskPool状态发生迁移
			type Result struct {
				done <-chan struct{}
				err  error
			}
			resultChan := make(chan Result)
			go func() {
				<-time.After(time.Millisecond)
				done, err := pool.Shutdown()
				resultChan <- Result{done: done, err: err}
			}()

			r := <-resultChan

			// 阻塞中的任务报错，证明处于TaskPool处于StateClosing状态
			assert.ErrorIs(t, <-submitErrChan, ErrTaskPoolIsClosing)

			// Shutdown调用成功
			assert.NoError(t, r.err)

			<-r.done
			// 等待状态迁移完成，并最终进入StateStopped状态
			<-time.After(100 * time.Millisecond)
			assert.Equal(t, StateStopped, pool.State())
		})

		t.Run("TaskPool状态由Running变为Stopped", func(t *testing.T) {
			t.Parallel()
			// 为了准确模拟，内部自定义一个pool
			pool, _ := NewBlockQueueTaskPool(1, 2)
			assert.Equal(t, StateCreated, pool.State())
			err := pool.Start()
			assert.NoError(t, err)
			assert.Equal(t, StateRunning, pool.State())

			// 模拟阻塞提交
			n := 5
			submitErrChan := make(chan error, 1)
			for i := 0; i < n; i++ {
				go func() {
					err := pool.Submit(context.Background(), &SlowTask{task: TaskFunc(func(ctx context.Context) error {
						<-time.After(10 * time.Millisecond)
						return nil
					})})
					if err != nil {
						submitErrChan <- err
					}
				}()
			}

			// 并发调用ShutdownNow
			type Result struct {
				tasks []Task
				err   error
			}
			result := make(chan Result, 1)
			go func() {
				<-time.After(time.Millisecond)
				tasks, err := pool.ShutdownNow()
				result <- Result{tasks: tasks, err: err}
			}()

			r := <-result
			assert.NoError(t, r.err)
			assert.NotNil(t, r.tasks)

			assert.ErrorIs(t, <-submitErrChan, ErrTaskPoolIsStopped)
			assert.Equal(t, StateStopped, pool.State())
		})

	})

}

func TestTaskPool_Shutdown(t *testing.T) {
	t.Parallel()

	pool, _ := NewBlockQueueTaskPool(1, 1)
	assert.Equal(t, StateCreated, pool.State())
	err := pool.Start()
	assert.Equal(t, StateRunning, pool.State())
	assert.NoError(t, err)

	// 第一次调用
	done, err := pool.Shutdown()
	assert.NoError(t, err)

	// 第二次调用
	select {
	case <-done:
		break
	default:
		done2, err2 := pool.Shutdown()
		assert.Equal(t, done2, done)
		assert.Equal(t, err2, err)
	}

	<-time.After(5 * time.Millisecond)
	assert.Equal(t, StateStopped, pool.State())

	// 第一个Shutdown将状态迁移至StateStopped
	// 第三次调用
	done, err = pool.Shutdown()
	assert.Nil(t, done)
	assert.ErrorIs(t, err, ErrTaskPoolIsStopped)
}

func TestTestPool_ShutdownNow(t *testing.T) {

	t.Parallel()

	pool, _ := NewBlockQueueTaskPool(1, 1)
	assert.Equal(t, StateCreated, pool.State())
	err := pool.Start()
	assert.Equal(t, StateRunning, pool.State())
	assert.NoError(t, err)

	type result struct {
		tasks []Task
		err   error
	}
	n := 3
	c := make(chan result, n)

	for i := 0; i < n; i++ {
		go func() {

			tasks, er := pool.ShutdownNow()
			c <- result{tasks: tasks, err: er}
		}()
	}

	for i := 0; i < n; i++ {
		r := <-c
		assert.Nil(t, r.tasks)
		assert.NoError(t, r.err)
		assert.Equal(t, StateStopped, pool.State())
	}
}

func TestTaskPool__Created_(t *testing.T) {
	t.Parallel()

	n, q := 1, 1
	pool, err := NewBlockQueueTaskPool(n, q)
	assert.NoError(t, err)
	assert.NotNil(t, pool)
	assert.Equal(t, StateCreated, pool.State())

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		t.Run("提交非法Task", func(t *testing.T) {
			t.Parallel()
			testSubmitInvalidTask(t, pool)
			assert.Equal(t, StateCreated, pool.State())
		})

		t.Run("正常提交Task", func(t *testing.T) {
			t.Parallel()
			testSubmitValidTask(t, pool)
			assert.Equal(t, StateCreated, pool.State())
		})

		t.Run("阻塞提交并导致超时", func(t *testing.T) {
			t.Parallel()
			// 为了准确模拟，内部自定义一个pool
			pool, _ := NewBlockQueueTaskPool(1, 1)
			assert.Equal(t, StateCreated, pool.State())
			testSubmitBlockingAndTimeout(t, pool)
			assert.Equal(t, StateCreated, pool.State())
		})
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Parallel()

		done, err := pool.Shutdown()
		assert.Nil(t, done)
		assert.ErrorIs(t, err, ErrTaskPoolIsNotRunning)
		assert.Equal(t, StateCreated, pool.State())
	})

	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()

		tasks, err := pool.ShutdownNow()
		assert.Nil(t, tasks)
		assert.ErrorIs(t, err, ErrTaskPoolIsNotRunning)
		assert.Equal(t, StateCreated, pool.State())
	})

}

func TestTaskPool__Running_(t *testing.T) {
	t.Parallel()

	pool, _ := NewBlockQueueTaskPool(1, 1)
	assert.Equal(t, StateCreated, pool.State())
	err := pool.Start()
	assert.Equal(t, StateRunning, pool.State())
	assert.NoError(t, err)

	t.Run("Start", func(t *testing.T) {
		t.Parallel()
		err = pool.Start()
		// todo: 调度器只启动一次
		assert.Equal(t, StateRunning, pool.State())
		assert.NoError(t, err)
	})

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		t.Run("提交非法Task", func(t *testing.T) {
			t.Parallel()
			testSubmitInvalidTask(t, pool)
			assert.Equal(t, StateRunning, pool.State())
		})

		t.Run("正常提交Task", func(t *testing.T) {
			t.Parallel()
			testSubmitValidTask(t, pool)
			assert.Equal(t, StateRunning, pool.State())
		})

		t.Run("阻塞提交并导致超时", func(t *testing.T) {
			t.Parallel()
			// 为了准确模拟，内部自定义一个pool
			pool, _ := NewBlockQueueTaskPool(1, 1)
			assert.Equal(t, StateCreated, pool.State())
			err := pool.Start()
			assert.Equal(t, StateRunning, pool.State())
			assert.NoError(t, err)

			testSubmitBlockingAndTimeout(t, pool)

			assert.Equal(t, StateRunning, pool.State())
		})
	})
}

func TestTaskPool__Closing_(t *testing.T) {

	t.Parallel()

	pool, _ := NewBlockQueueTaskPool(1, 0)
	assert.Equal(t, StateCreated, pool.State())
	err := pool.Start()
	assert.Equal(t, StateRunning, pool.State())
	assert.NoError(t, err)

	// 提交任务
	//num := 10
	//for i := 0; i < num; i++ {
	//	go func() {
	//		err := pool.Submit(context.Background(), &SlowTask{task: TaskFunc(func(ctx context.Context) error {
	//			<-time.After(10 * time.Millisecond)
	//			return nil
	//		})})
	//		t.Log(err)
	//	}()
	//}

	t.Run("Start", func(t *testing.T) {
		t.Parallel()

		pool, _ := NewBlockQueueTaskPool(1, 10)
		assert.Equal(t, StateCreated, pool.State())
		err := pool.Start()
		assert.Equal(t, StateRunning, pool.State())
		assert.NoError(t, err)

		// 提交任务
		num := 10
		for i := 0; i < num; i++ {
			go func() {
				err := pool.Submit(context.Background(), &SlowTask{task: TaskFunc(func(ctx context.Context) error {
					<-time.After(10 * time.Millisecond)
					return nil
				})})
				t.Log(err)
			}()
		}

		done, err := pool.Shutdown()
		assert.NoError(t, err)
		assert.ErrorIs(t, pool.Start(), ErrTaskPoolIsClosing)
		<-done
		<-time.After(10 * time.Millisecond)
		assert.Equal(t, StateStopped, pool.State())
	})

	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()

		pool, _ := NewBlockQueueTaskPool(1, 0)
		assert.Equal(t, StateCreated, pool.State())
		err := pool.Start()
		assert.Equal(t, StateRunning, pool.State())
		assert.NoError(t, err)

		// 提交任务
		num := 10
		for i := 0; i < num; i++ {
			go func() {
				err := pool.Submit(context.Background(), &SlowTask{task: TaskFunc(func(ctx context.Context) error {
					<-time.After(10 * time.Millisecond)
					return nil
				})})
				t.Log(err)
			}()
		}

		done, err := pool.Shutdown()
		assert.NoError(t, err)

		tasks, err := pool.ShutdownNow()
		assert.ErrorIs(t, err, ErrTaskPoolIsClosing)
		assert.Nil(t, tasks)

		<-done
		<-time.After(50 * time.Millisecond)
		assert.Equal(t, StateStopped, pool.State())
	})

}

func TestTestPool__Stopped_(t *testing.T) {
	t.Parallel()
	n := 2
	pool, _ := NewBlockQueueTaskPool(1, n)
	assert.Equal(t, StateCreated, pool.State())
	err := pool.Start()
	assert.Equal(t, StateRunning, pool.State())
	assert.NoError(t, err)

	// 模拟阻塞提交
	for i := 0; i < 3*n; i++ {
		go func() {
			err := pool.Submit(context.Background(), &FastTask{task: TaskFunc(func(ctx context.Context) error {
				<-time.After(2 * time.Millisecond)
				return nil
			})})
			t.Log(err)
			err = pool.Submit(context.Background(), &SlowTask{task: TaskFunc(func(ctx context.Context) error {
				<-time.After(10 * time.Millisecond)
				return nil
			})})
			t.Log(err)
		}()
	}

	tasks, err := pool.ShutdownNow()
	//assert.NotNil(t, tasks)
	assert.NoError(t, err)
	assert.Equal(t, StateStopped, pool.State())

	t.Run("Start", func(t *testing.T) {
		t.Parallel()
		assert.ErrorIs(t, pool.Start(), ErrTaskPoolIsStopped)
		assert.Equal(t, StateStopped, pool.State())
	})

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error { return nil }))
		assert.ErrorIs(t, err, ErrTaskPoolIsStopped)
		assert.Equal(t, StateStopped, pool.State())
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Parallel()

		done, err := pool.Shutdown()
		assert.Nil(t, done)
		assert.ErrorIs(t, err, ErrTaskPoolIsStopped)
		assert.Equal(t, StateStopped, pool.State())
	})

	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()
		// 多次调用返回相同结果
		tasks2, err := pool.ShutdownNow()
		assert.NoError(t, err)
		assert.Equal(t, tasks2, tasks)
	})
}

func testSubmitBlockingAndTimeout(t *testing.T, pool *BlockQueueTaskPool) {
	err := pool.Submit(context.Background(), &SlowTask{task: TaskFunc(func(ctx context.Context) error {
		<-time.After(2 * time.Millisecond)
		return nil
	})})
	assert.NoError(t, err)

	n := 2
	errChan := make(chan error, n)
	for i := 0; i < n; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()
			err := pool.Submit(ctx, &SlowTask{task: TaskFunc(func(ctx context.Context) error {
				<-time.After(2 * time.Millisecond)
				return nil
			})})
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
}

func testSubmitInvalidTask(t *testing.T, pool *BlockQueueTaskPool) {
	invalidTasks := map[string]Task{"*SlowTask": (*SlowTask)(nil), "*FastTask": (*FastTask)(nil), "nil": nil, "TaskFunc(nil)": TaskFunc(nil)}

	for name, task := range invalidTasks {
		t.Run(name, func(t *testing.T) {
			err := pool.Submit(context.Background(), task)
			assert.ErrorIs(t, err, ErrTaskIsInvalid)
		})
	}
}

func TestTaskExecutor(t *testing.T) {
	t.Parallel()

	/* todo:
	     [x]快慢任务分离
	     [x]快任务没有goroutine限制，提供方法检查个数
	     [x]慢任务占用固定个数goroutine，提供方法检查个数
	     [x]任务的panic处理
	     []任务的error处理
	todo: [x]Closing-优雅关闭
	        [x]返回一个chan,
	              供调用者监听，调用者从chan拿到消息即表明所有任务结束
	        [x]等待任务自然结束
	             [x]关闭chan——自动终止启动循环
	             [x]将队列中等待的任务启动执行
	             [x]等待未完成任务；
	       [x]Stop-强制关闭
	             [x] 关闭chan
	             [x] 终止任务启动循环
	             [x] 将队列中清空并未完成的任务返回
	*/

	t.Run("优雅关闭", func(t *testing.T) {

		t.Parallel()

		maxGo, numTasks := 2, 5
		n := 4 * numTasks

		ex := NewTaskExecutor(maxGo, n)
		ex.Start()

		// 注意：添加Task后需要调整否者会阻塞
		resultChan := make(chan struct{}, n)

		go func() {
			// chan may be closed
			defer func() {
				if r := recover(); r != nil {
					// 发送失败，也算执行了
					resultChan <- struct{}{}
					t.Log(fmt.Errorf("%w：%#v", ErrTaskPoolIsStopped, r))
				}
			}()
			for i := 0; i < numTasks; i++ {
				ex.SlowQueue() <- &SlowTask{task: TaskFunc(func(ctx context.Context) error {
					resultChan <- struct{}{}
					<-time.After(5 * time.Millisecond)
					return nil
				})}
				// panic slow task
				ex.SlowQueue() <- &SlowTask{task: TaskFunc(func(ctx context.Context) error {
					resultChan <- struct{}{}
					panic("slow task ")
				})}
			}
		}()
		go func() {
			// chan() may be closed
			defer func() {
				if r := recover(); r != nil {
					// 发送失败，也算执行了
					resultChan <- struct{}{}
					t.Log(fmt.Errorf("%w：%#v", ErrTaskPoolIsStopped, r))
				}
			}()
			for i := 0; i < numTasks; i++ {
				ex.FastQueue() <- &FastTask{task: TaskFunc(func(ctx context.Context) error {
					resultChan <- struct{}{}
					<-time.After(2 * time.Millisecond)
					return nil
				})}

				// panic fast task
				ex.FastQueue() <- &FastTask{task: TaskFunc(func(ctx context.Context) error {
					resultChan <- struct{}{}
					panic("fast task")
				})}
			}
		}()

		// 等待任务开始执行
		<-time.After(100 * time.Millisecond)

		<-ex.Close()
		assert.Equal(t, n, 4*numTasks)
		assert.Equal(t, int32(0), ex.NumRunningSlow())
		assert.Equal(t, int32(0), ex.NumRunningFast())
		close(resultChan)
		num := 0
		for r := range resultChan {
			if r == struct{}{} {
				num++
			}
		}
		assert.Equal(t, n, num)
	})

	t.Run("强制关闭", func(t *testing.T) {
		t.Parallel()

		maxGo, numTasks := 2, 5
		ex := NewTaskExecutor(maxGo, numTasks)
		ex.Start()

		// 注意：确保n = len(slowTasks) + len(fastTasks)
		n := 8
		resultChan := make(chan struct{}, n)

		slowTasks := []Task{
			&SlowTask{task: TaskFunc(func(ctx context.Context) error {
				resultChan <- struct{}{}
				<-time.After(5 * time.Millisecond)
				return nil
			})},
			// panic slow task
			&SlowTask{task: TaskFunc(func(ctx context.Context) error {
				resultChan <- struct{}{}
				panic("slow task ")
			})},
			&SlowTask{task: TaskFunc(func(ctx context.Context) error {
				resultChan <- struct{}{}
				<-time.After(5 * time.Millisecond)
				return nil
			})},
			// panic slow task
			&SlowTask{task: TaskFunc(func(ctx context.Context) error {
				resultChan <- struct{}{}
				panic("slow task ")
			})},
		}

		fastTasks := []Task{
			&FastTask{task: TaskFunc(func(ctx context.Context) error {
				resultChan <- struct{}{}
				<-time.After(2 * time.Millisecond)
				return nil
			})},
			&FastTask{task: TaskFunc(func(ctx context.Context) error {
				resultChan <- struct{}{}
				panic("fast task")
			})},
			&FastTask{task: TaskFunc(func(ctx context.Context) error {
				resultChan <- struct{}{}
				<-time.After(2 * time.Millisecond)
				return nil
			})},
			&FastTask{task: TaskFunc(func(ctx context.Context) error {
				resultChan <- struct{}{}
				panic("fast task")
			})},
		}
		go func() {
			// chan may be closed
			defer func() {
				if r := recover(); r != nil {
					// 发送任务时，chan被关闭，也算作执行中
					resultChan <- struct{}{}
				}
			}()
			for _, task := range slowTasks {
				ex.SlowQueue() <- task
			}
		}()
		go func() {
			// chan may be closed
			defer func() {
				if r := recover(); r != nil {
					// 发送任务时，chan被关闭，也算作执行中
					resultChan <- struct{}{}
					//t.Log(fmt.Errorf("%w：%#v", ErrTaskPoolIsStopped, r))
				}
			}()
			for _, task := range fastTasks {
				ex.FastQueue() <- task
			}
		}()

		// 等待任务开始执行
		<-time.After(100 * time.Millisecond)

		tasks := ex.Stop()

		// 等待任务执行并回传信号
		<-time.After(100 * time.Millisecond)

		// 统计执行的任务数
		for ex.NumRunningFast() != 0 || ex.NumRunningSlow() != 0 {
			time.Sleep(time.Millisecond)
		}
		close(resultChan)
		num := 0
		for r := range resultChan {
			if r == struct{}{} {
				num++
			}
		}

		assert.Equal(t, n, len(slowTasks)+len(fastTasks))
		assert.Equal(t, len(slowTasks)+len(fastTasks), num+len(tasks))
	})

}
