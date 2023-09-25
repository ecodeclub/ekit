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
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqlRowsScanner_New(t *testing.T) {
	t.Parallel()
	t.Run("当*sql.Rows为nil时，应该报错", func(t *testing.T) {
		t.Parallel()
		_, err := NewSQLRowsScanner(nil)
		require.ErrorIs(t, err, errInvalidArgument)
	})
	t.Run("当无法获取*sql.Rows列类型信息时，应该报错", func(t *testing.T) {
		t.Parallel()
		t.Run("*sql.Rows已关闭", func(t *testing.T) {
			t.Parallel()
			db, err := sql.Open("sqlite3", ":memory:")
			require.NoError(t, err)
			defer db.Close()
			rows, err := db.QueryContext(context.Background(), "")
			require.NoError(t, err)
			require.NoError(t, rows.Close())

			_, err = NewSQLRowsScanner(rows)
			assert.Error(t, err)
		})
		t.Run("*sql.Rows无列类型信息", func(t *testing.T) {
			t.Parallel()
			db, err := sql.Open("sqlite3", ":memory:")
			require.NoError(t, err)
			defer db.Close()
			rows, err := db.QueryContext(context.Background(), "")
			require.NoError(t, err)

			_, err = NewSQLRowsScanner(rows)
			assert.ErrorIs(t, err, errInvalidArgument)
		})
	})
}

