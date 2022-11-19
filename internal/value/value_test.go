// Package value 提供值相关的封装
package value

import (
	"errors"
	"math"
	"reflect"
	"testing"
)

func TestAnyValue_Int(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "normal case:",
			fields: fields{
				val: int(1),
			},
			want:    int(1),
			wantErr: false,
		},
		{
			name: "error case:",
			fields: fields{
				err: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "type error case:",
			fields: fields{
				val: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			got, err := av.Int()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyValue.Int() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AnyValue.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_IntOr(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	type args struct {
		def int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "normal case:",
			fields: fields{
				val: int(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			fields: fields{
				val: int(1),
				err: errors.New("error"),
			},
			args: args{
				def: int(2),
			},
			want: int(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			if got := a.IntOr(tt.args.def); got != tt.want {
				t.Errorf("AnyValue.IntOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Uint(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint
		wantErr bool
	}{
		{
			name: "normal case:",
			fields: fields{
				val: uint(1),
			},
			want:    uint(1),
			wantErr: false,
		},
		{
			name: "error case:",
			fields: fields{
				err: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "type error case:",
			fields: fields{
				val: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			got, err := av.Uint()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyValue.Uint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AnyValue.Uint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_UintOr(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	type args struct {
		def uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint
	}{
		{
			name: "normal case:",
			fields: fields{
				val: uint(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			fields: fields{
				val: uint(1),
				err: errors.New("error"),
			},
			args: args{
				def: uint(2),
			},
			want: uint(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			if got := a.UintOr(tt.args.def); got != tt.want {
				t.Errorf("AnyValue.UintOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Int32(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		want    int32
		wantErr bool
	}{
		{
			name: "normal case:",
			fields: fields{
				val: int32(1),
			},
			want:    int32(1),
			wantErr: false,
		},
		{
			name: "error case:",
			fields: fields{
				err: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "type error case:",
			fields: fields{
				val: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			got, err := av.Int32()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyValue.Int32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AnyValue.Int32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Int32Or(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	type args struct {
		def int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int32
	}{
		{
			name: "normal case:",
			fields: fields{
				val: int32(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			fields: fields{
				val: int32(1),
				err: errors.New("error"),
			},
			args: args{
				def: int32(2),
			},
			want: int32(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			if got := a.Int32Or(tt.args.def); got != tt.want {
				t.Errorf("AnyValue.Int32Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Uint32(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint32
		wantErr bool
	}{
		{
			name: "normal case:",
			fields: fields{
				val: uint32(1),
			},
			want:    uint32(1),
			wantErr: false,
		},
		{
			name: "error case:",
			fields: fields{
				err: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "type error case:",
			fields: fields{
				val: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			got, err := av.Uint32()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyValue.Uint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AnyValue.Uint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Uint32Or(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	type args struct {
		def uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		{
			name: "normal case:",
			fields: fields{
				val: uint32(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			fields: fields{
				val: uint32(1),
				err: errors.New("error"),
			},
			args: args{
				def: uint32(2),
			},
			want: uint32(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			if got := a.Uint32Or(tt.args.def); got != tt.want {
				t.Errorf("AnyValue.Uint32Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Int64(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		{
			name: "normal case:",
			fields: fields{
				val: int64(1),
			},
			want:    int64(1),
			wantErr: false,
		},
		{
			name: "error case:",
			fields: fields{
				err: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "type error case:",
			fields: fields{
				val: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			got, err := av.Int64()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyValue.Int64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AnyValue.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Int64Or(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	type args struct {
		def int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "normal case:",
			fields: fields{
				val: int64(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			fields: fields{
				val: int64(1),
				err: errors.New("error"),
			},
			args: args{
				def: int64(2),
			},
			want: int64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			if got := a.Int64Or(tt.args.def); got != tt.want {
				t.Errorf("AnyValue.Int64Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Uint64(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint64
		wantErr bool
	}{
		{
			name: "normal case:",
			fields: fields{
				val: uint64(1),
			},
			want:    uint64(1),
			wantErr: false,
		},
		{
			name: "error case:",
			fields: fields{
				err: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "type error case:",
			fields: fields{
				val: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			got, err := av.Uint64()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyValue.Uint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AnyValue.Uint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Uint64Or(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	type args struct {
		def uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint64
	}{
		{
			name: "normal case:",
			fields: fields{
				val: uint64(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			fields: fields{
				val: uint64(1),
				err: errors.New("error"),
			},
			args: args{
				def: uint64(2),
			},
			want: uint64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			if got := a.Uint64Or(tt.args.def); got != tt.want {
				t.Errorf("AnyValue.Uint64Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

const MIN = 1e-06

func isEqual(f1, f2 float64) bool {
	return math.Dim(f1, f2) < MIN
}

func TestAnyValue_Float32(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		want    float32
		wantErr bool
	}{
		{
			name: "normal case:",
			fields: fields{
				val: float32(1),
			},
			want:    float32(1),
			wantErr: false,
		},
		{
			name: "error case:",
			fields: fields{
				err: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "type error case:",
			fields: fields{
				val: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			got, err := av.Float32()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyValue.Float32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !isEqual(float64(got), float64(tt.want)) {
				t.Errorf("AnyValue.Float32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Float32Or(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	type args struct {
		def float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float32
	}{
		{
			name: "normal case:",
			fields: fields{
				val: float32(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			fields: fields{
				val: float32(1),
				err: errors.New("error"),
			},
			args: args{
				def: float32(2),
			},
			want: float32(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			if got := a.Float32Or(tt.args.def); !isEqual(float64(got), float64(tt.want)) {
				t.Errorf("AnyValue.Float32Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Float64(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		wantErr bool
	}{
		{
			name: "normal case:",
			fields: fields{
				val: float64(1),
			},
			want:    float64(1),
			wantErr: false,
		},
		{
			name: "error case:",
			fields: fields{
				err: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "type error case:",
			fields: fields{
				val: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			got, err := av.Float64()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyValue.Float64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !isEqual(got, tt.want) {
				t.Errorf("AnyValue.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Float64Or(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	type args struct {
		def float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "normal case:",
			fields: fields{
				val: float64(1),
			},
			want: 1,
		},
		{
			name: "default case:",
			fields: fields{
				val: float64(1),
				err: errors.New("error"),
			},
			args: args{
				def: float64(2),
			},
			want: float64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			if got := a.Float64Or(tt.args.def); !isEqual(got, tt.want) {
				t.Errorf("AnyValue.Float64Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_String(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "normal case:",
			fields: fields{
				val: "111",
			},
			want:    "111",
			wantErr: false,
		},
		{
			name: "error case:",
			fields: fields{
				err: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "type error case:",
			fields: fields{
				val: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			got, err := av.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyValue.String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AnyValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_StringOr(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	type args struct {
		def string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "normal case:",
			fields: fields{
				val: "111",
			},
			want: "111",
		},
		{
			name: "default case:",
			fields: fields{
				val: "111",
				err: errors.New("error"),
			},
			args: args{
				def: "222",
			},
			want: "222",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			if got := a.StringOr(tt.args.def); got != tt.want {
				t.Errorf("AnyValue.StringOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_Bytes(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "normal case:",
			fields: fields{
				val: []byte("111"),
			},
			want:    []byte("111"),
			wantErr: false,
		},
		{
			name: "error case:",
			fields: fields{
				err: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "type error case:",
			fields: fields{
				val: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			got, err := av.Bytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyValue.Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyValue.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyValue_BytesOr(t *testing.T) {
	type fields struct {
		val any
		err error
	}
	type args struct {
		def []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "normal case:",
			fields: fields{
				val: []byte("111"),
			},
			want: []byte("111"),
		},
		{
			name: "default case:",
			fields: fields{
				val: []byte("111"),
				err: errors.New("error"),
			},
			args: args{
				def: []byte("222"),
			},
			want: []byte("222"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyValue{
				val: tt.fields.val,
				err: tt.fields.err,
			}
			if got := a.BytesOr(tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyValue.BytesOr() = %v, want %v", got, tt.want)
			}
		})
	}
}
