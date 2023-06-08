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
			assert.Equal(t, multiMap.Keys(), []testData{})
			assert.Equal(t, multiMap.Values(), []int{})
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
				assert.Equal(t, linkedTreeMap.Keys(), []int{})
				assert.Equal(t, linkedTreeMap.Values(), []int{})
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
			name:   "change value of single key",
			keys:   []int{1, 1, 2, 3},
			values: []int{1, 11, 2, 3},

			wantKeys:   []int{1, 2, 3},
			wantValues: []int{11, 2, 3},
			wantErr:    nil,
		},
		{
			name:   "change value of multiple key",
			keys:   []int{1, 1, 2, 2, 3},
			values: []int{1, 11, 2, 22, 3},

			wantKeys:   []int{1, 2, 3},
			wantValues: []int{11, 22, 3},
			wantErr:    nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
			assert.NoError(t, err)
			for i := range tt.keys {
				err := linkedTreeMap.Put(tt.keys[i], tt.values[i])
				assert.Equal(t, tt.wantErr, err)
			}

			for i := range tt.wantKeys {
				v, b := linkedTreeMap.Get(tt.wantKeys[i])
				assert.Equal(t, true, b)
				assert.Equal(t, tt.wantValues[i], v)
			}

			assert.Equal(t, tt.wantKeys, linkedTreeMap.Keys())
			assert.Equal(t, tt.wantValues, linkedTreeMap.Values())
		})
	}
}

func TestLinkedMap_Get(t *testing.T) {
	testCases := []struct {
		name      string
		linkedMap func(t *testing.T) *LinkedMap[int, int]
		key       int

		wantValue int
		wantBool  bool
	}{
		{
			name: "can not find value in empty linked map",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				return linkedTreeMap
			},
			key: 1,

			wantValue: 0,
			wantBool:  false,
		},
		{
			name: "can not find value in linked map",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				_ = linkedTreeMap.Put(1, 1)
				return linkedTreeMap
			},
			key: 2,

			wantValue: 0,
			wantBool:  false,
		},
		{
			name: "found data",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				_ = linkedTreeMap.Put(1, 1)
				return linkedTreeMap
			},
			key: 1,

			wantValue: 1,
			wantBool:  true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			v, b := tt.linkedMap(t).Get(tt.key)
			assert.Equal(t, tt.wantBool, b)
			assert.Equal(t, tt.wantValue, v)
		})
	}
}

func TestLinkedMap_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		linkedMap func(t *testing.T) *LinkedMap[int, int]

		key int

		delValue  int
		wantBool  bool
		linkedKVs []struct {
			k, v int
		}
		wantKeys   []int
		wantValues []int
	}{
		{
			name: "delete key in empty linked map",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				return linkedTreeMap
			},

			key: 1,

			delValue:   0,
			wantBool:   false,
			linkedKVs:  nil,
			wantKeys:   []int{},
			wantValues: []int{},
		},
		{
			name: "delete unknown key in not empty linked map",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				return linkedTreeMap
			},

			key: 2,

			delValue: 0,
			wantBool: false,
			linkedKVs: []struct{ k, v int }{
				{
					k: 1,
					v: 1,
				},
			},
			wantKeys:   []int{1},
			wantValues: []int{1},
		},
		{
			name: "delete head in the data including one",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				return linkedTreeMap
			},
			key: 1,

			delValue:   1,
			wantBool:   true,
			linkedKVs:  nil,
			wantKeys:   []int{},
			wantValues: []int{},
		},
		{
			name: "delete head in the data including much",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				assert.NoError(t, linkedTreeMap.Put(2, 2))
				assert.NoError(t, linkedTreeMap.Put(3, 3))
				return linkedTreeMap
			},
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
			wantKeys:   []int{2, 3},
			wantValues: []int{2, 3},
		},
		{
			name: "delete tail in the data including much",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				assert.NoError(t, linkedTreeMap.Put(2, 2))
				assert.NoError(t, linkedTreeMap.Put(3, 3))
				return linkedTreeMap
			},
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
			wantKeys:   []int{1, 2},
			wantValues: []int{1, 2},
		},
		{
			name: "delete middle one in the data including much",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				assert.NoError(t, linkedTreeMap.Put(2, 2))
				assert.NoError(t, linkedTreeMap.Put(3, 3))
				return linkedTreeMap
			},
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
			wantKeys:   []int{1, 3},
			wantValues: []int{1, 3},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			linkedMap := tt.linkedMap(t)
			v, b := linkedMap.Delete(tt.key)
			assert.Equal(t, tt.wantBool, b)
			assert.Equal(t, tt.delValue, v)

			for i, cur := 0, linkedMap.head.next; i < linkedMap.length; i++ {
				assert.Equal(t, tt.linkedKVs[i].k, cur.key)
				assert.Equal(t, tt.linkedKVs[i].v, cur.value)
				cur = cur.next
			}
			assert.Equal(t, tt.wantKeys, linkedMap.Keys())
			assert.Equal(t, tt.wantValues, linkedMap.Values())
		})
	}
}

