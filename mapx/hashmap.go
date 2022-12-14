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

type Node struct {
	key   HashKey
	value any
	next  *Node
}

func NewNode(key HashKey, val any) *Node {
	return &Node{
		key:   key,
		value: val,
	}
}

type HashKey interface {
	Code() uint64
	Comparable(key any) bool
}

type MyHashMap[T HashKey] struct {
	hashmap map[uint64]*Node
}

func (m *MyHashMap[T]) Put(key T, val any) error {
	hash := key.Code()
	root, ok := m.hashmap[hash]
	if !ok {
		hash = key.Code()
		new_node := NewNode(key, val)
		m.hashmap[hash] = new_node
		return nil
	}
	var pre *Node
	for root != nil {
		if root.key.Comparable(key) {
			root.value = val
			return nil
		}
		pre = root
		root = root.next
	}
	new_node := NewNode(key, val)
	pre.next = new_node
	return nil
}

func (m *MyHashMap[T]) Get(key T) (any, bool) {
	hash := key.Code()
	root, ok := m.hashmap[hash]
	if !ok {
		return nil, false
	}
	for root != nil {
		if root.key.Comparable(key) {
			return root.value, true
		}
		root = root.next
	}
	return nil, false
}

func NewHashMap[T HashKey](size int) *MyHashMap[T] {
	return &MyHashMap[T]{
		hashmap: make(map[uint64]*Node, size),
	}
}
