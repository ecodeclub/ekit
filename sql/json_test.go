// Copyright 2021 gotomicro
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

package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonColumn_Value(t *testing.T) {
	js := JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}}
	value, err := js.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"Name":"Tom"}`), value)
	js = JsonColumn[User]{}
	value, err = js.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)
}

func TestJsonColumn_Scan(t *testing.T) {
	testCases := []struct {
		name    string
		src     any
		wantErr error
		wantVal User
	}{
		{
			name:    "nil",
			wantErr: errors.New("ekit：JsonColumn.Scan 不支持 src 类型 <nil>"),
		},
		{
			name:    "string",
			src:     `{"Name":"Tom"}`,
			wantVal: User{Name: "Tom"},
		},
		{
			name: "string pointer",
			src: func() string {
				return `{"Name":"Tom"}`
			}(),
			wantVal: User{Name: "Tom"},
		},
		{
			name:    "bytes",
			src:     []byte(`{"Name":"Tom"}`),
			wantVal: User{Name: "Tom"},
		},
		{
			name: "bytes pointer",
			src: func() *[]byte {
				res := []byte(`{"Name":"Tom"}`)
				return &res
			}(),
			wantVal: User{Name: "Tom"},
		},
		{
			name:    "sql.RawBytes",
			src:     sql.RawBytes(`{"Name":"Tom"}`),
			wantVal: User{Name: "Tom"},
		},
		{
			name: "sql.RawBytes pointer",
			src: func() *sql.RawBytes {
				res := sql.RawBytes(`{"Name":"Tom"}`)
				return &res
			}(),
			wantVal: User{Name: "Tom"},
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
			assert.Equal(t, tc.wantVal, js.Val)
			assert.True(t, js.Valid)
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
