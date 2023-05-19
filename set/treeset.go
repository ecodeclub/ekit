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

import (
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/mapx"
)

type TreeSet[T any] struct {
	treeMap *mapx.TreeMap[T, any]
}

func NewTreeSet[T any](compare ekit.Comparator[T]) (*TreeSet[T], error) {
	treeMap, err := mapx.NewTreeMap[T, any](compare)
	if err != nil {
		return nil, err
	}
	return &TreeSet[T]{
		treeMap: treeMap,
	}, nil
}

func (s *TreeSet[T]) Add(key T) {
	_ = s.treeMap.Put(key, nil)
}

func (s *TreeSet[T]) Delete(key T) {
	s.treeMap.Delete(key)
}

func (s *TreeSet[T]) Exist(key T) bool {
	_, isExist := s.treeMap.Get(key)
	return isExist
}

// Keys 方法返回的元素顺序不固定
func (s *TreeSet[T]) Keys() []T {
	return s.treeMap.Keys()
}

var _ Set[int] = (*TreeSet[int])(nil)
