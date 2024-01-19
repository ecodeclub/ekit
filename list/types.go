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

// List 接口
// 该接口只定义清楚各个方法的行为和表现
type List[T any] interface {
	// Get 返回对应下标的元素，
	// 在下标超出范围的情况下，返回错误
	Get(index int) (T, error)
	// Append 在末尾追加元素
	Append(ts ...T) error
	// Add 在特定下标处增加一个新元素
	// 如果下标不在[0, Len()]范围之内
	// 应该返回错误
	// 如果index == Len()则表示往List末端增加一个值
	Add(index int, t T) error
	// Set 重置 index 位置的值
	// 如果下标超出范围，应该返回错误
	Set(index int, t T) error
	// Delete 删除目标元素的位置，并且返回该位置的值
	// 如果 index 超出下标，应该返回错误
	Delete(index int) (T, error)
	// Len 返回长度
	Len() int
	// Cap 返回容量
	Cap() int
	// Range 遍历 List 的所有元素
	Range(fn func(index int, t T) error) error
	// AsSlice 将 List 转化为一个切片
	// 不允许返回nil，在没有元素的情况下，
	// 必须返回一个长度和容量都为 0 的切片
	// AsSlice 每次调用都必须返回一个全新的切片
	AsSlice() []T
}
