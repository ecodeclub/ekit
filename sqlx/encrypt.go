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
	"errors"
	"fmt"
	"io"
	"reflect"
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
	//获取块的大小
	e.blockSize = block.BlockSize()
	//使用cbc
	e.decrypt = cipher.NewCBCDecrypter(block, key[:e.blockSize])
	e.encrypt = cipher.NewCBCEncrypter(block, key[:e.blockSize])
	return e, nil
}

// EncryptColumn 代表一个加密的列
// 一般来说加密可以选择依赖于数据库进行加密
// EncryptColumn 并不打算使用极其难破解的加密算法
// 而是选择使用 AES GCM 模式。
// 如果你觉得安全性不够，那么你可以考虑自己实现类似的结构体.
func (e *EncryptColumn[T]) Value() (driver.Value, error) {
	relValue := reflect.ValueOf(&e.Val).Elem()
	data, err := valueToBytes(relValue, &e.Val)
	if err != nil {
		return nil, err
	}
	return aesEncryptWithSizeAndMode(data, e.blockSize, e.encrypt)
}

func (e *EncryptColumn[T]) Scan(src any) error {
	relValue := reflect.ValueOf(&e.Val).Elem()
	var decrBytes []byte
	switch value := src.(type) {
	case []byte:
		tmpBytes, err := aesDecryptWithSizeAndMode(value, e.blockSize, e.decrypt)
		if err != nil {
			return nil
		}
		decrBytes = tmpBytes
	default:
		return fmt.Errorf("ekit：EncryptColumn.Scan 不支持 src 类型 %v", src)
	}
	reader := bytes.NewReader(decrBytes)
	err := setVal[T](relValue, reader, decrBytes, &e.Val)
	if err != nil {
		return err
	}
	e.Valid = true
	return nil
}

func valueToBytes[T any](relValue reflect.Value, valT *T) ([]byte, error) {
	switch relValue.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128:
		buffer := new(bytes.Buffer)
		err := binary.Write(buffer, binary.BigEndian, *valT)
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	case reflect.Uint:
		val := relValue.Uint()
		buffer := new(bytes.Buffer)
		err := binary.Write(buffer, binary.BigEndian, val)
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	case reflect.Int:
		val := relValue.Int()
		buffer := new(bytes.Buffer)
		err := binary.Write(buffer, binary.BigEndian, val)
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	case reflect.Struct:
		jsonColumn := &JsonColumn[T]{
			Val:   *valT,
			Valid: true,
		}
		jsonByte, err := jsonColumn.Value()
		if err != nil {
			return nil, err
		}
		buffer := reflect.ValueOf(jsonByte)
		return buffer.Bytes(), nil
	case reflect.String:
		return []byte(relValue.String()), nil
	default:
		return nil, fmt.Errorf("ekit：EncryptColumn.Value 不支持 src 类型 %v", relValue.Kind())
	}
}

func setVal[T any](relValue reflect.Value, reader io.Reader, deData []byte, val *T) error {
	switch relValue.Kind() {
	case reflect.Int:
		tmp := int64(0)
		err := binary.Read(reader, binary.BigEndian, &tmp)
		if err != nil {
			return err
		}
		relValue.SetInt(tmp)
	case reflect.Uint:
		tmp := uint64(0)
		err := binary.Read(reader, binary.BigEndian, &tmp)
		if err != nil {
			return err
		}
		relValue.SetUint(tmp)
	case reflect.Struct:
		json := &JsonColumn[T]{}
		err := json.Scan(deData)
		if err != nil {
			return err
		}
		*val = json.Val
	case reflect.String:
		relValue.SetString(string(deData))
	default:
		err := binary.Read(reader, binary.BigEndian, val)
		if err != nil {
			return err
		}
	}
	return nil
}

// 加密过程：
//
//	1、处理数据，对数据进行填充，采用PKCS7（当密钥长度不够时，缺几位补几个几）的方式。
//	2、对数据进行加密，采用AES加密方法中CBC加密模式
//	3、对得到的加密数据，进行base64加密，得到字符串
//
// 解密过程相反
// 16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
// key不能泄露
// pkcs7Padding 填充
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

func aesDecryptWithSizeAndMode(data []byte, blockSize int, decrypt cipher.BlockMode) ([]byte, error) {
	crypted := make([]byte, len(data))
	decrypt.CryptBlocks(crypted, data)
	crypted, err := pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}
