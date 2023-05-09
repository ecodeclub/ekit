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
	"database/sql"
	"reflect"
)

func ScanRows(rows *sql.Rows) ([]any, error) {
	colsInfo, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	colsData := make([]any, 0, len(colsInfo))
	for _, colInfo := range colsInfo {
		typ := colInfo.ScanType()
		// 保险起见，循环的去除指针
		for typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}
		newData := reflect.New(typ).Interface()
		colsData = append(colsData, newData)
	}
	err = rows.Scan(colsData...)
	if err != nil {
		return nil, err
	}
	// 去掉reflect.New的指针
	for i := 0; i < len(colsData); i++ {
		colsData[i] = reflect.ValueOf(colsData[i]).Elem().Interface()
	}
	return colsData, nil
}

func ScanAll(rows *sql.Rows) ([][]any, error) {
	res := make([][]any, 0, 32)
	for rows.Next() {
		cols, err := ScanRows(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, cols)
	}
	return res, nil
}
