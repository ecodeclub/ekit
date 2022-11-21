// Package value 提供值相关的封装
package value

import (
	"reflect"

	"github.com/gotomicro/ekit/internal/errs"
)

type AnyValue struct {
	Val any
	Err error
}

// Int 返回 int 数据
func (av AnyValue) Int() (int, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(int)
	if !ok {
		return 0, errs.NewErrInvalidType("int", reflect.TypeOf(av.Val).Name())
	}
	return val, nil
}

// IntOr 返回 int 数据，或者默认值
func (a AnyValue) IntOr(def int) int {
	val, err := a.Int()
	if err != nil {
		return def
	}
	return val
}

// Uint 返回 uint 数据
func (av AnyValue) Uint() (uint, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(uint)
	if !ok {
		return 0, errs.NewErrInvalidType("uint", reflect.TypeOf(av.Val).Name())
	}
	return val, nil
}

// UintOr 返回 uint 数据，或者默认值
func (a AnyValue) UintOr(def uint) uint {
	val, err := a.Uint()
	if err != nil {
		return def
	}
	return val
}

// Int32 返回 int32 数据
func (av AnyValue) Int32() (int32, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(int32)
	if !ok {
		return 0, errs.NewErrInvalidType("int32", reflect.TypeOf(av.Val).Name())
	}
	return val, nil
}

// Int32Or 返回 int32 数据，或者默认值
func (a AnyValue) Int32Or(def int32) int32 {
	val, err := a.Int32()
	if err != nil {
		return def
	}
	return val
}

// Uint32 返回 uint32 数据
func (av AnyValue) Uint32() (uint32, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(uint32)
	if !ok {
		return 0, errs.NewErrInvalidType("uint32", reflect.TypeOf(av.Val).Name())
	}
	return val, nil
}

// Uint32Or 返回 uint32 数据，或者默认值
func (a AnyValue) Uint32Or(def uint32) uint32 {
	val, err := a.Uint32()
	if err != nil {
		return def
	}
	return val
}

// Int64 返回 int64 数据
func (av AnyValue) Int64() (int64, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(int64)
	if !ok {
		return 0, errs.NewErrInvalidType("int64", reflect.TypeOf(av.Val).Name())
	}
	return val, nil
}

// Int64Or 返回 int64 数据，或者默认值
func (a AnyValue) Int64Or(def int64) int64 {
	val, err := a.Int64()
	if err != nil {
		return def
	}
	return val
}

// Uint64 返回 uint64 数据
func (av AnyValue) Uint64() (uint64, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(uint64)
	if !ok {
		return 0, errs.NewErrInvalidType("uint64", reflect.TypeOf(av.Val).Name())
	}
	return val, nil
}

// Uint64Or 返回 uint64 数据，或者默认值
func (a AnyValue) Uint64Or(def uint64) uint64 {
	val, err := a.Uint64()
	if err != nil {
		return def
	}
	return val
}

// Float32 返回 float32 数据
func (av AnyValue) Float32() (float32, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(float32)
	if !ok {
		return 0, errs.NewErrInvalidType("float32", reflect.TypeOf(av.Val).Name())
	}
	return val, nil
}

// Float32Or 返回 float32 数据，或者默认值
func (a AnyValue) Float32Or(def float32) float32 {
	val, err := a.Float32()
	if err != nil {
		return def
	}
	return val
}

// Float64 返回 float64 数据
func (av AnyValue) Float64() (float64, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(float64)
	if !ok {
		return 0, errs.NewErrInvalidType("float64", reflect.TypeOf(av.Val).Name())
	}
	return val, nil
}

// Float64Or 返回 float64 数据，或者默认值
func (a AnyValue) Float64Or(def float64) float64 {
	val, err := a.Float64()
	if err != nil {
		return def
	}
	return val
}

// String 返回 string 数据
func (av AnyValue) String() (string, error) {
	if av.Err != nil {
		return "", av.Err
	}
	val, ok := av.Val.(string)
	if !ok {
		return "", errs.NewErrInvalidType("string", reflect.TypeOf(av.Val).Name())
	}
	return val, nil
}

// StringOr 返回 string 数据，或者默认值
func (a AnyValue) StringOr(def string) string {
	val, err := a.String()
	if err != nil {
		return def
	}
	return val
}

// Bytes 返回 []byte 数据
func (av AnyValue) Bytes() ([]byte, error) {
	if av.Err != nil {
		return nil, av.Err
	}
	val, ok := av.Val.([]byte)
	if !ok {
		return nil, errs.NewErrInvalidType("[]byte", reflect.TypeOf(av.Val).Name())
	}
	return val, nil
}

// BytesOr 返回 []byte 数据，或者默认值
func (a AnyValue) BytesOr(def []byte) []byte {
	val, err := a.Bytes()
	if err != nil {
		return def
	}
	return val
}
