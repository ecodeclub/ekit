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
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	StateCreated int32 = 1
	StateRunning int32 = 2
	StateClosing int32 = 3
	StateStopped int32 = 4

	ErrTaskPoolIsNotRunning = errors.New("pool: TaskPool未运行")
	ErrTaskPoolIsClosing    = errors.New("pool：TaskPool关闭中")
	ErrTaskPoolIsStopped    = errors.New("pool: TaskPool已停止")
	ErrTaskIsInvalid        = errors.New("pool: Task非法")

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

// FastTask 快任务，耗时较短，每个任务一个Goroutine
type FastTask struct{ task Task }

func (f *FastTask) Run(ctx context.Context) error { return f.task.Run(ctx) }

// SlowTask 慢任务，耗时较长，运行在固定个数Goroutine上
type SlowTask struct{ task Task }

func (s *SlowTask) Run(ctx context.Context) error { return s.task.Run(ctx) }

// BlockQueueTaskPool 并发阻塞的任务池
type BlockQueueTaskPool struct {
	state atomic.Int32
	numGo int
	lenQu int

	taskExecutor  *TaskExecutor
	fastTaskQueue chan<- Task
	slowTaskQueue chan<- Task

	// 缓存taskExecutor结果
	done   <-chan struct{}
	tasks  []Task
	mux    sync.RWMutex
	locked int32
}

// NewBlockQueueTaskPool 创建一个新的 BlockQueueTaskPool
// concurrency 是并发数，即最多允许多少个 goroutine 执行任务
// queueSize 是队列大小，即最多有多少个任务在等待调度
func NewBlockQueueTaskPool(concurrency int, queueSize int) (*BlockQueueTaskPool, error) {
	taskExecutor := NewTaskExecutor(concurrency, queueSize)
	b := &BlockQueueTaskPool{
		numGo:         concurrency,
		lenQu:         queueSize,
		done:          make(chan struct{}),
		taskExecutor:  taskExecutor,
		fastTaskQueue: taskExecutor.FastQueue(),
		slowTaskQueue: taskExecutor.SlowQueue(),
		locked:        int32(101),
	}
	b.state.Store(StateCreated)
	return b, nil
}

func (b *BlockQueueTaskPool) State() int32 {
	for {
		state := b.state.Load()
		if state != b.locked {
			return state
		}
	}
}

// Submit 提交一个任务
// 如果此时队列已满，那么将会阻塞调用者。
// 如果因为 ctx 的原因返回，那么将会返回 ctx.Err()
// 在调用 Start 前后都可以调用 Submit
func (b *BlockQueueTaskPool) Submit(ctx context.Context, task Task) error {
	if task == nil || reflect.ValueOf(task).IsNil() {
		return fmt.Errorf("%w", ErrTaskIsInvalid)
	}
	// todo: 用户未设置超时，可以考虑内部给个超时提交
	for {

		if b.state.Load() == StateClosing {
			return fmt.Errorf("%w", ErrTaskPoolIsClosing)
		}

		if b.state.Load() == StateStopped {
			return fmt.Errorf("%w", ErrTaskPoolIsStopped)
		}

		if b.state.CompareAndSwap(StateCreated, b.locked) {
			ok, err := b.submitTask(ctx, task, func() chan<- Task { return b.chanByTask(task) })
			if ok || err != nil {
				b.state.Swap(StateCreated)
				return err
			}
			b.state.Swap(StateCreated)
		}

		if b.state.CompareAndSwap(StateRunning, b.locked) {
			ok, err := b.submitTask(ctx, task, func() chan<- Task { return b.chanByTask(task) })
			if ok || err != nil {
				b.state.Swap(StateRunning)
				return err
			}
			b.state.Swap(StateRunning)
		}
	}
}

func (b *BlockQueueTaskPool) chanByTask(task Task) chan<- Task {
	switch task.(type) {
	case *SlowTask:
		return b.slowTaskQueue
	default:
		// FastTask, TaskFunc, 用户自定义类型实现Task接口
		return b.fastTaskQueue
	}
}

