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

package ekit

import "fmt"

type Triple struct {
	First  any
	Second any
	Third  any
}

// ToString 将 Triple 转为字符串，格式类似于 <key,value>
func (t Triple) ToString() string {
	return fmt.Sprint("<", t.First, ",", t.Second, ",", t.Third, ">")
}

// ToList 将 Triple 转为数组
func (t Triple) ToList() []any {
	return []any{t.First, t.Second, t.Third}
}

// Copy 传入一个 Triple 来修改对应位置的值
func (t Triple) Copy(toTriple Triple) (newTriple Triple) {
	newTriple = Triple{First: t.First, Second: t.Second}
	if toTriple.First != nil {
		newTriple.First = toTriple.First
	}
	if toTriple.Second != nil {
		newTriple.Second = toTriple.Second
	}
	if toTriple.Third != nil {
		newTriple.Third = toTriple.Third
	}
	return
}
