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

// Contains 判断 src 里面是否存在 dst
func Contains[T comparable](src []T, dst T) bool {
	for _, elem := range src {
		if elem == dst {
			return true
		}
	}
	return false
}

// ContainsFunc 判断 src 里面是否存在 dst
// 你应该优先使用 Contains
func ContainsFunc[T any](src []T, dst T, equal EqualFunc[T]) bool {
	for _, elem := range src {
		if equal(elem, dst) {
			return true
		}
	}
	return false
}

// ContainsAny 判断 src 里面是否存在 dst 中的任何一个元素
func ContainsAny[T comparable](src, dst []T) bool {
	mp := make(map[T]bool)
	for _, elem := range src {
		mp[elem] = true
	}
	for _, elem := range dst {
		if mp[elem] {
			return true
		}
	}

	return false
}

// ContainsAnyFunc 判断 src 里面是否存在 dst 中的任何一个元素
// 你应该优先使用 ContainsAny
func ContainsAnyFunc[T any](src, dst []T, equal EqualFunc[T]) bool {
	for _, d := range dst {
		for _, s := range src {
			if equal(d, s) {
				return true
			}
		}
	}
	return false
}

// ContainsAll 判断 src 里面是否存在 dst 中的所有元素
func ContainsAll[T comparable](src, dst []T) bool {
	if src == nil {
		return false
	}
	mp := make(map[T]bool)
	for _, elem := range src {
		mp[elem] = true
	}
	for _, elem := range dst {
		if mp[elem] == false {
			return false
		}
	}
	return true
}

// ContainsAllFunc 判断 src 里面是否存在 dst 中的所有元素
// 你应该优先使用 ContainsAllFunc
func ContainsAllFunc[T any](src, dst []T, equal EqualFunc[T]) bool {
	if src == nil {
		return false
	}
	for _, elem := range dst {
		if !ContainsFunc(src, elem, equal) {
			return false
		}
	}
	return true
}
