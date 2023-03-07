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

import "fmt"

type Pair struct {
	First  any
	Second any
}

// ToString 将 Pair 转为字符串，格式类似于 <key,value>
func (p Pair) ToString() string {
	return fmt.Sprint("<", p.First, ",", p.Second, ">")
}

// ToList 将 Pair 转为数组
func (p Pair) ToList() []any {
	return []any{p.First, p.Second}
}

// Copy 传入一个 Pair 来修改对应位置的值
func (p Pair) Copy(toPair Pair) (newPair Pair) {
	newPair = Pair{First: p.First, Second: p.Second}
	if toPair.First != nil {
		newPair.First = toPair.First
	}
	if toPair.Second != nil {
		newPair.Second = toPair.Second
	}
	return
}
