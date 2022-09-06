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

package slice

import (
	"log"
	"runtime"
)

// EqualFunc 比较两个元素是否相等
type EqualFunc[T any] func(src, dst T) bool

func (e EqualFunc[any]) safeEqual(src, dst any) (isPanic bool, result bool) {
	defer func() {
		if p := recover(); p != nil {
			isPanic = true
			log.Printf("ekit [PANIC]: %v", p)

			// 打印调用栈信息
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			log.Printf("ekit: panic stack info %s", buf[:n])
		}
	}()
	result = e(src, dst)
	return
}
