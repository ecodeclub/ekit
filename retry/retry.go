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

package retry

import (
	"sync"
	"time"
)

type EqualRetryStrategy struct {
	MaxRetries int           // 最大重试次数，如果是 0 或负数，表示无限重试
	Interval   time.Duration // 重试间隔时间
	retries    int           // 当前重试次数
	mu         sync.Mutex    // 用于保证并发安全
}

func NewEqualRetryStrategy(maxRetries int, interval time.Duration) *EqualRetryStrategy {
	return &EqualRetryStrategy{
		MaxRetries: maxRetries,
		Interval:   interval,
	}
}

func (s *EqualRetryStrategy) Next() (time.Duration, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.MaxRetries <= 0 || s.retries <= s.MaxRetries {
		s.retries++
		return s.Interval, true
	}
	return s.Interval, false
}
