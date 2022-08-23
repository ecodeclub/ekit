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

// ReflectCopier 基于反射的实现
// ReflectCopier 是浅拷贝
type ReflectCopier[Src any, Dst any] struct {
	// fieldMap Src 中字段下标到 Dst 字段下标的映射
	// 其中 key 是 Src 中字段下标使用连字符 - 连接起来
	// 不在 fieldMap 中字段则意味着被忽略
	fieldMap map[string][]int
}

func NewReflectCopier[Src any, Dst any]() *ReflectCopier[Src, Dst] {
	panic("implement me")
}

// CopyTo 执行复制
// 执行复制的逻辑是：
// 1. 按照字段的映射关系进行匹配
// 2. 如果 Src 和 Dst 中匹配的字段，其类型是基本类型（及其指针）或者内置类型（及其指针），并且类型一样，则直接用 Src 的值
// 3. 如果 Src 和 Dst 中匹配的字段，其类型都是结构体，或者都是结构体指针，那么会深入复制
// 4. 否则，返回类型不匹配的错误
// TODO: 支持不同类型之间的转换
func (r *ReflectCopier[Src, Dst]) CopyTo(src *Src, dst *Dst) error {
	panic("implement me")
}

func (r *ReflectCopier[Src, Dst]) Copy(src *Src) (*Dst, error) {
	dst := new(Dst)
	err := r.CopyTo(src, dst)
	return dst, err
}
