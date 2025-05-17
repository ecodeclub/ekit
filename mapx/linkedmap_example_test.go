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
	"fmt"
	"sort"
)

func ExampleLinkedMap_Iterate() {
	linkedMap := NewLinkedHashMap[testStringData, int](0)
	strArr := make([]string, 0)
	_ = linkedMap.Put(testStringData{data: "hello"}, 1)
	_ = linkedMap.Put(testStringData{data: "world"}, 2)
	_ = linkedMap.Put(testStringData{data: "ekit"}, 3)

	linkedMap.Iterate(
		func(key testStringData, val int) bool {
			strArr = append(strArr, key.data)
			return true
		})

	sort.Strings(strArr)
	for _, s := range strArr {
		fmt.Println(s)
	}

	// Output:
	// ekit
	// hello
	// world
}
