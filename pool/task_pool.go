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
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	stateCreated int32 = 1
	stateRunning int32 = 2
	stateClosing int32 = 3
	stateStopped int32 = 4
	stateLocked  int32 = 5

	errTaskPoolIsNotRunning = errors.New("ekit: TaskPool未运行")
	errTaskPoolIsClosing    = errors.New("ekit：TaskPool关闭中")
	errTaskPoolIsStopped    = errors.New("ekit: TaskPool已停止")
	errTaskPoolIsStarted    = errors.New("ekit：TaskPool已运行")
	errTaskIsInvalid        = errors.New("ekit: Task非法")
	errTaskRunningPanic     = errors.New("ekit: Task运行时异常")

	errInvalidArgument = errors.New("ekit: 参数非法")

	_            TaskPool = &BlockQueueTaskPool{}
	panicBuffLen          = 2048
)

// TaskPool 任务池
type TaskPool interface {
	// Submit 执行一个任务
	// 如果任务池提供了阻塞的功能，那么如果在 ctx 过期都没有提交成功，那么应该返回错误
	// 调用 Start 之后能否继续提交任务，则取决于具体的实现
	// 调用 Shutdown 或者 ShutdownNow 之后提交任务都会返回错误
	Submit(ctx context.Context, task Task) error

	// Start 开始调度任务执行。在调用 Start 之前，所有的任务都不会被调度执行。
	// Start 之后，能否继续调用 Submit 提交任务，取决于具体的实现
	Start() error

	// Shutdown 关闭任务池。如果此时尚未调用 Start 方法，那么将会立刻返回。
	// 任务池将会停止接收新的任务，但是会继续执行剩下的任务，
	// 在所有任务执行完毕之后，用户可以从返回的 chan 中得到通知
	// 任务池在发出通知之后会关闭 chan struct{}
	Shutdown() (<-chan struct{}, error)

	// ShutdownNow 立刻关闭线程池
	// 任务池能否中断当前正在执行的任务，取决于 TaskPool 的具体实现，以及 Task 的具体实现
	// 该方法会返回所有剩下的任务，剩下的任务是否包含正在执行的任务，也取决于具体的实现
	ShutdownNow() ([]Task, error)
}

// Task 代表一个任务
type Task interface {
	// Run 执行任务
	// 如果 ctx 设置了超时时间，那么实现者需要自己决定是否进行超时控制
	Run(ctx context.Context) error
}

// TaskFunc 一个可执行的任务
type TaskFunc func(ctx context.Context) error

// Run 执行任务
// 超时控制取决于衍生出 TaskFunc 的方法
func (t TaskFunc) Run(ctx context.Context) error { return t(ctx) }

// taskWrapper 是Task的装饰器
type taskWrapper struct {
	t Task
}

func (tw *taskWrapper) Run(ctx context.Context) (err error) {
	defer func() {
		// 处理 panic
		if r := recover(); r != nil {
			buf := make([]byte, panicBuffLen)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("%w：%s", errTaskRunningPanic, fmt.Sprintf("[PANIC]:\t%+v\n%s\n", r, buf))
		}
	}()
	return tw.t.Run(ctx)
}

// BlockQueueTaskPool 并发阻塞的任务池
type BlockQueueTaskPool struct {
	// TaskPool内部状态
	state int32

	queue chan Task
	token chan struct{}
	num   int32

	// 外部信号
	done chan struct{}
	// 内部中断信号
	ctx        context.Context
	cancelFunc context.CancelFunc
	// 缓存
	mux            sync.RWMutex
	submittedTasks []Task
}

// NewBlockQueueTaskPool 创建一个新的 BlockQueueTaskPool
// concurrency 是并发数，即最多允许多少个 goroutine 执行任务
// queueSize 是队列大小，即最多有多少个任务在等待调度
func NewBlockQueueTaskPool(concurrency int, queueSize int) (*BlockQueueTaskPool, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("%w：concurrency应该大于0", errInvalidArgument)
	}
	if queueSize < 0 {
		return nil, fmt.Errorf("%w：queueSize应该大于等于0", errInvalidArgument)
	}
	b := &BlockQueueTaskPool{
		queue: make(chan Task, queueSize),
		token: make(chan struct{}, concurrency),
		done:  make(chan struct{}),
	}
	b.ctx, b.cancelFunc = context.WithCancel(context.Background())
	atomic.StoreInt32(&b.state, stateCreated)
	return b, nil
}

// Submit 提交一个任务
// 如果此时队列已满，那么将会阻塞调用者。
// 如果因为 ctx 的原因返回，那么将会返回 ctx.Err()
// 在调用 Start 前后都可以调用 Submit
func (b *BlockQueueTaskPool) Submit(ctx context.Context, task Task) error {
	if task == nil {
		return fmt.Errorf("%w", errTaskIsInvalid)
	}
	// todo: 用户未设置超时，可以考虑内部给个超时提交
	for {

		if atomic.LoadInt32(&b.state) == stateClosing {
			return fmt.Errorf("%w", errTaskPoolIsClosing)
		}

		if atomic.LoadInt32(&b.state) == stateStopped {
			return fmt.Errorf("%w", errTaskPoolIsStopped)
		}

		task = &taskWrapper{t: task}

		ok, err := b.trySubmit(ctx, task, stateCreated)
		if ok || err != nil {
			return err
		}

		ok, err = b.trySubmit(ctx, task, stateRunning)
		if ok || err != nil {
			return err
		}
	}
}

