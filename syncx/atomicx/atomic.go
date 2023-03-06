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

package atomicx

import "sync/atomic"

// Value 是对 atomic.Value 的泛型封装
// 相比直接使用 atomic.Value，
// - Load 方法大概开销多了 0.5 ns
// - Store 方法多了不到 2 ns
// - Swap 方法多了 14 ns
// - CompareAndSwap 在失败的情况下，会多 2 ns，成功的时候多了 0.3 ns
// 使用 NewValue 或者 NewValueOf 来创建实例
type Value[T any] struct {
	val atomic.Value
}

// NewValue 会创建一个 Value 对象，里面存放着 T 的零值
// 注意，这个零值是带了类型的零值
func NewValue[T any]() *Value[T] {
	var t T
	return NewValueOf[T](t)
}

// NewValueOf 会使用传入的值来创建一个 Value 对象
func NewValueOf[T any](t T) *Value[T] {
	val := atomic.Value{}
	val.Store(t)
	return &Value[T]{
		val: val,
	}
}

func (v *Value[T]) Load() (val T) {
	data := v.val.Load()
	val = data.(T)
	return
}

func (v *Value[T]) Store(val T) {
	v.val.Store(val)
}

func (v *Value[T]) Swap(new T) (old T) {
	data := v.val.Swap(new)
	old = data.(T)
	return
}

func (v *Value[T]) CompareAndSwap(old, new T) (swapped bool) {
	return v.val.CompareAndSwap(old, new)
}
