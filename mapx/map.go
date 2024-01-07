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

import "fmt"

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

// ToMap 将会返回一个map[K]V
// 请保证传入的 keys 与 values 长度相同，长度均为n
// 长度不相同或者 keys 或者 values 为nil则会抛出异常
// 返回的 m map[K]V 保证对于所有的 0 <= i < n
// m[keys[i]] = values[i]
//
//	注意：
//	如果传入的数组中存在 0 <= i < j < n使得 keys[i] == keys[j]
//	则在返回的 m 中 m[keys[i]] = values[j]
//	如果keys和values的长度为0，则会返回一个空map
func ToMap[K comparable, V any](keys []K, values []V) (m map[K]V, err error) {
	if keys == nil || values == nil {
		return nil, fmt.Errorf("keys与values均不可为nil")
	}
	n := len(keys)
	if n != len(values) {
		return nil, fmt.Errorf("keys与values的长度不同, len(keys)=%d, len(values)=%d", n, len(values))
	}
	m = make(map[K]V, n)
	for i := 0; i < n; i++ {
		m[keys[i]] = values[i]
	}
	return
}
