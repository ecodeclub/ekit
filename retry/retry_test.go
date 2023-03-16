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
	"testing"
	"time"

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

func TestExponentialIntervalRetry_Next(t *testing.T) {
	testCases := []struct {
		name         string
		strategy     *ExponentialIntervalRetryStrategy
		wantStrategy *ExponentialIntervalRetryStrategy
		wantInterval time.Duration
		isContinue   bool
	}{
		{
			name: "test begin",
			strategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    3,
			},
			wantStrategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    3,
				retries:       1,
			},
			wantInterval: time.Second,
			isContinue:   true,
		},
		{
			name: "test normal",
			strategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    3,
				retries:       1,
			},
			wantStrategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      2 * time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    3,
				retries:       2,
			},
			wantInterval: 2 * time.Second,
			isContinue:   true,
		},
		{
			name: "test max interval",
			strategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      4 * time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    5,
				retries:       3,
			},
			wantStrategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      5 * time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    5,
				retries:       4,
			},
			wantInterval: 5 * time.Second,
			isContinue:   true,
		},
		{
			name: "test max retires",
			strategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      5 * time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    5,
				retries:       5,
			},
			wantStrategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      5 * time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    5,
				retries:       6,
			},
			wantInterval: 5 * time.Second,
			isContinue:   false,
		},
		{
			name: "test zero retires",
			strategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      5 * time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    0,
			},
			wantStrategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    0,
				retries:       1,
			},
			wantInterval: time.Second,
			isContinue:   true,
		},
		{
			name: "test zero retires",
			strategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      4 * time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    0,
				retries:       3,
			},
			wantStrategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				interval:      5 * time.Second,
				maxInterval:   5 * time.Second,
				maxRetries:    0,
				retries:       4,
			},
			wantInterval: 5 * time.Second,
			isContinue:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interval, isContinue := tc.strategy.Next()
			assert.Equal(t, tc.wantInterval, interval)
			assert.Equal(t, tc.wantStrategy, tc.strategy)
			assert.Equal(t, tc.isContinue, isContinue)
		})
	}
}

func TestExponentialIntervalRetryStrategy_New(t *testing.T) {
	testCases := []struct {
		name          string
		wantStrategy  *ExponentialIntervalRetryStrategy
		beginInterval time.Duration
		maxInterval   time.Duration
		maxRetries    int32
		wantErr       error
	}{
		{
			name: "test normal",
			wantStrategy: &ExponentialIntervalRetryStrategy{
				beginInterval: time.Second,
				maxInterval:   time.Minute,
				maxRetries:    5,
			},
			beginInterval: time.Second,
			maxInterval:   time.Minute,
			maxRetries:    5,
		},
		{
			name:          "test invalide beginInterval",
			beginInterval: time.Duration(-1),
			maxInterval:   time.Minute,
			maxRetries:    5,
			wantErr:       errs.NewErrInvalidIntervalValue(time.Duration(-1)),
		},
		{
			name:          "test invalide maxInterval",
			beginInterval: time.Second,
			maxInterval:   time.Duration(-1),
			maxRetries:    5,
			wantErr:       errs.NewErrInvalidIntervalValue(time.Duration(-1)),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strategy, err := NewExponentialIntervalRetryStrategy(tc.maxRetries, tc.beginInterval, tc.maxInterval)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantStrategy, strategy)
		})
	}
}
