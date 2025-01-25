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
	"sync"
)

// ConcurrentMultipleBytes 是 MultipleBytes 的线程安全装饰器
type ConcurrentMultipleBytes struct {
	mb   *MultipleBytes
	lock sync.Mutex
}

// NewConcurrentMultipleBytes 创建一个新的线程安全的 MultipleBytes 实例
// sliceCount 参数用于预分配内部切片数组的容量
func NewConcurrentMultipleBytes(sliceCount int) *ConcurrentMultipleBytes {
	return &ConcurrentMultipleBytes{
		mb: NewMultipleBytes(sliceCount),
	}
}

// Read 实现 io.Reader 接口
// 从当前位置读取数据到 p 中，如果没有数据可读返回 io.EOF
func (c *ConcurrentMultipleBytes) Read(p []byte) (n int, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.mb.Read(p)
}

// Write 实现 io.Writer 接口
// 将 p 中的数据写入到内部缓冲区
func (c *ConcurrentMultipleBytes) Write(p []byte) (n int, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.mb.Write(p)
}

// Reset 重置读取位置到开始处
func (c *ConcurrentMultipleBytes) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.mb.Reset()
}
