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
