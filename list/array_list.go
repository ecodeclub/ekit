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
	return &ArrayList[T]{vals: make([]T, 0, cap)}
}

// NewArrayListOf 直接使用 ts，而不会执行复制
func NewArrayListOf[T any](ts []T) *ArrayList[T] {
	return &ArrayList[T]{vals: ts}
}

func (a *ArrayList[T]) Get(index int) (T, error) {
	if len(a.vals) <= index || index < 0 {
		var zero T
		return zero, newErrIndexOutOfRange(len(a.vals), index)
	}
	return a.vals[index], nil
}

func (a *ArrayList[T]) Append(t T) error {
	a.vals = append(a.vals, t)
	return nil
}

func (a *ArrayList[T]) Add(index int, t T) error {
	// index 等于切片长度的时候 相当于append
	if len(a.vals) < index || index < 0 {
		return newErrIndexOutOfRange(len(a.vals), index)
	}
	origin := append([]T{}, a.vals...)
	a.vals = append(a.vals[:index], t)
	a.vals = append(a.vals, origin[index:]...)
	return nil
}

func (a *ArrayList[T]) Set(index int, t T) error {
	if len(a.vals) <= index || index < 0 {
		return newErrIndexOutOfRange(len(a.vals), index)
	}
	a.vals[index] = t
	return nil
}

func (a *ArrayList[T]) Delete(index int) (T, error) {
	if len(a.vals) <= index || index < 0 {
		var zero T
		return zero, newErrIndexOutOfRange(len(a.vals), index)
	}
	res := a.vals[index]
	a.vals = append(a.vals[:index], a.vals[index+1:]...)
	return res, nil
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
	res := make([]T, 0, a.Len())
	return append(res, a.vals...)
}
