// Package Value 提供值相关的封装
package value

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gotomicro/ekit/internal/errs"
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
			err: errs.NewErrInvalidType("int", reflect.TypeOf("").Name()),
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

func TestAnyValue_IntOr(t *testing.T) {
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
			assert.Equal(t, a.IntOr(tt.def), tt.want)
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
				Val: "",
			},
			err: errs.NewErrInvalidType("uint", reflect.TypeOf("").Name()),
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

func TestAnyValue_UintOr(t *testing.T) {
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
			assert.Equal(t, a.UintOr(tt.def), tt.want)
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
			err: errs.NewErrInvalidType("int32", reflect.TypeOf("").Name()),
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

func TestAnyValue_Int32Or(t *testing.T) {
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
			assert.Equal(t, a.Int32Or(tt.def), tt.want)
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
			err: errs.NewErrInvalidType("uint32", reflect.TypeOf("").Name()),
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

func TestAnyValue_Uint32Or(t *testing.T) {
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
			assert.Equal(t, a.Uint32Or(tt.def), tt.want)
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
			err: errs.NewErrInvalidType("int64", reflect.TypeOf("").Name()),
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

func TestAnyValue_Int64Or(t *testing.T) {
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
			assert.Equal(t, a.Int64Or(tt.def), tt.want)
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
			err: errs.NewErrInvalidType("uint64", reflect.TypeOf("").Name()),
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

func TestAnyValue_Uint64Or(t *testing.T) {
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
			assert.Equal(t, a.Uint64Or(tt.def), tt.want)
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
			err: errs.NewErrInvalidType("float32", reflect.TypeOf("").Name()),
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

func TestAnyValue_Float32Or(t *testing.T) {
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
			assert.Equal(t, a.Float32Or(tt.def), tt.want)
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
			err: errs.NewErrInvalidType("float64", reflect.TypeOf("").Name()),
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

func TestAnyValue_Float64Or(t *testing.T) {
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
			assert.Equal(t, a.Float64Or(tt.def), tt.want)
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
			err: errs.NewErrInvalidType("string", reflect.TypeOf(111).Name()),
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

func TestAnyValue_StringOr(t *testing.T) {
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
			assert.Equal(t, a.StringOr(tt.def), tt.want)
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
			err: errs.NewErrInvalidType("[]byte", reflect.TypeOf(111).Name()),
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

func TestAnyValue_BytesOr(t *testing.T) {
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
			assert.Equal(t, a.BytesOr(tt.def), tt.want)
		})
	}
}
