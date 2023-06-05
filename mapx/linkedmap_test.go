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
	"testing"

	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

func TestLinkedMap_NewLinkedHashMap(t *testing.T) {
	testCases := []struct {
		name string
		size int
	}{
		{
			name: "negative size",
			size: -1,
		},
		{
			name: "zero size",
			size: 0,
		},
		{
			name: "Positive size",
			size: 1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			multiMap := NewLinkedHashMap[testData, int](tt.size)
			assert.NotNil(t, multiMap)
		})
	}
}

func TestLinkedMap_NewLinkedTreeMap(t *testing.T) {
	testCases := []struct {
		name       string
		comparator ekit.Comparator[int]

		wantErr error
	}{
		{
			name:       "no error",
			comparator: ekit.ComparatorRealNumber[int],

			wantErr: nil,
		},
		{
			name:       "match errLinkedTreeMapComparatorIsNull error",
			comparator: nil,

			wantErr: errTreeMapComparatorIsNull,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			linkedTreeMap, err := NewLinkedTreeMap[int, int](tt.comparator)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				assert.Nil(t, linkedTreeMap)
			} else {
				assert.NotNil(t, linkedTreeMap)
			}
		})
	}
}

func TestLinkedMap_Put(t *testing.T) {
	testCases := []struct {
		name   string
		keys   []int
		values []int

		wantKeys   []int
		wantValues []int
		wantErr    error
	}{
		{
			name:   "put simple one",
			keys:   []int{1},
			values: []int{1},

			wantKeys:   []int{1},
			wantValues: []int{1},
			wantErr:    nil,
		},
		{
			name:   "put multiple",
			keys:   []int{1, 2, 3, 4},
			values: []int{1, 2, 3, 4},

			wantKeys:   []int{1, 2, 3, 4},
			wantValues: []int{1, 2, 3, 4},
			wantErr:    nil,
		},
		{
			name:   "the key include the same",
			keys:   []int{1, 1, 2, 3},
			values: []int{1, 1, 2, 3},

			wantKeys:   []int{1, 1, 2, 3},
			wantValues: []int{1, 1, 2, 3},
			wantErr:    nil,
		},
	}
	linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.keys {
				err := linkedTreeMap.Put(tt.keys[i], tt.values[i])
				assert.Equal(t, tt.wantErr, err)
			}

			for i := range tt.wantKeys {
				v, b := linkedTreeMap.Get(tt.wantKeys[i])
				assert.Equal(t, true, b)
				assert.Equal(t, tt.wantValues[i], v)
			}
		})
	}
}

func TestLinkedMap_Get(t *testing.T) {
	testCases := []struct {
		name      string
		linkedMap *LinkedMap[int, int]
		key       int

		wantValue int
		wantBool  bool
	}{
		{
			name: "not found (nil) in empty data",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				return linkedTreeMap
			}(),
			key: 1,

			wantValue: 0,
			wantBool:  false,
		},
		{
			name: "not found (nil) in data",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				return linkedTreeMap
			}(),
			key: 2,

			wantValue: 0,
			wantBool:  false,
		},
		{
			name: "found data",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				return linkedTreeMap
			}(),
			key: 1,

			wantValue: 1,
			wantBool:  true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			v, b := tt.linkedMap.Get(tt.key)
			assert.Equal(t, tt.wantBool, b)
			assert.Equal(t, tt.wantValue, v)
		})
	}
}

