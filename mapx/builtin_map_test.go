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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltinMap_Delete(t *testing.T) {
	testCases := []struct {
		name string
		data map[string]string

		key string

		wantVal     string
		wantDeleted bool
	}{
		{
			name: "deleted",
			data: map[string]string{
				"key1": "val1",
			},
			key: "key1",

			wantVal:     "val1",
			wantDeleted: true,
		},
		{
			name: "key not exist",
			data: map[string]string{
				"key1": "val1",
			},
			key: "key2",
		},
		{
			name: "nil map",
			key:  "key2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := builtinMapOf[string, string](tc.data)
			val, ok := m.Delete(tc.key)
			assert.Equal(t, tc.wantDeleted, ok)
			assert.Equal(t, tc.wantVal, val)
			_, ok = m.data[tc.key]
			assert.False(t, ok)
		})
	}
}

func TestBuiltinMap_Get(t *testing.T) {
	testCases := []struct {
		name string
		data map[string]string

		key string

		wantVal string
		found   bool
	}{
		{
			name: "found",
			data: map[string]string{
				"key1": "val1",
			},
			key: "key1",

			wantVal: "val1",
			found:   true,
		},
		{
			name: "key not exist",
			data: map[string]string{
				"key1": "val1",
			},
			key: "key2",
		},
		{
			name: "nil map",
			key:  "key2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := builtinMapOf[string, string](tc.data)
			val, ok := m.Get(tc.key)
			assert.Equal(t, tc.found, ok)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestBuiltinMap_Put(t *testing.T) {
	testCases := []struct {
		name string

		key string
		val string
		cap int

		wantErr error
	}{
		{
			name: "puted",
			key:  "key1",
			val:  "val1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := newBuiltinMap[string, string](tc.cap)
			err := m.Put(tc.key, tc.val)
			assert.Equal(t, tc.wantErr, err)
			v, ok := m.data[tc.key]
			assert.True(t, ok)
			assert.Equal(t, tc.val, v)
		})
	}
}

func TestBuiltinMap_Keys(t *testing.T) {
	testCases := []struct {
		name string
		data map[string]string

		wantKeys []string
	}{
		{
			name: "got keys",
			data: map[string]string{
				"key1": "val1",
				"key2": "val1",
				"key3": "val1",
				"key4": "val1",
			},
			wantKeys: []string{"key1", "key2", "key3", "key4"},
		},
		{
			name:     "empty map",
			data:     map[string]string{},
			wantKeys: []string{},
		},
		{
			name:     "nil map",
			wantKeys: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := builtinMapOf[string, string](tc.data)
			keys := m.Keys()
			assert.ElementsMatch(t, tc.wantKeys, keys)
		})
	}
}

func TestBuiltinMap_Values(t *testing.T) {
	testCases := []struct {
		name string
		data map[string]string

		wantValues []string
	}{
		{
			name: "got values",
			data: map[string]string{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
				"key4": "val4",
			},
			wantValues: []string{"val1", "val2", "val3", "val4"},
		},
		{
			name:       "empty map",
			data:       map[string]string{},
			wantValues: []string{},
		},
		{
			name:       "nil map",
			wantValues: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := builtinMapOf[string, string](tc.data)
			vals := m.Values()
			assert.ElementsMatch(t, tc.wantValues, vals)
		})
	}
}

func TestBuiltinMap_Len(t *testing.T) {
	testCases := []struct {
		name string
		data map[string]string

		wantLen int64
	}{
		{
			name: "got len",
			data: map[string]string{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
				"key4": "val4",
			},
			wantLen: 4,
		},
		{
			name:    "empty map",
			data:    map[string]string{},
			wantLen: 0,
		},
		{
			name:    "nil map",
			wantLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := builtinMapOf[string, string](tc.data)
			assert.Equal(t, tc.wantLen, m.Len())
		})
	}
}

func TestBuiltinMap_Iterate(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		expectedLen int
		inputStart  int
		inputEnd    int
		// 如果为true则遍历结束
		endConditionFunc func(key int) bool
	}{
		{
			name:        "builtinMap 为空",
			expectedLen: 0,
			inputStart:  1,
			inputEnd:    0,
			endConditionFunc: func(key int) bool {
				return false
			},
		},
		{
			name:        "builtinMap 有100个元素,遍历所有元素",
			expectedLen: 100,
			inputStart:  1,
			inputEnd:    100,
			endConditionFunc: func(key int) bool {
				return false
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			builtinMap := newBuiltinMap[int, int](0)
			for i := testCase.inputStart; i <= testCase.inputEnd; i++ {
				assert.Nil(t, builtinMap.Put(i, i))
			}
			arr := make([]int, 0)
			builtinMap.Iterate(func(key, value int) bool {
				if testCase.endConditionFunc(key) {
					return false
				}
				arr = append(arr, value)
				return true
			})
			assert.Equal(t, testCase.expectedLen, len(arr))
			sort.Ints(arr)
			for i := 0; i < testCase.expectedLen; i++ {
				assert.Equal(t, testCase.inputStart+i, arr[i])
			}
		})
	}
}

func TestBuiltinMap_Iterate_OnlyIterateHalf(t *testing.T) {
	builtinMap := newBuiltinMap[int, int](0)
	n := 100
	for i := 1; i <= n; i++ {
		assert.Nil(t, builtinMap.Put(i, i))
	}
	arr := make([]int, 0)
	builtinMap.Iterate(func(key, value int) bool {
		if len(arr) >= n/2 {
			return false
		}
		arr = append(arr, value)
		return true
	})
	assert.Equal(t, n/2, len(arr))
	sort.Ints(arr)
}

func builtinMapOf[K comparable, V any](data map[K]V) *builtinMap[K, V] {
	return &builtinMap[K, V]{data: data}
}
