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
	"math"
	"sync/atomic"
	"time"

	"github.com/ecodeclub/ekit/internal/errs"
)

type FixedIntervalRetryStrategy struct {
	maxRetries int32         // 最大重试次数，如果是 0 或负数，表示无限重试
	interval   time.Duration // 重试间隔时间
	retries    int32         // 当前重试次数
}

func NewFixedIntervalRetryStrategy(maxRetries int32, interval time.Duration) (*FixedIntervalRetryStrategy, error) {
	if interval <= 0 {
		return nil, errs.NewErrInvalidIntervalValue(interval)
	}
	return &FixedIntervalRetryStrategy{
		maxRetries: maxRetries,
		interval:   interval,
	}, nil
}

func (s *FixedIntervalRetryStrategy) Next() (time.Duration, bool) {
	retries := atomic.AddInt32(&s.retries, 1)
	if s.maxRetries <= 0 || retries <= s.maxRetries {
		return s.interval, true
	}
	return 0, false
}

type ExponentialBackoffRetryStrategy struct {
	// 初始重试间隔
	initialInterval time.Duration
	// 最大重试间隔
	maxInterval time.Duration
	// 最大重试次数
	maxRetries int32
	// 当前重试次数
	retries int32
}

func NewExponentialBackoffRetryStrategy(initialInterval, maxInterval time.Duration, maxRetries int32) (*ExponentialBackoffRetryStrategy, error) {
	if initialInterval <= 0 {
		return nil, errs.NewErrInvalidIntervalValue(initialInterval)
	} else if initialInterval > maxInterval {
		return nil, errs.NewErrInvalidMaxIntervalValue(maxInterval, initialInterval)
	}
	return &ExponentialBackoffRetryStrategy{
		initialInterval: initialInterval,
		maxInterval:     maxInterval,
		maxRetries:      maxRetries,
	}, nil
}

func (s *ExponentialBackoffRetryStrategy) Next() (time.Duration, bool) {
	retries := atomic.AddInt32(&s.retries, 1)
	if s.maxRetries <= 0 || retries <= s.maxRetries {
		interval := s.initialInterval * time.Duration(math.Pow(2, float64(retries-1)))
		// 无限重试的情况下，如果 interval 小于等于 0，说明溢出了，返回 maxInterval
		if interval <= 0 || interval > s.maxInterval {
			return s.maxInterval, true
		}
		return interval, true
	}
	return 0, false
}
