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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestHashMap(t *testing.T) {
	testKV := []struct {
		key testData
		val int
	}{
		{
			key: testData{
				id: 1,
			},
			val: 1,
		},
		{
			key: testData{
				id: 2,
			},
			val: 2,
		},
		{
			key: testData{
				id: 3,
			},
			val: 3,
		},
		{
			key: testData{
				id: 11,
			},
			val: 11,
		},
		{
			key: testData{
				id: 1,
			},
			val: 101,
		},
	}
	myhashmap := NewHashMap[testData, int](10)
	for _, kv := range testKV {
		err := myhashmap.Put(kv.key, kv.val)
		if err != nil {
			require.NoError(t, err)
		}
	}

	wantHashMap := NewHashMap[testData, int](10)
	wantHashMap.hashmap = map[uint64]*node[testData, int]{
		1: &node[testData, int]{
			key:   testData{id: 1},
			value: 101,
			next: &node[testData, int]{
				key:   testData{id: 11},
				value: 11,
			},
		},
		2: wantHashMap.newNode(newTestData(2), 2),
		3: wantHashMap.newNode(newTestData(3), 3),
	}

	assert.Equal(t, wantHashMap.hashmap, myhashmap.hashmap)
	getTestCases := []struct {
		name    string
		key     testData
		wantVal any
		isFound bool
	}{
		{
			name: "get normal val",
			key: testData{
				id: 1,
			},
			wantVal: 101,
			isFound: true,
		},
		{
			name: "hash conflicts",
			key: testData{
				id: 11,
			},
			wantVal: 11,
			isFound: true,
		},
		{
			name: "hash not Found",
			key: testData{
				id: 8,
			},
			isFound: false,
		},
		{
			name: "val not Found",
			key: testData{
				id: 21,
			},
			isFound: false,
		},
	}
	for _, tc := range getTestCases {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := myhashmap.Get(tc.key)
			assert.Equal(t, tc.isFound, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}

}
func TestHashMap_Delete(t *testing.T) {
	testcases := []struct {
		name        string
		key         testData
		setHashMap  func() map[uint64]*node[testData, int]
		wantHashMap func() map[uint64]*node[testData, int]
		wantVal     int
		isFound     bool
	}{
		{
			name: "hash not found",
			setHashMap: func() map[uint64]*node[testData, int] {
				return map[uint64]*node[testData, int]{}
			},
			isFound: false,
			key: testData{
				id: 1,
			},
		},
		{
			name: "key not found",
			setHashMap: func() map[uint64]*node[testData, int] {
				return map[uint64]*node[testData, int]{
					1: &node[testData, int]{
						key: testData{
							id: 1,
						},
						value: 1,
					},
				}
			},
			isFound: false,
			key: testData{
				id: 11,
			},
		},
		{
			name: "many link elements delete first",
			setHashMap: func() map[uint64]*node[testData, int] {
				return map[uint64]*node[testData, int]{
					1: &node[testData, int]{
						key: testData{
							id: 1,
						},
						value: 1,
						next: &node[testData, int]{
							key: testData{
								id: 11,
							},
							value: 11,
							next: &node[testData, int]{
								key: testData{
									id: 21,
								},
								value: 21,
							},
						},
					},
				}
			},
			isFound: true,
			key: testData{
				id: 1,
			},
			wantVal: 1,
			wantHashMap: func() map[uint64]*node[testData, int] {
				return map[uint64]*node[testData, int]{
					1: &node[testData, int]{
						key: testData{
							id: 11,
						},
						value: 11,
						next: &node[testData, int]{
							key: testData{
								id: 21,
							},
							value: 21,
						},
					},
				}
			},
		},
		{
			name: "delete only one link element",
			key: testData{
				id: 1,
			},
			setHashMap: func() map[uint64]*node[testData, int] {
				return map[uint64]*node[testData, int]{
					1: &node[testData, int]{
						key: testData{
							id: 1,
						},
						value: 1,
					},
				}
			},
			wantHashMap: func() map[uint64]*node[testData, int] {
				return map[uint64]*node[testData, int]{}
			},
			wantVal: 1,
			isFound: true,
		},
		{
			name: "many link elements delete middle",
			key: testData{
				id: 11,
			},
			setHashMap: func() map[uint64]*node[testData, int] {
				return map[uint64]*node[testData, int]{
					1: &node[testData, int]{
						key: testData{
							id: 1,
						},
						value: 1,
						next: &node[testData, int]{
							key: testData{
								id: 11,
							},
							value: 11,
							next: &node[testData, int]{
								key: testData{
									id: 21,
								},
								value: 21,
							},
						},
					},
				}
			},
			wantHashMap: func() map[uint64]*node[testData, int] {
				return map[uint64]*node[testData, int]{
					1: &node[testData, int]{
						key: testData{
							id: 1,
						},
						value: 1,
						next: &node[testData, int]{
							key: testData{
								id: 21,
							},
							value: 21,
						},
					},
				}
			},
			isFound: true,
			wantVal: 11,
		},
		{
			name: "many link elements delete end",
			key: testData{
				id: 21,
			},
			setHashMap: func() map[uint64]*node[testData, int] {
				return map[uint64]*node[testData, int]{
					1: &node[testData, int]{
						key: testData{
							id: 1,
						},
						value: 1,
						next: &node[testData, int]{
							key: testData{
								id: 11,
							},
							value: 11,
							next: &node[testData, int]{
								key: testData{
									id: 21,
								},
								value: 21,
							},
						},
					},
				}
			},
			wantHashMap: func() map[uint64]*node[testData, int] {
				return map[uint64]*node[testData, int]{
					1: &node[testData, int]{
						key: testData{
							id: 1,
						},
						value: 1,
						next: &node[testData, int]{
							key: testData{
								id: 11,
							},
							value: 11,
						},
					},
				}
			},
			isFound: true,
			wantVal: 21,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHashMap[testData, int](10)
			h.hashmap = tc.setHashMap()
			val, ok := h.Delete(tc.key)
			assert.Equal(t, tc.isFound, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.wantVal, val)
			assert.Equal(t, tc.wantHashMap(), h.hashmap)
		})
	}
}

