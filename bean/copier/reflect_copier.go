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
	"reflect"
)

// ReflectCopier 基于反射的实现
// ReflectCopier 是浅拷贝
type ReflectCopier[Src any, Dst any] struct {
	rootFiled fieldNode
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

	// 是否为根节点, 只有一个根节点
	isRoot bool
}

// NewReflectCopier 如果类型不匹配, 创建时直接检查报错.
func NewReflectCopier[Src any, Dst any]() (*ReflectCopier[Src, Dst], error) {
	src := new(Src)
	srcTyp := reflect.TypeOf(src).Elem()
	dst := new(Dst)
	dstTyp := reflect.TypeOf(dst).Elem()
	root := fieldNode{
		isRoot: true,
		isLeaf: false,
		fields: []fieldNode{},
	}
	if srcTyp.Kind() != reflect.Struct {
		return nil, newErrTypeError(srcTyp.Kind())
	}
	if dstTyp.Kind() != reflect.Struct {
		return nil, newErrTypeError(dstTyp.Kind())
	}
	if err := createFiledNodes(&root, srcTyp, dstTyp); err != nil {
		return nil, err
	}

	copier := &ReflectCopier[Src, Dst]{
		rootFiled: root,
	}
	return copier, nil
}

// createFiledNodes 递归创建 field 的前缀树, srcTyp 和 dstTyp 只能是结构体
func createFiledNodes(root *fieldNode, srcTyp, dstTyp reflect.Type) error {

	srcFieldNameID := make(map[string]int, 0)
	for i := 0; i < srcTyp.NumField(); i += 1 {
		fTyp := srcTyp.Field(i)
		if !fTyp.IsExported() {
			continue
		}
		srcFieldNameID[fTyp.Name] = i
	}

	// 构建当前节点的子节点
	for i := 0; i < dstTyp.NumField(); i += 1 {
		dstFieldTypStruct := dstTyp.Field(i)
		if !dstFieldTypStruct.IsExported() {
			continue
		}
		if srcIndex, ok := srcFieldNameID[dstFieldTypStruct.Name]; ok {
			srcFieldTypStruct := srcTyp.Field(srcIndex)
			dstFieldTypStruct := dstTyp.Field(i)

			if srcFieldTypStruct.Type.Kind() != dstFieldTypStruct.Type.Kind() {
				return newErrKindNotMatchError(srcFieldTypStruct.Type.Kind(), dstFieldTypStruct.Type.Kind(), dstFieldTypStruct.Name)
			}

			if srcFieldTypStruct.Type.Kind() == reflect.Pointer {
				if srcFieldTypStruct.Type.Elem().Kind() != dstFieldTypStruct.Type.Elem().Kind() {
					return newErrKindNotMatchError(srcFieldTypStruct.Type.Elem().Kind(), dstFieldTypStruct.Type.Elem().Kind(), dstFieldTypStruct.Name)
				}
				if srcFieldTypStruct.Type.Elem().Kind() == reflect.Pointer {
					return newErrMultiPointer(dstFieldTypStruct.Name)
				}
			}

			child := fieldNode{
				fields:   []fieldNode{},
				srcIndex: srcIndex,
				dstIndex: i,
				isRoot:   false,
				isLeaf:   false,
				name:     dstFieldTypStruct.Name,
			}

			fieldSrcTyp := srcFieldTypStruct.Type
			fieldDstTyp := dstFieldTypStruct.Type
			if fieldSrcTyp.Kind() == reflect.Pointer {
				fieldSrcTyp = fieldSrcTyp.Elem()
				fieldDstTyp = fieldDstTyp.Elem()
			}

			// 说明当前节点是叶子节点, 直接拷贝
			if isShadowCopyType(fieldSrcTyp.Kind()) {
				child.isLeaf = true
			} else if fieldSrcTyp.Kind() == reflect.Struct {
				if err := createFiledNodes(&child, fieldSrcTyp, fieldDstTyp); err != nil {
					return err
				}
			} else {
				// 不是我们能复制的类型, 直接跳过
				continue
			}

			root.fields = append(root.fields, child)
		}
	}
	return nil
}

func (r *ReflectCopier[Src, Dst]) copyToWithTree(src *Src, dst *Dst) error {
	srcTyp := reflect.TypeOf(src).Elem()
	dstTyp := reflect.TypeOf(dst).Elem()
	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()

	root := r.rootFiled

	return r.copyTreeNode(srcTyp, srcValue, dstTyp, dstValue, &root)
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
		child := root.fields[i]
		childSrcTyp := srcTyp.Field(child.srcIndex)
		childSrcValue := srcValue.Field(child.srcIndex)

		childDstTyp := dstType.Field(child.dstIndex)
		childDstValue := dstValue.Field(child.dstIndex)
		if err := r.copyTreeNode(childSrcTyp.Type, childSrcValue, childDstTyp.Type, childDstValue, &child); err != nil {
			return err
		}
	}
	return nil
}

