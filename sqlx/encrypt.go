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

package sqlx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// EncryptColumn 代表一个加密的列
// 一般来说加密可以选择依赖于数据库进行加密
// EncryptColumn 并不打算使用极其难破解的加密算法
// 而是选择使用 AES GCM 模式。
// 如果你觉得安全性不够，那么你可以考虑自己实现类似的结构体.
type EncryptColumn[T any] struct {
	Val   T
	Valid bool
	Key   string
}

var errInvalid = errors.New("ekit EncryptColumn无效")
var errKeyLengthInvalid = errors.New("ekit EncryptColumn仅支持 16/24/32 byte 的key")

// Value 返回加密后的值
// 如果 T 是基本类型，那么会对 T 进行直接加密
// 否则，将 T 按照 JSON 序列化之后进行加密，返回加密后的数据
func (e EncryptColumn[T]) Value() (driver.Value, error) {
	if !e.Valid {
		return nil, errInvalid
	}
	if len(e.Key) != 16 && len(e.Key) != 24 && len(e.Key) != 32 {
		return nil, errKeyLengthInvalid
	}
	var val any = e.Val
	var err error
	var b []byte
	switch valT := val.(type) {
	case string:
		b = []byte(valT)
	case []byte:
		b = valT
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64,
		float32, float64:
		buffer := new(bytes.Buffer)
		err = binary.Write(buffer, binary.BigEndian, val)
		b = buffer.Bytes()
	case int:
		tmp := int64(valT)
		buffer := new(bytes.Buffer)
		err = binary.Write(buffer, binary.BigEndian, tmp)
		b = buffer.Bytes()
	case uint:
		tmp := uint64(valT)
		buffer := new(bytes.Buffer)
		err = binary.Write(buffer, binary.BigEndian, tmp)
		b = buffer.Bytes()
	default:
		b, err = json.Marshal(e.Val)
	}
	if err != nil {
		return nil, err
	}
	return e.aesEncrypt(b)
}

// Scan 方法会把写入的数据转化进行解密，
// 并将解密后的数据进行反序列化，构造 T
func (e *EncryptColumn[T]) Scan(src any) error {
	var err error
	var b []byte
	switch value := src.(type) {
	case []byte:
		b, err = e.aesDecrypt(value)
	case string:
		b, err = e.aesDecrypt([]byte(value))
	default:
		return fmt.Errorf("ekit：EncryptColumn.Scan 不支持 src 类型 %v", src)
	}
	if err != nil {
		return err
	}
	err = e.setValAfterDecrypt(b)
	e.Valid = err == nil
	return err
}

func (e *EncryptColumn[T]) setValAfterDecrypt(deEncrypt []byte) error {
	var val any = &e.Val
	var err error
	switch valT := val.(type) {
	case *string:
		*valT = string(deEncrypt)
	case *[]byte:
		*valT = deEncrypt
	case *int8, *int16, *int32, *int64, *uint8, *uint16, *uint32, *uint64,
		*float32, *float64:
		reader := bytes.NewReader(deEncrypt)
		err = binary.Read(reader, binary.BigEndian, valT)
	case *int:
		tmp := new(int64)
		reader := bytes.NewReader(deEncrypt)
		err = binary.Read(reader, binary.BigEndian, tmp)
		*valT = int(*tmp)
	case *uint:
		tmp := new(uint64)
		reader := bytes.NewReader(deEncrypt)
		err = binary.Read(reader, binary.BigEndian, tmp)
		*valT = uint(*tmp)
	default:
		err = json.Unmarshal(deEncrypt, &e.Val)
	}
	return err
}

func (e *EncryptColumn[T]) aesEncrypt(data []byte) ([]byte, error) {
	newCipher, err := aes.NewCipher([]byte(e.Key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(newCipher)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}

func (e *EncryptColumn[T]) aesDecrypt(data []byte) ([]byte, error) {
	newCipher, err := aes.NewCipher([]byte(e.Key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(newCipher)
	if err != nil {
		return nil, err
	}
	nonce, cipherData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, cipherData, nil)
}
