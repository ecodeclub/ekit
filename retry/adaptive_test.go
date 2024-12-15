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
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAdaptiveTimeoutRetryStrategy_New(t *testing.T) {
	testCases := []struct {
		name          string
		threshold     int
		ringBufferLen int
		strategy      Strategy
		want          *AdaptiveTimeoutRetryStrategy
		wantErr       error
	}{
		{
			name:          "valid strategy and threshold",
			strategy:      &MockStrategy{}, // 假设有一个 MockStrategy 用于测试
			threshold:     50,
			ringBufferLen: 16,
			want: func() *AdaptiveTimeoutRetryStrategy {
				s, err := NewAdaptiveTimeoutRetryStrategy(&MockStrategy{}, WithThreshold(50), WithRingBufferLen(16))
				require.NoError(t, err)
				return s
			}(),
			wantErr: nil,
		},
		{
			name:          "threshold less than or equal to zero",
			strategy:      &MockStrategy{},
			ringBufferLen: 16,
			threshold:     0,
			want:          nil,
			wantErr:       fmt.Errorf("ekit: 失效比率阈值 [%d]", 0),
		},
		{
			name:          "ring buffer len less than or equal to zero",
			strategy:      &MockStrategy{},
			threshold:     10,
			ringBufferLen: 0,
			want:          nil,
			wantErr:       fmt.Errorf("ekit: 无效的滑动窗口长度 [%d]", 0),
		},
		{
			name:      "strategy is nil",
			strategy:  nil,
			threshold: 10,
			want:      nil,
			wantErr:   fmt.Errorf("ekit: strategy 不能为空"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewAdaptiveTimeoutRetryStrategy(tt.strategy, WithThreshold(tt.threshold), WithRingBufferLen(tt.ringBufferLen))
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, s)
		})
	}
}

func TestAdaptiveTimeoutRetryStrategy_Next(t *testing.T) {
	baseStrategy := &MockStrategy{}
	strategy, err := NewAdaptiveTimeoutRetryStrategy(baseStrategy, WithThreshold(50), WithRingBufferLen(16))
	require.NoError(t, err)

	tests := []struct {
		name      string
		err       error
		wantDelay time.Duration
		wantOk    bool
	}{
		{
			name:      "error below threshold",
			err:       errors.New("test error"),
			wantDelay: 1 * time.Second,
			wantOk:    true,
		},
		{
			name:      "error above threshold",
			err:       errors.New("test error"),
			wantDelay: 1 * time.Second,
			wantOk:    true,
		},
		{
			name:      "not retry",
			wantDelay: 0,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay, ok := strategy.Next(context.Background(), tt.err)
			assert.Equal(t, tt.wantDelay, delay)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

// 测试场景
// 阈值是50
// 2000个请求 有1500个成功的 有500个失败的 最后统计500个失败的有50个可以执行 有450个不能执行 1500成功的都能执行
func TestAdaptiveTimeoutRetryStrategy_Next_Concurrent(t *testing.T) {
	// 创建一个基础策略
	baseStrategy := &MockStrategy{}

	// 创建升级版自适应策略，设置阈值为50
	strategy, err := NewAdaptiveTimeoutRetryStrategy(baseStrategy,
		WithThreshold(50), WithRingBufferLen(16))
	assert.Nil(t, err)

	var wg sync.WaitGroup
	var successCount, errCount int64
	mockErr := errors.New("mock error")

	// 并发执行2000个请求
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			// 前1500个请求成功，后500个失败
			var err error
			if index >= 1500 {
				err = mockErr
			}
			_, allowed := strategy.Next(context.Background(), err)
			if err != nil {
				// 失败请求的统计
				if allowed {
					atomic.AddInt64(&successCount, 1)
				} else {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}(i)
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 验证结果：期望大约50个失败请求可以执行，450个被拒绝
	// 由于是环形缓冲区和并发执行，可能会有一些误差，这里使用一个合理的范围进行判断
	finalSuccessCount := int(atomic.LoadInt64(&successCount))
	finalErrCount := int(atomic.LoadInt64(&errCount))
	if finalSuccessCount < 45 || finalSuccessCount > 55 {
		t.Errorf("期望大约50个失败请求被允许执行，实际允许执行的失败请求数量为: %d", finalSuccessCount)
	}

	if finalErrCount < 445 || finalErrCount > 455 {
		t.Errorf("期望大约450个失败请求被拒绝执行，实际被拒绝的失败请求数量为: %d", finalErrCount)
	}
}

func ExampleAdaptiveTimeoutRetryStrategy_Next() {
	baseStrategy, err := NewExponentialBackoffRetryStrategy(time.Second, time.Second*5, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	strategy, err := NewAdaptiveTimeoutRetryStrategy(baseStrategy,
		WithThreshold(50), WithRingBufferLen(16))
	if err != nil {
		fmt.Println(err)
		return
	}
	nextErr := errors.New("test error")
	interval, ok := strategy.Next(context.Background(), nextErr)
	for ok {
		fmt.Println(interval)
		interval, ok = strategy.Next(context.Background(), nextErr)
	}
	// Output:
	// 1s
	// 2s
	// 4s
	// 5s
	// 5s
	// 5s
	// 5s
	// 5s
	// 5s
	// 5s
}

type MockStrategy struct {
}

func (m MockStrategy) Next(ctx context.Context, err error) (time.Duration, bool) {
	return 1 * time.Second, true
}
