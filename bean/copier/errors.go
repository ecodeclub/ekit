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
	"errors"
	"fmt"
	"reflect"
)

var (
	errConvertFieldTypeNotMatch = errors.New("ekit: 转化字段类型不匹配")
)

// newErrTypeError copier 不支持的类型
func newErrTypeError(typ reflect.Type) error {
	return fmt.Errorf("ekit: copier 入口只支持 Struct 不支持类型 %v, 种类 %v", typ, typ.Kind())
}

// newErrKindNotMatchError 字段类型不匹配
func newErrKindNotMatchError(src, dst reflect.Kind, field string) error {
	return fmt.Errorf("ekit: 字段 %s 的 Kind 不匹配, src: %v, dst: %v", field, src, dst)
}

// newErrTypeNotMatchError 字段不匹配
func newErrTypeNotMatchError(src, dst reflect.Type, field string) error {
	return fmt.Errorf("ekit: 字段 %s 的 Type 不匹配, src: %v, dst: %v", field, src, dst)
}

// newErrMultiPointer
func newErrMultiPointer(field string) error {
	return fmt.Errorf("ekit: 字段 %s 是多级指针", field)
}
