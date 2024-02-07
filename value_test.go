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

package ekit

import (
	"errors"
	"testing"

	"github.com/ecodeclub/ekit/internal/errs"
	"github.com/stretchr/testify/assert"
)

func TestAnyValue_Int(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want int
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: int(1),
			},
			want: int(1),
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			err: errs.NewErrInvalidType("int", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.Int()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_IntOrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  int
		want int
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: int(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: int(1),
				Err: errors.New("error"),
			},
			def:  int(2),
			want: int(2),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			def:  int(1),
			want: int(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, a.IntOrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Uint(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want uint
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: uint(1),
			},
			want: uint(1),
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []string{"111"},
			},
			err: errs.NewErrInvalidType("uint", []string{"111"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.Uint()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_UintOrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  uint
		want uint
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: uint(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: uint(1),
				Err: errors.New("error"),
			},
			def:  uint(2),
			want: uint(2),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			def:  uint(2),
			want: uint(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, a.UintOrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Int32(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want int32
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: int32(1),
			},
			want: int32(1),
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			err: errs.NewErrInvalidType("int32", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.Int32()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_Int32OrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  int32
		want int32
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: int32(1),
			},
			want: int32(1),
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: int32(1),
				Err: errors.New("error"),
			},
			def:  int32(2),
			want: int32(2),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			def:  int32(2),
			want: int32(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, a.Int32OrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Uint32(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want uint32
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: uint32(1),
			},
			want: uint32(1),
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			err: errs.NewErrInvalidType("uint32", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.Uint32()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_Uint32OrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  uint32
		want uint32
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: uint32(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: uint32(1),
				Err: errors.New("error"),
			},

			def:  uint32(2),
			want: uint32(2),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			def:  uint32(2),
			want: uint32(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, a.Uint32OrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Int64(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want int64
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: int64(1),
			},
			want: int64(1),
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			err: errs.NewErrInvalidType("int64", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.Int64()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_Int64OrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  int64
		want int64
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: int64(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: int64(1),
				Err: errors.New("error"),
			},
			def:  int64(2),
			want: int64(2),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			def:  int64(2),
			want: int64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, a.Int64OrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Uint64(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want uint64
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: uint64(1),
			},
			want: uint64(1),
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			err: errs.NewErrInvalidType("uint64", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.Uint64()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_Uint64OrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  uint64
		want uint64
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: uint64(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: uint64(1),
				Err: errors.New("error"),
			},
			def:  uint64(2),
			want: uint64(2),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			def:  uint64(2),
			want: uint64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, a.Uint64OrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Float32(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want float32
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: float32(1),
			},
			want: float32(1),
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			err: errs.NewErrInvalidType("float32", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.Float32()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_Float32OrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  float32
		want float32
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: float32(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: float32(1),
				Err: errors.New("error"),
			},
			def:  float32(2),
			want: float32(2),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			def:  float32(2),
			want: float32(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, a.Float32OrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Float64(t *testing.T) {

	tests := []struct {
		name string
		val  AnyValue
		want float64
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: float64(1),
			},
			want: float64(1),
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			err: errs.NewErrInvalidType("float64", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.Float64()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_Float64OrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  float64
		want float64
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: float64(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: float64(1),
				Err: errors.New("error"),
			},
			def:  float64(2),
			want: float64(2),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: "",
			},
			def:  float64(2),
			want: float64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, a.Float64OrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_String(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want string
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: "111",
			},
			want: "111",
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: 1,
			},
			err: errs.NewErrInvalidType("string", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.String()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_StringOrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  string
		want string
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: "111",
			},
			want: "111",
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: "111",
				Err: errors.New("error"),
			},
			def:  "222",
			want: "222",
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: 1,
			},
			def:  "222",
			want: "222",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, a.StringOrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Bytes(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want []byte
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: []byte("111"),
			},
			want: []byte("111"),
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: 1,
			},
			err: errs.NewErrInvalidType("[]byte", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.Bytes()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_BytesOrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  []byte
		want []byte
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: []byte("111"),
			},
			want: []byte("111"),
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: []byte("111"),
				Err: errors.New("error"),
			},
			def:  []byte("222"),
			want: []byte("222"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: 1,
			},
			def:  []byte("222"),
			want: []byte("222"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, a.BytesOrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Bool(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want bool
		err  error
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: true,
			},
			want: true,
			err:  nil,
		},
		{
			name: "error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: 1,
			},
			err: errs.NewErrInvalidType("bool", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			got, err := av.Bool()
			assert.Equal(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAnyValue_BoolOrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  bool
		want bool
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: true,
			},
			want: true,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: true,
				Err: errors.New("error"),
			},
			def:  false,
			want: false,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: 1,
			},
			def:  true,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, av.BoolOrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Int8OrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  int8
		want int8
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: int8(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: int8(0),
				Err: errors.New("error"),
			},
			def:  1,
			want: 1,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: true,
			},
			def:  10,
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, av.Int8OrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Int16OrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  int16
		want int16
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: int16(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: int16(0),
				Err: errors.New("error"),
			},
			def:  1,
			want: 1,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: true,
			},
			def:  10,
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, av.Int16OrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Uint8OrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  uint8
		want uint8
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: uint8(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: uint8(0),
				Err: errors.New("error"),
			},
			def:  1,
			want: 1,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: true,
			},
			def:  10,
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, av.Uint8OrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_Uint16OrDefault(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		def  uint16
		want uint16
	}{
		{
			name: "normal case:",
			val: AnyValue{
				Val: uint16(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			val: AnyValue{
				Val: uint16(0),
				Err: errors.New("error"),
			},
			def:  1,
			want: 1,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: true,
			},
			def:  10,
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				Val: tt.val.Val,
				Err: tt.val.Err,
			}
			assert.Equal(t, av.Uint16OrDefault(tt.def), tt.want)
		})
	}
}