type testData struct {
	id int
}

func (t testData) Code() uint64 {
	hash := t.id % 10
	return uint64(hash)
}

func (t testData) Equals(key any) bool {
	val, ok := key.(testData)
	if !ok {
		return false
	}
	if t.id != val.id {
		return false
	}
	return true
}

func newTestData(id int) testData {
	return testData{
		id: id,
	}
}

type hashInt uint64

func (h hashInt) Code() uint64 {
	return uint64(h)
}

func (h hashInt) Equals(key any) bool {
	switch keyVal := key.(type) {
	case hashInt:
		return keyVal == h
	default:
		return false
	}
}

func newHashInt(i int) hashInt {
	return hashInt(i)
}

// goos: linux
// goarch: amd64
// pkg: github.com/gotomicro/ekit/mapx
// cpu: Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz
// BenchmarkMyHashMap/hashmap_put-8                 4985634               374.1 ns/op            53 B/op          1 allocs/op
// BenchmarkMyHashMap/map_put-8                     5465565               235.5 ns/op            49 B/op          0 allocs/op
// BenchmarkMyHashMap/hashmap_get-8                 7080890               143.9 ns/op             5 B/op          0 allocs/op
// BenchmarkMyHashMap/map_get-8                    18534306                86.94 ns/op            0 B/op          0 allocs/op

func BenchmarkMyHashMap(b *testing.B) {
	hashmap := NewHashMap[hashInt, int](10)
	m := make(map[uint64]int, 10)
	b.Run("hashmap_put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = hashmap.Put(newHashInt(i), i)
		}
	})
	b.Run("map_put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m[uint64(i)] = i
		}
	})
	b.Run("hashmap_get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = hashmap.Get(newHashInt(i))
		}
	})

	b.Run("map_get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m[uint64(i)]
		}
	})

}
