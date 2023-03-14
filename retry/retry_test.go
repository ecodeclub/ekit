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

	"github.com/stretchr/testify/assert"
)

func TestEqualRetryStrategy_Next(t *testing.T) {

	testCases := []struct {
		name               string
		equalRetryStrategy *EqualRetryStrategy

		internal   time.Duration
		isContinue bool
	}{
		{
			name:               "retries = 0, MaxRetries = 3, interval = 1s",
			equalRetryStrategy: NewEqualRetryStrategy(3, time.Second),
			internal:           time.Second,
			isContinue:         true,
		},
		{
			name: "retries = 3, MaxRetries = 3, interval = 1s",
			equalRetryStrategy: &EqualRetryStrategy{
				MaxRetries: 3,
				Interval:   time.Second,
				retries:    3,
			},
			internal:   time.Second,
			isContinue: true,
		},
		{
			name: "retries = 4, MaxRetries = 3, interval = 1s",
			equalRetryStrategy: &EqualRetryStrategy{
				MaxRetries: 3,
				Interval:   time.Second,
				retries:    4,
			},
			internal:   time.Second,
			isContinue: false,
		},
		{
			name: "retries = 0, MaxRetries = 0, interval = 1s",
			equalRetryStrategy: &EqualRetryStrategy{
				MaxRetries: 3,
				Interval:   time.Second,
				retries:    3,
			},
			internal:   time.Second,
			isContinue: true,
		},
		{
			name: "retries = 0, MaxRetries = -1, interval = 1s",
			equalRetryStrategy: &EqualRetryStrategy{
				MaxRetries: 3,
				Interval:   time.Second,
				retries:    3,
			},
			internal:   time.Second,
			isContinue: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			interval, isContinue := tt.equalRetryStrategy.Next()
			assert.Equal(t, tt.internal, interval)
			assert.Equal(t, tt.isContinue, isContinue)
		})
	}
}
