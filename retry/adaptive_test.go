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
			want:          NewAdaptiveTimeoutRetryStrategy(&MockStrategy{}, 16, 50),
			wantErr:       nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAdaptiveTimeoutRetryStrategy(tt.strategy, tt.ringBufferLen, tt.threshold)
			assert.Equal(t, tt.want, s)
		})
	}
}

func TestAdaptiveTimeoutRetryStrategy_Next(t *testing.T) {
	baseStrategy := &MockStrategy{}
	strategy := NewAdaptiveTimeoutRetryStrategy(baseStrategy, 16, 50)

	tests := []struct {
		name      string
		wantDelay time.Duration
		wantOk    bool
	}{
		{
			name:      "error below threshold",
			wantDelay: 1 * time.Second,
			wantOk:    true,
		},
		{
			name:      "error above threshold",
			wantDelay: 1 * time.Second,
			wantOk:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay, ok := strategy.Next()
			assert.Equal(t, tt.wantDelay, delay)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func ExampleAdaptiveTimeoutRetryStrategy_Next() {
	baseStrategy, err := NewExponentialBackoffRetryStrategy(time.Second, time.Second*5, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	strategy := NewAdaptiveTimeoutRetryStrategy(baseStrategy, 16, 50)
	interval, ok := strategy.Next()
	for ok {
		fmt.Println(interval)
		interval, ok = strategy.Next()
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

func (m MockStrategy) Next() (time.Duration, bool) {
	return 1 * time.Second, true
}

func (m MockStrategy) Report(err error) Strategy {
	return m
}