func TestSqlRowsScanner_Scan(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test01.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()

	query := "DROP TABLE IF EXISTS t1; CREATE TABLE t1 (\n      id int primary key,\n      `int`  int,\n      `integer` integer,\n      `tinyint` TINYINT,\n      `smallint` smallint,\n      `MEDIUMINT` MEDIUMINT,\n      `BIGINT` BIGINT,\n      `UNSIGNED_BIG_INT` UNSIGNED BIG INT,\n      `INT2` INT2,\n      `INT8` INT8,\n      `VARCHAR` VARCHAR(20),\n  \t\t`CHARACTER` CHARACTER(20),\n  `VARYING_CHARACTER` VARYING_CHARACTER(20),\n  `NCHAR` NCHAR(23),\n  `TEXT` TEXT,\n  `CLOB` CLOB,\n  `REAL` REAL,\n  `DOUBLE` DOUBLE,\n  `DOUBLE_PRECISION` DOUBLE PRECISION,\n  `FLOAT` FLOAT,\n  `DATETIME` DATETIME \n    );"
	_, err = db.ExecContext(context.Background(), query)
	require.NoError(t, err)

	tests := []struct {
		name    string
		rows    *sql.Rows
		want    []any
		cleanup func()
	}{
		{
			name: "浮点类型",
			rows: func() *sql.Rows {
				res, er := db.Exec("INSERT INTO `t1` (`REAL`,`DOUBLE`,`DOUBLE_PRECISION`, `FLOAT`) VALUES (1.0,1.0,1.0,0);")
				require.NoError(t, er)
				id, _ := res.LastInsertId()
				q := "SELECT `REAL`,`DOUBLE`,`DOUBLE_PRECISION`,`FLOAT` FROM `t1` WHERE id=?;"
				rows, er := db.QueryContext(context.Background(), q, id)
				require.NoError(t, er)
				return rows
			}(),
			want: []any{sql.NullFloat64{Valid: true, Float64: 1.0}, sql.NullFloat64{Valid: true, Float64: 1.0}, sql.NullFloat64{Valid: true, Float64: 1.0}, sql.NullFloat64{Valid: true, Float64: 0}},
			cleanup: func() {
				_, er := db.Exec("DELETE FROM `t1`")
				require.NoError(t, er)
			},
		},
		{
			name: "整型",
			rows: func() *sql.Rows {
				res, er := db.Exec("INSERT INTO `t1` (`int`,`integer`,`tinyint`,`smallint`,`MEDIUMINT`,`BIGINT`,`UNSIGNED_BIG_INT`,`INT2`, `INT8`) VALUES (1,1,1,1,1,1,1,1,1);")
				require.NoError(t, er)
				q := "SELECT `int`,`integer`,`tinyint`,`smallint`,`MEDIUMINT`,`BIGINT`,`UNSIGNED_BIG_INT`,`INT2`,`INT8` FROM `t1` WHERE id=?;"
				id, _ := res.LastInsertId()
				rows, er := db.QueryContext(context.Background(), q, id)
				require.NoError(t, er)
				return rows
			}(),
			want: []any{sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: true, Int64: 1}},
			cleanup: func() {
				_, er := db.Exec("DELETE FROM `t1`")
				require.NoError(t, er)
			},
		},
		{
			name: "string类型",
			rows: func() *sql.Rows {
				res, er := db.Exec("INSERT INTO `t1` (`VARCHAR`,`CHARACTER`,`VARYING_CHARACTER`,`NCHAR`,`TEXT`) VALUES ('zwl','zwl','zwl','zwl','zwl');")
				require.NoError(t, er)
				id, _ := res.LastInsertId()
				q := "SELECT `VARCHAR`,`CHARACTER`,`VARYING_CHARACTER`,`NCHAR`,`TEXT`,`CLOB` FROM `t1` WHERE id=?;"
				rows, er := db.QueryContext(context.Background(), q, id)
				require.NoError(t, er)
				return rows
			}(),
			want: []any{sql.NullString{Valid: true, String: "zwl"}, sql.NullString{Valid: true, String: "zwl"}, sql.NullString{Valid: true, String: "zwl"}, sql.NullString{Valid: true, String: "zwl"}, sql.NullString{Valid: true, String: "zwl"}, sql.NullString{Valid: false, String: ""}},
			cleanup: func() {
				_, er := db.Exec("DELETE FROM `t1`")
				require.NoError(t, er)
			},
		},
		{
			name: "时间类型",
			rows: func() *sql.Rows {
				res, er := db.Exec("INSERT INTO `t1` (`DATETIME`) VALUES ('2022-01-01 12:00:00');")
				require.NoError(t, er)
				id, _ := res.LastInsertId()
				q := "SELECT `DATETIME` FROM `t1` WHERE id=?;"
				rows, er := db.QueryContext(context.Background(), q, id)
				require.NoError(t, er)
				return rows
			}(),
			want: []any{sql.NullTime{Valid: true, Time: func() time.Time {
				tim, er := time.ParseInLocation("2006-01-02 15:04:05", "2022-01-01 12:00:00", time.Local)
				require.NoError(t, er)
				return tim
			}()}},
			cleanup: func() {
				_, er := db.Exec("DELETE FROM `t1`")
				require.NoError(t, er)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewSQLRowsScanner(tt.rows)
			require.NoError(t, err)
			for {
				got, err := s.Scan()
				if err != nil && errors.Is(err, ErrNoMoreRows) {
					break
				}
				assert.NoError(t, err)
				assert.Equalf(t, tt.want, got, "ScanRows(%v)", tt.rows)
			}
			tt.cleanup()
		})
	}

	t.Run("迭代期间sql.Rows发生错误,Scan应该报错", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		expectedErr := errors.New("iteration error")
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "John").
			AddRow(2, "Jane").RowError(1, expectedErr))

		rows, err := db.Query("SELECT id, name FROM users")
		require.NoError(t, err)
		defer rows.Close()

		s, err := NewSQLRowsScanner(rows)
		require.NoError(t, err)

		values, err := s.Scan()
		assert.NoError(t, err)
		assert.Equal(t, []any{int64(1), "John"}, values)

		_, err = s.Scan()
		assert.Equal(t, expectedErr, err)
	})
}

