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
	errTypeNotSupported   = errors.New("ekit:不支持的类型")
	errLengthLessThanZero = errors.New("ekit:长度必须大于等于0")
)

type Type int

const (
	// TypeDigit 数字
	TypeDigit Type = 1
	// TypeLowerCase 小写字母
	TypeLowerCase Type = 1 << 1
	// TypeUpperCase 大写字母
	TypeUpperCase Type = 1 << 2
	// TypeSpecial 特殊符号
	TypeSpecial Type = 1 << 3
	// TypeMixed 混合类型
	TypeMixed = (TypeDigit | TypeUpperCase | TypeLowerCase | TypeSpecial)

	// CharsetDigit 数字字符组
	CharsetDigit = "0123456789"
	// CharsetLowerCase 小写字母字符组
	CharsetLowerCase = "abcdefghijklmnopqrstuvwxyz"
	// CharsetUpperCase 大写字母字符组
	CharsetUpperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// CharsetSpecial 特殊字符数组
	CharsetSpecial = " ~!@#$%^&*()_+-=[]{};'\\:\"|,./<>?"
)

var (
	// 只限于randx包内部使用
	typeCharsetPairs = []pair.Pair[Type, string]{
		pair.NewPair(TypeDigit, CharsetDigit),
		pair.NewPair(TypeLowerCase, CharsetLowerCase),
		pair.NewPair(TypeUpperCase, CharsetUpperCase),
		pair.NewPair(TypeSpecial, CharsetSpecial),
	}
)

// RandCode 根据传入的长度和类型生成随机字符串
// 请保证输入的 length >= 0，否则会返回 errLengthLessThanZero
// 请保证输入的 typ 的取值范围在 (0, type.MIXED] 内，否则会返回 errTypeNotSupported
func RandCode(length int, typ Type) (string, error) {
	if length < 0 {
		return "", errLengthLessThanZero
	}
	if length == 0 {
		return "", nil
	}
	if typ > TypeMixed {
		return "", errTypeNotSupported
	}
	charset := ""
	for _, p := range typeCharsetPairs {
		if (typ & p.Key) == p.Key {
			charset += p.Value
		}
	}
	return RandStrByCharset(length, charset)
}

// RandStrByCharset 根据传入的长度和字符集生成随机字符串
// 请保证输入的 length >= 0，否则会返回 errLengthLessThanZero
// 请保证输入的字符集不为空字符串，否则会返回 errTypeNotSupported
// 字符集内部字符可以无序或重复
func RandStrByCharset(length int, charset string) (string, error) {
	if length < 0 {
		return "", errLengthLessThanZero
	}
	if length == 0 {
		return "", nil
	}
	charsetSize := len(charset)
	if charsetSize == 0 {
		return "", errTypeNotSupported
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
