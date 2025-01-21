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
)

// MultipleBytes 是一个实现了 io.Reader 和 io.Writer 接口的结构体
// 它可以安全地在多个 goroutine 之间共享
type MultipleBytes struct {
	buf   []byte
	pos   int
	mutex sync.RWMutex
}

// NewMultipleBytes 创建一个新的 MultipleBytes 实例
// capacity 参数用于预分配内部缓冲区的容量
func NewMultipleBytes(capacity int) *MultipleBytes {
	return &MultipleBytes{
		buf: make([]byte, 0, capacity),
	}
}

// Read 实现 io.Reader 接口
// 从当前位置读取数据到 p 中，如果没有数据可读返回 io.EOF
func (m *MultipleBytes) Read(p []byte) (n int, err error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.pos >= len(m.buf) {
		return 0, io.EOF
	}

	n = copy(p, m.buf[m.pos:])
	m.pos += n
	return n, nil
}

// Write 实现 io.Writer 接口
// 将 p 中的数据写入到内部缓冲区
func (m *MultipleBytes) Write(p []byte) (n int, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.buf = append(m.buf, p...)
	return len(p), nil
}

// Len 返回当前缓冲区中的数据长度
func (m *MultipleBytes) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.buf)
}

// Cap 返回当前缓冲区的容量
func (m *MultipleBytes) Cap() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return cap(m.buf)
}

// Reset 重置读取位置到开始处
func (m *MultipleBytes) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.pos = 0
}

// Clear 清空缓冲区并重置读取位置
func (m *MultipleBytes) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.buf = m.buf[:0]
	m.pos = 0
}

// Bytes 返回内部缓冲区的副本
func (m *MultipleBytes) Bytes() []byte {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	res := make([]byte, len(m.buf))
	copy(res, m.buf)
	return res
}