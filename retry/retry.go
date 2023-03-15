package retry

import (
	"math"
	"sync"
	"time"
)

type Strategy interface {
	// Next 返回下一次重试的间隔，如果不需要继续重试，那么第二参数返回 false
	Next() (time.Duration, bool)
}

// FixIntervalRetry 固定间隔重试
type FixIntervalRetry struct {
	// 重试间隔
	Interval time.Duration
	// 最大次数
	Max int
	cnt int
	mu  sync.Mutex
}

func (f *FixIntervalRetry) Next() (time.Duration, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.Max <= 0 {
		return f.Interval, true
	}
	f.cnt++
	return f.Interval, f.cnt <= f.Max
}

// ExponentialIntervalRetry 指数间隔重试
type ExponentialIntervalRetry struct {
	// 初始重试间隔
	BeginInterval time.Duration
	interval      time.Duration
	// 最大重试间隔
	MaxInterval time.Duration
	// 最大次数
	Max int
	cnt int
	mu  sync.Mutex
}

func (e *ExponentialIntervalRetry) Next() (time.Duration, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cnt++
	e.interval = time.Duration(math.Exp2(float64(e.cnt))) * e.BeginInterval
	if e.interval > e.MaxInterval {
		e.interval = e.MaxInterval
	}
	if e.Max <= 0 {
		return e.interval, true
	}
	return e.interval, e.cnt <= e.Max
}
