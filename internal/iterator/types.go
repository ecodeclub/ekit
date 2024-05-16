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

import "errors"

var (
	ErrStructHasChange = errors.New("ekit: 迭代器内元素已经被修改")
	ErrNoSuchData      = errors.New("ekit:没有找到当前元素")
)

// Iterator 一个迭代器的接口，所有的容器类型都可以实现自己的迭代器,只需要继承当前接口即可
type Iterator[T any] interface {
	// Next 迭代器移动到下一个节点，并且返回当前指向的节点
	// 如果没有下一个节点，返回nil 和err
	Next() (T, error)

	// Get 获取迭代器当前所指向的节点的信息
	Get() (T, error)

	// HasNext 判断是否有后继节点
	HasNext() bool

	//// Valid 判断当前节点是否合法
	//Valid() bool
	// Err 主要是在 HasNext后获取错误
	Err() error

	// Delete 删除当前节点
	Delete() error
}

// IteratorAble Iterator的标记接口
type IteratorAble[T any] interface {
	Iterator() *Iterator[T]
}

// ModCount 修改计数
type ModCount interface {
	GetModCount() int
	Increment()
}

// UnModInRange 由于有的数据结构支持迭代器遍历的时候进行删除或者修改，所以单独设计了一个接口
type UnModInRange[T any] interface {
	CheckMod() error
}
