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

// Union 并集，只支持 comparable
func UnionSet[T comparable](src, dst []T) []T {
	srcMap, dstMap := setMapStruct[T](src), setMapStruct[T](dst)
	for key := range srcMap {
		dstMap[key] = struct{}{}
	}

	var ret = make([]T, 0, len(dstMap))
	for key := range dstMap {
		ret = append(ret, key)
	}

	return ret
}

// UnionByFunc 并集，支持任意类型
// 你应该优先使用 Union
func UnionSetByFunc[T any](src, dst []T, equal EqualFunc[T]) []T {
	// 双重循环检测
	var ret = make([]T, 0, len(src)+len(dst))
	ret = append(ret, dst...)
	ret = append(ret, src...)

	return deduplicateFunc[T](ret, equal)
}
