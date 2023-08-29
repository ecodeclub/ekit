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

var ERRTYPENOTSUPPORTTED = errors.New("ekit:不支持的类型")

type TYPE int

const (
	TYPE_DEFAULT TYPE = 0 //默认类型
	TYPE_DIGIT   TYPE = 1 //数字//
	TYPE_LETTER  TYPE = 2 //小写字母
	TYPE_CAPITAL TYPE = 3 //大写字母
	TYPE_MIXED   TYPE = 4 //数字+字母混合
)

// RandCode 根据传入的长度和类型生成随机字符串,这个方法目前可以生成数字、字母、数字+字母的随机字符串
func RandCode(length int, typ TYPE) (string, error) {
	switch typ {
	case TYPE_DEFAULT:
		fallthrough
	case TYPE_DIGIT:
		return generate("0123456789", length, 4), nil
	case TYPE_LETTER:
		return generate("abcdefghijklmnopqrstuvwxyz", length, 5), nil
	case TYPE_CAPITAL:
		return generate("ABCDEFGHIJKLMNOPQRSTUVWXYZ", length, 5), nil
	case TYPE_MIXED:
		return generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", length, 7), nil
	default:
		return "", ERRTYPENOTSUPPORTTED
	}
}

// generate 根据传入的随机源和长度生成随机字符串,一次随机，多次使用
func generate(source string, length, idxBits int) string {

	//掩码
	//例如： 使用低6位：0000 0000 --> 0011 1111
	idxMask := 1<<idxBits - 1

	// 63位最多可以使用多少次
	idxMax := 63 / idxBits

	result := make([]byte, length)

	//cache 随机位缓存
	//remain 当前还可以使用几次
	for i, cache, remain := 0, rand.Int63(), idxMax; i < length; {
		//如果使用次数剩余0，重新获取随机
		if remain == 0 {
			cache, remain = rand.Int63(), idxMax
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
