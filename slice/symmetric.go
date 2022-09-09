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

func SymmetricDiff[T comparable](src, dst []T) []T {
	srcMap, dstMap := setMapStruct[T](src), setMapStruct[T](dst)
	for dstKey := range dstMap {
		if _, exist := srcMap[dstKey]; exist {
			// 删除共同元素,两者剩余的并集即为对称差
			delete(dstMap, dstKey)
			delete(srcMap, dstKey)
		}
	}

	var ret = make([]T, 0, len(srcMap)+len(dstMap))
	for k := range srcMap {
		ret = append(ret, k)
	}
	for k := range dstMap {
		ret = append(ret, k)
	}

	return ret
}

func SymmetricDiffFunc[T any](src, dst []T, equal EqualFunc[T]) []T {
	// 双重循环检测
	var interSection = make([]T, 0, min(len(src), len(dst)))
	for _, valSrc := range src {
		for _, valDst := range dst {
			if equal(valSrc, valDst) {
				interSection = append(interSection, valSrc)
				break
			}
		}
	}

	ret := make([]T, 0, len(src)+len(dst)-len(interSection)*2)
	for _, v := range src {
		if !ContainsFunc[T](interSection, v, equal) {
			ret = append(ret, v)
		}
	}
	for _, v := range dst {
		if !ContainsFunc[T](interSection, v, equal) {
			ret = append(ret, v)
		}
	}

	return ret
}

func min(src, dst int) int {
	if src > dst {
		return dst
	}
	return src
}
