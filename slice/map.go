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

// Map 将一个切片转化为另外一个切片
func Map[Src any, Dst any](src []Src, m func(idx int, src Src) Dst) []Dst {
	return []Dst{}
}

// 构造map
func setMapStruct[T comparable](src []T) map[T]struct{} {
	var dataMap = make(map[T]struct{}, len(src))
	for _, v := range src {
		// 使用空结构体,减少内存消耗
		dataMap[v] = struct{}{}
	}
	return dataMap
}

func setMapIndex[T comparable](src []T) map[T][]int {
	var dataMap = make(map[T][]int, len(src))
	for k, v := range src {
		dataMap[v] = append(dataMap[v], k)
	}
	return dataMap
}
