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

// Queue 是队列的顶级接口
// 一个队列是否会阻塞调用取决于具体的实现
// 队列是否遵循 FIFO 也取决于具体实现
type Queue[T any] interface {
	// Put 将一个元素放入队列中
	// 对于阻塞队列来说，如果当前队列已满，那么调用者会被阻塞，直到 ctx 超时
	Put(ctx context.Context, t T) error
	// Poll 从队头移除一个元素，并返回该元素
	// 对于阻塞队列来说，如果队列为空，那么调用者会被阻塞，直到 ctx 超时
	Poll(ctx context.Context) (T, error)
}
