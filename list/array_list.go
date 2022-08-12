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

// ArrayList 基于切片的简单封装
type ArrayList[T any] struct {
	vals []T
}

func NewArrayList[T any](cap int) *ArrayList[T] {
	return &ArrayList[T]{
		vals: make([]T, 0, cap),
	}
}

// NewArrayListOf 直接使用 ts，而不会执行复制
func NewArrayListOf[T any](ts []T) *ArrayList[T] {
	return &ArrayList[T]{
		vals: ts,
	}
}

// Get 返回对应下标的元素，
// 在下标超出范围的情况下，返回错误
func (a *ArrayList[T]) Get(index int) (T, error) {
	if a.Len() <= index || index < 0 {
		// 笔记： T 不是一个确定的类型，最好不要返回 nil ，可以给一个默认的 T 的类型的初始值
		//例如： var i int = nil
		var defaultT T
		return defaultT, newErrIndexOutOfRange(a.Len(), index)
	}
	return a.vals[index], nil
}

// Append 在末尾追加元素
func (a *ArrayList[T]) Append(t T) error {
	a.vals = append(a.vals, t)
	return nil
}

// Add 在特定下标处增加一个新元素
// 如果下标超出范围，应该返回错误
func (a *ArrayList[T]) Add(index int, t T) error {
	if a.Len() < index || index < 0 {
		return newErrIndexOutOfRange(a.Len(), index)
	}
	// 如果 index == slice length ; 说明这是末尾添加一个 t
	if len(a.vals) == index {
		return a.Append(t)
	}
	source := append([]T{}, a.vals...)
	a.vals = append(a.vals[:index], t)
	a.vals = append(a.vals, source[index:]...)
	return nil
}

// Set 重置 index 位置的值
// 如果下标超出范围，应该返回错误
func (a *ArrayList[T]) Set(index int, t T) error {
	if a.Len() <= index || index < 0 {
		return newErrIndexOutOfRange(a.Len(), index)
	}
	a.vals[index] = t
	return nil
}

// Delete 删除目标元素的位置，并且返回该位置的值
// 如果 index 超出下标，应该返回错误
func (a *ArrayList[T]) Delete(index int) (T, error) {
	if a.Len() <= index || index < 0 {
		var defaultT T
		return defaultT, newErrIndexOutOfRange(a.Len(), index)
	}
	t := a.vals[index]
	a.vals = append(a.vals[:index], a.vals[index+1:]...)
	return t, nil
}

// Len 返回长度
func (a *ArrayList[T]) Len() int {
	return len(a.vals)
}

// Cap 返回容量
func (a *ArrayList[T]) Cap() int {
	return cap(a.vals)
}

// Range 遍历 List 的所有元素
func (a *ArrayList[T]) Range(fn func(index int, t T) error) error {
	for i, t := range a.vals {
		err := fn(i, t)
		if err != nil {
			return err
		}
	}
	return nil
}

// AsSlice 将 List 转化为一个切片
// 不允许返回nil，在没有元素的情况下，
// 必须返回一个长度和容量都为 0 的切片
// AsSlice 每次调用都必须返回一个全新的切片
func (a *ArrayList[T]) AsSlice() []T {
	return append([]T{}, a.vals...)
}
