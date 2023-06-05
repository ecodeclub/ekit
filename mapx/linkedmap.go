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

import "github.com/ecodeclub/ekit"

type LinkedMap[K any, V any] struct {
	m          mapi[K, *linkedKV[K, V]]
	head, tail *linkedKV[K, V]
}

type linkedKV[K any, V any] struct {
	key        K
	value      V
	prev, next *linkedKV[K, V]
}

func NewLinkedHashMap[K Hashable, V any](size int) *LinkedMap[K, V] {
	hashmap := NewHashMap[K, *linkedKV[K, V]](size)
	return &LinkedMap[K, V]{
		m: hashmap,
	}
}

func NewLinkedTreeMap[K any, V any](comparator ekit.Comparator[K]) (*LinkedMap[K, V], error) {
	treeMap, err := NewTreeMap[K, *linkedKV[K, V]](comparator)
	if err != nil {
		return nil, err
	}
	return &LinkedMap[K, V]{
		m: treeMap,
	}, nil
}

func (l *LinkedMap[K, V]) Put(key K, val V) error {
	lk := &linkedKV[K, V]{
		key:   key,
		value: val,
	}
	if err := l.m.Put(key, lk); err != nil {
		return err
	} else {
		if l.tail != nil {
			curTail := l.tail
			curTail.next = lk
			l.tail = lk
			l.tail.prev = curTail
		} else {
			l.head, l.tail = lk, lk
		}
	}
	return nil
}

func (l *LinkedMap[K, V]) Get(key K) (V, bool) {
	var v V
	if lk, ok := l.m.Get(key); ok {
		return lk.value, ok
	}
	return v, false
}

func (l *LinkedMap[K, V]) Delete(key K) (V, bool) {
	var v V
	if lk, ok := l.m.Delete(key); ok {
		prev := lk.prev
		next := lk.next
		if prev == nil {
			if next == nil {
				l.head, l.tail = nil, nil
			} else {
				l.head = next
				next.prev = nil
			}
		} else if next == nil {
			prev.next = nil
			l.tail = prev
		} else {
			prev.next = next
		}
		return lk.value, ok
	}
	return v, false
}

func (l *LinkedMap[K, V]) Keys() []K {
	keys := make([]K, 0)
	cur := l.head
	for cur != nil {
		keys = append(keys, cur.key)
		cur = cur.next
	}
	return keys
}

func (l *LinkedMap[K, V]) Values() []V {
	values := make([]V, 0)
	cur := l.head
	for cur != nil {
		values = append(values, cur.value)
		cur = cur.next
	}
	return values
}
