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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  int
		want int
	}{
		{
			name: "src nil",
			src:  nil,
			want: -1,
		},
		{
			name: "src empty",
			src:  []int{},
			want: -1,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3, 1, 2, 3},
			dst:  2,
			want: 1,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  4,
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Index(tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestIndexFunc(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  int
		want int
	}{
		{
			name: "src nil",
			src:  nil,
			want: -1,
		},
		{
			name: "src empty",
			src:  []int{},
			want: -1,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3, 1, 2, 3},
			dst:  2,
			want: 1,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  4,
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IndexFunc(tt.src, tt.dst, func(src int, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestLastIndex(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  int
		want int
	}{
		{
			name: "src nil",
			src:  nil,
			want: -1,
		},
		{
			name: "src empty",
			src:  []int{},
			want: -1,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3, 1, 2, 3},
			dst:  2,
			want: 4,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  4,
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := LastIndex(tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestLastIndexFunc(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  int
		want int
	}{
		{
			name: "src nil",
			src:  nil,
			want: -1,
		},
		{
			name: "src empty",
			src:  []int{},
			want: -1,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3, 1, 2, 3},
			dst:  2,
			want: 4,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  4,
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := LastIndexFunc(tt.src, tt.dst, func(src int, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestIndexAll(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  int
		want []int
	}{
		{
			name: "src nil",
			src:  nil,
			want: nil,
		},
		{
			name: "src empty",
			src:  []int{},
			want: nil,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3, 1, 2, 3},
			dst:  2,
			want: []int{1, 4},
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  4,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IndexAll(tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestIndexAllFunc(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  int
		want []int
	}{
		{
			name: "src nil",
			src:  nil,
			want: nil,
		},
		{
			name: "src empty",
			src:  []int{},
			want: nil,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3, 1, 2, 3},
			dst:  2,
			want: []int{1, 4},
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  4,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IndexAllFunc(tt.src, tt.dst, func(src int, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}

func ExampleIndex() {
	src := []int{1, 2, 3, 1, 2, 3}
	dst := 2
	index := Index(src, dst)
	fmt.Println(index)
	//Output: 1
}

func ExampleIndexFunc() {
	src := []int{1, 2, 3, 1, 2, 3}
	dst := 2
	equal := func(src int, dst int) bool {
		return src == dst
	}
	index := IndexFunc(src, dst, equal)
	fmt.Println(index)
	//Output: 1
}

func ExampleLastIndex() {
	src := []int{1, 2, 3, 1, 2, 3}
	dst := 2
	index := LastIndex(src, dst)
	fmt.Println(index)
	//Output: 4
}

func ExampleLastIndexFunc() {
	src := []int{1, 2, 3, 1, 2, 3}
	dst := 2
	equal := func(src int, dst int) bool {
		return src == dst
	}
	index := LastIndexFunc(src, dst, equal)
	fmt.Println(index)
	//Output: 4
}

func ExampleIndexAll() {
	src := []int{1, 2, 3, 1, 2, 3}
	dst := 2
	index := IndexAll(src, dst)
	fmt.Println(index)
	//Output: [1 4]
}

func ExampleIndexAllFunc() {
	src := []int{1, 2, 3, 1, 2, 3}
	dst := 2
	equal := func(src int, dst int) bool {
		return src == dst
	}
	index := IndexAllFunc(src, dst, equal)
	fmt.Println(index)
	//Output: [1 4]
}
