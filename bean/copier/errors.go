package copier

import (
	"fmt"
	"reflect"
)

// newErrTypeError copier 不支持的类型
func newErrTypeError(kind reflect.Kind) error {
	return fmt.Errorf("ekit: copier 入口只支持 Struct 不支持类型 %v", kind)
}

// newErrKindNotMatchError 字段类型不匹配
func newErrKindNotMatchError(src, dst reflect.Kind, field string) error {
	return fmt.Errorf("ekit: 字段 %s 的 Kind 不匹配, src: %v, dst: %v", field, src, dst)
}

// newErrMultiPointer
func newErrMultiPointer(field string) error {
	return fmt.Errorf("ekit: 字段 %s 是多级指针", field)
}
