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

package reflectx

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIsNilValue(t *testing.T) {
	type MyInterface interface{}
	// 空接口，用于测试
	var nilMI MyInterface
	// 非空接口，用于测试
	n := 1
	myInterface := MyInterface(n)

	// 空 channel
	var nilCh chan int
	// 非空 channel
	ch := make(chan int)

	// 空 pointer
	var nilPtr *int
	// 非空指针
	ptr := &n

	// 空 map
	var nilMp map[string]struct{}
	// 非空 map
	mp := make(map[string]struct{}, 1)

	// 空 slice
	var nilSlice []int
	slc := make([]int, 1)

	testCases := []struct {
		name string
		val  reflect.Value
		res  bool
	}{
		{
			name: "int 类型",
			val:  reflect.ValueOf(666),
			res:  false,
		},
		{
			name: "string 类型",
			val:  reflect.ValueOf("字符串类型"),
			res:  false,
		},
		{
			name: "bool 类型",
			val:  reflect.ValueOf(true),
			res:  false,
		},
		{
			name: "float 类型",
			val:  reflect.ValueOf(3.14),
			res:  false,
		},
		{
			name: "complex 类型",
			val:  reflect.ValueOf(complex(1, 1)),
			res:  false,
		},
		{
			name: "struct 类型",
			val:  reflect.ValueOf(struct{}{}),
			res:  false,
		},
		{
			name: "array 类型",
			val:  reflect.ValueOf([4]int{}),
			res:  false,
		},
		{
			name: "nil 非法值",
			val:  reflect.ValueOf(nil),
			res:  true,
		},
		{
			name: "interface 类型 - 非空",
			val:  reflect.ValueOf(myInterface),
			res:  false,
		},
		{
			name: "interface 类型 - 空",
			val:  reflect.ValueOf(nilMI),
			res:  true,
		},
		{
			name: "pointer 类型 - 非空",
			val:  reflect.ValueOf(ptr),
			res:  false,
		},
		{
			name: "pointer 类型 - 空",
			val:  reflect.ValueOf(nilPtr),
			res:  true,
		},
		{
			name: "channel 类型 - 非空",
			val:  reflect.ValueOf(ch),
			res:  false,
		},
		{
			name: "channel 类型 - 空",
			val:  reflect.ValueOf(nilCh),
			res:  true,
		},
		{
			name: "map 类型 - 非空",
			val:  reflect.ValueOf(mp),
			res:  false,
		},
		{
			name: "map 类型 - 空",
			val:  reflect.ValueOf(nilMp),
			res:  true,
		},
		{
			name: "slice 类型 - 非空",
			val:  reflect.ValueOf(slc),
			res:  false,
		},
		{
			name: "slice 类型 - 空",
			val:  reflect.ValueOf(nilSlice),
			res:  true,
		},
		{
			name: "func 类型 - 非空",
			val: reflect.ValueOf(func() func() {
				return func() {}
			}),
			res: false,
		},
		{
			name: "func 类型 - 空",
			val: reflect.ValueOf(func() func() {
				type MyFunc func()
				var myFunc MyFunc
				return myFunc
			}),
			res: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := IsNilValue(tc.val)
			assert.Equal(t, tc.res, res)
		})
	}
}
