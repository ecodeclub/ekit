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
	"math/rand"
)

var (
	ErrTypeNotSupported = errors.New("ekit:不支持的类型")
	// deprecated
	ERRTYPENOTSUPPORTTED = ErrTypeNotSupported
)

type TYPE int

const (
	// 数字
	TYPE_DIGIT TYPE = 1
	// 小写字母
	TYPE_LOWER  TYPE = 1 << 1
	TYPE_LETTER TYPE = TYPE_LOWER
	// 大写字母
	TYPE_UPPER   TYPE = 1 << 2
	TYPE_CAPITAL TYPE = TYPE_UPPER
	// 混合类型
	TYPE_MIXED = (TYPE_DIGIT | TYPE_UPPER | TYPE_LOWER)

	// 数字字符组
	CHARSET_DIGIT = "0123456789"
	// 小写字母字符组
	CHARSET_LOWER  = "abcdefghijklmnopqrstuvwxyz"
	CHARSET_LETTER = CHARSET_LOWER
	// 大写字母字符组
	CHARSET_UPPER   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CHARSET_CAPITAL = CHARSET_UPPER
)

// RandCode 根据传入的长度和类型生成随机字符串,这个方法目前可以生成数字、字母、数字+字母的随机字符串
func RandCode(length int, typ TYPE) (string, error) {
	charset := ""
	appendIfHas := func(typBase TYPE, baseCharset string) {
		if (typ & typBase) == typBase {
			charset += baseCharset
		}
	}
	appendIfHas(TYPE_DIGIT, CHARSET_DIGIT)
	appendIfHas(TYPE_UPPER, CHARSET_UPPER)
	appendIfHas(TYPE_LOWER, CHARSET_LOWER)

	charsetSize := len(charset)
	if charsetSize == 0 {
		return "", ErrTypeNotSupported
	}
	bits := 1
	for charsetSize > ((1 << bits) - 1) {
		bits++
	}
	return generate(charset, length, bits), nil
}

// generate 根据传入的随机源和长度生成随机字符串,一次随机，多次使用
func generate(source string, length, idxBits int) string {

	//掩码
	//例如： 使用低6位：0000 0000 --> 0011 1111
	idxMask := 1<<idxBits - 1

	// 63位最多可以使用多少次
	remain := 63 / idxBits

	//cache 随机位缓存
	cache := rand.Int63()

	result := make([]byte, length)

	for i := 0; i < length; {
		//如果使用次数剩余0，重新获取随机
		if remain == 0 {
			cache, remain = rand.Int63(), 63/idxBits
		}

		//利用掩码获取有效的随机数位
		if randIndex := int(cache & int64(idxMask)); randIndex < len(source) {
			result[i] = source[randIndex]
			i++
		}

		//使用下一组随机位
		cache >>= idxBits

		//扣减remain
		remain--

	}
	return string(result)

}
