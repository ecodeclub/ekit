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
	chain  []*Box[T]
	cursor int
	mutex  sync.Mutex
}

func NewArenaPool[T any]() *ArenaPool[T] {
	return &ArenaPool[T]{
		cursor: -1,
	}
}

func (a *ArenaPool[T]) newX() (*Box[T], error) {
	mem := arena.NewArena()
	box := &Box[T]{}
	box.Mem = mem
	box.object = arena.New[T](mem)
	return box, nil
}

func (a *ArenaPool[T]) Get() (*Box[T], error) {
	if a.cursor == -1 {
		X, err := a.newX()
		if err != nil {
			return nil, err
		}
		return X, nil
	}
	a.mutex.Lock()
	ret := a.chain[a.cursor]
	a.cursor--
	a.mutex.Unlock()
	return ret, nil
}

func (a *ArenaPool[T]) Put(X *Box[T]) error {
	a.mutex.Lock()
	a.chain = append(a.chain, X)
	a.cursor++
	a.mutex.Unlock()
	return nil
}

type Box[T any] struct {
	Mem    *arena.Arena
	object *T
}

func (b *Box[T]) Object() *T {
	return b.object
}

func (b *Box[T]) free() error {
	b.Mem.Free()
	return nil
}
