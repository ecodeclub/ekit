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
func Union[T comparable](src, dst []T) []T {
	if len(src) == 0 && len(dst) == 0 {
		return nil
	}
	result := src
	srcMap := map[T]struct{}{}
	for _, v := range src {
		srcMap[v] = struct{}{}
	}
	for _, v := range dst {
		if _, ok := srcMap[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}

// UnionByFunc 并集，支持任意类型
// 你应该优先使用 Union
func UnionByFunc[T any](src, dst []T, equal EqualFunc[T]) []T {
	if len(src) == 0 && len(dst) == 0 {
		return nil
	}
	result := src
	for _, dv := range dst {
		var contains bool
		for _, sv := range src {
			isPanic, isEqual := equal.safeEqual(sv, dv)
			if isPanic {
				return nil
			}
			if isEqual {
				contains = true
				break
			}
		}
		if !contains {
			result = append(result, dv)
		}
	}
	return result
}
