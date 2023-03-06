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

package tuple

import (
	. "github.com/ecodeclub/ekit/tuple/pair"
	. "github.com/ecodeclub/ekit/tuple/triple"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test_Example 期望的使用方法
func Test_Example(t *testing.T) {
	pair := Pair{First: 1, Second: "one"}
	_ = pair.First                       // 1
	_ = pair.Second                      // one
	_ = pair.ToString()                  // <1,one>
	_ = pair.ToList()                    // [1,"one"]
	pair = pair.Copy(Pair{First: "two"}) // <"two","one">
	pair = pair.Copy(Pair{Second: 2})    // <"two",2>
	pair = pair.Copy(Pair{
		First:  2,
		Second: "two",
	}) // <2,"two">

	triple := Triple{
		First:  3,
		Second: "Three",
		Third:  true,
	}
	_ = triple.First                            // 3
	_ = triple.Second                           // three
	_ = triple.Third                            // true
	_ = triple.ToString()                       // <3,"three",true>
	_ = triple.ToList()                         // [3,"three", true]
	triple = triple.Copy(Triple{First: "four"}) // <"four","three",true>
	triple = triple.Copy(Triple{Second: 4})     // <"four",4,true>
	triple = triple.Copy(Triple{Third: false})  //<"four",4,false>
	triple = triple.Copy(Triple{
		First:  4,
		Second: "four",
		Third:  nil,
	}) // <4,"four",false>
}

// Test_getValue 获取元素
func Test_getValue(t *testing.T) {
	testCases := []struct {
		name string
		// 用于区分是 Pair(2) 还是 Triple(3)
		class      int
		inputValue Tuple
		wantValue1 any
		wantValue2 any
		wantValue3 any
		wantErr    error
	}{
		{
			name:       "Pair:nil",
			class:      2,
			inputValue: Pair{},
			wantValue1: nil,
			wantValue2: nil,
			wantValue3: nil,
			wantErr:    nil,
		},
		{
			name:       "Triple:nil",
			class:      3,
			inputValue: Triple{},
			wantValue1: nil,
			wantValue2: nil,
			wantValue3: nil,
			wantErr:    nil,
		},
		{
			name:  "Pair:bool + int",
			class: 2,
			inputValue: Pair{
				First:  true,
				Second: 1,
			},
			wantValue1: true,
			wantValue2: int(1),
			wantValue3: nil,
			wantErr:    nil,
		},
		{
			name:  "Triple:bool + int + uint",
			class: 3,
			inputValue: Triple{
				First:  false,
				Second: int(0),
				Third:  uint(2),
			},
			wantValue1: false,
			wantValue2: int(0),
			wantValue3: uint(2),
			wantErr:    nil,
		},
		{
			name:  "Pair:rune + byte",
			class: 2,
			inputValue: Pair{
				First:  rune('a'),
				Second: byte('b'),
			},
			wantValue1: rune('a'),
			wantValue2: byte('b'),
			wantValue3: nil,
			wantErr:    nil,
		},
		{
			name:  "Triple:rune + byte + string",
			class: 3,
			inputValue: Triple{
				First:  rune('c'),
				Second: byte('d'),
				Third:  string("efg"),
			},
			wantValue1: rune('c'),
			wantValue2: byte('d'),
			wantValue3: string("efg"),
			wantErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.class {
			case 2:
				assert.Equal(t, tc.wantValue1, tc.inputValue.(Pair).First)
				assert.Equal(t, tc.wantValue2, tc.inputValue.(Pair).Second)
			case 3:
				assert.Equal(t, tc.wantValue1, tc.inputValue.(Triple).First)
				assert.Equal(t, tc.wantValue2, tc.inputValue.(Triple).Second)
				assert.Equal(t, tc.wantValue3, tc.inputValue.(Triple).Third)
			default:
			}
		})
	}
}

