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

func calCapacity(c, l int) int {
	if c <= 64 {
		return c
	}
	if c > 2048 && (c/l >= 2) {
		factor := 0.625
		return int(float32(c) * float32(factor))
	}
	if c <= 2048 && (c/l >= 4) {
		return c / 2
	}
	return c
}

func Shrink[T any](src []T) []T {
	c, l := cap(src), len(src)
	n := calCapacity(c, l)
	if n == c {
		return src
	}
	s := make([]T, 0, n)
	s = append(s, src...)
	return s
}
