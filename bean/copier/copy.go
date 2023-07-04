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

// Copier 复制数据
// 1. 深拷贝亦或是浅拷贝，取决于具体的实现。每个实现都要声明清楚这一点；
// 2. Src 和 Dst 都必须是普通的结构体，支持组合
// 3. 只复制公共字段
// 这种设计设计，即使用 *Src 和 *Dst 可能加剧内存逃逸
type Copier[Src any, Dst any] interface {
	// CopyTo 将 src 中的数据复制到 dst 中
	CopyTo(src *Src, dst *Dst, opts ...Option) error
	// Copy 将创建一个 Dst 的实例，并且将 Src 中的数据复制过去
	Copy(src *Src, opts ...Option) (*Dst, error)
}

// options 执行复制操作时的可选配置
type options struct {
	// ignoreFields 执行复制操作时，需要忽略的字段
	ignoreFields []string
}

// InIgnoreFields 判断 str 是不是在 ignoreFields 里面
func (r *options) InIgnoreFields(str string) bool {
	if len(r.ignoreFields) < 1 {
		return false
	}
	for _, s := range r.ignoreFields {
		if s == str {
			return true
		}
	}
	return false
}

type Option func(*options)

func IgnoreFields(fields ...string) Option {
	return func(opts *options) {
		opts.ignoreFields = append(opts.ignoreFields, fields...)
	}
}
