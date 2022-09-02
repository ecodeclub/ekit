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

// Map 将一个切片转化为另外一个切片,如果发生panic,或者src 为nil,会返回nil.
func Map[Src any, Dst any](src []Src, m func(idx int, src Src) Dst) []Dst {
	res := make([]Dst, 0, len(src))
	defer func() {
		if err := recover(); err != nil {
			res = nil
		}
	}()
	for index, elem := range src {
		res = append(res, m(index, elem))
	}
	return res
}
