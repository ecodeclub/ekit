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
	"regexp"
	"testing"
)

func TestRandCode_Digit(t *testing.T) {
	code, err := RandCode(6, TYPE_DIGIT)
	t.Log(code)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(code) != 6 {
		t.Errorf("expected 6 digits but got %d", len(code))
	}

	matched, err := regexp.MatchString("^[0-9]+$", code)
	if !matched || err != nil {
		t.Error("expected all digits code")
	}
}

func TestRandCode_Letter(t *testing.T) {

	// 类似的测试字母表生成
	code, err := RandCode(6, TYPE_LETTER)
	t.Log(code)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(code) != 6 {
		t.Errorf("expected 6 letters but got %d", len(code))
	}

	matched, err := regexp.MatchString("^[a-z]+$", code)
	if !matched || err != nil {
		t.Error("expected all letters code")
	}
}

func TestRandCode_Capital(t *testing.T) {

	// 类似的测试字母表生成
	code, err := RandCode(6, TYPE_CAPITAL)
	t.Log(code)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(code) != 6 {
		t.Errorf("expected 6 letters but got %d", len(code))
	}

	matched, err := regexp.MatchString("^[A-Z]+$", code)
	if !matched || err != nil {
		t.Error("expected all letters code")
	}
}

func TestRandCode_Mixed(t *testing.T) {

	// 类似的测试字母表生成
	code, err := RandCode(6, TYPE_MIXED)
	t.Log(code)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(code) != 6 {
		t.Errorf("expected 6 letters but got %d", len(code))
	}

	matched, err := regexp.MatchString("^[A-Za-z0-9]+$", code)
	if !matched || err != nil {
		t.Error("expected all letters code")
	}
}
