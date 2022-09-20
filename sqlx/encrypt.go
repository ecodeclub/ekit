// Copyright 2021 gotomicro
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

package sqlx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// EncryptColumn 代表一个加密的列
// 一般来说加密可以选择依赖于数据库进行加密
// EncryptColumn 并不打算使用极其难破解的加密算法
// 而是选择使用 AES GCM 模式。
// 如果你觉得安全性不够，那么你可以考虑自己实现类似的结构体.
type EncryptColumn[T any] struct {
	blockSize int
	decrypt   cipher.BlockMode
	encrypt   cipher.BlockMode
	Val       T
	Valid     bool
}

func NewEncryptColumn[T any](key []byte) (*EncryptColumn[T], error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	e := &EncryptColumn[T]{}
	e.blockSize = block.BlockSize()
	e.decrypt = cipher.NewCBCDecrypter(block, key[:e.blockSize])
	e.encrypt = cipher.NewCBCEncrypter(block, key[:e.blockSize])
	return e, nil
}

var notSpecifyInt = errors.New("ekit 请明确int/uint的长度，如int32/uint32，int8/uint8")

func notSupportType(s any) error {
	switch s.(type) {
	case map[any]any:
		return fmt.Errorf("ekit EncryptColumn不支持map类型")
	case []any:
		return fmt.Errorf("ekit EncryptColumn不支持slice类型")
	case *any, uintptr:
		return fmt.Errorf("ekit EncryptColumn不支持指针类型")
	default:
		return nil
	}
}

func (e *EncryptColumn[T]) Value() (driver.Value, error) {
	var val any = e.Val
	switch valT := val.(type) {
	case string:
		return aesEncryptWithSizeAndMode([]byte(valT), e.blockSize, e.encrypt)
	case []byte:
		return aesEncryptWithSizeAndMode(valT, e.blockSize, e.encrypt)
	case map[any]any, []any, *any, uintptr:
		return nil, notSupportType(val)
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64,
		float32, float64, complex64, complex128:
		buffer := new(bytes.Buffer)
		err := binary.Write(buffer, binary.BigEndian, val)
		if err != nil {
			return nil, err
		}
		return aesEncryptWithSizeAndMode(buffer.Bytes(), e.blockSize, e.encrypt)
	case int, uint:
		return nil, notSpecifyInt
	default:
		marshal, err := json.Marshal(e.Val)
		if err != nil {
			return nil, err
		}
		return aesEncryptWithSizeAndMode(marshal, e.blockSize, e.encrypt)
	}
}

func (e *EncryptColumn[T]) Scan(src any) error {
	switch value := src.(type) {
	case []byte:
		decrBytes, err := aesDecryptWithSizeAndMode(value, e.decrypt)
		if err != nil {
			return nil
		}
		return e.setVal(decrBytes)
	default:
		return fmt.Errorf("ekit：EncryptColumn.Scan 不支持 src 类型 %v", src)
	}
}

func (e *EncryptColumn[T]) setVal(deEncrypt []byte) error {
	var val any = e.Val
	switch val.(type) {
	case string:
		header := (*reflect.StringHeader)(unsafe.Pointer(&e.Val))
		s := string(deEncrypt)
		header.Len = (*reflect.StringHeader)(unsafe.Pointer(&s)).Len
		header.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
		return nil
	case []byte:
		header := (*reflect.SliceHeader)(unsafe.Pointer(&e.Val))
		header.Len = (*reflect.SliceHeader)(unsafe.Pointer(&deEncrypt)).Len
		header.Data = (*reflect.SliceHeader)(unsafe.Pointer(&deEncrypt)).Data
		header.Cap = (*reflect.SliceHeader)(unsafe.Pointer(&deEncrypt)).Cap
		return nil
	case map[any]any, []any, *any, uintptr:
		return notSupportType(val)
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64,
		float32, float64, complex64, complex128:
		reader := bytes.NewReader(deEncrypt)
		err := binary.Read(reader, binary.BigEndian, &e.Val)
		if err != nil {
			return err
		}
	case int, uint:
		return notSpecifyInt
	default:
		err := json.Unmarshal(deEncrypt, &e.Val)
		if err != nil {
			return err
		}
	}
	return nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func aesEncryptWithSizeAndMode(data []byte, blockSize int, encrypt cipher.BlockMode) ([]byte, error) {
	encryptBytes := pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	encrypt.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

func aesDecryptWithSizeAndMode(data []byte, decrypt cipher.BlockMode) ([]byte, error) {
	crypted := make([]byte, len(data))
	decrypt.CryptBlocks(crypted, data)
	crypted, err := pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}
