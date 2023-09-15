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

package stringx

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnsafeToBytes(t *testing.T) {
	testCase := []struct {
		name string
		val  string
		want []byte
	}{
		{
			name: "normal conversion",
			val:  "hello",
			want: []byte("hello"),
		},
		{
			name: "emoji coversion",
			val:  "😀!hello world",
			want: []byte("😀!hello world"),
		},
		{
			name: "chinese coversion",
			val:  "你好 世界！",
			want: []byte("你好 世界！"),
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			val := UnsafeToBytes(tt.val)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestUnsafeToString(t *testing.T) {
	testCase := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		val    func(t *testing.T) []byte
		want   string
	}{
		{
			name:   "normal conversion",
			before: func(t *testing.T) {},
			after:  func(t *testing.T) {},
			val: func(t *testing.T) []byte {
				return []byte("hello")
			},
			want: "hello",
		},
		{
			name:   "emoji coversion",
			before: func(t *testing.T) {},
			after:  func(t *testing.T) {},
			val: func(t *testing.T) []byte {
				return []byte("😀!hello world")
			},
			want: "😀!hello world",
		},
		{
			name:   "chinese coversion",
			before: func(t *testing.T) {},
			after:  func(t *testing.T) {},
			val: func(t *testing.T) []byte {
				return []byte("你好 世界！")
			},
			want: "你好 世界！",
		},
		{
			name: "file io.Reader coversion",
			before: func(t *testing.T) {
				create, err := os.Create("test_put.txt")
				require.NoError(t, err)
				defer create.Close()
				_, err = create.WriteString("the test file...")
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				require.NoError(t, os.Remove("test_put.txt"))
			},
			val: func(t *testing.T) []byte {
				open, err := os.Open("test_put.txt")
				require.NoError(t, err)
				defer open.Close()
				buf := bytes.Buffer{}
				_, err = buf.ReadFrom(open)
				require.NoError(t, err)
				return buf.Bytes()
			},
			want: "the test file...",
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			b := tt.val(t)
			val := UnsafeToString(b)
			assert.Equal(t, tt.want, val)
			tt.after(t)
		})
	}
}
