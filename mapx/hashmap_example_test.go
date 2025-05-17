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
	"hash/crc32"
	"sort"
)

type testStringData struct {
	data string
}

func (ts testStringData) Code() uint64 {
	return uint64(crc32.ChecksumIEEE([]byte(ts.data)))
}

func (ts testStringData) Equals(other any) bool {
	otherv, ok := other.(testStringData)
	if !ok {
		return false
	}
	return ts.data == otherv.data
}

func ExampleHashMap_Iterate() {
	hashMap := NewHashMap[testStringData, int](0)
	strArr := make([]string, 0)
	_ = hashMap.Put(testStringData{data: "hello"}, 1)
	_ = hashMap.Put(testStringData{data: "world"}, 2)
	_ = hashMap.Put(testStringData{data: "ekit"}, 3)

	hashMap.Iterate(
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
