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

package syncx

import (
	"sync/atomic"
)

// LimitPool 是对 Pool 的简单封装允许用户通过控制一段时间内对Pool的最大申请次数来间接控制Pool中对象的最大总量
type LimitPool[T any] struct {
	maxAttempts *atomic.Int32
	pool        *Pool[T]
}

// NewLimitPool 创建一个 LimitPool 实例
// maxAttempts 表示一段时间内的最大申请次数
// factory 必须返回 T 类型的值，并且不能返回 nil
func NewLimitPool[T any](maxAttempts int, factory func() T) *LimitPool[T] {
	var m atomic.Int32
	m.Add(int32(maxAttempts))
	return &LimitPool[T]{
		maxAttempts: &m,
		pool:        NewPool[T](factory),
	}
}

// Get 取出一个元素
func (l *LimitPool[T]) Get() T {
	for {
		currentAttempts := l.maxAttempts.Load()
		if currentAttempts <= 0 {
			var zero T
			return zero
		}
		if l.maxAttempts.CompareAndSwap(currentAttempts, currentAttempts-1) {
			return l.pool.Get()
		}
	}
}

// Put 放回去一个元素
func (l *LimitPool[T]) Put(t T) {
	l.pool.Put(t)
	l.maxAttempts.Add(1)
}
