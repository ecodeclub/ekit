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

//go:build goexperiment.arenas

package pool

import (
	"arena"
	"sync"
)

type ArenaPool[T any] struct {
	chain []*Arena[T]
	mutex sync.RWMutex
}

func NewArenaPool[T any]() *ArenaPool[T] {
	return &ArenaPool[T]{}
}

func (a *ArenaPool[T]) Get() (*Arena[T], error) {
	a.mutex.RLock()
	l := len(a.chain)
	a.mutex.RUnlock()
	if l == 0 {
		mem := arena.NewArena()
		obj := arena.New[T](mem)
		return &Arena[T]{arena: mem, obj: obj}, nil
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	l = len(a.chain)
	if l == 0 {
		mem := arena.NewArena()
		obj := arena.New[T](mem)
		return &Arena[T]{arena: mem, obj: obj}, nil
	}
	ret := a.chain[l-1]
	a.chain = a.chain[:l-1]
	return ret, nil
}

func (a *ArenaPool[T]) Put(X *Arena[T]) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.chain = append(a.chain, X)
	return nil
}

// Arena 二次封装
// 将来支持缩容，或者淘汰空闲很久的 Arena
type Arena[T any] struct {
	arena *arena.Arena
	obj   *T
}

// Obj 返回已有的对象
func (b *Arena[T]) Obj() *T {
	return b.obj
}
