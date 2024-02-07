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
	"encoding/json"
	"errors"
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
		return 0, errs.NewErrInvalidType("int", av.Val)
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
	return 0, errs.NewErrInvalidType("int", av.Val)
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
		return 0, errs.NewErrInvalidType("uint", av.Val)
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
	return 0, errs.NewErrInvalidType("uint", av.Val)
}

// UintOrDefault 返回 uint 数据，或者默认值
func (av AnyValue) UintOrDefault(def uint) uint {
	val, err := av.Uint()
	if err != nil {
		return def
	}
	return val
}

func (av AnyValue) Int8() (int8, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(int8)
	if !ok {
		return 0, errs.NewErrInvalidType("int", av.Val)
	}
	return val, nil
}

func (av AnyValue) AsInt8() (int8, error) {
	if av.Err != nil {
		return 0, av.Err
	}

	switch v := av.Val.(type) {
	case int8:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 64)
		return int8(res), err
	}
	return 0, errs.NewErrInvalidType("int8", av.Val)
}

func (av AnyValue) Int8OrDefault(def int8) int8 {
	val, err := av.Int8()
	if err != nil {
		return def
	}
	return val
}

func (av AnyValue) Uint8() (uint8, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(uint8)
	if !ok {
		return 0, errs.NewErrInvalidType("uint8", av.Val)
	}
	return val, nil
}

func (av AnyValue) AsUint8() (uint8, error) {
	if av.Err != nil {
		return 0, av.Err
	}

	switch v := av.Val.(type) {
	case uint8:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 8)
		return uint8(res), err
	}
	return 0, errs.NewErrInvalidType("uint8", av.Val)
}

func (av AnyValue) Uint8OrDefault(def uint8) uint8 {
	val, err := av.Uint8()
	if err != nil {
		return def
	}
	return val
}

func (av AnyValue) Int16() (int16, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(int16)
	if !ok {
		return 0, errs.NewErrInvalidType("int16", av.Val)
	}
	return val, nil
}

func (av AnyValue) AsInt16() (int16, error) {
	if av.Err != nil {
		return 0, av.Err
	}

	switch v := av.Val.(type) {
	case int16:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 16)
		return int16(res), err
	}
	return 0, errs.NewErrInvalidType("int16", av.Val)
}

func (av AnyValue) Int16OrDefault(def int16) int16 {
	val, err := av.Int16()
	if err != nil {
		return def
	}
	return val
}

func (av AnyValue) Uint16() (uint16, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(uint16)
	if !ok {
		return 0, errs.NewErrInvalidType("uint16", av.Val)
	}
	return val, nil
}

func (av AnyValue) AsUint16() (uint16, error) {
	if av.Err != nil {
		return 0, av.Err
	}

	switch v := av.Val.(type) {
	case uint16:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 16)
		return uint16(res), err
	}
	return 0, errs.NewErrInvalidType("uint16", av.Val)
}

func (av AnyValue) Uint16OrDefault(def uint16) uint16 {
	val, err := av.Uint16()
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
		return 0, errs.NewErrInvalidType("int32", av.Val)
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
	return 0, errs.NewErrInvalidType("int32", av.Val)
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
		return 0, errs.NewErrInvalidType("uint32", av.Val)
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
	return 0, errs.NewErrInvalidType("uint32", av.Val)
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
		return 0, errs.NewErrInvalidType("int64", av.Val)
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
	return 0, errs.NewErrInvalidType("int64", av.Val)
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
		return 0, errs.NewErrInvalidType("uint64", av.Val)
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
	return 0, errs.NewErrInvalidType("uint64", av.Val)
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
		return 0, errs.NewErrInvalidType("float32", av.Val)
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
	return 0, errs.NewErrInvalidType("float32", av.Val)
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
		return 0, errs.NewErrInvalidType("float64", av.Val)
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
	return 0, errs.NewErrInvalidType("float64", av.Val)
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
		return "", errs.NewErrInvalidType("string", av.Val)
	}
	return val, nil
}

func (av AnyValue) AsString() (string, error) {
	if av.Err != nil {
		return "", av.Err
	}

	var val string
	valueOf := reflect.ValueOf(av.Val)
	switch valueOf.Type().Kind() {
	case reflect.String:
		val = valueOf.String()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = strconv.FormatUint(valueOf.Uint(), 10)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = strconv.FormatInt(valueOf.Int(), 10)
	case reflect.Float32:
		val = strconv.FormatFloat(valueOf.Float(), 'f', 10, 32)
	case reflect.Float64:
		val = strconv.FormatFloat(valueOf.Float(), 'f', 10, 64)
	case reflect.Slice:
		if valueOf.Type().Elem().Kind() != reflect.Uint8 {
			return "", errs.NewErrInvalidType("[]byte", av.Val)
		}
		val = string(valueOf.Bytes())
	default:
		return "", errors.New("未兼容类型，暂时无法转换")
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
		return nil, errs.NewErrInvalidType("[]byte", av.Val)
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

	return []byte{}, errs.NewErrInvalidType("[]byte", av.Val)
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
		return false, errs.NewErrInvalidType("bool", av.Val)
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

// JSONScan 将 val 转化为一个对象
func (av AnyValue) JSONScan(val any) error {
	data, err := av.AsBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}
