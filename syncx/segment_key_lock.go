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
	"hash/fnv"
	"sync"
)

// SegmentKeysLock 部分key lock结构定义
type SegmentKeysLock struct {
	locks []*sync.RWMutex
}

// NewSegmentKeysLock 创建 SegmentKeysLock 示例
func NewSegmentKeysLock(size int) *SegmentKeysLock {
	locks := make([]*sync.RWMutex, size)
	for i := range locks {
		locks[i] = &sync.RWMutex{}
	}
	return &SegmentKeysLock{
		locks: locks,
	}
}

// hash 索引锁的hash函数
func (s *SegmentKeysLock) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// RLock 读锁加锁
func (s *SegmentKeysLock) RLock(key string) {
	hash := s.hash(key)
	lock := s.locks[hash%uint32(len(s.locks))]
	lock.RLock()
}

// RUnlock 读锁解锁
func (s *SegmentKeysLock) RUnlock(key string) {
	hash := s.hash(key)
	lock := s.locks[hash%uint32(len(s.locks))]
	lock.RUnlock()
}

// Lock 写锁加锁
func (s *SegmentKeysLock) Lock(key string) {
	hash := s.hash(key)
	lock := s.locks[hash%uint32(len(s.locks))]
	lock.Lock()
}

// Unlock 写锁解锁
func (s *SegmentKeysLock) Unlock(key string) {
	hash := s.hash(key)
	lock := s.locks[hash%uint32(len(s.locks))]
	lock.Unlock()
}
