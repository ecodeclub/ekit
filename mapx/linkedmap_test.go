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
	"errors"
	"testing"

	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

var (
	fakeErr = errors.New("fakeMap: put error")
)

type fakeMap[K any, V any] struct {
	*LinkedMap[K, V]
	count          int
	activeFirstErr bool
}

func (f *fakeMap[K, V]) Put(key K, val V) error {
	f.count++
	if f.activeFirstErr {
		f.activeFirstErr = false
		return fakeErr
	}
	if f.count == 3 {
		return fakeErr
	}
	if f.count == 5 {
		return fakeErr
	}
	return f.LinkedMap.Put(key, val)
}

func newLinkedFakeMap[K any, V any](activeFirstErr bool, comparator ekit.Comparator[K]) (*LinkedMap[K, V], error) {
	treeMap, err := NewLinkedTreeMap[K, *linkedKV[K, V]](comparator)
	if err != nil {
		return nil, err
	}
	fm := &fakeMap[K, *linkedKV[K, V]]{LinkedMap: treeMap, activeFirstErr: activeFirstErr}
	head := &linkedKV[K, V]{}
	tail := &linkedKV[K, V]{next: head, prev: head}
	head.prev, head.next = tail, tail
	return &LinkedMap[K, V]{
		m:    fm,
		head: head,
		tail: tail,
	}, nil
}

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
		name      string
		linkedMap func(t *testing.T) *LinkedMap[int, int]
		keys      []int
		values    []int

		wantKeys   []int
		wantValues []int
		wantErrs   []error
	}{
		{
			name: "put single key",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				return linkedTreeMap
			},
			keys:   []int{1},
			values: []int{1},

			wantKeys:   []int{1},
			wantValues: []int{1},
			wantErrs:   []error{nil},
		},
		{
			name: "put multiple keys",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				return linkedTreeMap
			},
			keys:   []int{1, 2, 3, 4},
			values: []int{1, 2, 3, 4},

			wantKeys:   []int{1, 2, 3, 4},
			wantValues: []int{1, 2, 3, 4},
			wantErrs:   []error{nil, nil, nil, nil},
		},
		{
			name: "change value of single key",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				return linkedTreeMap
			},
			keys:   []int{1, 1, 2, 3},
			values: []int{1, 11, 2, 3},

			wantKeys:   []int{1, 2, 3},
			wantValues: []int{11, 2, 3},
			wantErrs:   []error{nil, nil, nil, nil},
		},
		{
			name: "change value of multiple keys",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				return linkedTreeMap
			},
			keys:   []int{1, 1, 2, 2, 3},
			values: []int{1, 11, 2, 22, 3},

			wantKeys:   []int{1, 2, 3},
			wantValues: []int{11, 22, 3},
			wantErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name: "get error when put single key",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedFakeMap, err := newLinkedFakeMap[int, int](true, ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				return linkedFakeMap
			},
			keys:   []int{1},
			values: []int{1},

			wantKeys:   []int{},
			wantValues: []int{},
			wantErrs:   []error{fakeErr},
		},
		{
			name: "get multiple errors when put multiple keys",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedFakeMap, err := newLinkedFakeMap[int, int](true, ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				return linkedFakeMap
			},
			keys:   []int{1, 2, 3, 4, 5},
			values: []int{1, 2, 3, 4, 5},

			wantKeys:   []int{2, 4},
			wantValues: []int{2, 4},
			wantErrs:   []error{fakeErr, nil, fakeErr, nil, fakeErr},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			errs := make([]error, 0)
			linkedMap := tt.linkedMap(t)
			for i := range tt.keys {
				err := linkedMap.Put(tt.keys[i], tt.values[i])
				errs = append(errs, err)
			}

			for i := range tt.wantKeys {
				v, b := linkedMap.Get(tt.wantKeys[i])
				assert.Equal(t, true, b)
				assert.Equal(t, tt.wantValues[i], v)
			}

			assert.Equal(t, tt.wantKeys, linkedMap.Keys())
			assert.Equal(t, tt.wantValues, linkedMap.Values())
			assert.Equal(t, tt.wantErrs, errs)
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
				err = linkedTreeMap.Put(1, 1)
				assert.NoError(t, err)
				return linkedTreeMap
			},
			key: 2,

			wantValue: 0,
			wantBool:  false,
		},
		{
			name: "find value in linked map",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				err = linkedTreeMap.Put(1, 1)
				assert.NoError(t, err)
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

		delValue   int
		wantBool   bool
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

			delValue:   0,
			wantBool:   false,
			wantKeys:   []int{1},
			wantValues: []int{1},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			linkedMap := tt.linkedMap(t)
			v, b := linkedMap.Delete(tt.key)
			assert.Equal(t, tt.wantBool, b)
			assert.Equal(t, tt.delValue, v)

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
		{
			name: "put k1 → put k2 → put k3 → delete k2",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, err := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, err)
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				assert.NoError(t, linkedTreeMap.Put(2, 2))
				assert.NoError(t, linkedTreeMap.Put(3, 3))
				v, ok := linkedTreeMap.Delete(2)
				assert.Equal(t, 2, v)
				assert.Equal(t, true, ok)

				return linkedTreeMap
			},

			wantKeys:   []int{1, 3},
			wantValues: []int{1, 3},
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

func TestLinkedMap_Len(t *testing.T) {
	testCases := []struct {
		name      string
		linkedMap func(t *testing.T) *LinkedMap[int, int]

		wantLen int64
	}{
		{
			name: "empty linked map",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				return linkedTreeMap
			},

			wantLen: 0,
		},
		{
			name: "not empty linked map",
			linkedMap: func(t *testing.T) *LinkedMap[int, int] {
				linkedTreeMap, _ := NewLinkedTreeMap[int, int](ekit.ComparatorRealNumber[int])
				assert.NoError(t, linkedTreeMap.Put(1, 1))
				assert.NoError(t, linkedTreeMap.Put(2, 2))
				assert.NoError(t, linkedTreeMap.Put(3, 3))
				return linkedTreeMap
			},

			wantLen: 3,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantLen, tt.linkedMap(t).Len())
		})
	}
}
