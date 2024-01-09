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

package mapx

import (
	"github.com/ecodeclub/ekit"
)

// MultiMap 多映射的 Map
// 它可以将一个键映射到多个值上
type MultiMap[K any, V any] struct {
	m mapi[K, []V]
}

// NewMultiTreeMap 创建一个基于 TreeMap 的 MultiMap
// 注意：
// - comparator 不能为 nil
func NewMultiTreeMap[K any, V any](comparator ekit.Comparator[K]) (*MultiMap[K, V], error) {
	treeMap, err := NewTreeMap[K, []V](comparator)
	if err != nil {
		return nil, err
	}
	return &MultiMap[K, V]{
		m: treeMap,
	}, nil
}

// NewMultiHashMap 创建一个基于 HashMap 的 MultiMap
func NewMultiHashMap[K Hashable, V any](size int) *MultiMap[K, V] {
	var m mapi[K, []V] = NewHashMap[K, []V](size)
	return &MultiMap[K, V]{
		m: m,
	}
}

func NewMultiBuiltinMap[K comparable, V any](size int) *MultiMap[K, V] {
	var m mapi[K, []V] = newBuiltinMap[K, []V](size)
	return &MultiMap[K, V]{
		m: m,
	}
}

// Put 在 MultiMap 中添加键值对或向已有键 k 的值追加数据
func (m *MultiMap[K, V]) Put(k K, v V) error {
	return m.PutMany(k, v)
}

// PutMany 在 MultiMap 中添加键值对或向已有键 k 的值追加多个数据
func (m *MultiMap[K, V]) PutMany(k K, v ...V) error {
	val, _ := m.Get(k)
	val = append(val, v...)
	return m.m.Put(k, val)
}

// Get 从 MultiMap 中获取已有键 k 的值
// 如果键 k 不存在，则返回的 bool 值为 false
// 返回的切片是一个副本，你对该切片的修改不会影响原本的数据。
func (m *MultiMap[K, V]) Get(k K) ([]V, bool) {
	if v, ok := m.m.Get(k); ok {
		return append([]V{}, v...), ok
	}
	return nil, false
}

// Delete 从 MultiMap 中删除指定的键 k
func (m *MultiMap[K, V]) Delete(k K) ([]V, bool) {
	return m.m.Delete(k)
}

// Keys 返回 MultiMap 所有的键
func (m *MultiMap[K, V]) Keys() []K {
	return m.m.Keys()
}

// Values 返回 MultiMap 所有的值
func (m *MultiMap[K, V]) Values() [][]V {
	values := m.m.Values()
	copyValues := make([][]V, 0, len(values))
	for i := range values {
		copyValues = append(copyValues, append([]V{}, values[i]...))
	}
	return copyValues
}

// Len 返回 MultiMap 键值对的数量
func (m *MultiMap[K, V]) Len() int64 {
	return m.m.Len()
}
