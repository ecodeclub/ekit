package sqlx

import (
	"database/sql"
	"database/sql/driver"
	"testing"

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
		name    string
		encrypt any
		decrypt any
	}{
		{
			name:    "int (32位)",
			encrypt: &EncryptColumn[int32]{Val: int32(123), Valid: true, Key: key},
			decrypt: &EncryptColumn[int32]{Key: key},
		},
		{
			name:    "float (32位)",
			encrypt: &EncryptColumn[float32]{Val: float32(123.12), Valid: true, Key: key},
			decrypt: &EncryptColumn[float32]{Key: key},
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
			if err != nil {
				t.Error(err)
			}
			err = db.QueryRow(selectQuery).Scan(tc.decrypt)
			assert.Equal(t, tc.encrypt, tc.decrypt)
		})
	}
}

type Simple struct {
	Name string
	Age  int
}
