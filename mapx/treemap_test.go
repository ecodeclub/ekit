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
	"testing"

	"github.com/gotomicro/ekit"
	"github.com/gotomicro/ekit/internal/tree"
	"github.com/stretchr/testify/assert"
)

func TestBuildTreeMap(t *testing.T) {
	tests := []struct {
		name       string
		m          map[int]string
		comparable ekit.Comparator[int]
		want       bool
		wantErr    error
	}{
		{
			name:       "nil",
			m:          nil,
			comparable: nil,
			want:       false,
			wantErr:    errors.New("ekit: Comparator不能为nil"),
		},
		{
			name:       "empty",
			m:          map[int]string{},
			comparable: compare(),
			want:       true,
			wantErr:    nil,
		},
		{
			name: "single",
			m: map[int]string{
				0: "0",
			},
			comparable: compare(),
			want:       true,
			wantErr:    nil,
		},
		{
			name: "multiple",
			m: map[int]string{
				0: "0",
				1: "1",
				2: "2",
			},
			comparable: compare(),
			want:       true,
			wantErr:    nil,
		},
		{
			name: "disorder",
			m: map[int]string{
				1: "1",
				2: "2",
				0: "0",
				3: "3",
				5: "5",
				4: "4",
			},
			comparable: compare(),
			want:       true,
			wantErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treeMap, err := BuildTreeMap[int, string](tt.comparable, tt.m)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			assert.Equal(t, tt.want, tree.IsRedBlackTree[int, string](treeMap.Root()))

			treeNewTreeMap, err := NewTreeMap[int, string](tt.comparable)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			assert.Equal(t, tt.want, tree.IsRedBlackTree[int, string](treeNewTreeMap.Root()))
		})

	}
}

func TestPutAll(t *testing.T) {
	var tests = []struct {
		name    string
		m       map[int]int
		wantKey []int
		wantVal []int
	}{
		{
			name:    "empty-TreeMap",
			m:       map[int]int{},
			wantVal: nil,
			wantKey: nil,
		},
		{
			name: "single",
			m: map[int]int{
				1: 1,
			},
			wantVal: []int{1},
			wantKey: []int{1},
		},
		{
			name: "multiple",
			m: map[int]int{
				1: 1,
				2: 2,
				0: 0,
				3: 3,
				5: 5,
				4: 4,
			},
			wantVal: []int{0, 1, 2, 3, 4, 5},
			wantKey: []int{0, 1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treeMap, _ := NewTreeMap[int, int](compare())
			PutAll(treeMap, tt.m)
			if !tree.IsRedBlackTree(treeMap.Root()) {
				panic(errors.New("不是红黑树"))
			}
			keys, val := treeMap.keyValue()
			assert.ElementsMatch(t, tt.wantKey, keys)
			assert.ElementsMatch(t, tt.wantVal, val)

		})
	}
}

func TestTreeMap_Get(t *testing.T) {
	var tests = []struct {
		name     string
		m        map[int]int
		findKey  int
		wantVal  int
		wantBool bool
	}{
		{
			name:     "empty-TreeMap",
			m:        map[int]int{},
			findKey:  0,
			wantVal:  0,
			wantBool: false,
		},
		{
			name: "find",
			m: map[int]int{
				1: 1,
				2: 2,
				0: 0,
				3: 3,
				5: 5,
				4: 4,
			},
			findKey:  2,
			wantVal:  2,
			wantBool: true,
		},
		{
			name: "not-find",
			m: map[int]int{
				1: 1,
				2: 2,
				0: 0,
				3: 3,
				5: 5,
				4: 4,
			},
			findKey:  6,
			wantVal:  0,
			wantBool: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treeMap, _ := NewTreeMap[int, int](compare())
			PutAll(treeMap, tt.m)
			val, b := treeMap.Get(tt.findKey)
			assert.Equal(t, tt.wantBool, b)
			assert.Equal(t, tt.wantVal, val)
		})
	}
}

