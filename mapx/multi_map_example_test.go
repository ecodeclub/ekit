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

	"github.com/ecodeclub/ekit/tuple/pair"
)

func ExampleMultiMap_Iterate() {
	multiMap := NewMultiHashMap[testStringData, int](0)
	arr := make([]pair.Pair[string, int], 0)
	_ = multiMap.Put(testStringData{data: "hello"}, 1)
	_ = multiMap.Put(testStringData{data: "world"}, 2)
	_ = multiMap.Put(testStringData{data: "world"}, 3)
	_ = multiMap.Put(testStringData{data: "world"}, 4)
	_ = multiMap.Put(testStringData{data: "ekit"}, 3)

	multiMap.Iterate(
		func(key testStringData, val int) bool {
			arr = append(arr, pair.NewPair(key.data, val))
			return true
		})

	sort.Slice(arr, func(i, j int) bool {
		if arr[i].Key == arr[j].Key {
			return arr[i].Value < arr[j].Value
		}
		return arr[i].Key < arr[j].Key
	})

	for _, pa := range arr {
		fmt.Println(pa.Key, pa.Value)
	}
	// Output:
	// ekit 3
	// hello 1
	// world 2
	// world 3
	// world 4
}
