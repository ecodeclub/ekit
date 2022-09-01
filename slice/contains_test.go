package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
			name: "包含测试",
			args: args{[]int{1, 2, 3, 4}, 1},
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
	//简单用== 实现equal函数
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainsAny(tt.args.src, tt.args.dst), "ContainsAny(%v, %v)", tt.args.src, tt.args.dst)
		})
	}
}

func TestContainsAnyFunc(t *testing.T) {
	//简单用== 实现equal函数
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
	//简单用== 实现equal函数
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
		{name: "dst 是src 子集", args: args{[]int{5, 7, 8, 1, 2, 3, 4}, []int{1, 2, 3, 4}, f}, want: true},
		{name: "dst 不是src 子集", args: args{[]int{5, 7, 8, 1, 2, 3, 4}, []int{1, 2, 3, 4, 32}, f}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainsAllFunc(tt.args.src, tt.args.dst, tt.args.equal), "ContainsAllFunc(%v, %v, %v)", tt.args.src, tt.args.dst, tt.args.equal)
		})
	}
}
