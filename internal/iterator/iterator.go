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
	// 迭代器移动到下一个节点
	// 如果没有下一个节点，则迭代器所指向的位置变为非法位置，一般情况下为nil
	Next()

	// 获取迭代器当前所指向的节点的信息
	Get() (T, error)

	// 判断是否有后继节点
	HasNext() bool

	// 判断当前节点是否合法
	Valid() bool
}
