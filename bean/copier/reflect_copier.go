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

// 实际上不需要，但可以记在心里
// const (
//
//	bInvalid   = 1 << 0
//	bBool      = 1 << 1
//	bInt       = 1 << 2
//	bInt8      = 1 << 3
//	bInt16     = 1 << 4
//	bInt32     = 1 << 5
//	bInt64     = 1 << 6
//	bUint      = 1 << 7
//	bUint8     = 1 << 8
//	bUint16    = 1 << 9
//	bUint32    = 1 << 10
//	bUint64    = 1 << 11
//	bUintptr   = 1 << 12
//	bFloat32   = 1 << 13
//	bFloat64   = 1 << 14
//	bComplex64  = 1 << 15
//	bComplex128  = 1 << 16
//	bArray     = 1 << 17
//	bChan      = 1 << 18
//	bFunc      = 1 << 19
//	bInterface = 1 << 20
//	bMap       = 1 << 21
//	bPointer   = 1 << 22
//	bSlice     = 1 << 23
//	bString    = 1 << 24
//	bStruct    = 1 << 25
//	bUnsafePointer  = 1 << 26
//
// )
var errInvalidType = errors.New("only support struct")

type structFieldMap map[string][]int

var structFieldMapBuffer = make(map[string]structFieldMap)

// |除Invalid外的基本数据类型| Func(一种特殊的指针，应当直接复制) | string | map | slice
const basicKind = (1<<17 - 1) | (1 << 19) | (1 << 24) | (1 << 21) | (1 << 23)

// ReflectCopier 基于反射的实现
// ReflectCopier 是浅拷贝
type ReflectCopier[Src any, Dst any] struct {
	// fieldMap Src 中字段下标到 Dst 字段下标的映射
	// 其中 key 是 Src 中字段下标使用连字符 - 连接起来
	// 不在 fieldMap 中字段则意味着被忽略
	fieldMap map[string][]int
}

func NewReflectCopier[Src any, Dst any]() *ReflectCopier[Src, Dst] {
	return &ReflectCopier[Src, Dst]{}
}

func makeFieldMap(srcVal reflect.Value, dstVal reflect.Value) structFieldMap {
	//只支持src和dst同时为结构体指针
	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		return nil
	}
	fieldMap := make(structFieldMap)
	//确定目标字段中的key，并记录其偏移量，以供后面直接操作内存
	for i := 0; i < dstVal.NumField(); i++ {
		fileKey := dstVal.Type().Field(i).Name + "-" + dstVal.Field(i).Kind().String()
		//r.fieldMap[fileKey] = append(r.fieldMap[fileKey], i)
		fieldMap[fileKey] = append(fieldMap[fileKey], i)
	}
	for i := 0; i < srcVal.NumField(); i++ {
		fileKey := srcVal.Type().Field(i).Name + "-" + srcVal.Field(i).Kind().String()
		if _, ok := fieldMap[fileKey]; ok {
			fieldMap[fileKey] = append(fieldMap[fileKey], i)
		}
	}
	return fieldMap
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
	mapName := srcVal.Type().Name() + "-" + dstVal.Type().Name()
	if offsetMap, ok := structFieldMapBuffer[mapName]; ok {
		r.fieldMap = offsetMap
	} else {
		r.fieldMap = makeFieldMap(srcVal, dstVal)
		structFieldMapBuffer[mapName] = r.fieldMap
	}
	if r.fieldMap == nil {
		return errInvalidType
	}

	for _, pair := range r.fieldMap {
		if len(pair) == 2 {
			dstSettable := structOffsetValue(dstVal, pair[0])
			srcSettable := structOffsetValue(srcVal, pair[1])
			copyTo(srcSettable, dstSettable)
		}
	}
	return nil
}

func structOffsetValue(value reflect.Value, index int) reflect.Value {
	return reflect.NewAt(value.Field(index).Type(), unsafe.Pointer(uintptr(value.UnsafeAddr())+value.Type().Field(index).Offset)).Elem()
}

func copyTo(src reflect.Value, dst reflect.Value) {
	//如果可以Set，就直接赋值
	switch src.Kind() {
	case reflect.Pointer, reflect.UnsafePointer:
		if src.UnsafePointer() == nil {
			dst.Set(src)
		} else {
			newPtrV := reflect.New(dst.Type().Elem())
			dst.Set(newPtrV)
			copyTo(src.Elem(), dst.Elem())
		}
		break
	default:
		if dst.CanSet() {
			dst.Set(src)
			// 对于struct字段中的指针类型额外的处理
			// 目前只支持一级指针的深度拷贝
			if dst.Kind() == reflect.Struct {
				numFiled := dst.NumField()
				for i := 0; i < numFiled; i++ {
					value := dst.Field(i)
					if value.Kind() == reflect.Pointer || value.Kind() == reflect.UnsafePointer {
						//新建一个类型，并赋值到value上
						dstNewPtr := reflect.New(value.Type().Elem())
						value.Set(dstNewPtr)
						copyTo(dst.Field(i).Elem(), value.Elem())
					}
				}
			}
		}
	}

}

func isBasicKind(kind reflect.Kind) bool {
	return 1<<kind&basicKind > 0
}

func (r *ReflectCopier[Src, Dst]) Copy(src *Src) (*Dst, error) {
	dst := new(Dst)
	err := r.CopyTo(src, dst)
	return dst, err
}
