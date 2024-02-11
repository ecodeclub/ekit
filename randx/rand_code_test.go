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
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/ecodeclub/ekit/randx"
	"github.com/stretchr/testify/assert"
)

var (
	errTypeNotSupported   = errors.New("ekit:不支持的类型")
	errLengthLessThanZero = errors.New("ekit:长度必须大于等于0")
)

func TestRandCode(t *testing.T) {
	testCases := []struct {
		name      string
		length    int
		typ       randx.Type
		wantMatch string
		wantErr   error
	}{
		{
			name:      "数字验证码",
			length:    100,
			typ:       randx.TypeDigit,
			wantMatch: "^[0-9]+$",
			wantErr:   nil,
		},
		{
			name:      "小写字母验证码",
			length:    100,
			typ:       randx.TypeLowerCase,
			wantMatch: "^[a-z]+$",
			wantErr:   nil,
		},
		{
			name:      "数字+小写字母验证码",
			length:    100,
			typ:       randx.TypeDigit | randx.TypeLowerCase,
			wantMatch: "^[a-z0-9]+$",
			wantErr:   nil,
		},
		{
			name:      "数字+大写字母验证码",
			length:    100,
			typ:       randx.TypeDigit | randx.TypeUpperCase,
			wantMatch: "^[A-Z0-9]+$",
			wantErr:   nil,
		},
		{
			name:      "大写字母验证码",
			length:    100,
			typ:       randx.TypeUpperCase,
			wantMatch: "^[A-Z]+$",
			wantErr:   nil,
		},
		{
			name:      "大小写字母验证码",
			length:    100,
			typ:       randx.TypeUpperCase | randx.TypeLowerCase,
			wantMatch: "^[a-zA-Z]+$",
			wantErr:   nil,
		},
		{
			name:      "数字+大小写字母验证码",
			length:    100,
			typ:       randx.TypeDigit | randx.TypeUpperCase | randx.TypeLowerCase,
			wantMatch: "^[0-9a-zA-Z]+$",
			wantErr:   nil,
		},
		{
			name:      "所有类型验证",
			length:    100,
			typ:       randx.TypeMixed,
			wantMatch: "^[\\S\\s]+$",
			wantErr:   nil,
		},
		{
			name:      "特殊字符类型验证",
			length:    100,
			typ:       randx.TypeSpecial,
			wantMatch: "^[^0-9a-zA-Z]+$",
			wantErr:   nil,
		},
		{
			name:      "未定义类型(超过范围)",
			length:    100,
			typ:       randx.TypeMixed + 1,
			wantMatch: "",
			wantErr:   errTypeNotSupported,
		},
		{
			name:      "未定义类型(0)",
			length:    100,
			typ:       0,
			wantMatch: "",
			wantErr:   errTypeNotSupported,
		},
		{
			name:      "长度小于0",
			length:    -1,
			typ:       0,
			wantMatch: "",
			wantErr:   errLengthLessThanZero,
		},
		{
			name:      "长度等于0",
			length:    0,
			typ:       randx.TypeMixed,
			wantMatch: "",
			wantErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code, err := randx.RandCode(tc.length, tc.typ)
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.Len(t, code, tc.length)
			if tc.length > 0 {
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
			wantErr: errLengthLessThanZero,
		},
		{
			name:    "长度等于0",
			length:  0,
			charset: "123",
			wantErr: nil,
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
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}

			assert.Len(t, code, tc.length)
			if tc.length > 0 {
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
		res, err := randx.RandCode(n, randx.TypeMixed)
		b.StopTimer()
		assert.Nil(b, err)
		assert.Len(b, res, n)
	})
}
