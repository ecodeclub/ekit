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

package sql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JsonColumn 代表存储字段的 json 类型
// 主要用于没有提供默认 json 类型的数据库
// T 可以是结构体，也可以是切片或者 map
// 一切可以被 json 库所处理的类型都能被用作 T
type JsonColumn[T any] struct {
	Val   T
	Valid bool
}

// Value 返回一个 json 串。类型是 []byte
func (j JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return json.Marshal(j.Val)
}

// Scan 将 src 转化为对象
// src 的类型必须是 []byte, *[]byte, string, sql.RawBytes, *sql.RawBytes 之一
func (j *JsonColumn[T]) Scan(src any) error {
	var bs []byte
	switch val := src.(type) {
	case []byte:
		bs = val
	case *[]byte:
		bs = *val
	case string:
		bs = []byte(val)
	case sql.RawBytes:
		bs = val
	case *sql.RawBytes:
		bs = *val
	default:
		return fmt.Errorf("ekit：JsonColumn.Scan 不支持 src 类型 %v", src)
	}

	if err := json.Unmarshal(bs, &j.Val); err != nil {
		return err
	}
	j.Valid = true
	return nil
}