func TestLinkedMap_PutAndDelete(t *testing.T) {
	testCases := []struct {
		name      string
		linkedMap func(t *testing.T) *LinkedMap[int, int]

		wantKeys   []int
		wantValues []int
	}{
		{
			name: "put k1 → delete k1",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				v, ok := linkedTreeMap.Delete(1)
				assert.Equal(t, 1, v)
				assert.Equal(t, true, ok)
				return linkedTreeMap
			},

			wantKeys:   []int{},
			wantValues: []int{},
		},
		{
			name: "put k1 → delete k1 → put k1",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				v, ok := linkedTreeMap.Delete(1)
				assert.Equal(t, 1, v)
				assert.Equal(t, true, ok)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				return linkedTreeMap
			},

			wantKeys:   []int{1},
			wantValues: []int{1},
		},
		{
			name: "put k1 → delete k1 → put k2",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				v, ok := linkedTreeMap.Delete(1)
				assert.Equal(t, 1, v)
				assert.Equal(t, true, ok)
				assert.NoError(t, linkedTreeMap.Put(2, 2))
				return linkedTreeMap
			},

			wantKeys:   []int{2},
			wantValues: []int{2},
		},
		{
			name: "put k1 → put k1 → delete k1",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				assert.NoError(t, linkedTreeMap.Put(1, 2))
				v, ok := linkedTreeMap.Delete(1)
				assert.Equal(t, 2, v)
				assert.Equal(t, true, ok)
				return linkedTreeMap
			},

			wantKeys:   []int{},
			wantValues: []int{},
		},
		{
			name: "put k1 → put k2 → delete k1",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				assert.NoError(t, linkedTreeMap.Put(2, 2))
				v, ok := linkedTreeMap.Delete(1)
				assert.Equal(t, 1, v)
				assert.Equal(t, true, ok)
				return linkedTreeMap
			},

			wantKeys:   []int{2},
			wantValues: []int{2},
		},
		{
			name: "put k1 → put k2 → delete k2",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				assert.NoError(t, linkedTreeMap.Put(2, 2))
				v, ok := linkedTreeMap.Delete(2)
				assert.Equal(t, 2, v)
				assert.Equal(t, true, ok)
				return linkedTreeMap
			},

			wantKeys:   []int{1},
			wantValues: []int{1},
		},
		{
			name: "put k1 → delete k1 → put k2 → put k3",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				v, ok := linkedTreeMap.Delete(1)
				assert.Equal(t, 1, v)
				assert.Equal(t, true, ok)
				assert.NoError(t, linkedTreeMap.Put(2, 2))
				assert.NoError(t, linkedTreeMap.Put(3, 3))

				return linkedTreeMap
			},

			wantKeys:   []int{2, 3},
			wantValues: []int{2, 3},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			linkedMap := tt.linkedMap(t)
			for i := range tt.wantKeys {
				v, b := linkedMap.Get(tt.wantKeys[i])
				assert.Equal(t, true, b)
				assert.Equal(t, tt.wantValues[i], v)
			}
			assert.Equal(t, tt.wantKeys, linkedMap.Keys())
			assert.Equal(t, tt.wantValues, linkedMap.Values())
		})
	}
}
