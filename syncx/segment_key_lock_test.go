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

package syncx

import (
	"sync"
	"testing"
	"time"
)

func TestSegmentKeysLock(t *testing.T) {
	s := NewSegmentKeysLock(10)
	key := "test_key"

	var wg sync.WaitGroup
	wg.Add(2)

	// 写 goroutine
	go func() {
		defer wg.Done()
		s.Lock(key)
		defer s.Unlock(key)

		// 模拟写操作
		time.Sleep(100 * time.Millisecond)
	}()

	// 等待一段时间以确保写 goroutine 先获取锁
	time.Sleep(50 * time.Millisecond)

	// 读 goroutine
	go func() {
		defer wg.Done()
		s.RLock(key)
		defer s.RUnlock(key)

		// 如果读写锁工作正常，这个打印语句应该在写 goroutine 完成后才执行
		t.Log("Read operation executed")
	}()

	wg.Wait()
}
