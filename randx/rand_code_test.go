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

package randx_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/ecodeclub/ekit/randx"
	"github.com/stretchr/testify/assert"
)

func TestRandCode(t *testing.T) {
	testCases := []struct {
		name      string
		length    int
		typ       randx.TYPE
		wantMatch string
		wantErr   error
	}{
		{
			name:      "数字验证码",
			length:    100,
			typ:       randx.TYPE_DIGIT,
			wantMatch: "^[0-9]+$",
			wantErr:   nil,
		},
		{
			name:      "小写字母验证码",
			length:    100,
			typ:       randx.TYPE_LETTER,
			wantMatch: "^[a-z]+$",
			wantErr:   nil,
		},
		{
			name:      "数字+小写字母验证码",
			length:    100,
			typ:       randx.TYPE_DIGIT | randx.TYPE_LOWER,
			wantMatch: "^[a-z0-9]+$",
			wantErr:   nil,
		},
		{
			name:      "数字+大写字母验证码",
			length:    100,
			typ:       randx.TYPE_DIGIT | randx.TYPE_UPPER,
			wantMatch: "^[A-Z0-9]+$",
			wantErr:   nil,
		},
		{
			name:      "大写字母验证码",
			length:    100,
			typ:       randx.TYPE_CAPITAL,
			wantMatch: "^[A-Z]+$",
			wantErr:   nil,
		},
		{
			name:      "大小写字母验证码",
			length:    100,
			typ:       randx.TYPE_UPPER | randx.TYPE_LOWER,
			wantMatch: "^[a-zA-Z]+$",
			wantErr:   nil,
		},
		{
			name:      "数字+大小写字母验证码",
			length:    100,
			typ:       randx.TYPE_DIGIT | randx.TYPE_UPPER | randx.TYPE_LOWER,
			wantMatch: "^[0-9a-zA-Z]+$",
			wantErr:   nil,
		},
		{
			name:      "所有类型验证",
			length:    100,
			typ:       randx.TYPE_MIXED,
			wantMatch: "^[\\S\\s]*$",
			wantErr:   nil,
		},
		{
			name:      "未定义类型(超过范围)",
			length:    100,
			typ:       randx.TYPE_MIXED + 1,
			wantMatch: "",
			wantErr:   randx.ErrTypeNotSupported,
		},
		{
			name:      "未定义类型(0)",
			length:    100,
			typ:       0,
			wantMatch: "",
			wantErr:   randx.ErrTypeNotSupported,
		},
		{
			name:      "长度小于0",
			length:    -1,
			typ:       0,
			wantMatch: "",
			wantErr:   randx.ErrLengthLessThanZero,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code, err := randx.RandCode(tc.length, tc.typ)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
			} else {
				assert.Lenf(
					t,
					code,
					tc.length,
					"expected length: %d but got length:%d",
					tc.length, len(code))

				matched, err := regexp.MatchString(tc.wantMatch, code)
				assert.Nil(t, err)
				assert.Truef(t, matched, "expected %s but got %s", tc.wantMatch, code)
			}
		})
	}
}

func TestRandStrByCharset(t *testing.T) {
	matchFunc := func(str, charset string) bool {
		for _, c := range str {
			if !strings.Contains(charset, string(c)) {
				return false
			}
		}
		return true
	}
	testCases := []struct {
		name    string
		length  int
		charset string
		wantErr error
	}{
		{
			name:    "长度小于0",
			length:  -1,
			charset: "123",
			wantErr: randx.ErrLengthLessThanZero,
		},
		{
			name:    "随机字符串测试",
			length:  100,
			charset: "2rg248ry227t@@",
			wantErr: nil,
		},
		{
			name:    "随机字符串测试",
			length:  100,
			charset: "2rg248ry227t@&*($.!",
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code, err := randx.RandStrByCharset(tc.length, tc.charset)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
			} else {
				assert.Lenf(
					t,
					code,
					tc.length,
					"expected length: %d but got length:%d",
					tc.length, len(code))
				assert.True(t, matchFunc(code, tc.charset))
			}
		})
	}
}

// goos: linux
// goarch: amd64
// pkg: github.com/ecodeclub/ekit/randx
// cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
// BenchmarkRandCode_MIXED/length=1000000-8                1000000000               0.004584 ns/op        0 B/op          0 allocs/op
func BenchmarkRandCode_MIXED(b *testing.B) {
	b.Run("length=1000000", func(b *testing.B) {
		n := 1000000
		b.StartTimer()
		res, err := randx.RandCode(n, randx.TYPE_MIXED)
		b.StopTimer()
		assert.Nil(b, err)
		assert.Len(b, res, n)
	})
}
