package slice

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	type args struct {
		src []int
		m   func(idx int, src int) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "int转字符串",
			args: args{
				src: []int{1, 3, 4},
				m: func(idx int, src int) string {
					return strconv.Itoa(src)
				},
			},
			want: []string{`1`, `3`, `4`},
		},
		{
			name: "切片为nil",
			args: args{
				src: nil,
				m: func(idx int, src int) string {
					return strconv.Itoa(src)
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Map(tt.args.src, tt.args.m), "Map(%v, %v)", tt.args.src, tt.args.m)
		})
	}
}
