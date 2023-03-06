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

package list

import "sync"

var (
	_ List[any] = &ConcurrentList[any]{}
)

// ConcurrentList 用读写锁封装了对 List 的操作
// 达到线程安全的目标
type ConcurrentList[T any] struct {
	List[T]
	lock sync.RWMutex
}

func (c *ConcurrentList[T]) Get(index int) (T, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.List.Get(index)
}

func (c *ConcurrentList[T]) Append(ts ...T) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.List.Append(ts...)
}

func (c *ConcurrentList[T]) Add(index int, t T) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.List.Add(index, t)
}

func (c *ConcurrentList[T]) Set(index int, t T) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.List.Set(index, t)
}

func (c *ConcurrentList[T]) Delete(index int) (T, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.List.Delete(index)
}

func (c *ConcurrentList[T]) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.List.Len()
}

func (c *ConcurrentList[T]) Cap() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.List.Cap()
}

func (c *ConcurrentList[T]) Range(fn func(index int, t T) error) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.List.Range(fn)
}

func (c *ConcurrentList[T]) AsSlice() []T {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.List.AsSlice()
}
