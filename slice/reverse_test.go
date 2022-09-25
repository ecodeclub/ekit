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

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverseInt(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			want: []int{7, 5, 3, 1},
			src:  []int{1, 3, 5, 7},
			name: "normal test",
		},
		{
			src:  []int{},
			want: []int{},
			name: "length of src is 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Reverse[int](tt.src)
			assert.ElementsMatch(t, tt.want, res)
		})
	}
}

func ExampleReverseInt() {
	res := Reverse[int]([]int{1, 3, 2, 2, 4})
	fmt.Println(res)
	// Output:
	// [4,2,2,3,1]
}

func ExampleReverseString() {
	res := Reverse[string]([]string{"a", "b", "c", "d", "e"})
	fmt.Println(res)
	// Output:
	// ["e","d","c","b","a"]
}
