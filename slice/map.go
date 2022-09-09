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

package slice

// 构造map
func setMapStruct[T comparable](src []T) map[T]struct{} {
	var dataMap = make(map[T]struct{}, len(src))
	for _, v := range src {
		// 使用空结构体,减少内存消耗
		dataMap[v] = struct{}{}
	}
	return dataMap
}

func setMapIndexes[T comparable](src []T) map[T][]int {
	var dataMap = make(map[T][]int, len(src))
	for k, v := range src {
		dataMap[v] = append(dataMap[v], k)
	}
	return dataMap
}

func deduplicateFunc[T any](data []T, equal EqualFunc[T]) []T {
	var newData = make([]T, 0, len(data))
	for k, v := range data {
		if !ContainsFunc[T](data[k+1:], v, equal) {
			newData = append(newData, v)
		}
	}
	return newData
}

func deduplicateExist[T comparable](data []T) []T {
	dataMap := setMapStruct[T](data)
	var newData = make([]T, 0, len(dataMap))
	for key := range dataMap {
		newData = append(newData, key)
	}
	return newData
}