func (*BlockQueueTaskPool) submitTask(ctx context.Context, task Task, channel func() chan<- Task) (ok bool, err error) {
	// 此处channel() <- task不会出现panic——因为channel被关闭而panic
	// 代码执行到submit时TaskPool处于lock状态
	// 要关闭channel需要TaskPool处于RUNNING状态，Shutdown/ShutdownNow才能成功
	select {
	case <-ctx.Done():
		return false, fmt.Errorf("%w", ctx.Err())
	case channel() <- task:
		return true, nil
	default:
	}
	return false, nil
}

// Start 开始调度任务执行
// Start 之后，调用者可以继续使用 Submit 提交任务
func (b *BlockQueueTaskPool) Start() error {

	for {

		if b.state.Load() == StateClosing {
			return fmt.Errorf("%w", ErrTaskPoolIsClosing)
		}

		if b.state.Load() == StateStopped {
			return fmt.Errorf("%w", ErrTaskPoolIsStopped)
		}

		if b.state.Load() == StateRunning {
			// 重复调用，返回缓存结果
			return nil
		}

		if b.state.CompareAndSwap(StateCreated, StateRunning) {
			// todo: 启动task调度器，开始执行task
			b.taskExecutor.Start()
			return nil
		}
	}
}

// Shutdown 将会拒绝提交新的任务，但是会继续执行已提交任务
// 当执行完毕后，会往返回的 chan 中丢入信号
// Shutdown 会负责关闭返回的 chan
// Shutdown 无法中断正在执行的任务
func (b *BlockQueueTaskPool) Shutdown() (<-chan struct{}, error) {

	for {

		if b.state.Load() == StateCreated {
			return nil, fmt.Errorf("%w", ErrTaskPoolIsNotRunning)
		}

		if b.state.Load() == StateStopped {
			// 重复调用时，恰好前一个Shutdown调用将状态迁移为StateStopped
			// 这种情况与先调用ShutdownNow状态迁移为StateStopped再调用Shutdown效果一样
			return nil, fmt.Errorf("%w", ErrTaskPoolIsStopped)
		}

		if b.state.Load() == StateClosing {
			// 重复调用，返回缓存结果
			return b.done, nil
		}

		if b.state.CompareAndSwap(StateRunning, StateClosing) {
			// todo: 等待task完成，关闭b.done
			// 监听done信号，然后完成状态迁移StateClosing -> StateStopped
			b.done = b.taskExecutor.Close()
			go func() {
				<-b.done
				b.state.CompareAndSwap(StateClosing, StateStopped)
			}()
			return b.done, nil
		}

	}
}

// ShutdownNow 立刻关闭任务池，并且返回所有剩余未执行的任务（不包含正在执行的任务）
func (b *BlockQueueTaskPool) ShutdownNow() ([]Task, error) {

	for {

		if b.state.Load() == StateCreated {
			return nil, fmt.Errorf("%w", ErrTaskPoolIsNotRunning)
		}

		if b.state.Load() == StateClosing {
			return nil, fmt.Errorf("%w", ErrTaskPoolIsClosing)
		}

		if b.state.Load() == StateStopped {
			// 重复调用，返回缓存结果
			b.mux.RLock()
			tasks := append([]Task(nil), b.tasks...)
			b.mux.RUnlock()
			return tasks, nil
		}
		if b.state.CompareAndSwap(StateRunning, StateStopped) {
			b.mux.Lock()
			b.tasks = b.taskExecutor.Stop()
			tasks := append([]Task(nil), b.tasks...)
			b.mux.Unlock()
			return tasks, nil
		}
	}
}

type TaskExecutor struct {
	slowTasks chan Task
	fastTasks chan Task

	maxGo int32
	//
	done chan struct{}
	// 利用ctx充当内部信号
	ctx        context.Context
	cancelFunc context.CancelFunc
	// 统计
	numSlow atomic.Int32
	numFast atomic.Int32
}

