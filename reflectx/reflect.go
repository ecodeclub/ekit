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

package reflectx

import (
	"reflect"
)

// IsNilValue 对 IsNil 方法进一步封装
// 如果 val 的类型 为 map、chan、slice、interface 和 ptr 以及 func，可以执行 IsNil 方法
// 否则直接返回 false，避免执行 IsNil 方法时发生 panic
// 特别注意的是，如果 val 本身为非法的值时（例如 nil），需要先通过 IsValid 方法进行判断，避免后续操作发生 panic
func IsNilValue(val reflect.Value) bool {
	// 先判断 reflect.Value 本身是否为非法的值，例如 nil。避免后续获取 val.Type() 时 panic。
	if !val.IsValid() {
		return true
	}
	// 根据类型判断是否可以执行 IsNil 方法
	switch val.Type().Kind() {
	case reflect.Map, reflect.Chan, reflect.Slice, reflect.Interface, reflect.Ptr, reflect.Func, reflect.UnsafePointer:
		return val.IsNil()
	}
	return false
}
