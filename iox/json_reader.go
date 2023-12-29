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

package iox

import (
	"bytes"
	"encoding/json"
)

type JSONReader struct {
	val any
	bf  *bytes.Reader
}

func (j *JSONReader) Read(p []byte) (n int, err error) {
	if j.bf == nil {
		var data []byte
		data, err = json.Marshal(j.val)
		if err == nil {
			j.bf = bytes.NewReader(data)
		}
	}
	if err != nil {
		return
	}
	return j.bf.Read(p)
}

// NewJSONReader 用于解决将一个结构体序列化为 JSON 之后，再封装为 io.Reader 的场景。
// 该实现没有做任何输入检查。
// 也就是你需要自己确保 val 是一个可以被 json 正确处理的东西。
// 非线程安全。
// 如果你传入的是 nil，那么读到的结果应该是 null。务必小心。
func NewJSONReader(val any) *JSONReader {
	return &JSONReader{
		val: val,
	}
}
