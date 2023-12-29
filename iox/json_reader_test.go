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

package iox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONReader(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
		val   any

		wantRes []byte
		wantN   int
		wantErr error
	}{
		{
			name:    "正常读取",
			input:   make([]byte, 10),
			wantN:   10,
			val:     User{Name: "Tom"},
			wantRes: []byte(`{"name":"T`),
		},
		{
			name:    "输入 nil",
			input:   make([]byte, 7),
			wantN:   4,
			wantRes: append([]byte(`null`), 0, 0, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := NewJSONReader(tc.val)
			n, err := reader.Read(tc.input)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantN, n)
			assert.Equal(t, tc.wantRes, tc.input)
		})
	}
}

type User struct {
	Name string `json:"name"`
}
