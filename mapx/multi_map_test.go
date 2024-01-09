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

func getMultiTreeMap() *MultiMap[int, int] {
	multiTreeMap, _ := NewMultiTreeMap[int, int](ekit.ComparatorRealNumber[int])
	return multiTreeMap
}
func getMultiHashMap() *MultiMap[testData, int] {
	return NewMultiHashMap[testData, int](10)
}

func TestMultiMap_NewMultiHashMap(t *testing.T) {
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
			multiMap := NewMultiHashMap[testData, int](tt.size)
			assert.NotNil(t, multiMap)
		})
	}
}

func TestMultiMap_NewMultiTreeMap(t *testing.T) {
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
			name:       "match errMultiMapComparatorIsNull error",
			comparator: nil,

			wantErr: errTreeMapComparatorIsNull,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			multiMap, err := NewMultiTreeMap[int, int](tt.comparator)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				assert.Nil(t, multiMap)
			} else {
				assert.NotNil(t, multiMap)
			}
		})
	}
}

func TestNewMultiBuiltinMap(t *testing.T) {
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
			multiMap := NewMultiBuiltinMap[testData, int](tt.size)
			assert.NotNil(t, multiMap)
		})
	}
}

func TestMultiMap_Keys(t *testing.T) {
	testCases := []struct {
		name         string
		multiTreeMap *MultiMap[int, int]
		multiHashMap *MultiMap[testData, int]

		wantMultiTreeMapKeys []int
		wantMultiHashMapKeys []testData
	}{
		{
			name: "empty",
			multiTreeMap: func() *MultiMap[int, int] {
				return getMultiTreeMap()
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				return getMultiHashMap()
			}(),

			wantMultiTreeMapKeys: []int{},
			wantMultiHashMapKeys: []testData{},
		},
		{
			name: "single one",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				return multiHashMap
			}(),

			wantMultiTreeMapKeys: []int{1},
			wantMultiHashMapKeys: []testData{{id: 1}},
		},
		{
			name: "multiple",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				_ = multiTreeMap.Put(2, 2)
				_ = multiTreeMap.Put(3, 3)
				_ = multiTreeMap.Put(4, 4)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				_ = multiHashMap.Put(testData{id: 2}, 2)
				_ = multiHashMap.Put(testData{id: 3}, 3)
				_ = multiHashMap.Put(testData{id: 4}, 4)
				return multiHashMap
			}(),

			wantMultiTreeMapKeys: []int{1, 2, 3, 4},
			wantMultiHashMapKeys: []testData{
				{id: 1},
				{id: 2},
				{id: 3},
				{id: 4},
			},
		},
	}
	for _, tt := range testCases {
		t.Run("MultiTreeMap", func(t *testing.T) {
			assert.ElementsMatch(t, tt.wantMultiTreeMapKeys, tt.multiTreeMap.Keys())
		})

		t.Run("MultiHashMap", func(t *testing.T) {
			assert.ElementsMatch(t, tt.wantMultiHashMapKeys, tt.multiHashMap.Keys())
		})
	}
}

func TestMultiMap_Values(t *testing.T) {
	testCases := []struct {
		name         string
		multiTreeMap *MultiMap[int, int]
		multiHashMap *MultiMap[testData, int]

		wantValues [][]int
	}{
		{
			name: "empty",
			multiTreeMap: func() *MultiMap[int, int] {
				return getMultiTreeMap()
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				return getMultiHashMap()
			}(),

			wantValues: [][]int{},
		},
		{
			name: "single one",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				return multiHashMap
			}(),

			wantValues: [][]int{{1}},
		},
		{
			name: "multiple",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				_ = multiTreeMap.Put(2, 2)
				_ = multiTreeMap.Put(3, 3)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				_ = multiHashMap.Put(testData{id: 2}, 2)
				_ = multiHashMap.Put(testData{id: 3}, 3)
				return multiHashMap
			}(),

			wantValues: [][]int{{1}, {2}, {3}},
		},
	}
	t.Run("MultiTreeMap", func(t *testing.T) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				assert.ElementsMatch(t, tt.wantValues, tt.multiTreeMap.Values())
			})
		}
	})
	t.Run("MultiHashMap", func(t *testing.T) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				assert.ElementsMatch(t, tt.wantValues, tt.multiHashMap.Values())
			})
		}
	})

}

