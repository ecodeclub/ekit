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
	"fmt"
	"math/bits"
	"sync/atomic"
	"time"

	"github.com/ecodeclub/ekit/internal/errs"
)

var _ Strategy = (*AdaptiveTimeoutRetryStrategy)(nil)

type AdaptiveTimeoutRetryStrategy struct {
	strategy      Strategy // 基础重试策略
	threshold     int      // 超时比率阈值 (单位：比特数量)
	ringBuffer    []uint64 // 比特环（滑动窗口存储超时信息）
	reqCount      uint64   // 请求数量
	ringBufferLen int      // 滑动窗口长度
}

func (s *AdaptiveTimeoutRetryStrategy) Next() (time.Duration, bool) {
	failCount := s.getFailed()
	if failCount >= s.threshold {
		return 0, false
	}
	return s.strategy.Next()
}

func (s *AdaptiveTimeoutRetryStrategy) Report(err error) Strategy {
	if err == nil {
		s.markSuccess()
	} else {
		s.markFail()
	}
	return s
}

func (s *AdaptiveTimeoutRetryStrategy) markSuccess() {
	count := atomic.AddUint64(&s.reqCount, 1)
	count = count % (uint64(64) * uint64(len(s.ringBuffer)))
	// 对2^x进行取模或者整除运算时可以用位运算代替除法和取模
	// count / 64 可以转换成 count >> 6。 位运算会更高效。
	idx := count >> 6
	// count % 64 可以转换成 count & 63
	bitPos := count & 63
	for {
		old := atomic.LoadUint64(&s.ringBuffer[idx])
		// 检查 old 的第 bitPos 位是否为 1。如果结果为 0，表示该位为 0，即没有记录失败
		if old&(uint64(1)<<bitPos) == 0 {
			break
		}
		atomic.StoreUint64(&s.ringBuffer[idx], old&^(uint64(1)<<bitPos))
	}
}

func (s *AdaptiveTimeoutRetryStrategy) markFail() {
	count := atomic.AddUint64(&s.reqCount, 1)
	count = count % (uint64(64) * uint64(len(s.ringBuffer)))
	idx := count >> 6
	bitPos := count & 63
	for {
		old := atomic.LoadUint64(&s.ringBuffer[idx])
		// 检查 old 的第 bitPos 位是否为1。如果结果不等于0，表示该位已经被设置为1。
		if old&(uint64(1)<<bitPos) != 0 {
			// 已被设置为1
			break
		}
		// (uint64(1)<<bitPos) 将目标位设置为1
		atomic.StoreUint64(&s.ringBuffer[idx], old|(uint64(1)<<bitPos))
	}
}

func (s *AdaptiveTimeoutRetryStrategy) getFailed() int {
	var failCount int
	for i := 0; i < len(s.ringBuffer); i++ {
		v := atomic.LoadUint64(&s.ringBuffer[i])
		failCount += bits.OnesCount64(v)
	}
	return failCount
}

func NewAdaptiveTimeoutRetryStrategy(strategy Strategy, opts ...Option) (*AdaptiveTimeoutRetryStrategy, error) {
	if strategy == nil {
		return nil, fmt.Errorf("ekit: strategy 不能为空")
	}

	res := &AdaptiveTimeoutRetryStrategy{
		strategy: strategy,
	}
	for _, opt := range opts {
		opt(res)
	}

	if res.ringBufferLen <= 0 {
		return nil, fmt.Errorf("ekit: 无效的滑动窗口长度 [%d]", res.ringBufferLen)

	}

	if res.threshold <= 0 {
		return nil, errs.NewErrInvalidThresholdValue(res.threshold)
	}

	res.ringBuffer = make([]uint64, res.ringBufferLen)
	return res, nil
}

type Option func(*AdaptiveTimeoutRetryStrategy)

func WithRingBufferLen(l int) Option {
	return func(s *AdaptiveTimeoutRetryStrategy) {
		s.ringBufferLen = l
	}
}

func WithThreshold(t int) Option {
	return func(s *AdaptiveTimeoutRetryStrategy) {
		s.threshold = t
	}
}
