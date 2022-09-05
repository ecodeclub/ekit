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

			want: -1,
		},
		{
			name: "src empty",
			src:  []int{},

			want: -1,
		},
		{
			name: "src contains dst",
			src:  []int{1, 2, 3},
			dst:  3,

			want: 2,
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  4,

			want: -1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := Index[int](test.src, test.dst)
			assert.Equal(t, test.want, res)
		})
	}
}

func TestIndexFunc(t *testing.T) {
	tests := []struct {
		name  string
		src   []any
		dst   any
		equal EqualFunc[any]

		want int
	}{
		{
			name: "src nil",

			want: -1,
		},
		{
			name: "src empty",
			src:  []any{},

			want: -1,
		},
		{
			name: "src contains dst",
			src:  []any{1, 2, 3},
			dst:  3,

			want: 2,
		},
		{
			name: "src not contains dst",
			src:  []any{1, 2, 3},
			dst:  4,

			want: -1,
		},
		{
			name: "equal panic",
			src:  []any{1, 2, 3},
			dst:  1,
			equal: func(x, y any) bool {
				panic("panic test")
			},

			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func(x, y any) bool {
				return x == y
			}
			if tt.equal != nil {
				f = tt.equal
			}
			res := IndexFunc(tt.src, tt.dst, f)
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

			want: -1,
		},
		{
			name: "src empty",
			src:  []int{},

			want: -1,
		},
		{
			name: "src contains dst",
			src:  []int{1, 3, 3},
			dst:  3,

			want: 2,
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
			res := LastIndex[int](tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestLastIndexFunc(t *testing.T) {
	tests := []struct {
		name  string
		src   []any
		dst   any
		equal EqualFunc[any]

		want int
	}{
		{
			name: "src nil",

			want: -1,
		},
		{
			name: "src empty",
			src:  []any{},

			want: -1,
		},
		{
			name: "src contains dst",
			src:  []any{1, 3, 3},
			dst:  3,

			want: 2,
		},
		{
			name: "src not contains dst",
			src:  []any{1, 2, 3},
			dst:  4,

			want: -1,
		},
		{
			name: "equal panic",
			src:  []any{1, 2, 3},
			dst:  1,
			equal: func(x, y any) bool {
				panic("panic test")
			},

			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func(x, y any) bool {
				return x == y
			}
			if tt.equal != nil {
				f = tt.equal
			}
			res := LastIndexFunc[any](tt.src, tt.dst, f)
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

			want: nil,
		},
		{
			name: "src empty",
			src:  []int{},

			want: nil,
		},
		{
			name: "src contains dst",
			src:  []int{1, 3, 3},
			dst:  3,

			want: []int{1, 2},
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
			res := IndexAll[int](tt.src, tt.dst)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestIndexAllFunc(t *testing.T) {
	f := func(x, y int) bool {
		return x == y
	}
	tests := []struct {
		name  string
		src   []int
		dst   int
		equal EqualFunc[any]

		want []int
	}{
		{
			name: "src nil",

			want: nil,
		},
		{
			name: "src empty",
			src:  []int{},

			want: nil,
		},
		{
			name: "src contains dst",
			src:  []int{1, 3, 3},
			dst:  3,

			want: []int{1, 2},
		},
		{
			name: "src not contains dst",
			src:  []int{1, 2, 3},
			dst:  4,

			want: nil,
		},
		{
			name: "equal panic",
			src:  []int{1, 2, 3},
			dst:  1,
			equal: func(x, y any) bool {
				panic("panic test")
			},

			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IndexAllFunc[int](tt.src, tt.dst, f)
			assert.Equal(t, tt.want, res)
		})
	}
}

func ExampleIndex() {
	src := []int{1, 2, 3}
	dst := 1
	index := Index(src, dst)
	fmt.Println(index)
	// Output: 0
}

func ExampleIndexFunc() {
	equal := func(x, y int) bool {
		return x == y
	}
	src := []int{1, 2, 3}
	dst := 2
	index := IndexFunc(src, dst, equal)
	fmt.Println(index)
	// Output: 1
}

func ExampleLastIndex() {
	src := []int{1, 1, 1}
	dst := 1
	index := LastIndex(src, dst)
	fmt.Println(index)
	// Output: 2
}

func ExampleLastIndexFunc() {
	equal := func(x, y int) bool {
		return x == y
	}
	src := []int{1, 2, 2}
	dst := 2
	index := LastIndexFunc(src, dst, equal)
	fmt.Println(index)
	//Output: 2
}

func ExampleIndexAll() {
	src := []int{1, 2, 3, 2, 2}
	dst := 2
	result := IndexAll(src, dst)
	fmt.Println(result)
	//	Output: [1 3 4]
}

func ExampleIndexAllFunc() {
	equal := func(x, y int) bool {
		return x == y
	}
	src := []int{1, 2, 3, 2, 2}
	dst := 2
	result := IndexAllFunc(src, dst, equal)
	fmt.Println(result)
	//	Output: [1 3 4]
}
