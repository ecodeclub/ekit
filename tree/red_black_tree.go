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

package tree

import (
	"errors"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/internal/tree"
)

var (
	errRBTreeComparatorIsNull = errors.New("ekit: RBTree 的 Comparator 不能为 nil")
)

// RBTree 简单的封装一下红黑树
type RBTree[K any, V any] struct {
	rbTree *tree.RBTree[K, V] //红黑树本体
}

func NewRBTree[K any, V any](compare ekit.Comparator[K]) (*RBTree[K, V], error) {
	if nil == compare {
		return nil, errRBTreeComparatorIsNull
	}

	return &RBTree[K, V]{
		rbTree: tree.NewRBTree[K, V](compare),
	}, nil
}

// Add 增加节点
func (rb *RBTree[K, V]) Add(key K, value V) error {
	return rb.rbTree.Add(key, value)
}

// Delete 删除节点
func (rb *RBTree[K, V]) Delete(key K) (V, bool) {
	return rb.rbTree.Delete(key)
}

// Set 修改节点
func (rb *RBTree[K, V]) Set(key K, value V) error {
	return rb.rbTree.Set(key, value)
}

// Find 查找节点
func (rb *RBTree[K, V]) Find(key K) (V, error) {
	return rb.rbTree.Find(key)
}

// Size 返回红黑树结点个数
func (rb *RBTree[K, V]) Size() int {
	return rb.rbTree.Size()
}

// KeyValues 获取红黑树所有节点K,V
func (rb *RBTree[K, V]) KeyValues() ([]K, []V) {
	return rb.rbTree.KeyValues()
}
