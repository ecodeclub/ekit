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
			// 通过读取 file 文件 模拟 io.Reader 中存在的字节流 并将其转换为 string 检查他的正确性
			// 当然他必须是可控制的
			name: "file(io.Reader) read bytes stream coversion string",
			before: func(t *testing.T) {
				create, err := os.Create("/tmp/test_put.txt")
				require.NoError(t, err)
				defer create.Close()
				_, err = create.WriteString("the test file...")
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				require.NoError(t, os.Remove("/tmp/test_put.txt"))
			},
			val: func(t *testing.T) []byte {
				open, err := os.Open("/tmp/test_put.txt")
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
			defer tt.after(t)
			tt.before(t)
			b := tt.val(t)
			val := UnsafeToString(b)
			assert.Equal(t, tt.want, val)
		})
	}
}

func Benchmark_UnsafeToBytes(b *testing.B) {
	b.Run("safe to bytes", func(b *testing.B) {
		s := "hello ekit! hello golang! this is test benchmark"
		for i := 0; i < b.N; i++ {
			_ = []byte(s)
		}
	})

	b.Run("unsafe to bytes", func(b *testing.B) {
		s := "hello ekit! hello golang! this is test benchmark"
		for i := 0; i < b.N; i++ {
			_ = UnsafeToBytes(s)
		}
	})
}

func Benchmark_UnsafeToString(b *testing.B) {
	b.Run("safe to string", func(b *testing.B) {
		s := []byte("hello ekit! hello golang! this is test benchmark")
		for i := 0; i < b.N; i++ {
			_ = string(s)
		}
	})

	b.Run("unsafe to string", func(b *testing.B) {
		s := []byte("hello ekit! hello golang! this is test benchmark")
		for i := 0; i < b.N; i++ {
			_ = UnsafeToString(s)
		}
	})
}