func TestAnyValue_AsInt(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want int
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1",
			},
			want: 1,
		},
		{
			name: "normal int case:",
			val: AnyValue{
				Val: int(2),
			},
			want: 2,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("int", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Val: "",
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsInt()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsInt8(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want int8
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1",
			},
			want: 1,
		},
		{
			name: "normal int case:",
			val: AnyValue{
				Val: int8(2),
			},
			want: 2,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("int8", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Val: "",
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsInt8()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsInt16(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want int16
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1",
			},
			want: 1,
		},
		{
			name: "normal int16 case:",
			val: AnyValue{
				Val: int16(2),
			},
			want: 2,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("int16", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Val: "",
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsInt16()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsInt32(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want int32
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1",
			},
			want: 1,
		},
		{
			name: "normal int32 case:",
			val: AnyValue{
				Val: int32(2),
			},
			want: 2,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("int32", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Val: "",
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsInt32()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsInt64(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want int64
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1",
			},
			want: 1,
		},
		{
			name: "normal int64 case:",
			val: AnyValue{
				Val: int64(2),
			},
			want: 2,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("int64", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Val: "",
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsInt64()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsUint(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want uint
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1",
			},
			want: 1,
		},
		{
			name: "normal uint case:",
			val: AnyValue{
				Val: uint(2),
			},
			want: 2,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("uint", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Val: "",
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsUint()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsUint8(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want uint8
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1",
			},
			want: 1,
		},
		{
			name: "normal uint8 case:",
			val: AnyValue{
				Val: uint8(2),
			},
			want: 2,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("uint8", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Val: "",
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsUint8()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsUint16(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want uint16
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1",
			},
			want: 1,
		},
		{
			name: "normal uint16 case:",
			val: AnyValue{
				Val: uint16(2),
			},
			want: 2,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("uint16", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Val: "",
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsUint16()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsUint32(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want uint32
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1",
			},
			want: 1,
		},
		{
			name: "normal uint32 case:",
			val: AnyValue{
				Val: uint32(2),
			},
			want: 2,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("uint32", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Val: "",
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsUint32()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsUint64(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want uint64
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1",
			},
			want: 1,
		},
		{
			name: "normal uint64 case:",
			val: AnyValue{
				Val: uint64(2),
			},
			want: 2,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("uint64", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Val: "",
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsUint64()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsFloat32(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want float32
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "1.01",
			},
			want: 1.01,
		},
		{
			name: "normal float32 case:",
			val: AnyValue{
				Val: float32(2.44),
			},
			want: 2.44,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("float32", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsFloat32()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsFloat64(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want float64
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "100.0000000000",
			},
			want: 1e2,
		},
		{
			name: "normal float64 case:",
			val: AnyValue{
				Val: float64(2.44),
			},
			want: 2.44,
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			err: errs.NewErrInvalidType("float64", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			err: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsFloat64()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsBytes(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want []byte
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "hello",
			},
			want: []byte("hello"),
		},
		{
			name: "normal []byte case:",
			val: AnyValue{
				Val: []byte{1},
			},
			want: []byte{1},
		},
		{
			name: "type error case:",
			val: AnyValue{
				Val: []int{1},
			},
			want: []byte{},
			err:  errs.NewErrInvalidType("[]byte", []int{1}),
		},
		{
			name: "value exists error case:",
			val: AnyValue{
				Err: errors.New("error"),
			},
			want: []byte{},
			err:  errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsBytes()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_AsString(t *testing.T) {
	tests := []struct {
		name string
		val  AnyValue
		want string
		err  error
	}{
		{
			name: "normal string case:",
			val: AnyValue{
				Val: "hello ekit",
			},
			want: "hello ekit",
		},
		{
			name: "normal uint case:",
			val: AnyValue{
				Val: uint16(1231),
			},
			want: "1231",
		},
		{
			name: "normal int case:",
			val: AnyValue{
				Val: 1,
			},
			want: "1",
		},
		{
			name: "normal float case:",
			val: AnyValue{
				Val: 1e2,
			},
			want: "100.0000000000",
		},
		{
			name: "normal []byte case:",
			val: AnyValue{
				Val: []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33},
			},
			want: "Hello, World!",
		},
		{
			name: "type conversion failed",
			val: AnyValue{
				Val: []string{"h", "e", "llo"},
			},
			err: errs.NewErrInvalidType("[]byte", []string{"h", "e", "llo"}),
		},
		{
			name: "type conversion failed by int",
			val: AnyValue{
				Val: []int{1, 2, 3, 4, 5},
			},
			err: errs.NewErrInvalidType("[]byte", []int{1, 2, 3, 4, 5}),
		},
		{
			name: "unsupported type case:",
			val: AnyValue{
				Val: map[string]any{
					"test": 1,
					"hhh":  "sss",
				},
			},
			err: errors.New("未兼容类型，暂时无法转换"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.val.AsString()
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAnyValue_JSONScan(t *testing.T) {
	testCases := []struct {
		name string

		av AnyValue

		wantUser User
		wantErr  error
	}{
		{
			name: "OK",
			av: AnyValue{
				Val: `{"name": "Tom"}`,
			},
			wantUser: User{
				Name: "Tom",
			},
		},

		{
			name: "error",
			av: AnyValue{
				Err: errors.New("mock error"),
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var u User
			err := tc.av.JSONScan(&u)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

type User struct {
	Name string `json:"name"`
}
