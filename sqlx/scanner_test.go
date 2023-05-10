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
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanRows(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test01.db?cache=shared&mode=memory")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	query := "DROP TABLE IF EXISTS t1; CREATE TABLE t1 (\n      id int primary key,\n      `int`  int,\n      `integer` integer,\n      `tinyint` TINYINT,\n      `smallint` smallint,\n      `MEDIUMINT` MEDIUMINT,\n      `BIGINT` BIGINT,\n      `UNSIGNED_BIG_INT` UNSIGNED BIG INT,\n      `INT2` INT2,\n      `INT8` INT8,\n      `VARCHAR` VARCHAR(20),\n  \t\t`CHARACTER` CHARACTER(20),\n  `VARYING_CHARACTER` VARYING_CHARACTER(20),\n  `NCHAR` NCHAR(23),\n  `TEXT` TEXT,\n  `CLOB` CLOB,\n  `REAL` REAL,\n  `DOUBLE` DOUBLE,\n  `DOUBLE_PRECISION` DOUBLE PRECISION,\n  `FLOAT` FLOAT,\n  `DATETIME` DATETIME \n    );"
	_, err = db.ExecContext(context.Background(), query)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		rows  *sql.Rows
		want  []any
		after func()
	}{
		{
			name: "浮点类型",
			rows: func() *sql.Rows {
				res, er := db.Exec("INSERT INTO `t1` (`REAL`,`DOUBLE`,`DOUBLE_PRECISION`, `FLOAT`) VALUES (1.0,1.0,1.0,0);")
				require.NoError(t, er)
				id, _ := res.LastInsertId()
				q := "SELECT `REAL`,`DOUBLE`,`DOUBLE_PRECISION`,`FLOAT` FROM `t1` where id=?;"
				rows, er := db.QueryContext(context.Background(), q, id)
				require.NoError(t, er)
				return rows
			}(),
			want: []any{sql.NullFloat64{Valid: true, Float64: 1.0}, sql.NullFloat64{Valid: true, Float64: 1.0}, sql.NullFloat64{Valid: true, Float64: 1.0}, sql.NullFloat64{Valid: true, Float64: 0}},
			after: func() {
				_, er := db.Exec("delete from `t1`")
				require.NoError(t, er)
			},
		},
		{
			name: "整型",
			rows: func() *sql.Rows {
				res, er := db.Exec("INSERT INTO `t1` (`int`,`integer`,`tinyint`,`smallint`,`MEDIUMINT`,`BIGINT`,`UNSIGNED_BIG_INT`,`INT2`, `INT8`) VALUES (1,1,1,1,1,1,1,1,1);")
				require.NoError(t, er)
				q := "SELECT `int`,`integer`,`tinyint`,`smallint`,`MEDIUMINT`,`BIGINT`,`UNSIGNED_BIG_INT`,`INT2`,`INT8` FROM `t1` where id=?;"
				id, _ := res.LastInsertId()
				rows, er := db.QueryContext(context.Background(), q, id)
				require.NoError(t, er)
				return rows
			}(),
			want: []any{sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}},
			after: func() {
				_, er := db.Exec("delete from `t1`")
				require.NoError(t, er)
			},
		},
		{
			name: "string类型",
			rows: func() *sql.Rows {
				res, er := db.Exec("INSERT INTO `t1` (`VARCHAR`,`CHARACTER`,`VARYING_CHARACTER`,`NCHAR`,`TEXT`) VALUES ('zwl','zwl','zwl','zwl','zwl');")
				require.NoError(t, er)
				id, _ := res.LastInsertId()
				q := "SELECT `VARCHAR`,`CHARACTER`,`VARYING_CHARACTER`,`NCHAR`,`TEXT`,`CLOB` FROM `t1` where id=?;"
				rows, er := db.QueryContext(context.Background(), q, id)
				require.NoError(t, er)
				return rows
			}(),
			want: []any{sql.NullString{Valid: true, String: "zwl"}, sql.NullString{Valid: true, String: "zwl"}, sql.NullString{Valid: true, String: "zwl"}, sql.NullString{Valid: true, String: "zwl"}, sql.NullString{Valid: true, String: "zwl"}, sql.NullString{Valid: false, String: ""}},
			after: func() {
				_, er := db.Exec("delete from `t1`")
				require.NoError(t, er)
			},
		},
		{
			name: "时间类型",
			rows: func() *sql.Rows {
				res, er := db.Exec("INSERT INTO `t1` (`DATETIME`) VALUES ('2022-01-01 12:00:00');")
				require.NoError(t, er)
				id, _ := res.LastInsertId()
				q := "SELECT `DATETIME` FROM `t1` where id=?;"
				rows, er := db.QueryContext(context.Background(), q, id)
				require.NoError(t, er)
				return rows
			}(),
			want: []any{sql.NullTime{Valid: true, Time: func() time.Time {
				tim, er := time.ParseInLocation("2006-01-02 15:04:05", "2022-01-01 12:00:00", time.Local)
				require.NoError(t, er)
				return tim
			}()}},
			after: func() {
				_, er := db.Exec("delete from `t1`")
				require.NoError(t, er)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := tt.rows
			for rows.Next() {
				got, er := ScanRows(rows)
				assert.Nil(t, er)
				assert.Equalf(t, tt.want, got, "ScanRows(%v)", tt.rows)
			}
			tt.after()
		})
	}
}

