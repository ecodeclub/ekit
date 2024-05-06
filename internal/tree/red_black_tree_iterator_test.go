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
	"math/rand"
	"testing"

	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

func TestIteratorToVisitFullRBTree(t *testing.T) {
	t.Parallel()
	n := 10000
	arr := generateArray(n)
	rbTree := NewRBTree[int, int](ekit.ComparatorRealNumber[int])
	for _, v := range arr {
		assert.Nil(t, rbTree.Add(v, v))
	}

	arrVisit := make([]int, n)
	id := 0
	for iter := rbTree.Begin(); iter.Valid(); iter.Next() {
		pa, err := iter.Get(), iter.Err()
		assert.Nil(t, err)
		arrVisit[id] = pa.Key
		assert.Equal(t, id, pa.Key)
		id++
	}
	assert.Equal(t, n, id)
}

func TestIteratorFind(t *testing.T) {
	t.Run("查找存在的节点", func(t *testing.T) {
		t.Parallel()
		rbt := NewRBTree[int, int](ekit.ComparatorRealNumber[int])
		assert.Nil(t, rbt.Add(1, 101))
		assert.Nil(t, rbt.Add(-100, 102))
		assert.Nil(t, rbt.Add(100, 103))
		it, err := rbt.FindIt(-100)
		assert.Nil(t, err)
		assert.Equal(t, 102, it.Get().Value)
	})

	t.Run("查找不存在的节点", func(t *testing.T) {
		t.Parallel()
		rbt := NewRBTree[int, int](ekit.ComparatorRealNumber[int])
		assert.Nil(t, rbt.Add(1, 101))
		assert.Nil(t, rbt.Add(-100, 102))
		assert.Nil(t, rbt.Add(100, 103))
		it, err := rbt.FindIt(2)
		assert.Equal(t, ErrRBTreeNotRBNode, err)
		assert.Nil(t, it)
	})

	t.Run("查找存在的节点，删除后不存在", func(t *testing.T) {
		t.Parallel()
		rbt := NewRBTree[int, int](ekit.ComparatorRealNumber[int])
		assert.Nil(t, rbt.Add(1, 101))
		assert.Nil(t, rbt.Add(-100, 102))
		assert.Nil(t, rbt.Add(100, 103))
		it, err := rbt.FindIt(-100)
		assert.Nil(t, err)
		assert.Equal(t, 102, it.Get().Value)
		assert.Nil(t, rbt.DeleteIt(it))
		it, err = rbt.FindIt(-100)
		assert.Equal(t, ErrRBTreeNotRBNode, err)
		assert.Nil(t, it)
	})

	t.Run("查找不存在的节点,增加后存在", func(t *testing.T) {
		t.Parallel()
		rbt := NewRBTree[int, int](ekit.ComparatorRealNumber[int])
		assert.Nil(t, rbt.Add(1, 101))
		assert.Nil(t, rbt.Add(-100, 102))
		assert.Nil(t, rbt.Add(100, 103))
		it, err := rbt.FindIt(2)
		assert.Equal(t, ErrRBTreeNotRBNode, err)
		assert.Nil(t, it)
		rbt.Add(2, 104)
		it, err = rbt.FindIt(2)
		assert.Nil(t, err)
		assert.Equal(t, 104, it.Get().Value)
	})
}

func TestIteratorDelete(t *testing.T) {
	t.Run("删除不存在的节点", func(t *testing.T) {
		t.Parallel()
		rbt := NewRBTree[int, int](ekit.ComparatorRealNumber[int])
		assert.Nil(t, rbt.Add(1, 101))
		assert.Nil(t, rbt.Add(-100, 102))
		assert.Nil(t, rbt.Add(100, 103))
		assert.Equal(t, ErrRBTreeIteratorInvalid, rbt.DeleteIt(nil))
	})
	t.Run("重复删除某个节点", func(t *testing.T) {
		t.Parallel()
		rbt := NewRBTree[int, int](ekit.ComparatorRealNumber[int])
		assert.Nil(t, rbt.Add(1, 101))
		assert.Nil(t, rbt.Add(-100, 102))
		assert.Nil(t, rbt.Add(100, 103))
		it, err := rbt.FindIt(-100)
		assert.Nil(t, err)
		assert.Equal(t, nil, rbt.DeleteIt(it))
		assert.Equal(t, ErrRBTreeIteratorInvalid, rbt.DeleteIt(it))
	})
}

func generateArray(n int) []int {
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = i
	}
	rand.Shuffle(n, func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})
	return res
}
