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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	bizErr := errors.New("biz error")
	testCases := []struct {
		name      string
		biz       func() error
		strategy  Strategy
		wantError error
	}{
		{
			name: "第一次就成功",
			biz: func() error {
				t.Log("模拟业务")
				return nil
			},
			strategy: func() Strategy {
				res, _ := NewFixedIntervalRetryStrategy(time.Second, 3)
				return res
			}(),
		},
		{
			name: "重试最终失败",
			biz: func() error {
				return bizErr
			},
			strategy: func() Strategy {
				res, _ := NewFixedIntervalRetryStrategy(time.Second, 3)
				return res
			}(),
			wantError: bizErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			err := Retry(ctx, tc.strategy, tc.biz)
			assert.ErrorIs(t, err, tc.wantError)
		})
	}
}

func ExampleRetry() {
	// 这是你的业务
	bizFunc := func() error {
		fmt.Print("hello, world")
		return nil
	}
	strategy, _ := NewFixedIntervalRetryStrategy(time.Millisecond*100, 3)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := Retry(ctx, strategy, bizFunc)
	if err != nil {
		fmt.Println("error:", err)
	}
	// Output:
	// hello, world
}