func NewTaskExecutor(maxGo int, queueSize int) *TaskExecutor {
	t := &TaskExecutor{maxGo: int32(maxGo), done: make(chan struct{})}
	t.ctx, t.cancelFunc = context.WithCancel(context.Background())
	t.slowTasks = make(chan Task, queueSize)
	t.fastTasks = make(chan Task, queueSize)
	return t
}

func (t *TaskExecutor) Start() {
	go t.startSlowTasks()
	go t.startFastTasks()
}

func (t *TaskExecutor) startFastTasks() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case task := <-t.fastTasks:
			// handle close(t.fastTasks)
			if task == nil {
				return
			}
			go func() {
				t.numFast.Add(1)
				// log.Println("fast N", t.numFast.Add(1))
				defer func() {
					// 恢复统计
					t.numFast.Add(-1)

					// handle panic
					if r := recover(); r != nil {
						buf := make([]byte, panicBuffLen)
						buf = buf[:runtime.Stack(buf, false)]
						fmt.Printf("[PANIC]:\t%+v\n%s\n", r, buf)
					}
				}()
				// todo: handle err
				fmt.Println(task.Run(t.ctx))
			}()
		}
	}
}

func (t *TaskExecutor) startSlowTasks() {

	for {
		n := atomic.AddInt32(&t.maxGo, -1)
		if n < 0 {
			atomic.AddInt32(&t.maxGo, 1)
			continue
		}
		// log.Println("maxGo=", n)
		select {
		case <-t.ctx.Done():
			return
		case task := <-t.slowTasks:
			// handle close(t.slowTasks)
			if task == nil {
				return
			}
			go func() {
				t.numSlow.Add(1)
				// log.Println("slow N=", t.numSlow.Add(1))
				defer func() {
					// 恢复
					atomic.AddInt32(&t.maxGo, 1)
					t.numSlow.Add(-1)

					// handle panic
					if r := recover(); r != nil {
						buf := make([]byte, panicBuffLen)
						buf = buf[:runtime.Stack(buf, false)]
						fmt.Printf("[PANIC]:\t%+v\n%s\n", r, buf)
					}
				}()
				// todo: handle err
				fmt.Println(task.Run(t.ctx))
			}()
		}
	}
}

func (t *TaskExecutor) FastQueue() chan<- Task {
	return t.fastTasks
}

func (t *TaskExecutor) SlowQueue() chan<- Task {
	return t.slowTasks
}

func (t *TaskExecutor) NumRunningSlow() int32 {
	return t.numSlow.Load()
}

func (t *TaskExecutor) NumRunningFast() int32 {
	return t.numFast.Load()
}

// Close 优雅关闭
// 目标：不但希望正在运行中的任务自然退出，还希望队列中等待的任务也能启动执行并自然退出
// 策略：先将所有队列中的任务启动并执行（清空队列），再等待全部运行中的任务自然退出。
func (t *TaskExecutor) Close() <-chan struct{} {

	// 先关闭等待队列不再允许提交
	// 同时任务启动循环能够通过Task==nil来终止循环
	close(t.slowTasks)
	close(t.fastTasks)

	go func() {

		// 检查三次是因为可能出现：
		// 两队列中有任务且正在创建启动任务尚未执行计数，恰巧此时正在运行中的任务为0
		for i := 0; i < 3; i++ {

			// 确保所有运行中任务也自然退出
			for t.numFast.Load() != 0 || t.numSlow.Load() != 0 {
				time.Sleep(time.Second)
			}
		}

		// 通知外部调用者
		close(t.done)
	}()

	return t.done
}

// Stop 强制关闭
// 目标：立刻关闭并且返回所有剩下未执行的任务
// 策略：关闭等待队列不再接受新任务，中断任务启动循环，清空等待队列并保存返回
func (t *TaskExecutor) Stop() []Task {

	close(t.fastTasks)
	close(t.slowTasks)

	// 发送中断信号，中断任务启动循环
	t.cancelFunc()

	// 清空队列并保存
	var tasks []Task
	for task := range t.fastTasks {
		tasks = append(tasks, task)
	}
	for task := range t.slowTasks {
		tasks = append(tasks, task)
	}
	return tasks
}
