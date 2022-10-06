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

import "context"

// BlockingQueue 阻塞队列
// 参考 Queue 普通队列
// 一个阻塞队列是否遵循 FIFO 取决于具体实现
type BlockingQueue[T any] interface {
	// Enqueue 将元素放入队列。如果此时队列已经满了，那么调用者会被阻塞，直到 ctx 超时
	Enqueue(ctx context.Context, t T) error
	// Dequeue 从队首获得一个元素
	// 如果此时队列里面没有元素，那么调用者会被阻塞，直到 ctx 超时
	Dequeue(ctx context.Context) (T, error)
}

// Queue 普通队列
// 参考 BlockingQueue 阻塞队列
// 一个队列是否遵循 FIFO 取决于具体实现
type Queue[T any] interface {
	// Enqueue 将元素放入队列，如果此时队列已经满了，那么返回错误
	Enqueue(t T) error
	// Dequeue 从队首获得一个元素
	// 如果此时队列里面没有元素，那么返回错误
	Dequeue() (T, error)
}
