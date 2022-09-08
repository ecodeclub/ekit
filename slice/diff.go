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

// Diff 差集，只支持 comparable 类型
func Diff[T comparable](src, dst []T) []T {
	srcMap, dstMap := setMapStruct[T](src), setMapStruct[T](dst)
	for key := range dstMap {
		if _, exist := srcMap[key]; exist {
			delete(srcMap, key)
		}
	}
	
	var ret = make([]T, 0, len(srcMap))
	for key := range srcMap {
		ret = append(ret, key)
	}
	
	return ret
}

// DiffFunc 差集
// 你应该优先使用 Diff
func DiffFunc[T any](src, dst []T, equal EqualFunc[T]) []T {
	var ret = make([]T, 0, len(src))
	for _, val := range src {
		if !ContainsFunc[T](dst, val, equal) {
			ret = append(ret, val)
		}
	}
	return removeExistFunc[T](ret, equal)
}
