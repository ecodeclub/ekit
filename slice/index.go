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

// Index 返回和 dst 相等的第一个元素下标
// -1 表示没找到
func Index[T comparable](src []T, dst T) int {
	return IndexFunc[T](src, dst, func(x, y T) bool {
		return x == y
	})
}

// IndexFunc 返回和 dst 相等的第一个元素下标
// -1 表示没找到
// 你应该优先使用 Index
func IndexFunc[T any](src []T, dst T, equal EqualFunc[T]) int {
	for i, v := range src {
		isPanic, result := equal.safeEqual(v, dst)
		if isPanic {
			return -1
		}
		if result {
			return i
		}
	}
	return -1
}

// LastIndex 返回和 dst 相等的最后一个元素下标
// -1 表示没找到
func LastIndex[T comparable](src []T, dst T) int {
	return LastIndexFunc[T](src, dst, func(x, y T) bool {
		return x == y
	})
}

// LastIndexFunc 返回和 dst 相等的最后一个元素下标
// -1 表示没找到
// 你应该优先使用 LastIndex
func LastIndexFunc[T any](src []T, dst T, equal EqualFunc[T]) int {
	for i := len(src) - 1; i >= 0; i-- {
		isPanic, result := equal.safeEqual(src[i], dst)
		if isPanic {
			return -1
		}
		if result {
			return i
		}
	}
	return -1
}

// IndexAll 返回和 dst 相等的所有元素的下标
func IndexAll[T comparable](src []T, dst T) []int {
	return IndexAllFunc[T](src, dst, func(x, y T) bool {
		return x == y
	})
}

// IndexAllFunc 返回和 dst 相等的所有元素的下标
// 你应该优先使用 IndexAll
func IndexAllFunc[T any](src []T, dst T, equal EqualFunc[T]) []int {
	var ans []int
	for i, v := range src {
		isPanic, result := equal.safeEqual(v, dst)
		if isPanic {
			return ans
		}
		if result {
			ans = append(ans, i)
		}
	}
	return ans
}
