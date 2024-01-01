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
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ecodeclub/ekit/reflectx"
)

type Pair[FirstType any, SecondType any] struct {
	first  FirstType
	second SecondType
}

type pairJson[FirstType any, SecondType any] struct {
	First  FirstType  `json:"first"`
	Second SecondType `json:"second"`
}

func (pair *Pair[FirstType, SecondType]) First() FirstType {
	return pair.first
}

func (pair *Pair[FirstType, SecondType]) SetFirst(v FirstType) {
	pair.first = v
}

func (pair *Pair[FirstType, SecondType]) Second() SecondType {
	return pair.second
}

func (pair *Pair[FirstType, SecondType]) SetSecond(v SecondType) {
	pair.second = v
}

func (pair *Pair[FirstType, SecondType]) ToString() string {
	return fmt.Sprintf("<%#v, %#v>", pair.first, pair.second)
}

func (pair *Pair[FirstType, SecondType]) ToArray() []any {
	return []any{pair.first, pair.second}
}

func (pair *Pair[FirstType, SecondType]) ToJson() ([]byte, error) {
	return json.Marshal(pairJson[FirstType, SecondType]{
		First:  pair.first,
		Second: pair.second,
	})
}

func (pair *Pair[FirstType, SecondType]) FromJson(jsonByte []byte) (err error) {
	pairJsonObj := pairJson[FirstType, SecondType]{}
	if err = json.Unmarshal(jsonByte, &pairJsonObj); err != nil {
		return err
	}
	pair.first = pairJsonObj.First
	pair.second = pairJsonObj.Second
	return
}

// 使用other的内容覆盖原先的pair, ignoreNil == true会忽略other中为nil的值，否则的话会直接覆盖。
func (pair *Pair[FirstType, SecondType]) MergeFrom(
	other Pair[FirstType, SecondType],
	ignoreNil bool,
) {
	if ignoreNil {
		if !reflectx.IsNilValue(reflect.ValueOf(other.first)) {
			pair.first = other.first
		}
		if !reflectx.IsNilValue(reflect.ValueOf(other.second)) {
			pair.second = other.second
		}
	} else {
		pair.first = other.first
		pair.second = other.second
	}
}

// 复制出一个与原来完全一样的Pair
func (pair *Pair[FirstType, SecondType]) Copy() Pair[FirstType, SecondType] {
	return NewPair(pair.values())
}

// 内部函数
func (pair Pair[FirstType, SecondType]) values() (FirstType, SecondType) {
	return pair.first, pair.second
}

func NewPair[FirstType any, SecondType any](
	first FirstType,
	second SecondType,
) Pair[FirstType, SecondType] {
	return Pair[FirstType, SecondType]{
		first:  first,
		second: second,
	}
}

func NewEmptyPair[FirstType any, SecondType any]() Pair[FirstType, SecondType] {
	return Pair[FirstType, SecondType]{}
}
