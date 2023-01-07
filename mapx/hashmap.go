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

import "github.com/gotomicro/ekit/syncx"

type node[T Hashable, ValType any] struct {
	key   Hashable
	value ValType
	next  *node[T, ValType]
}

func (m *HashMap[T, ValType]) newNode(key Hashable, val ValType) *node[T, ValType] {
	newNode := m.nodePool.Get()
	newNode.value = val
	newNode.key = key
	return newNode
}

type Hashable interface {
	Code() uint64
	Equals(key any) bool
}

type HashMap[T Hashable, ValType any] struct {
	hashmap  map[uint64]*node[T, ValType]
	nodePool *syncx.Pool[*node[T, ValType]]
}

func (m *HashMap[T, ValType]) Put(key T, val ValType) error {
	hash := key.Code()
	root, ok := m.hashmap[hash]
	if !ok {
		hash = key.Code()
		newNode := m.newNode(key, val)
		m.hashmap[hash] = newNode
		return nil
	}
	pre := root
	for root != nil {
		if root.key.Equals(key) {
			root.value = val
			return nil
		}
		pre = root
		root = root.next
	}
	newNode := m.newNode(key, val)
	pre.next = newNode
	return nil
}

func (m *HashMap[T, ValType]) Get(key T) (ValType, bool) {
	hash := key.Code()
	root, ok := m.hashmap[hash]
	var val ValType
	if !ok {
		return val, false
	}
	for root != nil {
		if root.key.Equals(key) {
			return root.value, true
		}
		root = root.next
	}
	return val, false
}

// Keys 返回 Hashmap 里面的所有的 key。
// 注意：key 的顺序是随机的。
func (m *HashMap[T, ValType]) Keys() []Hashable {
	res := make([]Hashable, 0)
	for _, bucketNode := range m.hashmap {
		res = append(res, bucketNode.key)
	}
	return res
}

// Values 返回 Hashmap 里面的所有的 value。
// 注意：value 的顺序是随机的。
func (m *HashMap[T, ValType]) Values() []ValType {
	res := make([]ValType, 0)
	for _, bucketNode := range m.hashmap {
		res = append(res, bucketNode.value)
	}
	return res
}

func NewHashMap[T Hashable, ValType any](size int) *HashMap[T, ValType] {
	return &HashMap[T, ValType]{
		nodePool: syncx.NewPool[*node[T, ValType]](func() *node[T, ValType] {
			return &node[T, ValType]{}
		}),
		hashmap: make(map[uint64]*node[T, ValType], size),
	}
}

type mapi[T any, ValType any] interface {
	Put(key T, val ValType) error
	Get(key T) (ValType, bool)
}

var _ mapi[Hashable, any] = (*HashMap[Hashable, any])(nil)