func TestMultiMap_Put(t *testing.T) {
	testCases := []struct {
		name   string
		keys   []int
		values []int

		wantKeys   []int
		wantValues [][]int
		wantErr    error
	}{
		{
			name:   "put simple one",
			keys:   []int{1},
			values: []int{1},

			wantKeys:   []int{1},
			wantValues: [][]int{{1}},
			wantErr:    nil,
		},
		{
			name:   "put multiple",
			keys:   []int{1, 2, 3, 4},
			values: []int{1, 2, 3, 4},

			wantKeys:   []int{1, 2, 3, 4},
			wantValues: [][]int{{1}, {2}, {3}, {4}},
			wantErr:    nil,
		},
		{
			name:   "the key include the same",
			keys:   []int{1, 2, 1, 4},
			values: []int{1, 2, 3, 4},

			wantKeys: []int{1, 2, 4},
			wantValues: [][]int{
				{1, 3},
				{2},
				{4},
			},
			wantErr: nil,
		},
	}
	for _, tt := range testCases {
		t.Run("MultiTreeMap", func(t *testing.T) {
			multiTreeMap, _ := NewMultiTreeMap[int, int](ekit.ComparatorRealNumber[int])
			for i := range tt.keys {
				err := multiTreeMap.Put(tt.keys[i], tt.values[i])
				assert.Equal(t, tt.wantErr, err)
			}

			for i := range tt.wantKeys {
				v, b := multiTreeMap.Get(tt.wantKeys[i])
				assert.Equal(t, true, b)
				assert.Equal(t, tt.wantValues[i], v)
			}
		})

		t.Run("MultiHashMap", func(t *testing.T) {
			multiHashMap := NewMultiHashMap[testData, int](10)
			for i := range tt.keys {
				err := multiHashMap.Put(testData{id: tt.keys[i]}, tt.values[i])
				assert.Equal(t, tt.wantErr, err)
			}

			for i := range tt.wantKeys {
				v, b := multiHashMap.Get(testData{id: tt.wantKeys[i]})
				assert.Equal(t, true, b)
				assert.Equal(t, tt.wantValues[i], v)
			}
		})
	}
}

func TestMultiMap_Get(t *testing.T) {
	testCases := []struct {
		name         string
		multiTreeMap *MultiMap[int, int]
		multiHashMap *MultiMap[testData, int]
		key          int

		wantValue []int
		wantBool  bool
	}{
		{
			name: "not found (nil) in empty data",
			multiTreeMap: func() *MultiMap[int, int] {
				return getMultiTreeMap()
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				return getMultiHashMap()
			}(),
			key: 1,

			wantValue: nil,
			wantBool:  false,
		},
		{
			name: "not found (nil) in data",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				_ = multiTreeMap.Put(2, 2)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				_ = multiHashMap.Put(testData{id: 2}, 2)
				return multiHashMap
			}(),
			key: 3,

			wantValue: nil,
			wantBool:  false,
		},
		{
			name: "found data",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				return multiHashMap
			}(),
			key: 1,

			wantValue: []int{1},
			wantBool:  true,
		},
	}
	for _, tt := range testCases {
		t.Run("MultiTreeMap", func(t *testing.T) {
			v, b := tt.multiTreeMap.Get(tt.key)
			assert.Equal(t, tt.wantBool, b)
			assert.ElementsMatch(t, tt.wantValue, v)
		})

		t.Run("MultiHashMap", func(t *testing.T) {
			v2, b2 := tt.multiHashMap.Get(testData{id: tt.key})
			assert.Equal(t, tt.wantBool, b2)
			assert.ElementsMatch(t, tt.wantValue, v2)
		})
	}
}

func TestMultiMap_Delete(t *testing.T) {
	testCases := []struct {
		name         string
		multiTreeMap *MultiMap[int, int]
		multiHashMap *MultiMap[testData, int]

		key int

		delValue []int
		wantBool bool
	}{
		{
			name: "not found in empty data",
			multiTreeMap: func() *MultiMap[int, int] {
				return getMultiTreeMap()
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				return getMultiHashMap()
			}(),

			key: 1,

			delValue: nil,
			wantBool: false,
		},
		{
			name: "not found in data",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				return multiHashMap
			}(),

			key: 2,

			delValue: nil,
			wantBool: false,
		},
		{
			name: "found and deleted",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				_ = multiTreeMap.Put(2, 2)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				_ = multiHashMap.Put(testData{id: 2}, 2)
				return multiHashMap
			}(),
			key: 1,

			delValue: []int{1},
			wantBool: true,
		},
	}
	for _, tt := range testCases {
		t.Run("MultiTreeMap", func(t *testing.T) {
			v, b := tt.multiTreeMap.Delete(tt.key)
			assert.Equal(t, tt.wantBool, b)
			assert.ElementsMatch(t, tt.delValue, v)
		})
		t.Run("MultiHashMap", func(t *testing.T) {
			v, b := tt.multiHashMap.Delete(testData{id: tt.key})
			assert.Equal(t, tt.wantBool, b)
			assert.ElementsMatch(t, tt.delValue, v)
		})
	}
}

