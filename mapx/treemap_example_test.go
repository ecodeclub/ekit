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

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/mapx"
)

func ExampleNewTreeMap() {
	m, _ := mapx.NewTreeMap[int, int](ekit.ComparatorRealNumber[int])
	_ = m.Put(1, 11)
	val, _ := m.Get(1)
	fmt.Println(val)
	// Output:
	// 11
}

func ExampleTreeMap_Iterate() {
	m, _ := mapx.NewTreeMap[int, int](ekit.ComparatorRealNumber[int])
	_ = m.Put(1, 11)
	_ = m.Put(-1, 12)
	_ = m.Put(100, 13)
	_ = m.Put(-100, 14)
	_ = m.Put(-101, 15)

	m.Iterate(func(key, value int) bool {
		if key > 1 {
			return false
		}
		fmt.Println(key, value)
		return true
	})

	// Output:
	// -101 15
	// -100 14
	// -1 12
	// 1 11
}
