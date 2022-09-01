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
