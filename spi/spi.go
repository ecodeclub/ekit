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

package spi

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/pkg/errors"
)

// LoadService 加载 dir 下面的所有的实现了 T 接口的类型
// 举个例子来说，如果你有一个叫做 UserService 的接口
// 而后你将所有的实现都放到了 /ext/user_service 目录下
// 并且所有的实现，虽然在不同的包，但是都叫做 UserService
// 那么我可以执行 LoadService("/ext/user_service", "UserService")
// 加载到所有的实现
// LoadService 加载 dir 下面的所有的实现了 T 接口的类型

var (
	ErrDirNotFound        = errors.New("ekit: 目录不存在")
	ErrSymbolNameIsEmpty  = errors.New("ekit: 结构体名不能为空")
	ErrOpenPluginFailed   = errors.New("ekit: 打开插件失败")
	ErrSymbolNameNotFound = errors.New("ekit: 从插件中查找对象失败")
	ErrInvalidSo          = errors.New("ekit: 插件非该接口类型")
)

func LoadService[T any](dir string, symName string) ([]T, error) {
	var services []T
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w", ErrDirNotFound)
	}
	if symName == "" {
		return nil, fmt.Errorf("%w", ErrSymbolNameIsEmpty)
	}
	// 遍历目录下的所有 .so 文件
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".so" {
			// 打开插件
			p, err := plugin.Open(path)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrOpenPluginFailed, err)
			}
			// 查找变量
			sym, err := p.Lookup(symName)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrSymbolNameNotFound, err)
			}

			// 尝试将符号断言为接口类型
			service, ok := sym.(T)
			if !ok {
				return fmt.Errorf("%w", ErrInvalidSo)
			}
			// 收集服务
			services = append(services, service)
		}
		return nil
	})
	return services, err
}