func TestLinkedMap_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		linkedMap *LinkedMap[int, int]

		key int

		delValue  int
		wantBool  bool
		linkedKVs []struct {
			k, v int
		}
	}{
		{
			name: "not found in empty data",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				return linkedTreeMap
			}(),

			key: 1,

			delValue:  0,
			wantBool:  false,
			linkedKVs: nil,
		},
		{
			name: "not found in data",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				return linkedTreeMap
			}(),

			key: 2,

			delValue: 0,
			wantBool: false,
			linkedKVs: []struct{ k, v int }{
				{
					k: 1,
					v: 1,
				},
			},
		},
		{
			name: "delete head | tail in one data",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				return linkedTreeMap
			}(),
			key: 1,

			delValue:  1,
			wantBool:  true,
			linkedKVs: nil,
		},
		{
			name: "delete head in many data",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				_ = linkedTreeMap.Put(2, 2)
				_ = linkedTreeMap.Put(3, 3)
				return linkedTreeMap
			}(),
			key: 1,

			delValue: 1,
			wantBool: true,
			linkedKVs: []struct{ k, v int }{
				{
					k: 2,
					v: 2,
				},
				{
					k: 3,
					v: 3,
				},
			},
		},
		{
			name: "delete tail in many data",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				_ = linkedTreeMap.Put(2, 2)
				_ = linkedTreeMap.Put(3, 3)
				return linkedTreeMap
			}(),
			key: 3,

			delValue: 3,
			wantBool: true,
			linkedKVs: []struct{ k, v int }{
				{
					k: 1,
					v: 1,
				},
				{
					k: 2,
					v: 2,
				},
			},
		},
		{
			name: "delete middle one in many data",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				_ = linkedTreeMap.Put(2, 2)
				_ = linkedTreeMap.Put(3, 3)
				return linkedTreeMap
			}(),
			key: 2,

			delValue: 2,
			wantBool: true,
			linkedKVs: []struct{ k, v int }{
				{
					k: 1,
					v: 1,
				},
				{
					k: 3,
					v: 3,
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			v, b := tt.linkedMap.Delete(tt.key)
			assert.Equal(t, tt.wantBool, b)
			assert.Equal(t, tt.delValue, v)

			idx := 0
			cur := tt.linkedMap.head
			for cur != nil {
				assert.Equal(t, tt.linkedKVs[idx].k, cur.key)
				assert.Equal(t, tt.linkedKVs[idx].v, cur.value)
				cur = cur.next
				idx++
			}
		})
	}
}

func TestLinkedMap_Keys(t *testing.T) {
	testCases := []struct {
		name      string
		linkedMap *LinkedMap[int, int]

		wantLinkedMapMapKeys []int
	}{
		{
			name: "empty",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				return linkedTreeMap
			}(),

			wantLinkedMapMapKeys: []int{},
		},
		{
			name: "single one",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				return linkedTreeMap
			}(),

			wantLinkedMapMapKeys: []int{1},
		},
		{
			name: "multiple",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				_ = linkedTreeMap.Put(2, 2)
				_ = linkedTreeMap.Put(3, 3)
				_ = linkedTreeMap.Put(4, 4)
				return linkedTreeMap
			}(),

			wantLinkedMapMapKeys: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			keys := tt.linkedMap.Keys()
			for i := range tt.wantLinkedMapMapKeys {
				assert.Equal(t, tt.wantLinkedMapMapKeys[i], keys[i])
			}
		})

	}
}

func TestLinkedMap_Values(t *testing.T) {
	testCases := []struct {
		name      string
		linkedMap *LinkedMap[int, int]

		wantLinkedMapMapValues []int
	}{
		{
			name: "empty",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				return linkedTreeMap
			}(),

			wantLinkedMapMapValues: []int{},
		},
		{
			name: "single one",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				return linkedTreeMap
			}(),

			wantLinkedMapMapValues: []int{1},
		},
		{
			name: "multiple",
			linkedMap: func() *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				_ = linkedTreeMap.Put(1, 1)
				_ = linkedTreeMap.Put(2, 2)
				_ = linkedTreeMap.Put(3, 3)
				_ = linkedTreeMap.Put(4, 4)
				return linkedTreeMap
			}(),

			wantLinkedMapMapValues: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			keys := tt.linkedMap.Keys()
			for i := range tt.wantLinkedMapMapValues {
				assert.Equal(t, tt.wantLinkedMapMapValues[i], keys[i])
			}
		})

	}
}
