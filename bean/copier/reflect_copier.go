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

package copier

import (
	"errors"
	"reflect"
	"unsafe"
)

var errInvalidType = errors.New("only support struct")
var errToDoPointer = errors.New("doesn't support inside multi-pointer yet")

type structFieldMap map[string][]int

// ReflectCopier 基于反射的实现
// ReflectCopier 是浅拷贝
type ReflectCopier[Src any, Dst any] struct {
	// fieldMap Src 中字段下标到 Dst 字段下标的映射
	// 其中 key 是 Src 中字段下标使用连字符 - 连接起来
	// 不在 fieldMap 中字段则意味着被忽略
	fieldMap map[string][]int

	structHelperMap map[string]*structOffsets
}

func NewReflectCopier[Src any, Dst any](structHelpMap map[string]*structOffsets) (*ReflectCopier[Src, Dst], error) {
	src := new(Src)
	dst := new(Dst)
	fieldMap, err := makeFieldMap(reflect.TypeOf(src).Elem(), reflect.TypeOf(dst).Elem())
	if err != nil {
		return nil, err
	}
	return &ReflectCopier[Src, Dst]{
		fieldMap:        fieldMap,
		structHelperMap: structHelpMap,
	}, nil
}

func makeFieldMap(srcType reflect.Type, dstType reflect.Type) (structFieldMap, error) {
	//只支持srcVal和dstVal同时为结构体
	if srcType.Kind() != reflect.Struct {
		return nil, newErrTypeError(srcType)
	}
	if dstType.Kind() != reflect.Struct {
		return nil, newErrTypeError(dstType)
	}
	fieldMap := make(structFieldMap)
	//确定目标字段中的key，并记录其偏移量，以供后面直接操作内存
	for i := 0; i < dstType.NumField(); i++ {
		fileKey := dstType.Field(i).Name + "-" + dstType.Field(i).Type.String()
		fieldMap[fileKey] = append(fieldMap[fileKey], i)
	}
	for i := 0; i < srcType.NumField(); i++ {
		fileKey := srcType.Field(i).Name + "-" + srcType.Field(i).Type.String()
		if _, ok := fieldMap[fileKey]; ok {
			switch srcType.Kind() {
			case reflect.Map:
				keyKind := srcType.Key().Kind()
				valueKind := srcType.Elem().Kind()
				if keyKind != srcType.Key().Kind() || valueKind != srcType.Elem().Kind() {
					continue
				}
			case reflect.Slice:
				if srcType.Elem().Kind() != dstType.Elem().Kind() {
					continue
				}
			}
			fieldMap[fileKey] = append(fieldMap[fileKey], i)
		}
	}
	for key, pair := range fieldMap {
		if len(pair) < 2 {
			delete(fieldMap, key)
		}
	}
	return fieldMap, nil
}

// CopyTo 执行复制
// 执行复制的逻辑是：
// 1. 按照字段的映射关系进行匹配
// 2. 如果 Src 和 Dst 中匹配的字段，其类型是基本类型（及其指针）或者内置类型（及其指针），并且类型一样，则直接用 Src 的值
// 3. 如果 Src 和 Dst 中匹配的字段，其类型都是结构体，或者都是结构体指针，那么会深入复制
// 4. 否则，返回类型不匹配的错误
//
// 需求疑惑的点,关于匹配的粒度问题, 目前仅对完全匹配（名称-类型）的保证深度复制
// 对于结构体内部的完全匹配，暂时不予考虑，详细见测试用例 Complicated Sub Struct
func (r *ReflectCopier[Src, Dst]) CopyTo(src *Src, dst *Dst) error {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()
	for _, pair := range r.fieldMap {
		dstSettable := structOffsetValue(dstVal, pair[0])
		srcSettable := structOffsetValue(srcVal, pair[1])
		if err := r.copyTo(srcSettable, dstSettable); err != nil {
			return err
		}
	}
	return nil
}

func structOffsetValue(value reflect.Value, index int) reflect.Value {
	return reflect.NewAt(value.Field(index).Type(), unsafe.Pointer(uintptr(value.UnsafeAddr())+value.Type().Field(index).Offset)).Elem()
}

func checkValid(srcVal reflect.Value, dstVal reflect.Value) error {
	if srcVal.Type().Kind() != dstVal.Type().Kind() {
		panic("srcVal 和 dstVal一定要属于相同类型")
	}
	if srcVal.Type().Kind() == reflect.UnsafePointer || dstVal.Type().Kind() == reflect.Pointer {
		return errToDoPointer
	}

	if srcVal.Type().Kind() != reflect.Struct {
		return errInvalidType
	}
	return nil
}

func (r *ReflectCopier[Src, Dst]) copyTo(srcVal reflect.Value, dstVal reflect.Value) error {

	switch srcVal.Type().Kind() {
	case reflect.Pointer, reflect.UnsafePointer:
		if srcVal.Pointer() != 0 {
			subType := srcVal.Elem().Type()
			newValue := reflect.New(subType)
			r.copyTo(srcVal.Elem(), newValue.Elem())
			dstVal.Set(newValue)
		}
		return nil
	//case reflect.Map:
	//case reflect.Slice:
	//case reflect.Array:
	case reflect.Struct:
		dstVal.Set(srcVal)
		structOffsets, err := findTypeOffsets(srcVal.Type(), r.structHelperMap)
		if err != nil {
			return err
		}
		srcAddr := srcVal.UnsafeAddr()
		dstAddr := dstVal.UnsafeAddr()
		// 解决普通的指针结构
		for _, h := range structOffsets.helper {
			subSrcValue := reflect.NewAt(h.typ, unsafe.Pointer(srcAddr+h.ptrOffset)).Elem()
			subDstValue := reflect.NewAt(h.typ, unsafe.Pointer(dstAddr+h.ptrOffset)).Elem()
			if subSrcValue.Pointer() != 0 {
				newValue := reflect.New(h.typ.Elem())
				err = r.copyTo(subSrcValue.Elem(), newValue.Elem())
				if err != nil {
					return err
				}
				subDstValue.Set(newValue)
			} else {
				subDstValue.SetPointer(nil)
			}
		}
		// 解决slice、map等特殊需要深度拷贝的结构

	default:
		dstVal.Set(srcVal)
	}
	return nil
}

func (r *ReflectCopier[Src, Dst]) Copy(src *Src) (*Dst, error) {
	dst := new(Dst)
	err := r.CopyTo(src, dst)
	return dst, err
}
