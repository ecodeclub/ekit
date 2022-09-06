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
	for dstKey := range dstMap {
		if _, exist := srcMap[dstKey]; exist {
			// 删除共同元素,两者剩余的并集即为差集
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

// DiffFunc 差集
// 你应该优先使用 Diff
func DiffFunc[T any](src, dst []T, equal EqualFunc[T]) []T {
	// 双重循环检测
	var sameData = make([]T, 0, min(len(src), len(dst)))
	for i := 0; i < len(src); i++ {
		for j := 0; j < len(dst); j++ {
			// 保持sameData中的数据唯一
			if equal(src[i], dst[j]) && ContainsFunc(sameData, src[i], equal) {
				sameData = append(sameData, src[i])
				break
			}
		}
	}
	
	var ret = make([]T, 0)
	for i := 0; i < len(src); i++ {
		if !ContainsFunc(sameData, src[i], equal) {
			ret = append(ret, src[i])
		}
	}
	for i := 0; i < len(dst); i++ {
		if !ContainsFunc(sameData, dst[i], equal) {
			ret = append(ret, dst[i])
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
