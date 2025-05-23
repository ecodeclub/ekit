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

package mapx_test

import (
	"fmt"

	"github.com/ecodeclub/ekit/mapx"
)

func ExampleNewHashMap() {
	m := mapx.NewHashMap[MockKey, int](10)
	_ = m.Put(MockKey{}, 123)
	val, _ := m.Get(MockKey{})
	fmt.Println(val)
	// Output:
	// 123
}

type MockKey struct {
	values []int
}

func (m MockKey) Code() uint64 {
	res := 3
	for _, v := range m.values {
		res += v * 7
	}
	return uint64(res)
}

func (m MockKey) Equals(key any) bool {
	k, ok := key.(MockKey)
	if !ok {
		return false
	}
	if len(k.values) != len(m.values) {
		return false
	}
	if k.values == nil && m.values != nil {
		return false
	}
	if k.values != nil && m.values == nil {
		return false
	}
	for i, v := range m.values {
		if v != k.values[i] {
			return false
		}
	}
	return true
}

func ExampleMerge() {
	m1 := map[int]int{1: 1, 2: 2, 3: 3}
	m2 := map[int]int{4: 4, 5: 5, 6: 6}
	got := mapx.Merge(m1, m2)
	fmt.Println(got)

	m3 := map[int]int{1: 1, 2: 2, 3: 3}
	m4 := map[int]int{1: 5, 2: 6, 3: 7}
	got = mapx.Merge(m3, m4)
	fmt.Println(got)

	var m map[int]int
	got = mapx.Merge(m)
	fmt.Println(got == nil) // 不会返回 nil map

	// Output:
	// map[1:1 2:2 3:3 4:4 5:5 6:6]
	// map[1:5 2:6 3:7]
	// false
}

func ExampleMergeFunc() {
	m1 := map[int]int{1: 1, 2: 2, 3: 3}
	m2 := map[int]int{1: 2, 2: 3, 3: 4}
	got := mapx.MergeFunc(func(val1, val2 int) int {
		return val1 + val2
	}, m1, m2)
	fmt.Println(got)

	// Output:
	// map[1:3 2:5 3:7]
}
