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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExponentialBackoffRetryStrategy_New(t *testing.T) {
	testCases := []struct {
		name            string
		initialInterval time.Duration
		maxInterval     time.Duration
		maxRetries      int32
		want            *ExponentialBackoffRetryStrategy
		wantErr         error
	}{
		{
			name:            "no error",
			initialInterval: 2 * time.Second,
			maxInterval:     2 * time.Minute,
			maxRetries:      5,
			want: func() *ExponentialBackoffRetryStrategy {
				s, err := NewExponentialBackoffRetryStrategy(2*time.Second, 2*time.Minute, 5)
				require.NoError(t, err)
				return s
			}(),
			wantErr: nil,
		},
		{
			name:            "return error, initialInterval equals 0",
			initialInterval: 0 * time.Second,
			maxInterval:     2 * time.Minute,
			maxRetries:      5,
			want:            nil,
			wantErr:         fmt.Errorf("ekit: 无效的间隔时间 %d, 预期值应大于 0", 0*time.Second),
		},
		{
			name:            "return error, initialInterval equals -60",
			initialInterval: -1 * time.Second,
			maxInterval:     2 * time.Minute,
			maxRetries:      5,
			want:            nil,
			wantErr:         fmt.Errorf("ekit: 无效的间隔时间 %d, 预期值应大于 0", -1*time.Second),
		},
		{
			name:            "return error, maxInternal > initialInterval",
			initialInterval: 5 * time.Second,
			maxInterval:     1 * time.Second,
			maxRetries:      5,
			want:            nil,
			wantErr:         fmt.Errorf("ekit: 最大重试间隔的时间 [%d] 应大于等于初始重试的间隔时间 [%d] ", 1*time.Second, 5*time.Second),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewExponentialBackoffRetryStrategy(tt.initialInterval, tt.maxInterval, tt.maxRetries)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, s)
		})
	}
}

func TestExponentialBackoffRetryStrategy_Next(t *testing.T) {
	testCases := []struct {
		name     string
		strategy *ExponentialBackoffRetryStrategy

		wantIntervals []time.Duration
	}{
		{
			name: "stop if retries reaches maxRetries",
			strategy: func() *ExponentialBackoffRetryStrategy {
				s, err := NewExponentialBackoffRetryStrategy(1*time.Second, 10*time.Second, 3)
				require.NoError(t, err)
				return s
			}(),

			wantIntervals: []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second},
		},
		{
			name: "initialInterval over maxInterval",
			strategy: func() *ExponentialBackoffRetryStrategy {
				s, err := NewExponentialBackoffRetryStrategy(1*time.Second, 4*time.Second, 5)
				require.NoError(t, err)
				return s
			}(),

			wantIntervals: []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second, 4 * time.Second, 4 * time.Second},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			intervals := make([]time.Duration, 0)
			for {
				if interval, ok := tt.strategy.Next(); ok {
					intervals = append(intervals, interval)
				} else {
					break
				}
			}
			assert.Equal(t, tt.wantIntervals, intervals)
		})
	}
}

// 指数退避重试策略子测试函数，无限重试
func TestExponentialBackoffRetryStrategy_Next4InfiniteRetry(t *testing.T) {
	t.Run("maxRetries equals 0", func(t *testing.T) {
		testNext4InfiniteRetry(t, 0)
	})

	t.Run("maxRetries equals -1", func(t *testing.T) {
		testNext4InfiniteRetry(t, -1)
	})
}

func ExampleExponentialBackoffRetryStrategy_Next() {
	// 注意，因为在例子里面我们设置初始的重试间隔是 1s，最大重试间隔是 5s
	// 所以在前面四次，重试间隔都是在增长的，每次变为原来的2倍。
	// 在触及到了最大重试间隔之后，就一直以最大重试间隔来重试。
	retry, err := NewExponentialBackoffRetryStrategy(time.Second, time.Second*5, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	interval, ok := retry.Next()
	for ok {
		fmt.Println(interval)
		interval, ok = retry.Next()
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
