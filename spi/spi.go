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

// LoadService 加载 dir 下面的所有的实现了 T 接口的类型
// 举个例子来说，如果你有一个叫做 UserService 的接口
// 而后你将所有的实现都放到了 /ext/user_service 目录下
// 并且所有的实现，虽然在不同的包，但是都叫做 UserService
// 那么我可以执行 LoadService("/ext/user_service", "UserService")
// 加载到所有的实现
func LoadService[T any](dir string, symName string) ([]T, error) {
	panic("implement me")
}
