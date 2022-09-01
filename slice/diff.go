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
	if src == nil {
		return dst
	}
	if dst == nil {
		return src
	}

	diff := make([]T, 0)
	for i := 0; i < len(src); i++ {
		found := false
		for j := 0; j < len(dst); j++ {
			if src[i] == dst[j] {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, src[i])
		}
	}
	for i := 0; i < len(dst); i++ {
		found := false
		for j := 0; j < len(src); j++ {
			if dst[i] == src[j] {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, dst[i])
		}
	}
	return diff
}

// DiffFunc 差集
// 你应该优先使用 Diff
func DiffFunc[T any](src, dst []T, equal EqualFunc[T]) []T {
	if src == nil {
		return dst
	}
	if dst == nil {
		return src
	}

	diff := make([]T, 0)
	for i := 0; i < len(src); i++ {
		found := false
		for j := 0; j < len(dst); j++ {
			if equal(src[i], dst[j]) {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, src[i])
		}
	}
	for i := 0; i < len(dst); i++ {
		found := false
		for j := 0; j < len(src); j++ {
			if equal(dst[i], src[j]) {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, dst[i])
		}
	}
	return diff
}
