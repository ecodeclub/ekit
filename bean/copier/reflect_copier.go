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

package copier

import (
	"reflect"
	"time"

	"github.com/ecodeclub/ekit/set"

	"github.com/ecodeclub/ekit/bean/option"
)

var defaultAtomicTypes = []reflect.Type{
	reflect.TypeOf(time.Time{}),
}

// ReflectCopier 基于反射的实现
// ReflectCopier 是浅拷贝
type ReflectCopier[Src any, Dst any] struct {

	// rootField 字典树的根节点
	rootField fieldNode

	// options 执行复制操作时的可选配置
	// 如果默认配置和Copy()/CopyTo()中的配置同名,会替换defaultOptions同名内容
	// 初始化时的默认配置,仅作为记录,执行时会拷贝到options中
	defaultOptions options

	atomicTypes []reflect.Type
}

// fieldNode 字段的前缀树
type fieldNode struct {
	// 当前节点的名字
	name string

	// 当前 Struct 的子节点, 如果为叶子节点, 则没有子节点
	fields []fieldNode

	// 在 source 的 index
	srcIndex int

	// 在 dst 的 index
	dstIndex int

	// 是否为叶子节点, 如果为叶子节点, 应该直接进行拷贝该字段
	isLeaf bool
}

// NewReflectCopier 如果类型不匹配, 创建时直接检查报错.
func NewReflectCopier[Src any, Dst any](opts ...option.Option[options]) (*ReflectCopier[Src, Dst], error) {
	src := new(Src)
	srcTyp := reflect.TypeOf(src).Elem()
	dst := new(Dst)
	dstTyp := reflect.TypeOf(dst).Elem()
	root := fieldNode{
		isLeaf: false,
		fields: []fieldNode{},
	}
	if srcTyp.Kind() != reflect.Struct {
		return nil, newErrTypeError(srcTyp)
	}
	if dstTyp.Kind() != reflect.Struct {
		return nil, newErrTypeError(dstTyp)
	}

	copier := &ReflectCopier[Src, Dst]{
		atomicTypes: defaultAtomicTypes,
	}

	if err := copier.createFieldNodes(&root, srcTyp, dstTyp); err != nil {
		return nil, err
	}
	copier.rootField = root

	defaultOpts := newOptions()
	option.Apply(&defaultOpts, opts...)
	copier.defaultOptions = defaultOpts
	return copier, nil
}

// createFieldNodes 递归创建 field 的前缀树, srcTyp 和 dstTyp 只能是结构体
func (r *ReflectCopier[Src, Dst]) createFieldNodes(root *fieldNode, srcTyp, dstTyp reflect.Type) error {

	fieldMap := map[string]int{}
	for i := 0; i < srcTyp.NumField(); i++ {
		srcFieldTypStruct := srcTyp.Field(i)
		if !srcFieldTypStruct.IsExported() {
			continue
		}
		fieldMap[srcFieldTypStruct.Name] = i
	}

	for dstIndex := 0; dstIndex < dstTyp.NumField(); dstIndex++ {

		dstFieldTypStruct := dstTyp.Field(dstIndex)
		if !dstFieldTypStruct.IsExported() {
			continue
		}
		srcIndex, ok := fieldMap[dstFieldTypStruct.Name]
		if !ok {
			continue
		}
		srcFieldTypStruct := srcTyp.Field(srcIndex)

		if srcFieldTypStruct.Type.Kind() == reflect.Pointer && srcFieldTypStruct.Type.Elem().Kind() == reflect.Pointer {
			return newErrMultiPointer(srcFieldTypStruct.Name)
		}
		if dstFieldTypStruct.Type.Kind() == reflect.Pointer && dstFieldTypStruct.Type.Elem().Kind() == reflect.Pointer {
			return newErrMultiPointer(dstFieldTypStruct.Name)
		}

		child := fieldNode{
			fields:   []fieldNode{},
			srcIndex: srcIndex,
			dstIndex: dstIndex,
			isLeaf:   false,
			name:     dstFieldTypStruct.Name,
		}

		fieldSrcTyp := srcFieldTypStruct.Type
		fieldDstTyp := dstFieldTypStruct.Type
		if fieldSrcTyp.Kind() == reflect.Pointer {
			fieldSrcTyp = fieldSrcTyp.Elem()
		}

		if fieldDstTyp.Kind() == reflect.Pointer {
			fieldDstTyp = fieldDstTyp.Elem()
		}

		if isShadowCopyType(fieldSrcTyp.Kind()) {
			// 内置类型，但不匹配，如别名、map和slice
			// 说明当前节点是叶子节点, 直接拷贝
			child.isLeaf = true
		} else if r.isAtomicType(fieldSrcTyp) {
			// 指定可作为一个整体的类型,不用递归
			// 同上，当当前节点是叶子节点时, 直接拷贝
			child.isLeaf = true
		} else if fieldSrcTyp.Kind() == reflect.Struct {
			if err := r.createFieldNodes(&child, fieldSrcTyp, fieldDstTyp); err != nil {
				return err
			}
		} else {
			// 不是我们能复制的类型, 直接跳过
			continue
		}

		root.fields = append(root.fields, child)
	}
	return nil
}

