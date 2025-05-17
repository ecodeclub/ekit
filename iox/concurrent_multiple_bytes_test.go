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

package iox

import (
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentMultipleBytes(t *testing.T) {
	t.Run("基本读写功能", func(t *testing.T) {
		cmb := NewConcurrentMultipleBytes(2)
		data := []byte{1, 2, 3, 4}

		// 写入数据
		n, err := cmb.Write(data)
		assert.Equal(t, len(data), n)
		assert.Nil(t, err)

		// 读取数据
		read := make([]byte, 4)
		n, err = cmb.Read(read)
		assert.Equal(t, len(data), n)
		assert.Nil(t, err)
		assert.Equal(t, data, read[:n])
	})

	t.Run("并发读写", func(t *testing.T) {
		cmb := NewConcurrentMultipleBytes(3)
		var wg sync.WaitGroup

		// 并发写入
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(val byte) {
				defer wg.Done()
				n, err := cmb.Write([]byte{val})
				assert.Equal(t, 1, n)
				assert.Nil(t, err)
			}(byte(i + 1))
		}
		wg.Wait()

		// 并发读取
		results := make([][]byte, 3)
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				read := make([]byte, 1)
				n, err := cmb.Read(read)
				if err != nil && err != io.EOF {
					assert.Nil(t, err)
					return
				}
				results[idx] = read[:n]
			}(i)
		}
		wg.Wait()

		// 验证总读取字节数
		total := 0
		for _, res := range results {
			total += len(res)
		}
		assert.Equal(t, 3, total)
	})

	t.Run("边界场景", func(t *testing.T) {
		cmb := NewConcurrentMultipleBytes(1)

		// 空切片写入
		n, err := cmb.Write([]byte{})
		assert.Equal(t, 0, n)
		assert.Nil(t, err)

		// 空切片读取
		read := make([]byte, 1)
		n, err = cmb.Read(read)
		assert.Equal(t, 0, n)
		assert.Equal(t, io.EOF, err)
	})

	t.Run("Reset功能", func(t *testing.T) {
		cmb := NewConcurrentMultipleBytes(1)
		data := []byte{1, 2}

		// 写入数据
		n, err := cmb.Write(data)
		assert.Equal(t, len(data), n)
		assert.Nil(t, err)

		// 读取一部分
		read := make([]byte, 1)
		n, err = cmb.Read(read)
		assert.Equal(t, 1, n)
		assert.Nil(t, err)

		// 重置
		cmb.Reset()

		// 重新读取
		read = make([]byte, 2)
		n, err = cmb.Read(read)
		assert.Equal(t, 2, n)
		assert.Nil(t, err)
		assert.Equal(t, data, read[:n])
	})
}
