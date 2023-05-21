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
	"errors"

	"github.com/ecodeclub/ekit"
)

var (
	errMultiMapComparatorIsNull = errors.New("ekit: Comparator 不能为 nil")
)

// MultiMap 多映射的 Map
// 它可以将一个键映射到多个值上
type MultiMap[K any, V any] struct {
	m mapi[K, []V]
}

// NewMultiTreeMap MultiTreeMap 的构造方法，创建一个新的 MultiTreeMap
// 注意：
// - comparator 不能为 nil
func NewMultiTreeMap[K any, V any](comparator ekit.Comparator[K]) (*MultiMap[K, V], error) {
	if comparator == nil {
		return nil, errMultiMapComparatorIsNull
	}

	treeMap, err := NewTreeMap[K, []V](comparator)
	if err != nil {
		return nil, err
	}
	return &MultiMap[K, V]{
		m: treeMap,
	}, nil
}

// NewMultiHashMap MultiHashMap 的构造方法，创建一个新的 MultiHashMap
func NewMultiHashMap[K Hashable, V any](size int) *MultiMap[K, V] {
	var m mapi[K, []V] = NewHashMap[K, []V](size)
	return &MultiMap[K, V]{
		m: m,
	}
}

// Put 在 MultiMap 中添加键值对或向已有键 k 的值追加数据
func (m *MultiMap[K, V]) Put(k K, v V) error {
	var val []V
	var ok bool
	if val, ok = m.Get(k); ok {
		val = append(val, v)
	} else {
		val = []V{v}
	}
	return m.m.Put(k, val)
}

// Get 从 MultiMap 中获取已有键 k 的值
// 如果键 k 不存在，则返回的 bool 值为 false
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
