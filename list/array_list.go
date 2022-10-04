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

package list

import (
	"context"

	"github.com/gotomicro/ekit/internal/errs"
	"github.com/gotomicro/ekit/internal/slice"
)

var (
	_ List[any] = &ArrayList[any]{}
)

// ArrayList 基于切片的简单封装
type ArrayList[T any] struct {
	vals []T
}

// NewArrayList 初始化一个len为0，cap为cap的ArrayList
func NewArrayList[T any](cap int) *ArrayList[T] {
	return &ArrayList[T]{vals: make([]T, 0, cap)}
}

// NewArrayListOf 直接使用 ts，而不会执行复制
func NewArrayListOf[T any](ts []T) *ArrayList[T] {
	return &ArrayList[T]{
		vals: ts,
	}
}

func (a *ArrayList[T]) Put(ctx context.Context, t T) error {
	return a.Append(t)
}

func (a *ArrayList[T]) Poll(ctx context.Context) (T, error) {
	return a.Delete(0)
}

func (a *ArrayList[T]) Get(index int) (t T, e error) {
	l := a.Len()
	if index < 0 || index >= l {
		return t, errs.NewErrIndexOutOfRange(l, index)
	}
	return a.vals[index], e
}

// Append 往ArrayList里追加数据
func (a *ArrayList[T]) Append(ts ...T) error {
	a.vals = append(a.vals, ts...)
	return nil
}

// Add 在ArrayList下标为index的位置插入一个元素
// 当index等于ArrayList长度等同于append
func (a *ArrayList[T]) Add(index int, t T) error {
	if index < 0 || index > len(a.vals) {
		return errs.NewErrIndexOutOfRange(len(a.vals), index)
	}
	a.vals = append(a.vals, t)
	copy(a.vals[index+1:], a.vals[index:])
	a.vals[index] = t
	return nil
}

// Set 设置ArrayList里index位置的值为t
func (a *ArrayList[T]) Set(index int, t T) error {
	length := len(a.vals)
	if index >= length || index < 0 {
		return errs.NewErrIndexOutOfRange(length, index)
	}
	a.vals[index] = t
	return nil
}

// Delete 方法会在必要的时候引起缩容，其缩容规则是：
// - 如果容量 > 2048，并且长度小于容量一半，那么就会缩容为原本的 5/8
// - 如果容量 (64, 2048]，如果长度是容量的 1/4，那么就会缩容为原本的一半
// - 如果此时容量 <= 64，那么我们将不会执行缩容。在容量很小的情况下，浪费的内存很少，所以没必要消耗 CPU去执行缩容
func (a *ArrayList[T]) Delete(index int) (T, error) {
	res, t, err := slice.Delete(a.vals, index)
	if err != nil {
		return t, err
	}
	a.vals = res
	a.shrink()
	return t, nil
}

// shrink 数组缩容
func (a *ArrayList[T]) shrink() {
	var newCap int
	c, l := a.Cap(), a.Len()
	if c <= 64 {
		return
	}
	if c > 2048 && (c/l >= 2) {
		newCap = int(float32(c) * float32(0.625))
	} else if c <= 2048 && (c/l >= 4) {
		newCap = c / 2
	} else {
		// 不满足缩容
		return
	}
	newSlice := make([]T, 0, newCap)
	newSlice = append(newSlice, a.vals...)
	a.vals = newSlice
}

func (a *ArrayList[T]) Len() int {
	return len(a.vals)
}

func (a *ArrayList[T]) Cap() int {
	return cap(a.vals)
}

func (a *ArrayList[T]) Range(fn func(index int, t T) error) error {
	for key, value := range a.vals {
		e := fn(key, value)
		if e != nil {
			return e
		}
	}
	return nil
}

func (a *ArrayList[T]) AsSlice() []T {
	res := make([]T, len(a.vals))
	copy(res, a.vals)
	return res
}
