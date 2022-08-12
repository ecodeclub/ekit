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
	"fmt"
	"math"
)

// ArrayList 基于切片的简单封装
type ArrayList[T any] struct {
	data []T
	size int
}

const (
	_defaultCap = 10
	// _capacityReductionThreshold 数组缩容阈值，当size大于这个值才触发缩容
	_capacityReductionThreshold = 8
)

func NewArrayList[T any](cap int) *ArrayList[T] {
	if cap < 0 || cap > math.MaxInt {
		panic(fmt.Sprintf("invalid cap: %d", cap))
	}
	return &ArrayList[T]{
		data: make([]T, cap, cap),
		size: 0,
	}
}

// NewArrayListOf 直接使用 ts，而不会执行复制
func NewArrayListOf[T any](ts []T) *ArrayList[T] {
	if ts == nil {
		ts = make([]T, _defaultCap)
	}
	return &ArrayList[T]{
		data: ts,
		size: len(ts),
	}
}

// Get 根据搜索获取元素
func (a *ArrayList[T]) Get(index int) (T, error) {
	var t T
	if index < 0 || index >= len(a.data) {
		return t, newErrIndexOutOfRange(a.size, index)
	}
	return a.data[index], nil
}

// Append 在List的尾部添加元素
func (a *ArrayList[T]) Append(t T) error {
	return a.Add(a.size, t)
}

// 手动对内部的slice扩缩容
func (a *ArrayList[T]) resize(newCapacity int) {
	newData := make([]T, newCapacity)
	copy(newData, a.data)
	a.data = newData
}

// Add 在 index 位置插入一个新的元素 t
func (a *ArrayList[T]) Add(index int, t T) error {
	if index < 0 || index > a.Len() {
		return newErrIndexOutOfRange(a.Len(), index)
	}
	if a.size == a.Len() {
		a.resize(2 * a.Len())
	}
	for i := a.size - 1; i >= index; i-- {
		a.data[i+1] = a.data[i]
	}
	a.data[index] = t
	a.size++
	return nil
}

// Set 修改 index 位置的元素的值为 t
func (a *ArrayList[T]) Set(index int, t T) error {
	if index < 0 || index >= a.size {
		return newErrIndexOutOfRange(a.Len(), index)
	}
	a.data[index] = t
	return nil
}

// Delete 删除index位置的元素，并返回删除的元素
func (a *ArrayList[T]) Delete(index int) (T, error) {
	var t T
	if index < 0 || index >= a.size {
		return t, newErrIndexOutOfRange(a.Len(), index)
	}
	ret := a.data[index]
	for i := index + 1; i < a.size; i++ {
		a.data[i-1] = a.data[i]
	}
	a.size--
	a.data[a.size] = t
	// 避免性能抖动
	if a.size == a.Len()/4 && a.Len() > _capacityReductionThreshold && a.Len()/2 != 0 {
		a.resize(a.size / 2)
	}
	return ret, nil
}

func (a *ArrayList[T]) Len() int {
	return len(a.data)
}

func (a *ArrayList[T]) Cap() int {
	return cap(a.data)
}

// Range 遍历数组
func (a *ArrayList[T]) Range(fn func(index int, t T) error) error {
	var err error
	for index, val := range a.data {
		err = fn(index, val)
		if err != nil {
			break
		}
	}
	return err
}

// AsSlice 将List转换成 Slice, 需要将内部的 data 拷贝一份，避免返回的 Slice被外部改动从而影响内部
func (a *ArrayList[T]) AsSlice() []T {
	s := make([]T, a.size)
	copy(s, a.data)
	return s
}
