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

package set

type Set[T comparable] interface {
	Add(key T)
	Delete(key T)
	// Exist 返回是否存在这个元素
	Exist(key T) bool
	Keys() []T
}

type MapSet[T comparable] struct {
	m map[T]struct{}
}

func NewMapSet[T comparable](size int) *MapSet[T] {
	return &MapSet[T]{
		m: make(map[T]struct{}, size),
	}
}
func (s *MapSet[T]) Add(val T) {
	s.m[val] = struct{}{}
}

func (s *MapSet[T]) Delete(key T) {
	delete(s.m, key)
}

func (s *MapSet[T]) Exist(key T) bool {
	_, ok := s.m[key]
	return ok
}

// Keys 方法返回的元素顺序不固定
func (s *MapSet[T]) Keys() []T {
	ans := make([]T, 0, len(s.m))
	for key := range s.m {
		ans = append(ans, key)
	}
	return ans
}
