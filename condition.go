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

package ekit

// IfThenElse 根据条件返回对应的泛型结果
// 注意避免结果的空指针问题
func IfThenElse[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// IfThenElseFunc 根据条件执行对应的函数并返回泛型结果
func IfThenElseFunc[T any](condition bool, trueFunc, falseFunc func() (T, error)) (T, error) {
	if condition {
		return trueFunc()
	}
	return falseFunc()
}