func TestTreeMap_Put(t *testing.T) {
	tests := []struct {
		name    string
		k       []int
		v       []string
		wantKey []int
		wantVal []string
		wantErr error
	}{
		{
			name:    "nil",
			k:       []int{0},
			v:       nil,
			wantKey: []int{0},
			wantVal: nil,
		},
		{
			name:    "single",
			k:       []int{0},
			v:       []string{"0"},
			wantKey: []int{0},
			wantVal: []string{"0"},
			wantErr: nil,
		},
		{
			name:    "multiple",
			k:       []int{0, 1, 2},
			v:       []string{"0", "1", "2"},
			wantKey: []int{0, 1, 2},
			wantVal: []string{"0", "1", "2"},
			wantErr: nil,
		},
		{
			name:    "same",
			k:       []int{0, 0, 1, 2, 2, 3},
			v:       []string{"0", "999", "1", "998", "2", "3"},
			wantKey: []int{0, 1, 2, 3},
			wantVal: []string{"999", "1", "2", "3"},
			wantErr: nil,
		},
		{
			name:    "same",
			k:       []int{0, 0},
			v:       []string{"0", "999"},
			wantKey: []int{0},
			wantVal: []string{"999"},
			wantErr: nil,
		},
		{
			name:    "disorder",
			k:       []int{1, 2, 0, 3, 5, 4},
			v:       []string{"1", "2", "0", "3", "5", "4"},
			wantKey: []int{0, 1, 2, 3, 4, 5},
			wantVal: []string{"0", "1", "2", "3", "4", "5"},
			wantErr: nil,
		},
		{
			name:    "disorder-same",
			k:       []int{1, 3, 2, 0, 2, 3},
			v:       []string{"1", "2", "998", "0", "3", "997"},
			wantKey: []int{0, 1, 2, 3},
			wantVal: []string{"0", "1", "3", "997"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.v == nil {
				treeMapNil, _ := NewTreeMap[int, []string](compare())
				for i := 0; i < len(tt.k); i++ {
					err := treeMapNil.Put(tt.k[i], tt.v)
					if err != nil {
						assert.Equal(t, tt.wantErr, err)
						return
					}
				}
				keys, val := treeMapNil.keyValue()
				assert.ElementsMatch(t, tt.wantKey, keys)
				assert.ElementsMatch(t, tt.wantVal, val[0])
			} else {
				treeMap, _ := NewTreeMap[int, string](compare())
				for i := 0; i < len(tt.k); i++ {
					err := treeMap.Put(tt.k[i], tt.v[i])
					if err != nil {
						assert.Equal(t, tt.wantErr, err)
						return
					}
				}
				keys, val := treeMap.keyValue()
				assert.ElementsMatch(t, tt.wantKey, keys)
				assert.ElementsMatch(t, tt.wantVal, val)
			}

		})
	}
}

func TestTreeMap_Remove(t *testing.T) {
	var tests = []struct {
		name     string
		m        map[int]int
		delKey   int
		wantVal  int
		wantBool bool
	}{
		{
			name:    "empty-TreeMap",
			m:       map[int]int{},
			delKey:  0,
			wantVal: 0,
		},
		{
			name: "find",
			m: map[int]int{
				1: 1,
				2: 2,
				0: 0,
				3: 3,
				5: 5,
				4: 4,
			},
			delKey:  2,
			wantVal: 0,
		},
		{
			name: "not-find",
			m: map[int]int{
				1: 1,
				2: 2,
				0: 0,
				3: 3,
				5: 5,
				4: 4,
			},
			delKey:  6,
			wantVal: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treeMap, _ := NewTreeMap[int, int](compare())
			treeMap.Remove(tt.delKey)
			assert.Equal(t, true, tree.IsRedBlackTree[int, int](treeMap.Root()))
			val, err := treeMap.Get(tt.delKey)
			assert.Equal(t, tt.wantBool, err)
			assert.Equal(t, tt.wantVal, val)
		})
	}
}

func compare() ekit.Comparator[int] {
	return ekit.ComparatorRealNumber[int]
}

type kv[Key any, Val any] struct {
	ks   []Key
	vals []Val
}

func (treeMap *TreeMap[int, string]) keyValue() ([]int, []string) {
	treeNode := treeMap.Root()
	var m = &kv[int, string]{}
	if treeMap.Size() > 0 {
		midOrder(treeNode, m)
		return m.ks, m.vals
	}
	return nil, nil
}

func midOrder[Key any, Val any](node *tree.RBNode[Key, Val], m *kv[Key, Val]) {
	// 先遍历左子树
	if node.Left() != nil {
		midOrder(node.Left(), m)
	}
	// 再遍历自己
	if node != nil {
		m.ks = append(m.ks, node.Key)
		m.vals = append(m.vals, node.Value)
	}
	// 最后遍历右子树
	if node.Right() != nil {
		midOrder(node.Right(), m)
	}
}
