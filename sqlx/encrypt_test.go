package sqlx

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptColumn_Basic(t *testing.T) {

	testCases := []struct {
		name      string
		input     any // 因为泛型的缘故我们这里只能使用 any
		output    any
		wantEnErr error
		wantDeErr error
	}{
		{
			name:   "int",
			input:  &EncryptColumn[int32]{Key: "ABCDABCDABCDABCDABCDABCDABCDABCD", Val: 123, Valid: true},
			output: &EncryptColumn[int32]{Key: "ABCDABCDABCDABCDABCDABCDABCDABCD"},
		},
		{
			name:   "int",
			input:  &EncryptColumn[int]{Key: "ABCDABCDABCDABCD", Val: 123, Valid: true},
			output: &EncryptColumn[int]{Key: "ABCDABCDABCDABCD"},
		},
		{
			name:   "string",
			input:  &EncryptColumn[string]{Key: "ABCDABCDABCDABCD", Val: "adsnfjkenfjkndjsknfjenjfknsadnfkjejfn", Valid: true},
			output: &EncryptColumn[string]{Key: "ABCDABCDABCDABCD"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypt, err := tc.input.(driver.Valuer).Value()
			assert.Equal(t, tc.wantEnErr, err)
			err = tc.output.(sql.Scanner).Scan(encrypt)
			assert.Equal(t, tc.wantDeErr, err)
			assert.Equal(t, tc.input, tc.output)
		})
	}
}

func TestEncryptColumn_Sql(t *testing.T) {
	db, err := sql.Open("sqlite3", "./sqlx_test.db")
	if err != nil {
		t.Error(err)
	}

	sqlTable := `
    CREATE TABLE IF NOT EXISTS product(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sale TEXT NOT NULL UNIQUE, 
        num  TEXT NOT NULL,
        json TEXT NOT NULL
    );`

	insertQuery := `INSERT INTO product (id, sale, num, json) VALUES (?, ?, ?, ?)`
	updateQuery := `UPDATE product SET %s = ? WHERE id = 1`
	selectQuery := `SELECT %s FROM product  WHERE id = 1`

	stmt, err := db.Prepare(sqlTable)
	res, err := stmt.Exec()
	_, err = res.RowsAffected()
	if err != nil {
		return
	}
	stmt, err = db.Prepare(insertQuery)
	_, err = stmt.Exec("1232d", "ddd3sdf", "ec3adc", "123df")

	key := "ABCDABCDABCDABCD"

	testCases := []struct {
		name    string
		valName string
		encrypt any
		decrypt any
	}{
		{
			name:    "int (32位)",
			valName: "num",
			encrypt: &EncryptColumn[int32]{Val: int32(123), Valid: true, Key: key},
			decrypt: &EncryptColumn[int32]{Key: key},
		},
		{
			name:    "float (32位)",
			valName: "sale",
			encrypt: &EncryptColumn[float32]{Val: float32(123.12), Valid: true, Key: key},
			decrypt: &EncryptColumn[float32]{Key: key},
		},
		{
			name:    "int tiny ",
			valName: "num",
			encrypt: &EncryptColumn[int]{Val: 123, Valid: true, Key: key},
			decrypt: &EncryptColumn[int]{Key: key},
		},
		{
			name:    "int small ",
			valName: "num",
			encrypt: &EncryptColumn[int]{Val: 1<<16 + 1, Valid: true, Key: key},
			decrypt: &EncryptColumn[int]{Key: key},
		},
		{
			name:    "int mid ",
			valName: "num",
			encrypt: &EncryptColumn[int64]{Val: 1<<32 + 1, Valid: true, Key: key}, //虽然结构体不完全相等，但在数值上是一样的
			decrypt: &EncryptColumn[int]{Key: key},
		},
		{
			name:    "int big",
			valName: "num",
			encrypt: &EncryptColumn[int64]{Val: 1<<63 - 1, Valid: true, Key: key}, //虽然结构体不完全相等，但在数值上是一样的
			decrypt: &EncryptColumn[int]{Key: key},
		},
		{
			name:    "int huge",
			valName: "num",
			encrypt: &EncryptColumn[int64]{Val: 1<<63 - 1, Valid: true, Key: key}, //虽然结构体不完全相等，但在数值上是一样的
			decrypt: &EncryptColumn[int]{Key: key},
		},
		{
			name:    "int huge",
			valName: "num",
			encrypt: &EncryptColumn[int64]{Val: 1<<63 - 1, Valid: true, Key: key}, //虽然结构体不完全相等，但在数值上是一样的
			decrypt: &EncryptColumn[int]{Key: key},
		},
		{
			name:    "bool",
			valName: "json",
			encrypt: &EncryptColumn[bool]{Val: true, Valid: true, Key: key}, //虽然结构体不完全相等，但在数值上是一样的
			decrypt: &EncryptColumn[bool]{Key: key},
		},
		{
			name:    "struct",
			valName: "json",
			encrypt: &EncryptColumn[Simple]{Val: Simple{"大明", 99}, Valid: true, Key: key}, //虽然结构体不完全相等，但在数值上是一样的
			decrypt: &EncryptColumn[Simple]{Key: key},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateQ := fmt.Sprintf(updateQuery, tc.valName)
			_, err = db.Exec(updateQ, tc.encrypt)
			if err != nil {
				t.Error(err)
			}
			selectQ := fmt.Sprintf(selectQuery, tc.valName)
			err = db.QueryRow(selectQ).Scan(tc.decrypt)
			assert.Equal(t, tc.encrypt, tc.decrypt)
		})
	}
}

type Simple struct {
	Name string
	Age  int
}