func TestSqlRowsScanner_ScanAll(t *testing.T) {
	t.Parallel()
	t.Run("迭代期间sql.Rows没有错误,ScanAll正常结束", func(t *testing.T) {
		t.Parallel()
		db, err := sql.Open("sqlite3", "file:test01.db?cache=shared&mode=memory")
		require.NoError(t, err)
		defer db.Close()

		query := "DROP TABLE IF EXISTS t1; CREATE TABLE t1 " +
			"(id int primary key," +
			"`name` VARCHAR(20), " +
			"`intro` TEXT, " +
			"`create_time` DATETIME);"
		_, err = db.ExecContext(context.Background(), query)
		require.NoError(t, err)

		t1, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-02-01 19:00:01", time.UTC)
		t2, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-04-01 11:00:00", time.UTC)
		t3, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-02-02 09:00:23", time.UTC)
		t4, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-02-04 15:00:00", time.UTC)

		_, err = db.Exec("INSERT INTO `t1` (`id`, `name`, `intro`, `create_time`) VALUES " +
			"(1, 'zhangsan','这是一段中文介绍', \"2023-02-01 19:00:01\"), " +
			"(2, 'lisi','这是一段中文介绍', \"2023-04-01 11:00:00\"), " +
			"(3, 'wangwu','this is English introduction', \"2023-02-02 09:00:23\"), " +
			"(4, 'zhaoliu','this is English introduction', \"2023-02-04 15:00:00\");")
		require.NoError(t, err)

		expected := [][]any{
			{sql.NullInt64{Valid: true, Int64: 1}, sql.NullString{Valid: true, String: "zhangsan"}, sql.NullString{Valid: true, String: "这是一段中文介绍"}, sql.NullTime{Valid: true, Time: t1}},
			{sql.NullInt64{Valid: true, Int64: 2}, sql.NullString{Valid: true, String: "lisi"}, sql.NullString{Valid: true, String: "这是一段中文介绍"}, sql.NullTime{Valid: true, Time: t2}},
			{sql.NullInt64{Valid: true, Int64: 3}, sql.NullString{Valid: true, String: "wangwu"}, sql.NullString{Valid: true, String: "this is English introduction"}, sql.NullTime{Valid: true, Time: t3}},
			{sql.NullInt64{Valid: true, Int64: 4}, sql.NullString{Valid: true, String: "zhaoliu"}, sql.NullString{Valid: true, String: "this is English introduction"}, sql.NullTime{Valid: true, Time: t4}},
		}

		rows, err := db.QueryContext(context.Background(), "SELECT * FROM `t1`;")
		require.NoError(t, err)
		defer rows.Close()

		s, err := NewSQLRowsScanner(rows)
		require.NoError(t, err)

		actual, err := s.ScanAll()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("迭代期间sql.Rows发生错误,ScanAll应该报错", func(t *testing.T) {
		t.Parallel()
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		expectedErr := errors.New("iteration error")

		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "John").
			AddRow(2, "Jane").RowError(1, expectedErr))

		rows, err := db.Query("SELECT id, name FROM users")
		require.NoError(t, err)
		defer rows.Close()

		s, err := NewSQLRowsScanner(rows)
		require.NoError(t, err)

		_, err = s.ScanAll()
		assert.Equal(t, expectedErr, err)
	})
}

func TestSqlRowsScanner_NextResultSet(t *testing.T) {
	t.Parallel()
	t.Run("没有更多 ResultSet", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		mock.ExpectQuery("SELECT .*").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
		rows, err := db.Query("SELECT id, name FROM users")
		require.NoError(t, err)
		scanner, err := NewSQLRowsScanner(rows)
		require.NoError(t, err)
		assert.False(t, scanner.NextResultSet())
	})
	t.Run("还有 ResultSet", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		mock.ExpectQuery("SELECT .*").
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name"}),
				sqlmock.NewRows([]string{"id", "name"}),
				sqlmock.NewRows([]string{"id", "name"}))
		rows, err := db.Query("SELECT id, name FROM users")
		require.NoError(t, err)
		scanner, err := NewSQLRowsScanner(rows)
		require.NoError(t, err)
		assert.True(t, scanner.NextResultSet())
		assert.True(t, scanner.NextResultSet())
		assert.False(t, scanner.NextResultSet())
	})
}