func (r *ReflectCopier[Src, Dst]) Copy(src *Src, opts ...option.Option[options]) (*Dst, error) {
	dst := new(Dst)
	err := r.CopyTo(src, dst, opts...)
	return dst, err
}

// CopyTo 执行复制
// 执行复制的逻辑是：
// 1. 按照字段的映射关系进行匹配
// 2. 如果 Src 和 Dst 中匹配的字段，其类型是基本类型（及其指针）或者内置类型（及其指针），并且类型一样，则直接用 Src 的值
// 3. 如果 Src 和 Dst 中匹配的字段，其类型都是结构体，或者都是结构体指针，则会深入复制
// 4. 否则，忽略字段
func (r *ReflectCopier[Src, Dst]) CopyTo(src *Src, dst *Dst, opts ...option.Option[options]) error {
	localOption := r.copyDefaultOptions()
	option.Apply(&localOption, opts...)
	return r.copyToWithTree(src, dst, localOption)
}

// copyDefaultOptions 复制默认配置
func (r *ReflectCopier[Src, Dst]) copyDefaultOptions() options {
	localOption := newOptions()
	// 复制ignoreFields default配置
	if r.defaultOptions.ignoreFields != nil {
		ignoreFields := set.NewMapSet[string](8)
		for _, key := range r.defaultOptions.ignoreFields.Keys() {
			ignoreFields.Add(key)
		}
		localOption.ignoreFields = ignoreFields
	}

	// 复制convertFields default配置
	for field, convert := range r.defaultOptions.convertFields {
		if localOption.convertFields == nil {
			localOption.convertFields = make(map[string]converterWrapper, 8)
		}
		localOption.convertFields[field] = convert
	}
	return localOption
}

func (r *ReflectCopier[Src, Dst]) copyToWithTree(src *Src, dst *Dst, opts options) error {
	srcTyp := reflect.TypeOf(src)
	dstTyp := reflect.TypeOf(dst)
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	return r.copyTreeNode(srcTyp, srcValue, dstTyp, dstValue, &r.rootField, opts)
}

func (r *ReflectCopier[Src, Dst]) copyTreeNode(srcTyp reflect.Type, srcValue reflect.Value,
	dstType reflect.Type, dstValue reflect.Value, root *fieldNode, opts options) error {
	originSrcVal := srcValue
	originDstVal := dstValue
	if srcValue.Kind() == reflect.Pointer {
		if srcValue.IsNil() {
			return nil
		}
		srcValue = srcValue.Elem()
		srcTyp = srcTyp.Elem()
	}

	if dstValue.Kind() == reflect.Pointer {
		if dstValue.IsNil() {
			dstValue.Set(reflect.New(dstType.Elem()))
		}
		dstValue = dstValue.Elem()
		dstType = dstType.Elem()
	}

	// 执行拷贝
	if root.isLeaf {
		convert, ok := opts.convertFields[root.name]
		if !dstValue.CanSet() {
			return nil
		}
		// 获取convert失败,就需要检测类型是否匹配,类型匹配就直接set
		if !ok {
			if srcTyp != dstType {
				return newErrTypeNotMatchError(srcTyp, dstType, root.name)
			}
			if srcValue.IsZero() {
				return nil
			}
			dstValue.Set(srcValue)
			return nil
		}

		// 字段执行转换函数时,需要用到原始类型进行判断,set的时候也是根据原始value设置
		if !originDstVal.CanSet() {
			return nil
		}
		srcConv, err := convert(originSrcVal.Interface())
		if err != nil {
			return err
		}

		srcConvType := reflect.TypeOf(srcConv)
		srcConvVal := reflect.ValueOf(srcConv)
		// 待设置的value和转换获取的value类型不匹配
		if srcConvType != originDstVal.Type() {
			return newErrTypeNotMatchError(srcConvType, originDstVal.Type(), root.name)
		}

		originDstVal.Set(srcConvVal)
		return nil
	}

	for i := range root.fields {
		child := &root.fields[i]

		// 只要结构体属性的名字在需要忽略的字段里面，就不走下面的复制逻辑
		if opts.InIgnoreFields(child.name) {
			continue
		}

		childSrcTyp := srcTyp.Field(child.srcIndex)
		childSrcValue := srcValue.Field(child.srcIndex)

		childDstTyp := dstType.Field(child.dstIndex)
		childDstValue := dstValue.Field(child.dstIndex)
		if err := r.copyTreeNode(childSrcTyp.Type, childSrcValue, childDstTyp.Type, childDstValue, child, opts); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReflectCopier[Src, Dst]) isAtomicType(typ reflect.Type) bool {
	for _, dt := range r.atomicTypes {
		if dt == typ {
			return true
		}
	}
	return false
}

func isShadowCopyType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String,
		reflect.Slice,
		reflect.Map,
		reflect.Chan,
		reflect.Array:
		return true
	}
	return false
}
