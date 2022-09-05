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
	var result []T
	srcMap := map[T]struct{}{}
	for _, v := range src {
		srcMap[v] = struct{}{}
	}
	for _, v := range dst {
		if _, ok := srcMap[v]; ok {
			result = append(result, v)
		}
	}
	return result
}

// IntersectByFunc 支持任意类型
// 你应该优先使用 Intersect
func IntersectByFunc[T any](src []T, dst []T, equal EqualFunc[T]) []T {
	var result []T
	for _, sv := range src {
		for _, dv := range dst {
			isPanic, isEqual := equal.safeEqual(sv, dv)
			if isPanic {
				return nil
			}
			if isEqual {
				result = append(result, sv)
			}
		}
	}
	return result
}