func TestMultiMap_PutMany(t *testing.T) {
	testCases := []struct {
		name   string
		keys   []int
		values [][]int

		wantKeys   []int
		wantValues [][]int
		wantErr    error
	}{
		{
			name:   "one to one",
			keys:   []int{1},
			values: [][]int{{1}},

			wantKeys:   []int{1},
			wantValues: [][]int{{1}},
			wantErr:    nil,
		},
		{
			name:   "many [one to one]",
			keys:   []int{1, 2, 3},
			values: [][]int{{1}, {2}, {3}},

			wantKeys:   []int{1, 2, 3},
			wantValues: [][]int{{1}, {2}, {3}},
			wantErr:    nil,
		},
		{
			name:   "one to many",
			keys:   []int{1},
			values: [][]int{{1, 2, 3}},

			wantKeys: []int{1},
			wantValues: [][]int{
				{1, 2, 3},
			},
			wantErr: nil,
		},
		{
			name:   "many [one to many]",
			keys:   []int{1, 2, 3},
			values: [][]int{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}},

			wantKeys: []int{1, 2, 3},
			wantValues: [][]int{
				{1, 2, 3},
				{1, 2, 3},
				{1, 2, 3},
			},
			wantErr: nil,
		},
		{
			name:   "the key include the same for append one",
			keys:   []int{1, 1},
			values: [][]int{{1, 2, 3, 4, 5}, {6}},

			wantKeys: []int{1},
			wantValues: [][]int{
				{1, 2, 3, 4, 5, 6},
			},
			wantErr: nil,
		},
		{
			name:   "the key include the same for append many",
			keys:   []int{1, 1},
			values: [][]int{{1}, {2, 3, 4, 5, 6}},

			wantKeys: []int{1},
			wantValues: [][]int{
				{1, 2, 3, 4, 5, 6},
			},
			wantErr: nil,
		},
	}
	for _, tt := range testCases {
		t.Run("MultiTreeMap", func(t *testing.T) {
			multiTreeMap, _ := NewMultiTreeMap[int, int](ekit.ComparatorRealNumber[int])
			for i := range tt.keys {
				err := multiTreeMap.PutMany(tt.keys[i], tt.values[i]...)
				assert.Equal(t, tt.wantErr, err)
			}

			for i := range tt.wantKeys {
				v, b := multiTreeMap.Get(tt.wantKeys[i])
				assert.Equal(t, true, b)
				assert.Equal(t, tt.wantValues[i], v)
			}
		})

		t.Run("MultiHashMap", func(t *testing.T) {
			multiHashMap := NewMultiHashMap[testData, int](10)
			for i := range tt.keys {
				err := multiHashMap.PutMany(testData{id: tt.keys[i]}, tt.values[i]...)
				assert.Equal(t, tt.wantErr, err)
			}

			for i := range tt.wantKeys {
				v, b := multiHashMap.Get(testData{id: tt.wantKeys[i]})
				assert.Equal(t, true, b)
				assert.Equal(t, tt.wantValues[i], v)
			}
		})
	}
}

func TestMultiMap_Len(t *testing.T) {
	testCases := []struct {
		name         string
		multiTreeMap *MultiMap[int, int]
		multiHashMap *MultiMap[testData, int]
		wantLen      int64
	}{
		{
			name:         "empty",
			multiTreeMap: getMultiTreeMap(),
			multiHashMap: getMultiHashMap(),

			wantLen: 0,
		},
		{
			name: "single",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				return multiHashMap
			}(),

			wantLen: 1,
		},
		{
			name: "multiple",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				_ = multiTreeMap.Put(2, 2)
				_ = multiTreeMap.Put(3, 3)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				_ = multiHashMap.Put(testData{id: 2}, 2)
				_ = multiHashMap.Put(testData{id: 3}, 3)
				return multiHashMap
			}(),

			wantLen: 3,
		},
		{
			name: "multiple with same key",
			multiTreeMap: func() *MultiMap[int, int] {
				multiTreeMap := getMultiTreeMap()
				_ = multiTreeMap.Put(1, 1)
				_ = multiTreeMap.Put(1, 2)
				_ = multiTreeMap.Put(1, 3)
				return multiTreeMap
			}(),
			multiHashMap: func() *MultiMap[testData, int] {
				multiHashMap := getMultiHashMap()
				_ = multiHashMap.Put(testData{id: 1}, 1)
				_ = multiHashMap.Put(testData{id: 1}, 2)
				_ = multiHashMap.Put(testData{id: 1}, 3)
				return multiHashMap
			}(),
			wantLen: 1,
		},
	}
	for _, tt := range testCases {
		t.Run("MultiTreeMap", func(t *testing.T) {
			assert.Equal(t, tt.wantLen, tt.multiTreeMap.Len())
		})

		t.Run("MultiHashMap", func(t *testing.T) {
			assert.Equal(t, tt.wantLen, tt.multiHashMap.Len())
		})
	}
}