// Test_changeValue 改变内部值，只能从 非nil 变成 nil，不能从 nil 变成 非nil
func Test_changeValue(t *testing.T) {
	runes := []rune{'a', 'b', 'c'}
	bytes := []byte{'d', 'e', 'f'}
	testCases := []struct {
		name string
		// 用于区分是 Pair(2) 还是 Triple(3)
		class       int
		inputValue  Tuple
		changeValue Tuple
		wantValue1  any
		wantValue2  any
		wantValue3  any
		wantErr     error
	}{
		{
			name:       "Pair:nil -> all filed",
			class:      2,
			inputValue: Pair{},
			changeValue: Pair{
				First:  0,
				Second: "zero",
			},
			wantValue1: 0,
			wantValue2: "zero",
			wantValue3: nil,
			wantErr:    nil,
		},
		{
			name:       "Triple:nil -> all filed",
			class:      3,
			inputValue: Triple{},
			changeValue: Triple{
				First:  3,
				Second: "three",
				Third:  true,
			},
			wantValue1: 3,
			wantValue2: "three",
			wantValue3: true,
			wantErr:    nil,
		},
		{
			name:  "Pair:bool + int -> the same type",
			class: 2,
			inputValue: Pair{
				First:  true,
				Second: 1,
			},
			changeValue: Pair{First: false},
			wantValue1:  false,
			wantValue2:  int(1),
			wantValue3:  nil,
			wantErr:     nil,
		},
		{
			name:  "Triple:bool + int + uint -> the same type",
			class: 3,
			inputValue: Triple{
				First:  false,
				Second: int(0),
				Third:  uint(2),
			},
			changeValue: Triple{
				Third: uint(20)},
			wantValue1: false,
			wantValue2: int(0),
			wantValue3: uint(20),
			wantErr:    nil,
		},
		{
			name:  "Pair:rune + byte -> different type",
			class: 2,
			inputValue: Pair{
				First:  rune('a'),
				Second: byte('b'),
			},
			changeValue: Pair{
				First:  float32(3.14),
				Second: complex64(314),
			},
			wantValue1: float32(3.14),
			wantValue2: complex64(314),
			wantValue3: nil,
			wantErr:    nil,
		},
		{
			name:  "Triple:rune + byte + string -> different type",
			class: 3,
			inputValue: Triple{
				First:  rune('c'),
				Second: byte('d'),
				Third:  string("efg"),
			},
			changeValue: Triple{
				First:  uintptr(0x9876543210),
				Second: runes,
				Third:  bytes,
			},
			wantValue1: uintptr(0x9876543210),
			wantValue2: runes,
			wantValue3: bytes,
			wantErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.class {
			case 2:
				output := tc.inputValue.(Pair).Copy(tc.changeValue.(Pair))
				assert.Equal(t, tc.wantValue1, output.First)
				assert.Equal(t, tc.wantValue2, output.Second)
			case 3:
				output := tc.inputValue.(Triple).Copy(tc.changeValue.(Triple))
				assert.Equal(t, tc.wantValue1, output.First)
				assert.Equal(t, tc.wantValue2, output.Second)
				assert.Equal(t, tc.wantValue3, output.Third)
			default:
			}
		})
	}
}

// Test_TransInChan 尝试通过一组 chan 来传输
func Test_TransInChan(t *testing.T) {
	ch := make(chan Tuple, 1)
	testCases := []struct {
		name    string
		input   Tuple
		want    Tuple
		wantErr error
	}{
		{
			name:    "pair + nil",
			input:   Pair{},
			want:    Pair{},
			wantErr: nil,
		},
		{
			name:    "triple + nil",
			input:   Triple{},
			want:    Triple{},
			wantErr: nil,
		},
		{
			name: "pair",
			input: Pair{
				First:  0,
				Second: "zero",
			},
			want: Pair{
				First:  0,
				Second: "zero",
			},
			wantErr: nil,
		},
		{
			name: "triple",
			input: Triple{
				First:  1,
				Second: "one",
				Third:  "first",
			},
			want: Triple{
				First:  1,
				Second: "one",
				Third:  "first",
			},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			go func() {
				ch <- tc.input
			}()
			v := <-ch
			assert.Equal(t, tc.want, v)
		})
	}
}
