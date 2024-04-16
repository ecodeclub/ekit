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

// LimitPool 是对 Pool 的简单封装允许用户通过控制一段时间内对Pool的令牌申请次数来间接控制Pool中对象的内存总占用量
type LimitPool[T any] struct {
	pool   *Pool[T]
	tokens *atomic.Int32
}

// NewLimitPool 创建一个 LimitPool 实例
// maxTokens 表示一段时间内的允许发放的最大令牌数
// factory 必须返回 T 类型的值，并且不能返回 nil
func NewLimitPool[T any](maxTokens int, factory func() T) *LimitPool[T] {
	var tokens atomic.Int32
	tokens.Add(int32(maxTokens))
	return &LimitPool[T]{
		pool:   NewPool[T](factory),
		tokens: &tokens,
	}
}

// Get 取出一个元素
// 如果返回值是 true，则代表确实从 Pool 里面取出来了一个
// 否则是新建了一个
func (l *LimitPool[T]) Get() (T, bool) {
	if l.tokens.Add(-1) < 0 {
		l.tokens.Add(1)
		var zero T
		return zero, false
	}
	return l.pool.Get(), true
}

// Put 放回去一个元素
func (l *LimitPool[T]) Put(t T) {
	l.pool.Put(t)
	l.tokens.Add(1)
}
