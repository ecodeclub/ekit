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
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
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
			name:      "wrong length key",
			input:     &EncryptColumn[string]{Key: "ABC", Val: "abc", Valid: true},
			wantEnErr: errKeyLengthInvalid,
		},
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
		{
			name:      "complex64",
			input:     &EncryptColumn[complex64]{Key: "ABCDABCDABCDABCD", Val: complex(1, 2), Valid: true},
			wantEnErr: &json.UnsupportedTypeError{Type: reflect.TypeOf(complex64(complex(1, 2)))},
		},
		{
			name:      "complex128",
			input:     &EncryptColumn[complex128]{Key: "ABCDABCDABCDABCD", Val: complex(1, 2), Valid: true},
			wantEnErr: &json.UnsupportedTypeError{Type: reflect.TypeOf(complex(1, 2))},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypt, err := tc.input.(driver.Valuer).Value()
			assert.Equal(t, tc.wantEnErr, err)
			if err == nil {
				err = tc.output.(sql.Scanner).Scan(encrypt)
				assert.Equal(t, tc.wantDeErr, err)
				assert.Equal(t, tc.input, tc.output)
			}
		})
	}
}

func TestEncryptColumn_Sql(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Error(err)
	}

	sqlTable := `
	DROP TABLE IF EXISTS product;
    CREATE TABLE IF NOT EXISTS product(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        encrypt TEXT NOT NULL
    );`

	insertQuery := `INSERT INTO product (id, encrypt) VALUES (1, "13adfdf")`
	updateQuery := `UPDATE product SET encrypt = ? WHERE id = 1`
	selectQuery := `SELECT encrypt FROM product  WHERE id = 1`

	_, err = db.Exec(sqlTable)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec(insertQuery)
	if err != nil {
		t.Error(err)
	}

	key := "ABCDABCDABCDABCD"
	testCases := []struct {
		name      string
		encrypt   any
		decrypt   any
		wantError error
	}{
		{
			name:    "int8",
			encrypt: &EncryptColumn[int8]{Val: int8(123), Valid: true, Key: key},
			decrypt: &EncryptColumn[int8]{Key: key},
		},
		{
			name:    "int16",
			encrypt: &EncryptColumn[int16]{Val: int16(330), Valid: true, Key: key},
			decrypt: &EncryptColumn[int16]{Key: key},
		},
		{
			name:    "int32",
			encrypt: &EncryptColumn[int32]{Val: int32(65550), Valid: true, Key: key},
			decrypt: &EncryptColumn[int32]{Key: key},
		},
		{
			name:    "int64",
			encrypt: &EncryptColumn[int64]{Val: int64(4294967300), Valid: true, Key: key},
			decrypt: &EncryptColumn[int64]{Key: key},
		},
		{
			name:    "uint8",
			encrypt: &EncryptColumn[uint8]{Val: uint8(123), Valid: true, Key: key},
			decrypt: &EncryptColumn[uint8]{Key: key},
		},
		{
			name:    "uint16",
			encrypt: &EncryptColumn[uint16]{Val: uint16(330), Valid: true, Key: key},
			decrypt: &EncryptColumn[uint16]{Key: key},
		},
		{
			name:    "uint32",
			encrypt: &EncryptColumn[uint32]{Val: uint32(65550), Valid: true, Key: key},
			decrypt: &EncryptColumn[uint32]{Key: key},
		},
		{
			name:    "uint64",
			encrypt: &EncryptColumn[uint64]{Val: uint64(4294967300), Valid: true, Key: key},
			decrypt: &EncryptColumn[uint64]{Key: key},
		},
		{
			name:    "int tiny ",
			encrypt: &EncryptColumn[int]{Val: 123, Valid: true, Key: key},
			decrypt: &EncryptColumn[int]{Key: key},
		},
		{
			name:    "int small ",
			encrypt: &EncryptColumn[int]{Val: 1<<16 + 1, Valid: true, Key: key},
			decrypt: &EncryptColumn[int]{Key: key},
		},
		{
			name:    "uint tiny ",
			encrypt: &EncryptColumn[uint]{Val: 123, Valid: true, Key: key},
			decrypt: &EncryptColumn[uint]{Key: key},
		},
		{
			name:    "uint small ",
			encrypt: &EncryptColumn[uint]{Val: 1<<16 + 1, Valid: true, Key: key},
			decrypt: &EncryptColumn[uint]{Key: key},
		},
		{
			name:    "float32",
			encrypt: &EncryptColumn[float32]{Val: float32(123.12), Valid: true, Key: key},
			decrypt: &EncryptColumn[float32]{Key: key},
		},
		{
			name:    "float64",
			encrypt: &EncryptColumn[float64]{Val: 1212321412321323.12222221322, Valid: true, Key: key},
			decrypt: &EncryptColumn[float64]{Key: key},
		},
		{
			name: "map string string",
			encrypt: &EncryptColumn[map[string]string]{Val: map[string]string{
				"A": "B",
				"C": "D",
			}, Valid: true, Key: key},
			decrypt: &EncryptColumn[map[string]string]{Key: key},
		},
		{
			name: "map int string",
			encrypt: &EncryptColumn[map[int]string]{Val: map[int]string{
				1: "B",
				2: "D",
				3: "E",
			}, Valid: true, Key: key},
			decrypt: &EncryptColumn[map[int]string]{Key: key},
		},
		{
			name: "slice string",
			encrypt: &EncryptColumn[[]string]{Val: []string{
				"B",
				"D",
				"E",
			}, Valid: true, Key: key},
			decrypt: &EncryptColumn[[]string]{Key: key},
		},
		{
			name:    "bytes",
			encrypt: &EncryptColumn[[]byte]{Val: []byte("hello"), Valid: true, Key: key},
			decrypt: &EncryptColumn[[]byte]{Key: key},
		},
		{
			name:    "bool",
			encrypt: &EncryptColumn[bool]{Val: true, Valid: true, Key: key},
			decrypt: &EncryptColumn[bool]{Key: key},
		},
		{
			name:    "struct",
			encrypt: &EncryptColumn[Simple]{Val: Simple{"大明", 99}, Valid: true, Key: key},
			decrypt: &EncryptColumn[Simple]{Key: key},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err = db.Exec(updateQuery, tc.encrypt)
			require.Nil(t, err)
			err = db.QueryRow(selectQuery).Scan(tc.decrypt)
			assert.Equal(t, tc.encrypt, tc.decrypt)
		})
	}
}

type Simple struct {
	Name string
	Age  int
}
