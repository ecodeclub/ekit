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
	"math/bits"
	"sync/atomic"
	"time"
)

var _ Strategy = (*AdaptiveTimeoutRetryStrategy)(nil)

type AdaptiveTimeoutRetryStrategy struct {
	strategy   Strategy // 基础重试策略
	threshold  int      // 超时比率阈值 (单位：比特数量)
	ringBuffer []uint64 // 比特环（滑动窗口存储超时信息）
	reqCount   uint64   // 请求数量
	bufferLen  int      // 滑动窗口长度
	bitCnt     uint64   // 比特位总数
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
	count = count % s.bitCnt
	// 对2^x进行取模或者整除运算时可以用位运算代替除法和取模
	// count / 64 可以转换成 count >> 6。 位运算会更高效。
	idx := count >> 6
	// count % 64 可以转换成 count & 63
	bitPos := count & 63
	old := atomic.LoadUint64(&s.ringBuffer[idx])
	atomic.StoreUint64(&s.ringBuffer[idx], old&^(uint64(1)<<bitPos))
}

func (s *AdaptiveTimeoutRetryStrategy) markFail() {
	count := atomic.AddUint64(&s.reqCount, 1)
	count = count % s.bitCnt
	idx := count >> 6
	bitPos := count & 63
	old := atomic.LoadUint64(&s.ringBuffer[idx])
	// (uint64(1)<<bitPos) 将目标位设置为1
	atomic.StoreUint64(&s.ringBuffer[idx], old|(uint64(1)<<bitPos))
}

func (s *AdaptiveTimeoutRetryStrategy) getFailed() int {
	var failCount int
	for i := 0; i < len(s.ringBuffer); i++ {
		v := atomic.LoadUint64(&s.ringBuffer[i])
		failCount += bits.OnesCount64(v)
	}
	return failCount
}

func NewAdaptiveTimeoutRetryStrategy(strategy Strategy, bufferLen, threshold int) *AdaptiveTimeoutRetryStrategy {
	return &AdaptiveTimeoutRetryStrategy{
		strategy:   strategy,
		threshold:  threshold,
		bufferLen:  bufferLen,
		ringBuffer: make([]uint64, bufferLen),
		bitCnt:     uint64(64) * uint64(bufferLen),
	}
}
