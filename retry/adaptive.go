package retry

import (
	"context"
	"fmt"
	"math/bits"
	"sync/atomic"
	"time"

	"github.com/ecodeclub/ekit/internal/errs"
)

type AdaptiveTimeoutRetryStrategy struct {
	strategy      Strategy // 基础重试策略
	threshold     int      // 超时比率阈值 (单位：比特数量)
	ringBuffer    []uint64 // 比特环（滑动窗口存储超时信息）
	reqCount      uint64   // 当前滑动窗口内超时的数量
	ringBufferLen int      // 滑动窗口长度
}

func (s *AdaptiveTimeoutRetryStrategy) Next(ctx context.Context, err error) (time.Duration, bool) {
	if err == nil {
		s.markSuccess()
		return 0, false
	}
	failCount := s.getFailed()
	s.markFail()
	if failCount >= s.threshold {
		return 0, false
	}
	return s.strategy.Next(ctx, err)
}

func (s *AdaptiveTimeoutRetryStrategy) markSuccess() {
	count := atomic.AddUint64(&s.reqCount, 1)
	count = count % (uint64(64) * uint64(len(s.ringBuffer)))
	idx := count / 64
	bitPos := count % 64
	for {
		old := atomic.LoadUint64(&s.ringBuffer[idx])
		// 检查 old 的第 bitPos 位是否为 1。如果结果为 0，表示该位为 0，即没有记录失败
		if old&(uint64(1)<<bitPos) == 0 {
			break
		}
		if atomic.CompareAndSwapUint64(&s.ringBuffer[idx], old, old&^(uint64(1)<<bitPos)) {
			break
		}
	}
}

func (s *AdaptiveTimeoutRetryStrategy) markFail() {
	count := atomic.AddUint64(&s.reqCount, 1)
	count = count % (uint64(64) * uint64(len(s.ringBuffer)))
	idx := count / 64
	bitPos := count % 64
	for {
		old := atomic.LoadUint64(&s.ringBuffer[idx])
		// 检查 old 的第 bitPos 位是否为1。如果结果不等于0，表示该位已经被设置为1。
		if old&(uint64(1)<<bitPos) != 0 {
			// 已被设置为1
			break
		}
		// (uint64(1)<<bitPos) 将目标位设置为1
		if atomic.CompareAndSwapUint64(&s.ringBuffer[idx], old, old|(uint64(1)<<bitPos)) {
			break
		}
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
