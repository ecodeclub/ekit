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
	"github.com/ecodeclub/ekit/bean/copier/converter"
	"github.com/ecodeclub/ekit/bean/option"
	"github.com/ecodeclub/ekit/set"
)

// Copier 复制数据
// 1. 深拷贝亦或是浅拷贝，取决于具体的实现。每个实现都要声明清楚这一点；
// 2. Src 和 Dst 都必须是普通的结构体，支持组合
// 3. 只复制公共字段
// 这种设计设计，即使用 *Src 和 *Dst 可能加剧内存逃逸
type Copier[Src any, Dst any] interface {
	// CopyTo 将 src 中的数据复制到 dst 中
	CopyTo(src *Src, dst *Dst, opts ...option.Option[options]) error
	// Copy 将创建一个 Dst 的实例，并且将 Src 中的数据复制过去
	Copy(src *Src, opts ...option.Option[options]) (*Dst, error)
}

// options 执行复制操作时的可选配置
type options struct {
	// ignoreFields 执行复制操作时，需要忽略的字段
	ignoreFields *set.MapSet[string]
	// convertFields 执行转换的field和转化接口的泛型包装
	convertFields map[string]converterWrapper
}

type converterWrapper func(src any) (any, error)

func newOptions() options {
	return options{}
}

// InIgnoreFields 判断 str 是不是在 ignoreFields 里面
func (r *options) InIgnoreFields(str string) bool {
	// 如果没有设置过忽略的字段的话，ignoreFields 就有可能是 nil，这里需要判断一下
	if r.ignoreFields == nil {
		return false
	}
	return r.ignoreFields.Exist(str)
}

// IgnoreFields 设置复制时要忽略的字段（option 设计模式）
func IgnoreFields(fields ...string) option.Option[options] {
	return func(opt *options) {
		if len(fields) < 1 {
			return
		}
		// 需要用的时候再延迟初始化 ignoreFields
		if opt.ignoreFields == nil {
			opt.ignoreFields = set.NewMapSet[string](len(fields))
		}
		for i := 0; i < len(fields); i++ {
			opt.ignoreFields.Add(fields[i])
		}
	}
}

func ConvertField[Src any, Dst any](field string, converter converter.Converter[Src, Dst]) option.Option[options] {
	return func(opt *options) {
		if field == "" || converter == nil {
			return
		}
		if opt.convertFields == nil {
			opt.convertFields = make(map[string]converterWrapper, 8)
		}
		opt.convertFields[field] = func(src any) (any, error) {
			var dst Dst
			srcVal, ok := src.(Src)
			if !ok {
				return dst, errConvertFieldTypeNotMatch
			}
			return converter.Convert(srcVal)
		}
	}
}
