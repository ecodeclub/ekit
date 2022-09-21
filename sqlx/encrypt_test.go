package sqlx

import (
	"database/sql"
	"database/sql/driver"
	"github.com/stretchr/testify/assert"
	"reflect"
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
			input:  EncryptColumn[int32]{Key: "ABCDABCDABCDABCD", Val: 123, Valid: true},
			output: &EncryptColumn[int32]{Key: "ABCDABCDABCDABCD"},
		},
		{
			name:   "int",
			input:  EncryptColumn[int]{Key: "ABCDABCDABCDABCD", Val: 123, Valid: true},
			output: &EncryptColumn[int]{Key: "ABCDABCDABCDABCD"},
		},
		{
			name:   "string",
			input:  EncryptColumn[string]{Key: "ABCDABCDABCDABCD", Val: "adsnfjkenfjkndjsknfjenjfknsadnfkjejfn", Valid: true},
			output: &EncryptColumn[string]{Key: "ABCDABCDABCDABCD"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypt, err := tc.input.(driver.Valuer).Value()
			assert.Equal(t, tc.wantEnErr, err)
			err = tc.output.(sql.Scanner).Scan(encrypt)
			assert.Equal(t, tc.wantDeErr, err)
			assert.Equal(t, reflect.ValueOf(tc.input).FieldByName("Val"), reflect.ValueOf(tc.output).Elem().FieldByName("Val"))
		})
	}
}
