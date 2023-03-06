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
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonColumn_Value(t *testing.T) {
	testCases := []struct {
		name    string
		valuer  driver.Valuer
		wantRes any
		wantErr error
	}{
		{
			name:    "user",
			valuer:  JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}},
			wantRes: []byte(`{"Name":"Tom"}`),
		},
		{
			name:   "invalid",
			valuer: JsonColumn[User]{},
		},
		{
			name:   "nil",
			valuer: JsonColumn[*User]{},
		},
		{
			name:    "nil but valid",
			valuer:  JsonColumn[*User]{Valid: true},
			wantRes: []uint8("null"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.valuer.Value()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, value)
		})
	}
}

func TestJsonColumn_Scan(t *testing.T) {
	testCases := []struct {
		name      string
		src       any
		wantErr   error
		wantValid bool
		wantVal   User
	}{
		{
			name:    "nil",
			wantVal: User{},
		},
		{
			name:      "string",
			src:       `{"Name":"Tom"}`,
			wantVal:   User{Name: "Tom"},
			wantValid: true,
		},
		{
			name:      "bytes",
			src:       []byte(`{"Name":"Tom"}`),
			wantVal:   User{Name: "Tom"},
			wantValid: true,
		},
		{
			name:    "int",
			src:     123,
			wantErr: errors.New("ekit：JsonColumn.Scan 不支持 src 类型 123"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			js := &JsonColumn[User]{}
			err := js.Scan(tc.src)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantValid, js.Valid)
			if !js.Valid {
				return
			}
			assert.Equal(t, tc.wantVal, js.Val)
		})
	}
}

func TestJsonColumn_ScanTypes(t *testing.T) {
	jsSlice := JsonColumn[[]string]{}
	err := jsSlice.Scan(`["a", "b", "c"]`)
	assert.Nil(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, jsSlice.Val)
	val, err := jsSlice.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`["a","b","c"]`), val)

	jsMap := JsonColumn[map[string]string]{}
	err = jsMap.Scan(`{"a":"a value"}`)
	assert.Nil(t, err)
	val, err = jsMap.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"a":"a value"}`), val)
}

type User struct {
	Name string
}

func ExampleJsonColumn_Value() {
	js := JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}}
	value, err := js.Value()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(value.([]byte)))
	// Output:
	// {"Name":"Tom"}
}

func ExampleJsonColumn_Scan() {
	js := JsonColumn[User]{}
	err := js.Scan(`{"Name":"Tom"}`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(js.Val)
	// Output:
	// {Tom}
}
