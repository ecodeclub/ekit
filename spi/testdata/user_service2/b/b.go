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

//go:generate go build -race --buildmode=plugin   -o ../b.so ./b.go
package main

// 测试用

type UserService struct{}

// GetName returns the name of the service
func (u UserService) Get() string {
	return "B"
}

// 导出对象
var UserSvc UserService