func Test_ScanAll(t *testing.T) {
	db, err1 := sql.Open("sqlite3", "file:test01.db?cache=shared&mode=memory")
	if err1 != nil {
		t.Error(err1)
	}
	defer db.Close()

	query := "DROP TABLE IF EXISTS t1; CREATE TABLE t1 " +
		"(id int primary key," +
		"`name` VARCHAR(20), " +
		"`intro` TEXT, " +
		"`create_time` DATETIME);"
	_, err2 := db.ExecContext(context.Background(), query)
	if err2 != nil {
		t.Fatal(err2)
	}

	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-02-01 19:00:01", time.UTC)
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-04-01 11:00:00", time.UTC)
	t3, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-02-02 09:00:23", time.UTC)
	t4, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-02-04 15:00:00", time.UTC)

	_, err3 := db.Exec("INSERT INTO `t1` (`id`, `name`, `intro`, `create_time`) VALUES " +
		"(1, 'zhangsan','这是一段中文介绍', \"2023-02-01 19:00:01\"), " +
		"(2, 'lisi','这是一段中文介绍', \"2023-04-01 11:00:00\"), " +
		"(3, 'wangwu','this is English introduction', \"2023-02-02 09:00:23\"), " +
		"(4, 'zhaoliu','this is English introduction', \"2023-02-04 15:00:00\");")
	assert.Nil(t, err3)

	rows, err4 := db.QueryContext(context.Background(), "SELECT * FROM `t1`;")
	assert.Nil(t, err4)

	got, err5 := ScanAll(rows)
	assert.Nil(t, err5)

	assert.Equal(t, [][]any{
		{sql.NullInt64{Valid: true, Int64: 1}, sql.NullString{Valid: true, String: "zhangsan"}, sql.NullString{Valid: true, String: "这是一段中文介绍"}, sql.NullTime{Valid: true, Time: t1}},
		{sql.NullInt64{Valid: true, Int64: 2}, sql.NullString{Valid: true, String: "lisi"}, sql.NullString{Valid: true, String: "这是一段中文介绍"}, sql.NullTime{Valid: true, Time: t2}},
		{sql.NullInt64{Valid: true, Int64: 3}, sql.NullString{Valid: true, String: "wangwu"}, sql.NullString{Valid: true, String: "this is English introduction"}, sql.NullTime{Valid: true, Time: t3}},
		{sql.NullInt64{Valid: true, Int64: 4}, sql.NullString{Valid: true, String: "zhaoliu"}, sql.NullString{Valid: true, String: "this is English introduction"}, sql.NullTime{Valid: true, Time: t4}},
	}, got)
}
