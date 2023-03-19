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

// ExponentialBackoffRetryStrategy 指数间隔重试
type ExponentialBackoffRetryStrategy struct {
	// 初始重试间隔
	initInterval time.Duration
	interval     time.Duration
	// 最大重试间隔
	maxInterval time.Duration
	// 最大重试次数
	maxRetries int32
	retries    int32
}

func (e *ExponentialBackoffRetryStrategy) Next() (time.Duration, bool) {
	retries := atomic.AddInt32(&e.retries, 1)
	expRetries := math.Pow(2, float64(retries-1))
	if expRetries > math.MaxInt32 {
		e.interval = e.maxInterval
	} else {
		e.interval = time.Duration(expRetries) * e.initInterval
	}

	if e.interval > e.maxInterval {
		e.interval = e.maxInterval
	}
	if e.maxRetries <= 0 || e.retries <= e.maxRetries {
		return e.interval, true
	}
	return 0, false
}
func NewExponentialIntervalRetryStrategy(maxRetries int32, initInterval time.Duration, maxInterval time.Duration) (*ExponentialBackoffRetryStrategy, error) {
	if initInterval <= 0 {
		return nil, errs.NewErrInvalidIntervalValue(initInterval)
	}

	if maxInterval <= 0 {
		return nil, errs.NewErrInvalidIntervalValue(maxInterval)
	}
	if initInterval > maxInterval {
		return nil, errs.NewErrInvalidIntervalValue(initInterval)
	}
	return &ExponentialBackoffRetryStrategy{
		initInterval: initInterval,
		maxInterval:  maxInterval,
		maxRetries:   maxRetries,
	}, nil
}
