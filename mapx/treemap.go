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

package mapx

import (
	"errors"

	"github.com/gotomicro/ekit"
	"github.com/gotomicro/ekit/internal/tree"
)

const (
	Red   = false
	Black = true
)

var (
	errTreeMapComparatorIsNull = errors.New("ekit: Comparator不能为nil")
)

// TreeMap 是基于红黑树实现的Map
// 需要注意TreeMap是有序的所以必须传入比较器
// compare	比较器
// root	根节点
type TreeMap[T any, V any] struct {
	*tree.RBTree[T, V]
}

// BuildTreeMap TreeMap构造方法
// 支持传入compare比较器，并根据传入的m构建TreeMap
// 需注意比较器compare不能为nil
func BuildTreeMap[T comparable, V any](compare ekit.Comparator[T], m map[T]V) (*TreeMap[T, V], error) {
	treeMap, err := NewTreeMap[T, V](compare)
	if err != nil {
		return treeMap, err
	}
	err = PutAll(treeMap, m)
	if err != nil {
		return nil, err
	}
	return treeMap, nil
}

// NewTreeMap TreeMap构造方法,创建一个的TreeMap
func NewTreeMap[T any, V any](compare ekit.Comparator[T]) (*TreeMap[T, V], error) {
	if compare == nil {
		return nil, errTreeMapComparatorIsNull
	}
	return &TreeMap[T, V]{
		RBTree: tree.NewRBTree[T, V](compare),
	}, nil
}

// PutAll 将map传入TreeMap
// 需注意如果map中的Key已存在TreeMap将被替换
func PutAll[T comparable, V any](treeMap *TreeMap[T, V], m map[T]V) error {
	if len(m) != 0 {
		keys, values := KeysValues[T, V](m)
		for i := 0; i < len(keys); i++ {
			err := treeMap.Put(keys[i], values[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Put 在TreeMap插入指定Key的节点
// 需注意如果TreeMap已存在该Key那么原值会被替换
// 错误：
// TreeMap中比较器为nil将会返回error
func (treeMap *TreeMap[T, V]) Put(k T, v V) error {
	oldNode := treeMap.Find(k)
	if oldNode == nil {
		node := tree.NewRBNode[T, V](k, v)
		return treeMap.Add(node)
	}
	oldNode.SetValue(v)
	return nil
}

// Get 在TreeMap找到指定Key的节点,返回Val
// 错误：
// TreeMap未找到指定Key将会返回error
// TreeMap中比较器为nil将会返回error
func (treeMap *TreeMap[T, V]) Get(k T) (V, bool) {
	var defaultV V
	node := treeMap.Find(k)
	for node != nil {
		return node.Value, true
	}
	return defaultV, false
}

// Remove TreeMap中删除对应K的节点
func (treeMap *TreeMap[T, V]) Remove(k T) {
	treeMap.Delete(k)
}

var _ mapi[any, any] = (*TreeMap[any, any])(nil)
