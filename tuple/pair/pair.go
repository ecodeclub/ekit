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

package pair

import (
	"fmt"
)

type Pair[K any, V any] struct {
	Key   K
	Value V
}

func (pair *Pair[K, V]) String() string {
	return fmt.Sprintf("<%#v, %#v>", pair.Key, pair.Value)
}

// Split 方法将Key, Value作为返回参数传出。
func (pair *Pair[K, V]) Split() (K, V) {
	return pair.Key, pair.Value
}

func NewPair[K any, V any](
	key K,
	value V,
) Pair[K, V] {
	return Pair[K, V]{
		Key:   key,
		Value: value,
	}
}

// NewPairs 需要传入两个长度相同并且均不为nil的数组 keys 和 values，
// 设keys长度为n，返回一个长度为n的pair数组。
// 保证：
//
//	返回的pair数组满足条件（设pair数组为p）:
//		对于所有的 0 <= i < n
//		p[i].Key == keys[i] 并且 p[i].Value == values[i]
//
//	如果传入的keys或者values为nil，会返回error
//
//	如果传入的keys长度与values长度不同，会返回error
func NewPairs[K any, V any](
	keys []K,
	values []V,
) ([]Pair[K, V], error) {
	if keys == nil || values == nil {
		return nil, fmt.Errorf("keys与values均不可为nil")
	}
	n := len(keys)
	if n != len(values) {
		return nil, fmt.Errorf("keys与values的长度不同, len(keys)=%d, len(values)=%d", n, len(values))
	}
	pairs := make([]Pair[K, V], n)
	for i := 0; i < n; i++ {
		pairs[i] = NewPair(keys[i], values[i])
	}
	return pairs, nil
}

// SplitPairs 需要传入一个[]Pair[K, V]，数组可以为nil。
// 设pairs数组的长度为n，返回两个长度均为n的数组keys, values。
// 如果pairs数组是nil, 则返回的keys与values也均为nil。
func SplitPairs[K any, V any](pairs []Pair[K, V]) (keys []K, values []V) {
	if pairs == nil {
		return nil, nil
	}
	n := len(pairs)
	keys = make([]K, n)
	values = make([]V, n)
	for i, pair := range pairs {
		keys[i], values[i] = pair.Split()
	}
	return
}

// FlattenPairs 需要传入一个[]Pair[K, V]，数组可以为nil
// 如果pairs数组为nil，则返回的flatPairs数组也为nil
//
//	设pairs数组长度为n，保证返回的flatPairs数组长度为2 * n且满足:
//		对于所有的 0 <= i < n
//		flatPairs[i * 2] == pairs[i].Key
//		flatPairs[i * 2 + 1] == pairs[i].Value
func FlattenPairs[K any, V any](pairs []Pair[K, V]) (flatPairs []any) {
	if pairs == nil {
		return nil
	}
	n := len(pairs)
	flatPairs = make([]any, 0, n*2)
	for _, pair := range pairs {
		flatPairs = append(flatPairs, pair.Key, pair.Value)
	}
	return
}

// PackPairs 需要传入一个长度为2 * n的数组flatPairs，数组可以为nil。
//
//	函数将会返回一个长度为n的pairs数组，pairs满足
//		对于所有的 0 <= i < n
//		pairs[i].Key == flatPairs[i * 2]
//		pairs[i].Value == flatPairs[i * 2 + 1]
//	如果flatPairs为nil,则返回的pairs也为nil
//
//	入参flatPairs需要满足以下条件：
//		对于所有的 0 <= i < n
//		flatPairs[i * 2] 的类型为 K
//		flatPairs[i * 2 + 1] 的类型为 V
//	否则会panic
func PackPairs[K any, V any](flatPairs []any) (pairs []Pair[K, V]) {
	if flatPairs == nil {
		return nil
	}
	n := len(flatPairs) / 2
	pairs = make([]Pair[K, V], n)
	for i := 0; i < n; i++ {
		pairs[i] = NewPair(flatPairs[i*2].(K), flatPairs[i*2+1].(V))
	}
	return
}
