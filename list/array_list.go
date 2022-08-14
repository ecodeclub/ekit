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
	panic("implement me")
}

// NewArrayListOf 直接使用 ts，而不会执行复制
func NewArrayListOf[T any](ts []T) *ArrayList[T] {
	return &ArrayList[T]{
		vals: ts,
	}
}

func (a *ArrayList[T]) Get(index int) (t T, e error) {
	l := a.Len()
	if index < 0 || index >= l {
		return t, newErrIndexOutOfRange(l, index)
	}
	return a.vals[index], e
}

func (a *ArrayList[T]) Append(t T) error {
	// TODO implement me
	panic("implement me")
}

// Add 在ArrayList下标为index的位置插入一个元素
// 当index等于ArrayList长度等同于append
func (a *ArrayList[T]) Add(index int, t T) error {
	if index < 0 || index > len(a.vals) {
		return newErrIndexOutOfRange(len(a.vals), index)
	}
	a.vals = append(a.vals, t)
	copy(a.vals[index+1:], a.vals[index:])
	a.vals[index] = t
	return nil
}

func (a *ArrayList[T]) Set(index int, t T) error {
	// TODO implement me
	panic("implement me")
}

func (a *ArrayList[T]) Delete(index int) (T, error) {
	// TODO implement me
	panic("implement me")
}

func (a *ArrayList[T]) Len() int {
	if a == nil {
		return 0
	}
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
	slice := make([]T, len(a.vals))
	copy(slice, a.vals)
	return slice
}
