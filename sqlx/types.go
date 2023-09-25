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

import "database/sql"

// 因为 sql 包里面缺乏顶级接口定义，而在研发一些中间件的时候，又必须用到不同的实现
// 因此在这里提前定义一些顶级接口
// 一般来说，如果你不是设计一些和数据库有关的中间件，你是用不上这些接口的

var _ Rows = (*sql.Rows)(nil)

type Rows interface {
	Next() bool
	NextResultSet() bool
	Err() error
	Columns() ([]string, error)
	// ColumnTypes 还是返回了原本的 sql.ColumnType
	// 因为 ColumnType 同样不是一个接口，所以为了兼容 sql.Rows，
	// 就只有保持这个设计
	ColumnTypes() ([]*sql.ColumnType, error)
	Scan(dest ...any) error
	Close() error
}
