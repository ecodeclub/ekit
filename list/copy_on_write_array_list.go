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

import (
	"sync"

	"github.com/ecodeclub/ekit/internal/errs"
	"github.com/ecodeclub/ekit/internal/slice"
)

var (
	_ List[any] = &CopyOnWriteArrayList[any]{}
)

// CopyOnWriteArrayList 基于切片的简单封装 写时加锁，读不加锁，适合于读多写少场景
type CopyOnWriteArrayList[T any] struct {
	vals  []T
	mutex *sync.Mutex
}

// NewCopyOnWriteArrayList
func NewCopyOnWriteArrayList[T any]() *CopyOnWriteArrayList[T] {
	m := &sync.Mutex{}
	return &CopyOnWriteArrayList[T]{
		vals:  make([]T, 0),
		mutex: m,
	}
}

// NewCopyOnWriteArrayListOf 直接使用 ts，会执行复制
func NewCopyOnWriteArrayListOf[T any](ts []T) *CopyOnWriteArrayList[T] {
	items := make([]T, len(ts))
	copy(items, ts)
	m := &sync.Mutex{}
	return &CopyOnWriteArrayList[T]{
		vals:  items,
		mutex: m,
	}
}

func (a *CopyOnWriteArrayList[T]) Get(index int) (t T, e error) {
	l := a.Len()
	if index < 0 || index >= l {
		return t, errs.NewErrIndexOutOfRange(l, index)
	}
	return a.vals[index], e
}

// Append 往CopyOnWriteArrayList里追加数据
func (a *CopyOnWriteArrayList[T]) Append(ts ...T) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	n := len(a.vals)
	newItems := make([]T, n, n+len(ts))
	copy(newItems, a.vals)
	newItems = append(newItems, ts...)
	a.vals = newItems
	return nil
}

// Add 在CopyOnWriteArrayList下标为index的位置插入一个元素
// 当index等于CopyOnWriteArrayList长度等同于append
func (a *CopyOnWriteArrayList[T]) Add(index int, t T) (err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	n := len(a.vals)
	newItems := make([]T, n, n+1)
	copy(newItems, a.vals)
	newItems, err = slice.Add(newItems, t, index)
	a.vals = newItems
	return
}

// Set 设置CopyOnWriteArrayList里index位置的值为t
func (a *CopyOnWriteArrayList[T]) Set(index int, t T) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	n := len(a.vals)
	if index >= n || index < 0 {
		return errs.NewErrIndexOutOfRange(n, index)
	}
	newItems := make([]T, n)
	copy(newItems, a.vals)
	newItems[index] = t
	a.vals = newItems
	return nil
}

// 这里不涉及缩容，每次都是当前内容长度申请的数组容量
func (a *CopyOnWriteArrayList[T]) Delete(index int) (T, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	var ret T
	n := len(a.vals)
	if index >= n || index < 0 {
		return ret, errs.NewErrIndexOutOfRange(n, index)
	}
	newItems := make([]T, len(a.vals)-1)
	item := 0
	for i, v := range a.vals {
		if i == index {
			ret = v
			continue
		}
		newItems[item] = v
		item++
	}
	a.vals = newItems
	return ret, nil
}

func (a *CopyOnWriteArrayList[T]) Len() int {
	return len(a.vals)
}

func (a *CopyOnWriteArrayList[T]) Cap() int {
	return cap(a.vals)
}

func (a *CopyOnWriteArrayList[T]) Range(fn func(index int, t T) error) error {
	for key, value := range a.vals {
		e := fn(key, value)
		if e != nil {
			return e
		}
	}
	return nil
}

func (a *CopyOnWriteArrayList[T]) AsSlice() []T {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	res := make([]T, len(a.vals))
	copy(res, a.vals)
	return res
}
