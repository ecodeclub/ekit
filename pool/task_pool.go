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

import "context"

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
func (t TaskFunc) Run(ctx context.Context) error {
	return t(ctx)
}

// BlockQueueTaskPool 并发阻塞的任务池
type BlockQueueTaskPool struct {
}

// NewBlockQueueTaskPool 创建一个新的 BlockQueueTaskPool
// concurrency 是并发数，即最多允许多少个 goroutine 执行任务
// queueSize 是队列大小，即最多有多少个任务在等待调度
func NewBlockQueueTaskPool(concurrency int, queueSize int) (*BlockQueueTaskPool, error) {
	return &BlockQueueTaskPool{}, nil
}

// Submit 提交一个任务
// 如果此时队列已满，那么将会阻塞调用者。
// 如果因为 ctx 的原因返回，那么将会返回 ctx.Err()
// 在调用 Start 前后都可以调用 Submit
func (b *BlockQueueTaskPool) Submit(ctx context.Context, task func()) error {
	// TODO implement me
	panic("implement me")
}

// Start 开始调度任务执行
// Start 之后，调用者可以继续使用 Submit 提交任务
func (b *BlockQueueTaskPool) Start() error {
	// TODO implement me
	panic("implement me")
}

// Shutdown 将会拒绝提交新的任务，但是会继续执行已提交任务
// 当执行完毕后，会往返回的 chan 中丢入信号
// Shutdown 会负责关闭返回的 chan
// Shutdown 无法中断正在执行的任务
func (b *BlockQueueTaskPool) Shutdown() (<-chan struct{}, error) {
	// TODO implement me
	panic("implement me")
}

// ShutdownNow 立刻关闭任务池，并且返回所有剩余未执行的任务（不包含正在执行的任务）
func (b *BlockQueueTaskPool) ShutdownNow() ([]Task, error) {
	// TODO implement me
	panic("implement me")
}
