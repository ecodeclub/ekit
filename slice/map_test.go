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

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		want []string
	}{
		{
			name: "src nil",
			want: []string{},
		},
		{
			name: "src empty",
			src:  []int{},
			want: []string{},
		},
		{
			name: "src has element",
			src:  []int{1, 2, 3},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Map(tt.src, func(idx int, src int) string {
				return strconv.Itoa(src)
			})
			assert.Equal(t, res, tt.want)
		})
	}
}

func ExampleMap() {
	src := []int{1, 2, 3}
	dst := Map(src, func(idx int, src int) string {
		return strconv.Itoa(src)
	})
	fmt.Println(dst)
	// Output: [1 2 3]
}

func TestFilterMap(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		want []string
	}{
		{
			name: "src nil",
			want: []string{},
		},
		{
			name: "src empty",
			src:  []int{},
			want: []string{},
		},
		{
			name: "src has element",
			src:  []int{1, -2, 3},
			want: []string{"1", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := FilterMap(tt.src, func(idx int, src int) (string, bool) {
				return strconv.Itoa(src), src >= 0
			})
			assert.Equal(t, res, tt.want)
		})
	}
}

func ExampleFilterMap() {
	src := []int{1, -2, 3}
	dst := FilterMap[int, string](src, func(idx int, src int) (string, bool) {
		return strconv.Itoa(src), src >= 0
	})
	fmt.Println(dst)
	// Output: [1 3]
}
