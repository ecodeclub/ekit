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

import "reflect"

// CopyTo 复制结构体, 纯递归实现. src 和 dst 都必须是结构体的指针
func CopyTo(src any, dst any) error {
	srcPtrTyp := reflect.TypeOf(src)
	if srcPtrTyp.Kind() != reflect.Pointer {
		return newErrTypeError(srcPtrTyp)
	}
	srcTyp := srcPtrTyp.Elem()
	if srcTyp.Kind() != reflect.Struct {
		return newErrTypeError(srcTyp)
	}
	dstPtrTyp := reflect.TypeOf(dst)
	if dstPtrTyp.Kind() != reflect.Pointer {
		return newErrTypeError(dstPtrTyp)
	}
	dstTyp := dstPtrTyp.Elem()
	if dstTyp.Kind() != reflect.Struct {
		return newErrTypeError(dstTyp)
	}

	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()

	return copyStruct(srcTyp, srcValue, dstTyp, dstValue)
}

func copyStruct(srcTyp reflect.Type, srcValue reflect.Value, dstTyp reflect.Type, dstValue reflect.Value) error {
	srcFieldNameIndex := make(map[string]int, 0)
	for i := 0; i < srcTyp.NumField(); i += 1 {
		fTyp := srcTyp.Field(i)
		if !fTyp.IsExported() {
			continue
		}
		srcFieldNameIndex[fTyp.Name] = i
	}

	for i := 0; i < dstTyp.NumField(); i += 1 {
		fTyp := dstTyp.Field(i)
		if !fTyp.IsExported() {
			continue
		}
		if idx, ok := srcFieldNameIndex[fTyp.Name]; ok {
			if err := copyStructField(srcTyp, srcValue, dstTyp, dstValue, idx, i); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyStructField(
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
		return copyData(srcFieldType.Type.Elem(), srcFiledValue.Elem(), dstFieldType.Type.Elem(), dstFiledValue.Elem(), srcFieldType.Name)
	}

	return copyData(srcFieldType.Type, srcFiledValue, dstFieldType.Type, dstFiledValue, srcFieldType.Name)
}

func copyData(
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
		return copyStruct(srcTyp, srcValue, dstTyp, dstValue)
	}
	return nil
}
