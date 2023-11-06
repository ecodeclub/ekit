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

package list

import (
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/internal/list"
)

func NewSkipList[T any](compare ekit.Comparator[T]) *SkipList[T] {
	pq := &SkipList[T]{}
	pq.skiplist = list.NewSkipList[T](compare)
	return pq
}

type SkipList[T any] struct {
	skiplist *list.SkipList[T]
}

func (sl *SkipList[T]) Search(target T) bool {
	return sl.skiplist.Search(target)
}

func (sl *SkipList[T]) AsSlice() []T {
	return sl.skiplist.AsSlice()
}

func (sl *SkipList[T]) Len() int {
	return sl.skiplist.Len()
}

func (sl *SkipList[T]) Cap() int {
	return sl.Len()
}

func (sl *SkipList[T]) Insert(Val T) {
	sl.skiplist.Insert(Val)
}

func (sl *SkipList[T]) DeleteElement(target T) bool {
	return sl.skiplist.DeleteElement(target)
}
