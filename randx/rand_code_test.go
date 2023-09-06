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

package randx

import (
	"errors"
	"regexp"
	"testing"
)

func TestRandCode(t *testing.T) {
	testCases := []struct {
		name      string
		length    int
		typ       TYPE
		wantMatch string
		wantErr   error
	}{
		{
			name:      "默认类型",
			length:    8,
			typ:       TYPE_DEFAULT,
			wantMatch: "^[0-9]+$",
			wantErr:   nil,
		},
		{
			name:      "数字验证码",
			length:    8,
			typ:       TYPE_DIGIT,
			wantMatch: "^[0-9]+$",
			wantErr:   nil,
		}, {
			name:      "小写字母验证码",
			length:    8,
			typ:       TYPE_LETTER,
			wantMatch: "^[a-z]+$",
			wantErr:   nil,
		}, {
			name:      "大写字母验证码",
			length:    8,
			typ:       TYPE_CAPITAL,
			wantMatch: "^[A-Z]+$",
			wantErr:   nil,
		}, {
			name:      "混合验证码",
			length:    8,
			typ:       TYPE_MIXED,
			wantMatch: "^[0-9a-zA-Z]+$",
			wantErr:   nil,
		}, {
			name:      "未定义类型",
			length:    8,
			typ:       9,
			wantMatch: "",
			wantErr:   ERRTYPENOTSUPPORTTED,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code, err := RandCode(tc.length, tc.typ)
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				//长度检验
				if len(code) != tc.length {
					t.Errorf("expected length: %d but got length:%d  ", tc.length, len(code))
				}
				//模式检验
				matched, _ := regexp.MatchString(tc.wantMatch, code)
				if !matched {
					t.Errorf("expected %s but got %s", tc.wantMatch, code)
				}
			}
		})
	}

}
