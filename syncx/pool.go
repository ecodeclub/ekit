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

import "sync"

// Pool 是对 sync.Pool 的简单封装
// 会有一些性能损耗，但是基本可以忽略不计。担忧性能问题的可以参考
type Pool[T any] struct {
	p sync.Pool
}

// NewPool 创建一个 Pool 实例
// factory 必须返回 T 类型的值，并且不能返回 nil
func NewPool[T any](factory func() T) *Pool[T] {
	return &Pool[T]{
		p: sync.Pool{
			New: func() any {
				return factory()
			},
		},
	}
}

// Get 取出一个元素
func (p *Pool[T]) Get() T {
	return p.p.Get().(T)
}

// Put 放回去一个元素
func (p *Pool[T]) Put(t T) {
	p.p.Put(t)
}
