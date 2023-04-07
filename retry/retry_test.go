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

	"github.com/stretchr/testify/require"

	"github.com/ecodeclub/ekit/internal/errs"

	"github.com/stretchr/testify/assert"
)

func TestFixedIntervalRetryStrategy_Next(t *testing.T) {

	testCases := []struct {
		name     string
		s        *FixedIntervalRetryStrategy
		interval time.Duration

		isContinue bool
	}{
		{
			name: "init case, retries 0",
			s: &FixedIntervalRetryStrategy{
				maxRetries: 3,
				interval:   time.Second,
			},

			interval:   time.Second,
			isContinue: true,
		},
		{
			name: "retries equals to MaxRetries 3 after the increase",
			s: &FixedIntervalRetryStrategy{
				maxRetries: 3,
				interval:   time.Second,
				retries:    2,
			},
			interval:   time.Second,
			isContinue: true,
		},
		{
			name: "retries over MaxRetries after the increase",
			s: &FixedIntervalRetryStrategy{
				maxRetries: 3,
				interval:   time.Second,
				retries:    3,
			},
			interval:   0,
			isContinue: false,
		},
		{
			name: "MaxRetries equals to 0",
			s: &FixedIntervalRetryStrategy{
				maxRetries: 0,
				interval:   time.Second,
			},
			interval:   time.Second,
			isContinue: true,
		},
		{
			name: "negative MaxRetries",
			s: &FixedIntervalRetryStrategy{
				maxRetries: -1,
				interval:   time.Second,
				retries:    0,
			},
			interval:   time.Second,
			isContinue: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			interval, isContinue := tt.s.Next()
			assert.Equal(t, tt.interval, interval)
			assert.Equal(t, tt.isContinue, isContinue)
		})
	}
}

func TestFixedIntervalRetryStrategy_New(t *testing.T) {
	testCases := []struct {
		name       string
		maxRetries int32
		interval   time.Duration

		want    *FixedIntervalRetryStrategy
		wantErr error
	}{
		{
			name:       "no error",
			maxRetries: 5,
			interval:   time.Second,

			want: &FixedIntervalRetryStrategy{
				maxRetries: 5,
				interval:   time.Second,
			},
			wantErr: nil,
		},
		{
			name:       "returns error, interval equals to 0",
			maxRetries: 5,
			interval:   0,

			want:    nil,
			wantErr: errs.NewErrInvalidIntervalValue(0),
		},
		{
			name:       "returns error, interval equals to -1",
			maxRetries: 5,
			interval:   -1,

			want:    nil,
			wantErr: errs.NewErrInvalidIntervalValue(-1),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFixedIntervalRetryStrategy(tt.maxRetries, tt.interval)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

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
			want: &ExponentialBackoffRetryStrategy{
				initialInterval: 2 * time.Second,
				maxInterval:     2 * time.Minute,
				maxRetries:      5,
			},
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
		n := 100

		s, err := NewExponentialBackoffRetryStrategy(1*time.Second, 4*time.Second, 0)
		require.NoError(t, err)

		wantIntervals := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
		length := n - len(wantIntervals)
		for i := 0; i < length; i++ {
			wantIntervals = append(wantIntervals, 4*time.Second)
		}

		intervals := make([]time.Duration, 0, n)
		for i := 0; i < n; i++ {
			res, _ := s.Next()
			intervals = append(intervals, res)
		}
		assert.Equal(t, wantIntervals, intervals)
	})

	t.Run("maxRetries equals -1", func(t *testing.T) {
		n := 100

		s, err := NewExponentialBackoffRetryStrategy(1*time.Second, 4*time.Second, -1)
		require.NoError(t, err)

		wantIntervals := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
		length := n - len(wantIntervals)
		for i := 0; i < length; i++ {
			wantIntervals = append(wantIntervals, 4*time.Second)
		}

		intervals := make([]time.Duration, 0, n)
		for i := 0; i < n; i++ {
			res, _ := s.Next()
			intervals = append(intervals, res)
		}
		assert.Equal(t, wantIntervals, intervals)
	})
}
