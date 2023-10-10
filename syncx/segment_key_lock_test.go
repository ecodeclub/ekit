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
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegmentKeysLock(t *testing.T) {
	s := NewSegmentKeysLock(10)
	key := "test_key"
	var wg sync.WaitGroup
	wg.Add(2)
	var writeDone atomic.Bool
	var readStarted atomic.Bool
	val := false
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()

	// 写 goroutine
	go func() {
		defer wg.Done()
		s.Lock(key)
		val = true       // 加写锁写
		s.Unlock(key)
		writeDone.Store(true)
		cond.Broadcast()
	}()

	// 读 goroutine
	go func() {
		defer wg.Done()
		cond.L.Lock()
		defer cond.L.Unlock()

		// 等待写操作完成
		for !writeDone.Load() {
			cond.Wait()
		}

		readStarted.Store(true)
		cond.Broadcast()
		s.RLock(key)
		assert.Equal(t, true, val, "Read lock err")  // 加读锁读
		defer s.RUnlock(key)
	}()

	// 等待读操作开始
	for !readStarted.Load() {
		cond.Wait()
	}

	// 检查写操作是否已完成，防止意外情况导致读优先写发生
	assert.Equal(t, true, writeDone.Load(), "Write operation did not complete before read operation started")

	cond.L.Unlock()
	wg.Wait()
}
