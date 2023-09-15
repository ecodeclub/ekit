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

package stringx

import (
	"unsafe"
)

// 确保传入的字符串和字节切片的生命周期足够长，不会在转换后被释放或修改。
// 确保传入的字符串和字节切片的长度和容量是一致的，否则可能导致访问越界。
// 不要对转换后的字节切片或字符串进行修改，因为它们可能与原始的字符串或字节切片共享底层的内存。

// UnsafeToBytes 非安全 string 转 []byte 他必须遵守上述规则
func UnsafeToBytes(val string) []byte {
	sh := (*[2]uintptr)(unsafe.Pointer(&val))
	bh := [3]uintptr{sh[0], sh[1], sh[1]}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// UnsafeToString 非安全 []byte 转 string 他必须遵守上述规则
func UnsafeToString(val []byte) string {
	bh := (*[3]uintptr)(unsafe.Pointer(&val))
	sh := [2]uintptr{bh[0], bh[1]}
	return *(*string)(unsafe.Pointer(&sh))
}
