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
	"sync"
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

// ExponentialIntervalRetryStrategy 指数间隔重试
type ExponentialIntervalRetryStrategy struct {
	// 初始重试间隔
	beginInterval time.Duration
	interval      time.Duration
	// 最大重试间隔
	maxInterval time.Duration
	// 最大重试次数
	maxRetries int32
	retries    int32
	mu         sync.Mutex
}

func (e *ExponentialIntervalRetryStrategy) Next() (time.Duration, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.interval = time.Duration(math.Exp2(float64(e.retries))) * e.beginInterval
	atomic.AddInt32(&e.retries, 1)

	if e.interval > e.maxInterval {
		e.interval = e.maxInterval
	}
	if e.maxRetries <= 0 || e.retries <= e.maxRetries {
		return e.interval, true
	}
	return e.maxInterval, false
}
func NewExponentialIntervalRetryStrategy(maxRetries int32, beginInterval time.Duration, maxInterval time.Duration) (*ExponentialIntervalRetryStrategy, error) {
	if beginInterval <= 0 {
		return nil, errs.NewErrInvalidIntervalValue(beginInterval)
	}

	if maxInterval <= 0 {
		return nil, errs.NewErrInvalidIntervalValue(maxInterval)
	}
	return &ExponentialIntervalRetryStrategy{
		beginInterval: beginInterval,
		maxInterval:   maxInterval,
		maxRetries:    maxRetries,
	}, nil
}
