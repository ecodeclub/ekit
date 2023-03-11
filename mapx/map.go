// Copyright 2021 ecodeclub
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

package mapx

// Keys 返回 map 里面的所有的 key。
// 需要注意：这些 key 的顺序是随机。
func Keys[K comparable, V any](m map[K]V) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

// Values 返回 map 里面的所有的 value。
// 需要注意：这些 value 的顺序是随机。
func Values[K comparable, V any](m map[K]V) []V {
	res := make([]V, 0, len(m))
	for k := range m {
		res = append(res, m[k])
	}
	return res
}

// KeysValues 返回 map 里面的所有的 key,value。
// 需要注意：这些 (key,value) 的顺序是随机,相对顺序是一致的。
func KeysValues[K comparable, V any](m map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))
	for k := range m {
		keys = append(keys, k)
		values = append(values, m[k])
	}
	return keys, values
}
