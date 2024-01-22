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

package slice

import "github.com/ecodeclub/ekit/internal/slice"

// Add 在index处添加元素
// index 范围应为[0, len(src)]
// 如果index == len(src) 则表示往末尾添加元素
func Add[Src any](src []Src, element Src, index int) ([]Src, error) {
	res, err := slice.Add[Src](src, element, index)
	return res, err
}
