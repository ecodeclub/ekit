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

	"github.com/ecodeclub/ekit/tuple/pair"
)

var (
	ErrTypeNotSupported   = errors.New("ekit:不支持的类型")
	ErrLengthLessThanZero = errors.New("ekit:长度必须大于0")
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
	// 特殊符号
	TYPE_SPECIAL TYPE = 1 << 3
	// 混合类型
	TYPE_MIXED = (TYPE_DIGIT | TYPE_UPPER | TYPE_LOWER | TYPE_SPECIAL)

	// 数字字符组
	CHARSET_DIGIT = "0123456789"
	// 小写字母字符组
	CHARSET_LOWER  = "abcdefghijklmnopqrstuvwxyz"
	CHARSET_LETTER = CHARSET_LOWER
	// 大写字母字符组
	CHARSET_UPPER   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CHARSET_CAPITAL = CHARSET_UPPER
	// 特殊字符数组
	CHARSET_SPECIAL = " ~!@#$%^&*()_+-=[]{};'\\:\"|,./<>?"
)

var (
	// 只限于randx包内部使用
	typeCharSetPair = []pair.Pair[TYPE, string]{
		pair.NewPair(TYPE_DIGIT, CHARSET_DIGIT),
		pair.NewPair(TYPE_LOWER, CHARSET_LOWER),
		pair.NewPair(TYPE_UPPER, CHARSET_UPPER),
		pair.NewPair(TYPE_SPECIAL, CHARSET_SPECIAL),
	}
)

// RandCode 根据传入的长度和类型生成随机字符串
func RandCode(length int, typ TYPE) (string, error) {
	if typ > TYPE_MIXED {
		return "", ErrTypeNotSupported
	}
	charset := ""
	for _, p := range typeCharSetPair {
		if (typ & p.Key) == p.Key {
			charset += p.Value
		}
	}
	return RandStrByCharset(length, charset)
}

// 根据传入的长度和字符集生成随机字符串
func RandStrByCharset(length int, charset string) (string, error) {
	if length < 0 {
		return "", ErrLengthLessThanZero
	}
	charsetSize := len(charset)
	if charsetSize == 0 {
		return "", ErrTypeNotSupported
	}
	return generate(charset, length, getFirstMask(charsetSize)), nil
}

func getFirstMask(charsetSize int) int {
	bits := 0
	for charsetSize > ((1 << bits) - 1) {
		bits++
	}
	return bits
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
