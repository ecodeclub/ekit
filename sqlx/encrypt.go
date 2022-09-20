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
	"crypto/rand"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
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

func NewEncryptColumn[T any](key string) *EncryptColumn[T] {
	return &EncryptColumn[T]{
		Key:   key,
		Valid: true,
	}
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

// Value 返回加密后的值
func (e *EncryptColumn[T]) Value() (driver.Value, error) {
	var val any = e.Val
	switch valT := val.(type) {
	case string:
		return e.aesEncrypt([]byte(valT))
	case []byte:
		return e.aesEncrypt(valT)
	case map[any]any, []any, *any, uintptr:
		return nil, notSupportType(val)
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64,
		float32, float64, complex64, complex128:
		buffer := new(bytes.Buffer)
		err := binary.Write(buffer, binary.BigEndian, val)
		if err != nil {
			return nil, err
		}
		return e.aesEncrypt(buffer.Bytes())
	case int, uint:
		return nil, notSpecifyInt
	default:
		marshal, err := json.Marshal(e.Val)
		if err != nil {
			return nil, err
		}
		return e.aesEncrypt(marshal)
	}
}

func (e *EncryptColumn[T]) Scan(src any) error {
	switch value := src.(type) {
	case []byte:
		decrBytes, err := e.aesDecrypt(value)
		if err != nil {
			return nil
		}
		return e.setVal(decrBytes)
	case *[]byte:
		decrBytes, err := e.aesDecrypt(*value)
		if err != nil {
			return nil
		}
		return e.setVal(decrBytes)
	case string:
		decrBytes, err := e.aesDecrypt([]byte(value))
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
