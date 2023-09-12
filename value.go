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
	"reflect"
	"strconv"

	"github.com/ecodeclub/ekit/internal/errs"
)

// AnyValue 类型转换结构定义
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
		return 0, errs.NewErrInvalidType("int", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

func (av AnyValue) AsInt() (int, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	switch v := av.Val.(type) {
	case int:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 64)
		return int(res), err
	}
	return 0, errs.NewErrInvalidType("int", reflect.TypeOf(av.Val).String())
}

// IntOrDefault 返回 int 数据，或者默认值
func (av AnyValue) IntOrDefault(def int) int {
	val, err := av.Int()
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
		return 0, errs.NewErrInvalidType("uint", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

func (av AnyValue) AsUint() (uint, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	switch v := av.Val.(type) {
	case uint:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 64)
		return uint(res), err
	}
	return 0, errs.NewErrInvalidType("uint", reflect.TypeOf(av.Val).String())
}

// UintOrDefault 返回 uint 数据，或者默认值
func (av AnyValue) UintOrDefault(def uint) uint {
	val, err := av.Uint()
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
		return 0, errs.NewErrInvalidType("int32", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

func (av AnyValue) AsInt32() (int32, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	switch v := av.Val.(type) {
	case int32:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 32)
		return int32(res), err
	}
	return 0, errs.NewErrInvalidType("int32", reflect.TypeOf(av.Val).String())
}

// Int32OrDefault 返回 int32 数据，或者默认值
func (av AnyValue) Int32OrDefault(def int32) int32 {
	val, err := av.Int32()
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
		return 0, errs.NewErrInvalidType("uint32", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

func (av AnyValue) AsUint32() (uint32, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	switch v := av.Val.(type) {
	case uint32:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 32)
		return uint32(res), err
	}
	return 0, errs.NewErrInvalidType("uint32", reflect.TypeOf(av.Val).String())
}

// Uint32OrDefault 返回 uint32 数据，或者默认值
func (av AnyValue) Uint32OrDefault(def uint32) uint32 {
	val, err := av.Uint32()
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
		return 0, errs.NewErrInvalidType("int64", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

func (av AnyValue) AsInt64() (int64, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	switch v := av.Val.(type) {
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	}
	return 0, errs.NewErrInvalidType("int64", reflect.TypeOf(av.Val).String())
}

// Int64OrDefault 返回 int64 数据，或者默认值
func (av AnyValue) Int64OrDefault(def int64) int64 {
	val, err := av.Int64()
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
		return 0, errs.NewErrInvalidType("uint64", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

func (av AnyValue) AsUint64() (uint64, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	switch v := av.Val.(type) {
	case uint64:
		return v, nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	}
	return 0, errs.NewErrInvalidType("uint64", reflect.TypeOf(av.Val).String())
}

// Uint64OrDefault 返回 uint64 数据，或者默认值
func (av AnyValue) Uint64OrDefault(def uint64) uint64 {
	val, err := av.Uint64()
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
		return 0, errs.NewErrInvalidType("float32", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

func (av AnyValue) AsFloat32() (float32, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	switch v := av.Val.(type) {
	case float32:
		return v, nil
	case string:
		res, err := strconv.ParseFloat(v, 32)
		return float32(res), err
	}
	return 0, errs.NewErrInvalidType("float32", reflect.TypeOf(av.Val).String())
}

// Float32OrDefault 返回 float32 数据，或者默认值
func (av AnyValue) Float32OrDefault(def float32) float32 {
	val, err := av.Float32()
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
		return 0, errs.NewErrInvalidType("float64", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

func (av AnyValue) AsFloat64() (float64, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	switch v := av.Val.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	}
	return 0, errs.NewErrInvalidType("float64", reflect.TypeOf(av.Val).String())
}

// Float64OrDefault 返回 float64 数据，或者默认值
func (av AnyValue) Float64OrDefault(def float64) float64 {
	val, err := av.Float64()
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
		return "", errs.NewErrInvalidType("string", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// StringOrDefault 返回 string 数据，或者默认值
func (av AnyValue) StringOrDefault(def string) string {
	val, err := av.String()
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
		return nil, errs.NewErrInvalidType("[]byte", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

func (av AnyValue) AsBytes() ([]byte, error) {
	if av.Err != nil {
		return []byte{}, av.Err
	}
	switch v := av.Val.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	}

	return []byte{}, errs.NewErrInvalidType("[]byte", reflect.TypeOf(av.Val).String())
}

// BytesOrDefault 返回 []byte 数据，或者默认值
func (av AnyValue) BytesOrDefault(def []byte) []byte {
	val, err := av.Bytes()
	if err != nil {
		return def
	}
	return val
}

// Bool 返回 bool 数据
func (av AnyValue) Bool() (bool, error) {
	if av.Err != nil {
		return false, av.Err
	}
	val, ok := av.Val.(bool)
	if !ok {
		return false, errs.NewErrInvalidType("bool", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// BoolOrDefault 返回 bool 数据，或者默认值
func (av AnyValue) BoolOrDefault(def bool) bool {
	val, err := av.Bool()
	if err != nil {
		return def
	}
	return val
}
