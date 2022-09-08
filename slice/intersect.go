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

// Intersect 取交集，只支持 comparable 类型
// 返回值永远不为 nil
func Intersect[T comparable](src []T, dst []T) []T {
	srcMap, dstMap := setMapStruct(src), setMapStruct(dst)
	var ret = make([]T, 0, len(src))
	// 交集小于等于两个集合中的任意一个
	for dstKey := range dstMap {
		if _, exist := srcMap[dstKey]; exist {
			ret = append(ret, dstKey)
		}
	}
	return removeExist[T](ret)
}

// IntersectByFunc 支持任意类型
// 你应该优先使用 Intersect
func IntersectByFunc[T any](src []T, dst []T, equal EqualFunc[T]) []T {
	// 双重循环检测
	var ret = make([]T, 0, len(src))
	for _, valSrc := range src {
		for _, valDst := range dst {
			if equal(valDst, valSrc) {
				ret = append(ret, valSrc)
				break
			}
		}
	}
	return removeExistFunc[T](ret, equal)
}
