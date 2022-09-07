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
	type args struct {
		src []int
		dst int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "切片为nil",
			args: args{nil, 1},
			want: false,
		},
		{
			name: "切片为空",
			args: args{[]int{}, 1},
			want: false,
		},
		{
			name: "包含测试",
			args: args{[]int{2, 3, 4}, 2},
			want: true,
		},
		{
			name: "不包含测试",
			args: args{[]int{1, 2, 3, 4}, 5},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Contains(tt.args.src, tt.args.dst), "Contains(%v, %v)", tt.args.src, tt.args.dst)
		})
	}
}

func TestContainsFunc(t *testing.T) {
	f := func(a, b int) bool {
		return a == b
	}

	type args struct {
		src   []int
		dst   int
		equal EqualFunc[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "切片为nil",
			args: args{nil, 1, f},
			want: false,
		},
		{
			name: "切片为空",
			args: args{[]int{}, 1, f},
			want: false,
		},
		{
			name: "包含测试",
			args: args{[]int{1, 2, 3, 4}, 1, f},
			want: true,
		},
		{
			name: "不包含测试",
			args: args{[]int{1, 2, 3, 4}, 5, f},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainsFunc(tt.args.src, tt.args.dst, tt.args.equal), "ContainsFunc(%v, %v, %v)", tt.args.src, tt.args.dst, tt.args.equal)
		})
	}
}

func TestContainsAny(t *testing.T) {
	type args struct {
		src []int
		dst []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"dst和src为nil", args{nil, nil}, false},
		{"存在交集", args{[]int{1, 2, 3}, []int{5, 8, 2}}, true},
		{"不存在交集", args{[]int{1, 2, 3}, []int{5, 8, 23}}, false},
		{"dst为nil", args{[]int{1, 2, 3}, nil}, false},
		{"src为nil", args{nil, []int{1, 2, 3}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainsAny(tt.args.src, tt.args.dst), "ContainsAny(%v, %v)", tt.args.src, tt.args.dst)
		})
	}
}

func TestContainsAnyFunc(t *testing.T) {
	f := func(a, b int) bool {
		return a == b
	}
	type args struct {
		src   []int
		dst   []int
		equal EqualFunc[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"dst和src为nil", args{nil, nil, f}, false},
		{"存在交集", args{[]int{1, 2, 3}, []int{5, 8, 2}, f}, true},
		{"不存在交集", args{[]int{1, 2, 3}, []int{5, 8, 23}, f}, false},
		{"dst为nil", args{[]int{1, 2, 3}, nil, f}, false},
		{"src为nil", args{nil, []int{1, 2, 3}, f}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainsAnyFunc(tt.args.src, tt.args.dst, tt.args.equal), "ContainsAnyFunc(%v, %v, %v)", tt.args.src, tt.args.dst, tt.args.equal)
		})
	}
}

func TestContainsAll(t *testing.T) {
	type args struct {
		src []int
		dst []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "dst和src为nil", args: args{nil, nil}, want: false},
		{name: "dst为nil", args: args{[]int{1, 2, 3}, nil}, want: true},
		{name: "src为{}和dst为nil", args: args{[]int{}, nil}, want: false},
		{name: "dst 是src 子集", args: args{[]int{5, 7, 8, 1, 2, 3, 4}, []int{1, 2, 3, 4}}, want: true},
		{name: "dst 不是src 子集", args: args{[]int{5, 7, 8, 1, 2, 3, 4}, []int{1, 2, 3, 4, 32}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainsAll(tt.args.src, tt.args.dst), "ContainsAll(%v, %v)", tt.args.src, tt.args.dst)
		})
	}
}

func TestContainsAllFunc(t *testing.T) {
	f := func(a, b int) bool {
		return a == b
	}
	type args struct {
		src   []int
		dst   []int
		equal EqualFunc[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "dst和src为nil", args: args{nil, nil, f}, want: false},
		{name: "dst为nil", args: args{[]int{1, 2, 3}, nil, f}, want: true},
		{name: "src为{}和dst为nil", args: args{[]int{}, nil, f}, want: false},
		{name: "dst 是src 子集", args: args{[]int{5, 7, 8, 1, 2, 3, 4}, []int{1, 2, 3, 4}, f}, want: true},
		{name: "dst 不是src 子集", args: args{[]int{5, 7, 8, 1, 2, 3, 4}, []int{1, 2, 3, 4, 32}, f}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainsAllFunc(tt.args.src, tt.args.dst, tt.args.equal), "ContainsAllFunc(%v, %v, %v)", tt.args.src, tt.args.dst, tt.args.equal)
		})
	}
}
func ExampleContains() {
	contains := Contains([]string{`php`, `python`, `golang`, `java`}, `java`)
	fmt.Println(contains)
	// Output: true
}
func ExampleContainsFunc() {
	contains := ContainsFunc(
		[]string{`php`, `python`, `golang`, `java`},
		`java`,
		func(src, dst string) bool {
			return src == dst
		},
	)
	fmt.Println(contains)
	// Output: true
}
func ExampleContainsAny() {
	containsAny := ContainsAny([]string{`php`, `python`, `golang`, `java`}, []string{`C#`, `C++`, `golang`, `java`})
	fmt.Println(containsAny)
	// Output: true
}
func ExampleContainsAnyFunc() {
	containsAny := ContainsAnyFunc([]string{`php`, `python`, `golang`, `java`}, []string{`C#`, `C++`, `golang`, `java`}, func(src, dst string) bool {
		return src == dst
	})
	fmt.Println(containsAny)
	// Output: true
}

func ExampleContainsAll() {
	all := ContainsAll([]string{`php`, `python`, `golang`, `java`}, []string{`C#`, `C++`, `golang`, `java`})
	fmt.Println(all)
	// Output: false
}

func ExampleContainsAllFunc() {
	all := ContainsAllFunc([]string{`php`, `python`, `golang`, `java`}, []string{`C#`, `C++`, `golang`, `java`}, func(src, dst string) bool {
		return src == dst
	})
	fmt.Println(all)
	// Output: false
}
