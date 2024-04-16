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

package sqlx

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNoMoreRows              = errors.New("ekit: 已读取完")
	errInvalidArgument         = errors.New("ekit: 参数非法")
	_                  Scanner = &sqlRowsScanner{}
)

// Scanner 用于简化sql.Rows包中的Scan操作
// Scanner 不会关闭sql.Rows，用户需要对此负责
type Scanner interface {
	Scan() (values []any, err error)
	// ScanAll 扫描当前结果集的全部数据
	ScanAll() (allValues [][]any, err error)
	// NextResultSet 移动到下一个结果集
	NextResultSet() bool
}

type sqlRowsScanner struct {
	sqlRows             Rows
	columnValuePointers []any
}

// NewSQLRowsScanner 返回一个Scanner
func NewSQLRowsScanner(r Rows) (Scanner, error) {
	if r == nil {
		return nil, fmt.Errorf("%w *sql.Rows不能为nil", errInvalidArgument)
	}
	columnTypes, err := r.ColumnTypes()
	if err != nil || len(columnTypes) < 1 {
		return nil, fmt.Errorf("%w 无法获取*sql.Rows列类型信息: %v", errInvalidArgument, err)
	}
	columnValuePointers := make([]any, len(columnTypes))
	for i, columnType := range columnTypes {
		typ := columnType.ScanType()
		for typ.Kind() == reflect.Pointer {
			// 兼容 sqlite，理论上来说其他 driver 不应该命中这个分支
			typ = typ.Elem()
		}
		columnValuePointers[i] = reflect.New(typ).Interface()
	}
	return &sqlRowsScanner{sqlRows: r, columnValuePointers: columnValuePointers}, nil
}

func (s *sqlRowsScanner) NextResultSet() bool {
	return s.sqlRows.NextResultSet()
}

// Scan 返回一行
func (s *sqlRowsScanner) Scan() ([]any, error) {
	if !s.sqlRows.Next() {
		if err := s.sqlRows.Err(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%w", ErrNoMoreRows)
	}
	err := s.sqlRows.Scan(s.columnValuePointers...)
	if err != nil {
		return nil, err
	}
	return s.columnValues(), nil
}

func (s *sqlRowsScanner) columnValues() []any {
	values := make([]any, len(s.columnValuePointers))
	for i := 0; i < len(s.columnValuePointers); i++ {
		val := reflect.ValueOf(s.columnValuePointers[i]).Elem().Interface()
		// sql.RawBytes 存在内存共享的问题，所以需要执行复制
		if rawBytes, ok := val.(sql.RawBytes); ok {
			val = sql.RawBytes(bytes.Clone(rawBytes))
		}
		values[i] = val
	}
	return values
}

// ScanAll 返回所有行
func (s *sqlRowsScanner) ScanAll() ([][]any, error) {
	all := make([][]any, 0, 32)
	for {
		columnValues, err := s.Scan()
		if err != nil {
			if errors.Is(err, ErrNoMoreRows) {
				break
			}
			return nil, err
		}
		all = append(all, columnValues)
	}
	return all, nil
}