func (b *BlockQueueTaskPool) trySubmit(ctx context.Context, task Task, state int32) (bool, error) {
	// 进入临界区
	if atomic.CompareAndSwapInt32(&b.state, state, stateLocked) {
		defer atomic.CompareAndSwapInt32(&b.state, stateLocked, state)

		// 此处b.queue <- task不会因为b.queue被关闭而panic
		// 代码执行到trySubmit时TaskPool处于lock状态
		// 要关闭b.queue需要TaskPool处于RUNNING状态，Shutdown/ShutdownNow才能成功
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("%w", ctx.Err())
		case b.queue <- task:
			return true, nil
		default:
			// 不能阻塞在临界区
		}
		return false, nil
	}
	return false, nil
}

// Start 开始调度任务执行
// Start 之后，调用者可以继续使用 Submit 提交任务
func (b *BlockQueueTaskPool) Start() error {

	for {

		if atomic.LoadInt32(&b.state) == stateClosing {
			return fmt.Errorf("%w", errTaskPoolIsClosing)
		}

		if atomic.LoadInt32(&b.state) == stateStopped {
			return fmt.Errorf("%w", errTaskPoolIsStopped)
		}

		if atomic.LoadInt32(&b.state) == stateRunning {
			return fmt.Errorf("%w", errTaskPoolIsStarted)
		}

		if atomic.CompareAndSwapInt32(&b.state, stateCreated, stateRunning) {
			go b.startTasks()
			return nil
		}
	}
}

func (b *BlockQueueTaskPool) startTasks() {
	defer close(b.token)

	for {
		select {
		case <-b.ctx.Done():
			return
		case b.token <- struct{}{}:

			task, ok := <-b.queue
			if !ok {
				return
			}

			go func() {

				atomic.AddInt32(&b.num, 1)
				defer func() {
					atomic.AddInt32(&b.num, -1)
					<-b.token
				}()

				// todo: handle err
				err := task.Run(b.ctx)
				if err != nil {
					return
				}
			}()
		}
	}
}

// Shutdown 将会拒绝提交新的任务，但是会继续执行已提交任务
// 当执行完毕后，会往返回的 chan 中丢入信号
// Shutdown 会负责关闭返回的 chan
// Shutdown 无法中断正在执行的任务
func (b *BlockQueueTaskPool) Shutdown() (<-chan struct{}, error) {

	for {

		if atomic.LoadInt32(&b.state) == stateCreated {
			return nil, fmt.Errorf("%w", errTaskPoolIsNotRunning)
		}

		if atomic.LoadInt32(&b.state) == stateStopped {
			// 重复调用时，恰好前一个Shutdown调用将状态迁移为StateStopped
			// 这种情况与先调用ShutdownNow状态迁移为StateStopped再调用Shutdown效果一样
			return nil, fmt.Errorf("%w", errTaskPoolIsStopped)
		}

		if atomic.LoadInt32(&b.state) == stateClosing {
			return nil, fmt.Errorf("%w", errTaskPoolIsClosing)
		}

		if atomic.CompareAndSwapInt32(&b.state, stateRunning, stateClosing) {
			// 目标：不但希望正在运行中的任务自然退出，还希望队列中等待的任务也能启动执行并自然退出
			// 策略：先将队列中的任务启动并执行（清空队列），再等待全部运行中的任务自然退出

			// 先关闭等待队列不再允许提交
			// 同时任务启动循环能够通过Task==nil来终止循环
			close(b.queue)

			go func() {
				// 等待运行中的Task自然结束
				for atomic.LoadInt32(&b.num) != 0 {
					time.Sleep(time.Second)
				}
				// 通知外部调用者
				close(b.done)
				// 完成最终的状态迁移
				atomic.CompareAndSwapInt32(&b.state, stateClosing, stateStopped)
			}()
			return b.done, nil
		}

	}
}

// ShutdownNow 立刻关闭任务池，并且返回所有剩余未执行的任务（不包含正在执行的任务）
func (b *BlockQueueTaskPool) ShutdownNow() ([]Task, error) {

	for {

		if atomic.LoadInt32(&b.state) == stateCreated {
			return nil, fmt.Errorf("%w", errTaskPoolIsNotRunning)
		}

		if atomic.LoadInt32(&b.state) == stateClosing {
			return nil, fmt.Errorf("%w", errTaskPoolIsClosing)
		}

		if atomic.LoadInt32(&b.state) == stateStopped {
			return nil, fmt.Errorf("%w", errTaskPoolIsStopped)
		}

		if atomic.CompareAndSwapInt32(&b.state, stateRunning, stateStopped) {
			// 目标：立刻关闭并且返回所有剩下未执行的任务
			// 策略：关闭等待队列不再接受新任务，中断任务启动循环，清空等待队列并保存返回

			close(b.queue)

			// 发送中断信号，中断任务启动循环
			b.cancelFunc()

			b.mux.Lock()
			// 清空队列并保存
			var tasks []Task
			for task := range b.queue {
				b.submittedTasks = append(b.submittedTasks, task)
				tasks = append(tasks, task)
			}
			b.mux.Unlock()

			return tasks, nil
		}
	}
}

// internalState 用于查看TaskPool状态
func (b *BlockQueueTaskPool) internalState() int32 {
	for {
		state := atomic.LoadInt32(&b.state)
		if state != stateLocked {
			return state
		}
	}
}

func (b *BlockQueueTaskPool) NumGo() int32 {
	return atomic.LoadInt32(&b.num)
}
