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

	"github.com/ecodeclub/ekit/bean/option"
)

// ReflectCopier 基于反射的实现
// ReflectCopier 是浅拷贝
type ReflectCopier[Src any, Dst any] struct {

	// rootField 字典树的根节点
	rootField fieldNode

	// options 执行复制操作时的可选配置
	options *options
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
func NewReflectCopier[Src any, Dst any]() (*ReflectCopier[Src, Dst], error) {
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
	if err := createFieldNodes(&root, srcTyp, dstTyp); err != nil {
		return nil, err
	}

	copier := &ReflectCopier[Src, Dst]{
		rootField: root,
	}
	return copier, nil
}

// createFieldNodes 递归创建 field 的前缀树, srcTyp 和 dstTyp 只能是结构体
func createFieldNodes(root *fieldNode, srcTyp, dstTyp reflect.Type) error {

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
		if srcFieldTypStruct.Type.Kind() != dstFieldTypStruct.Type.Kind() {
			return newErrKindNotMatchError(srcFieldTypStruct.Type.Kind(), dstFieldTypStruct.Type.Kind(), dstFieldTypStruct.Name)
		}

		if srcFieldTypStruct.Type.Kind() == reflect.Pointer {
			if srcFieldTypStruct.Type.Elem().Kind() != dstFieldTypStruct.Type.Elem().Kind() {
				return newErrKindNotMatchError(srcFieldTypStruct.Type.Kind(), dstFieldTypStruct.Type.Kind(), dstFieldTypStruct.Name)
			}
			if srcFieldTypStruct.Type.Elem().Kind() == reflect.Pointer {
				return newErrMultiPointer(dstFieldTypStruct.Name)
			}
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
			fieldDstTyp = fieldDstTyp.Elem()
		}

		if isShadowCopyType(fieldSrcTyp.Kind()) {
			// 内置类型，但不匹配，如别名、map和slice
			if fieldSrcTyp != fieldDstTyp {
				return newErrTypeNotMatchError(srcFieldTypStruct.Type, dstFieldTypStruct.Type, dstFieldTypStruct.Name)
			}
			// 说明当前节点是叶子节点, 直接拷贝
			child.isLeaf = true
		} else if fieldSrcTyp.Kind() == reflect.Struct {
			if err := createFieldNodes(&child, fieldSrcTyp, fieldDstTyp); err != nil {
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
	opt := newOptions()
	option.Apply(opt, opts...)
	r.options = opt

	return r.copyToWithTree(src, dst)
}

func (r *ReflectCopier[Src, Dst]) copyToWithTree(src *Src, dst *Dst) error {
	srcTyp := reflect.TypeOf(src)
	dstTyp := reflect.TypeOf(dst)
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	return r.copyTreeNode(srcTyp, srcValue, dstTyp, dstValue, &r.rootField)
}

func (r *ReflectCopier[Src, Dst]) copyTreeNode(srcTyp reflect.Type, srcValue reflect.Value, dstType reflect.Type, dstValue reflect.Value, root *fieldNode) error {
	if srcValue.Kind() == reflect.Pointer {
		if srcValue.IsNil() {
			return nil
		}
		if dstValue.IsNil() {
			dstValue.Set(reflect.New(dstType.Elem()))
		}
		srcValue = srcValue.Elem()
		srcTyp = srcTyp.Elem()

		dstValue = dstValue.Elem()
		dstType = dstType.Elem()
	}
	// 执行拷贝
	if root.isLeaf {
		if dstValue.CanSet() {
			dstValue.Set(srcValue)
		}
		return nil
	}

	for i := range root.fields {
		child := &root.fields[i]

		// 只要结构体属性的名字在需要忽略的字段里面，就不走下面的复制逻辑
		if r.options.InIgnoreFields(child.name) {
			continue
		}

		childSrcTyp := srcTyp.Field(child.srcIndex)
		childSrcValue := srcValue.Field(child.srcIndex)

		childDstTyp := dstType.Field(child.dstIndex)
		childDstValue := dstValue.Field(child.dstIndex)
		if err := r.copyTreeNode(childSrcTyp.Type, childSrcValue, childDstTyp.Type, childDstValue, child); err != nil {
			return err
		}
	}
	return nil
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
