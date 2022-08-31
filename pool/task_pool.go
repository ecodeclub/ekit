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
)

// NewBlockQueueTaskPool 创建一个新的 BlockQueueTaskPool
// concurrency 是并发数，即最多允许多少个 goroutine 执行任务
// queueSize 是队列大小，即最多有多少个任务在等待调度
func NewBlockQueueTaskPool(concurrency int, queueSize int) (*BlockQueueTaskPool, error) {
	if queueSize < 0 || concurrency < 0 {
		return nil, errData
	}

	b := &BlockQueueTaskPool{
		waitTask:  make(chan Task, queueSize),
		doingTask: make(chan struct{}, concurrency),
	}

	b.ctx, b.ctxCancel = context.WithCancel(context.Background())
	b.changePoolStatus(statusNew)

	return b, nil
}

// Submit 提交一个任务
// 如果此时队列已满，那么将会阻塞调用者。
// 如果因为 ctx 的原因返回，那么将会返回 ctx.Err()
// 在调用 Start 前后都可以调用 Submit
func (b *BlockQueueTaskPool) Submit(ctx context.Context, task Task) error {
	if task == nil {
		return errTaskEmpty
	}
	if b.status != statusOpen {
		return errPoolStatus
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", ctx.Err())
	case b.waitTask <- task:
		return nil
	}
}

// Start 开始调度任务执行
// Start 之后，调用者可以继续使用 Submit 提交任务
func (b *BlockQueueTaskPool) Start() error {
	if b.status == statusOpen {
		return errPoolOpened
	}
	if b.status != statusNew {
		return errPoolStatus
	}
	b.changePoolStatus(statusOpen)
	go func() {
		defer close(b.doingTask)
		wg := sync.WaitGroup{}
		for {
			select {
			case <-b.ctx.Done():
				return
			case b.doingTask <- struct{}{}:
				t, ok := <-b.waitTask
				if !ok {
					// 可能通道被关闭了或者发生其他问题
					b.changePoolStatus(statusClose)
					wg.Wait()
					return
				}
				wg.Add(1)
				go func() {
					defer func() {
						<-b.doingTask
						wg.Done()
					}()
					err := t.Run(b.ctx)
					if err != nil {
						return
					}
				}()
			}
		}
	}()
	return nil
}

// Shutdown 将会拒绝提交新的任务，但是会继续执行已提交任务
// 当执行完毕后，会往返回的 chan 中丢入信号
// Shutdown 会负责关闭返回的 chan
// Shutdown 无法中断正在执行的任务
func (b *BlockQueueTaskPool) Shutdown() (<-chan struct{}, error) {
	for {
		if b.status != statusOpen && b.status != statusNew {
			return nil, errPoolStatus
		}
		b.changePoolStatus(statusClose)
		tmp := make(chan struct{})
		close(b.waitTask)
		return tmp, nil
	}
}

// ShutdownNow 立刻关闭任务池，并且返回所有剩余未执行的任务（不包含正在执行的任务）
func (b *BlockQueueTaskPool) ShutdownNow() ([]Task, error) {
	for {
		if b.status != statusOpen && b.status != statusNew {
			return nil, errPoolStatus
		}
		b.changePoolStatus(statusClose)

		taskList := b.getAllTasks()
		close(b.waitTask)
		b.ctxCancel()

		return taskList, nil
	}
}

func (b *BlockQueueTaskPool) getAllTasks() []Task {
	var tasks []Task
	l := len(b.waitTask)
	if l <= 0 {
		return tasks
	}
	b.mu.Lock()
	for task := range b.waitTask {
		tasks = append(tasks, task)
	}
	b.mu.Unlock()
	return tasks
}

func (b *BlockQueueTaskPool) changePoolStatus(wantStatus int) {
	b.mu.Lock()
	b.status = wantStatus
	b.mu.Unlock()
}

func (b *BlockQueueTaskPool) showPoolStatus() int {
	b.mu.Lock()
	s := b.status
	b.mu.Unlock()
	return s
}

// TaskFunc 一个可执行的任务
type TaskFunc func(ctx context.Context) error

// Run 执行任务
// 超时控制取决于衍生出 TaskFunc 的方法
func (t TaskFunc) Run(ctx context.Context) error {
	return t(ctx)
}
