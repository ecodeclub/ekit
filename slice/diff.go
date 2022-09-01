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
	var res []T
	mp := make(map[T]bool)
	for _, elem := range dst {
		mp[elem] = true
	}
	for _, elem := range src {
		if mp[elem] == false {
			res = append(res, elem)
		}
	}
	return res
}

// DiffFunc 差集
// 你应该优先使用 Diff
func DiffFunc[T any](src, dst []T, equal EqualFunc[T]) []T {
	// 双重循环检测
	return nil
}
