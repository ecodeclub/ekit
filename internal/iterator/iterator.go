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

package iterator

// 一个迭代器的接口，所有的容器类型都可以实现自己的迭代器
// 只需要继承当前接口即可
type Iterator[T any] interface {
	// 迭代器移动到下一个节点，如果成功的话就返回true
	// 如果没有下一个节点，则迭代器所指向的位置会变为非法，一般为nil，并且返回false
	Next() bool

	// 获取迭代器当前所指向的节点的信息
	Get() T

	// 获取error
	Err() error

	// 判断是否有后继节点
	HasNext() bool

	// 判断当前节点是否合法
	Valid() bool

	// 删除当前节点
	Delete()
}
