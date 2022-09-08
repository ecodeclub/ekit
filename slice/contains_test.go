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

func TestContains(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  int
		want bool
	}{
		{
			name: "src nil",
			src:  nil,
			want: false,
		},
		{
			name: "src empty",
			src:  []int{},
			want: false,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3},
			dst:  3,
			want: true,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  4,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Contains(tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestContainsFunc(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  int
		want bool
	}{
		{
			name: "src nil",
			src:  nil,
			want: false,
		},
		{
			name: "src empty",
			src:  []int{},
			want: false,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3},
			dst:  3,
			want: true,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  4,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ContainsFunc(tt.src, tt.dst, func(src int, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestContainsAny(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want bool
	}{
		{
			name: "src and dst nil",
			src:  nil,
			dst:  nil,
			want: false,
		},
		{
			name: "only src nil",
			src:  nil,
			dst:  []int{1, 2, 3},
			want: false,
		},
		{
			name: "only dst nil",
			src:  []int{1, 2, 3},
			dst:  nil,
			want: false,
		},
		{
			name: "src and dst empty",
			src:  []int{},
			dst:  []int{},
			want: false,
		},
		{
			name: "src empty",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: false,
		},
		{
			name: "dst empty",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: false,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3},
			dst:  []int{1, 2, 3},
			want: true,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ContainsAny(tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestContainsAnyFunc(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want bool
	}{
		{
			name: "src and dst nil",
			src:  nil,
			dst:  nil,
			want: false,
		},
		{
			name: "only src nil",
			src:  nil,
			dst:  []int{1, 2, 3},
			want: false,
		},
		{
			name: "only dst nil",
			src:  []int{1, 2, 3},
			dst:  nil,
			want: false,
		},
		{
			name: "src and dst empty",
			src:  []int{},
			dst:  []int{},
			want: false,
		},
		{
			name: "src empty",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: false,
		},
		{
			name: "dst empty",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: false,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3},
			dst:  []int{1, 2, 3},
			want: true,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ContainsAnyFunc(tt.src, tt.dst, func(src int, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestContainsAll(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want bool
	}{
		{
			name: "src and dst nil",
			src:  nil,
			dst:  nil,
			want: true,
		},
		{
			name: "only src nil",
			src:  nil,
			dst:  []int{1, 2, 3},
			want: false,
		},
		{
			name: "only dst nil",
			src:  []int{1, 2, 3},
			dst:  nil,
			want: true,
		},
		{
			name: "src and dst empty",
			src:  []int{},
			dst:  []int{},
			want: true,
		},
		{
			name: "src empty",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: false,
		},
		{
			name: "dst empty",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: true,
		},
		{
			name: "src contains all dst",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{1, 2, 3},
			want: true,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			want: false,
		},
		{
			name: "src not contains all dst",
			src:  []int{1, 2, 3},
			dst:  []int{2, 3, 4},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ContainsAll(tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestContainsAllFunc(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		dst  []int
		want bool
	}{
		{
			name: "src and dst nil",
			src:  nil,
			dst:  nil,
			want: true,
		},
		{
			name: "only src nil",
			src:  nil,
			dst:  []int{1, 2, 3},
			want: false,
		},
		{
			name: "only dst nil",
			src:  []int{1, 2, 3},
			dst:  nil,
			want: true,
		},
		{
			name: "src and dst empty",
			src:  []int{},
			dst:  []int{},
			want: true,
		},
		{
			name: "src empty",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: false,
		},
		{
			name: "dst empty",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: true,
		},
		{
			name: "src contains all dst",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{1, 2, 3},
			want: true,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			want: false,
		},
		{
			name: "src not contains all dst",
			src:  []int{1, 2, 3},
			dst:  []int{2, 3, 4},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ContainsAllFunc(tt.src, tt.dst, func(src int, dst int) bool {
				return src == dst
			})
			assert.Equal(t, tt.want, res)
		})
	}
}

func ExampleContains() {
	src := []int{1, 2, 3}
	dst := 1
	contains := Contains(src, dst)
	fmt.Println(contains)
	//Output: true
}

func ExampleContainsFunc() {
	src := []int{1, 2, 3}
	dst := 1
	equal := func(src int, dst int) bool {
		return src == dst
	}
	contains := ContainsFunc(src, dst, equal)
	fmt.Println(contains)
	//Output: true
}

func ExampleContainsAny() {
	src := []int{1, 2, 3}
	dst := []int{3, 4}
	contains := ContainsAny(src, dst)
	fmt.Println(contains)
	//Output: true
}

func ExampleContainsAnyFunc() {
	src := []int{1, 2, 3}
	dst := []int{3, 4}
	equal := func(src int, dst int) bool {
		return src == dst
	}
	contains := ContainsAnyFunc(src, dst, equal)
	fmt.Println(contains)
	//Output: true
}

func ExampleContainsAll() {
	src := []int{1, 2, 3}
	dst := []int{2, 3}
	contains := ContainsAll(src, dst)
	fmt.Println(contains)
	//Output: true
}

func ExampleContainsAllFunc() {
	src := []int{1, 2, 3}
	dst := []int{2, 3}
	equal := func(src int, dst int) bool {
		return src == dst
	}
	contains := ContainsAllFunc(src, dst, equal)
	fmt.Println(contains)
	//Output: true
}