// CopyTo 执行复制
// 执行复制的逻辑是：
// 1. 按照字段的映射关系进行匹配
// 2. 如果 Src 和 Dst 中匹配的字段，其类型是基本类型（及其指针）或者内置类型（及其指针），并且类型一样，则直接用 Src 的值
// 3. 如果 Src 和 Dst 中匹配的字段，其类型都是结构体，或者都是结构体指针，那么会深入复制
// 4. 否则，返回类型不匹配的错误
func (r *ReflectCopier[Src, Dst]) CopyTo(src *Src, dst *Dst) error {
	return r.copyToWithTree(src, dst)
}

// copyWithRuntime 是不使用字典树的复制
func (r *ReflectCopier[Src, Dst]) copyWithRuntime(src *Src, dst *Dst) error {
	srcTyp := reflect.TypeOf(src).Elem()
	if srcTyp.Kind() != reflect.Struct {
		return newErrTypeError(srcTyp.Kind())
	}
	dstTyp := reflect.TypeOf(dst).Elem()
	if dstTyp.Kind() != reflect.Struct {
		return newErrTypeError(dstTyp.Kind())
	}

	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()

	return r.copyStruct(srcTyp, srcValue, dstTyp, dstValue)
}

func (r *ReflectCopier[Src, Dst]) copyStruct(srcTyp reflect.Type, srcValue reflect.Value, dstTyp reflect.Type, dstValue reflect.Value) error {
	srcFieldNameID := make(map[string]int, 0)
	//dstFiledNameID := make(map[string]int32, 0)
	for i := 0; i < srcTyp.NumField(); i += 1 {
		fTyp := srcTyp.Field(i)
		if !fTyp.IsExported() {
			continue
		}
		srcFieldNameID[fTyp.Name] = i
	}

	for i := 0; i < dstTyp.NumField(); i += 1 {
		fTyp := dstTyp.Field(i)
		if !fTyp.IsExported() {
			continue
		}
		if idx, ok := srcFieldNameID[fTyp.Name]; ok {
			if err := r.copyStructField(srcTyp, srcValue, dstTyp, dstValue, idx, i); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *ReflectCopier[Src, Dst]) copyStructField(
	srcTyp reflect.Type,
	srcValue reflect.Value,
	dstTyp reflect.Type,
	dstValue reflect.Value,
	srcFiledIndex int,
	dstFiledIndex int) error {

	srcFieldType := srcTyp.Field(srcFiledIndex)
	dstFieldType := dstTyp.Field(dstFiledIndex)
	if srcFieldType.Type.Kind() != dstFieldType.Type.Kind() {
		return newErrKindNotMatchError(srcFieldType.Type.Kind(), dstFieldType.Type.Kind(), srcFieldType.Name)
	}
	srcFiledValue := srcValue.Field(srcFiledIndex)
	dstFiledValue := dstValue.Field(dstFiledIndex)

	if srcFieldType.Type.Kind() == reflect.Pointer {
		if srcFiledValue.IsNil() {
			return nil
		}
		if dstFiledValue.IsNil() {
			dstFiledValue.Set(reflect.New(dstFieldType.Type.Elem()))
		}
		return r.copyData(srcFieldType.Type.Elem(), srcFiledValue.Elem(), dstFieldType.Type.Elem(), dstFiledValue.Elem(), srcFieldType.Name)
	}

	return r.copyData(srcFieldType.Type, srcFiledValue, dstFieldType.Type, dstFiledValue, srcFieldType.Name)
}

func (r *ReflectCopier[Src, Dst]) copyData(
	srcTyp reflect.Type,
	srcValue reflect.Value,
	dstTyp reflect.Type,
	dstValue reflect.Value,
	fieldName string,
) error {
	if srcTyp.Kind() == reflect.Pointer {
		return newErrMultiPointer(fieldName)
	}
	if srcTyp.Kind() != dstTyp.Kind() {
		return newErrKindNotMatchError(srcTyp.Kind(), dstTyp.Kind(), fieldName)
	}

	if isShadowCopyType(srcTyp.Kind()) {
		if dstValue.CanSet() {
			dstValue.Set(srcValue)
		}
	} else if srcTyp.Kind() == reflect.Struct {
		return r.copyStruct(srcTyp, srcValue, dstTyp, dstValue)
	}
	return nil
}

func (r *ReflectCopier[Src, Dst]) Copy(src *Src) (*Dst, error) {
	dst := new(Dst)
	err := r.CopyTo(src, dst)
	return dst, err
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
