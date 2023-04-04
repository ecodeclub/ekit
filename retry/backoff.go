package retry

import (
	"sync/atomic"
	"time"
)

type BackoffIntervalRetryStrategy struct {
	initialInterval time.Duration // 初始重试间隔
	maxInterval     time.Duration // 最大重试间隔
	maxRetries      int32         // 最大重试次数，如果是 0 或负数，表示无限重试
	multiplier      float64       // 指数增长

	interval time.Duration
	retries  int32
}

func NewBackoffIntervalRetryStrategyWithInitialInterval(initialInterval, maxInterval time.Duration, maxRetries int32, multiplier float64) *BackoffIntervalRetryStrategy {
	if maxInterval == 0 {
		maxInterval = 30 * time.Second
	}
	if multiplier < 1 {
		multiplier = 2
	}

	return &BackoffIntervalRetryStrategy{
		initialInterval: initialInterval,
		maxInterval:     maxInterval,
		maxRetries:      maxRetries,
		multiplier:      multiplier,
		interval:        initialInterval,
	}
}

func (bo *BackoffIntervalRetryStrategy) Next() (time.Duration, bool) {
	atomic.StoreInt64((*int64)(&bo.interval), int64(time.Duration(float64(bo.interval)*bo.multiplier)))
	if bo.maxRetries > 0 && bo.interval > bo.maxInterval {
		return 0, false
	}

	retries := atomic.AddInt32(&bo.retries, 1)
	if bo.maxRetries > 0 && retries > bo.maxRetries {
		return 0, false
	}

	return bo.interval, true
}
