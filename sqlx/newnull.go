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
	"time"
)

//这一个系列的方法，会在数据为零值时，将valid设置为false；否则设置为true；

func NewNullString(val string) sql.NullString {
	return sql.NullString{String: val, Valid: val != ""}
}

func NewNullInt64(val int64) sql.NullInt64 {
	return sql.NullInt64{Int64: val, Valid: val != 0}
}

func NewNullFloat64(val float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: val, Valid: val != 0}
}

func NewNullBool(val bool) sql.NullBool {
	return sql.NullBool{Bool: val, Valid: val}
}

func NewNullTime(val time.Time) sql.NullTime {
	return sql.NullTime{Time: val, Valid: !val.IsZero()}
}

func NewNullBytes(val []byte) sql.NullString {
	return sql.NullString{String: string(val), Valid: len(val) > 0}
}
