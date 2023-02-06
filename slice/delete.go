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

import "github.com/gotomicro/ekit/internal/slice"

// Delete 删除 index 处的元素
func Delete[Src any](src []Src, index int) ([]Src, error) {
	res, _, err := slice.Delete[Src](src, index)
	return res, err
}

// FilterDelete 删除符合条件的元素
// 考虑到性能问题，所有操作都会在原切片上进行
// 被删除元素之后的元素会往前移动，有且只会移动一次
func FilterDelete[Src any](src []Src, m func(idx int, src Src) bool) []Src {
	// 记录被删除的元素位置，也称空缺的位置
	emptyPos := -1
	for idx := range src {
		// 判断是否满足删除的条件
		if m(idx, src[idx]) {
			// 如果在此之前没有被删除的元素，记录为空缺位置，等待后续元素往前移填补
			if emptyPos == -1 {
				emptyPos = idx
			}
		} else if emptyPos != -1 { // 如果前面有空缺的位置（也就是有元素被删除了）
			// 交换（填补）元素
			src[emptyPos] = src[idx]
			// 空缺位置往后移
			emptyPos++
		}
	}
	// 如果一个元素都没有删除，返回原切片
	if emptyPos == -1 {
		return src
	}
	return src[:emptyPos]
}
