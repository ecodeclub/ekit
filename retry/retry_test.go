package retry

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestExponentialIntervalRetry_Next(t *testing.T) {
	testCases := []struct {
		name         string
		strategy     *ExponentialIntervalRetry
		wantRetryCnt int
		wantInterval []time.Duration
	}{
		{
			name: "test max cnt",
			strategy: &ExponentialIntervalRetry{
				BeginInterval: time.Second,
				MaxInterval:   5 * time.Minute,
				Max:           5,
			},
			wantRetryCnt: 5,
			wantInterval: func() []time.Duration {
				var res []time.Duration
				beginInterval := time.Second
				for i := 1; i <= 5; i++ {
					res = append(res, time.Duration(math.Exp2(float64(i)))*beginInterval)
				}
				return res
			}(),
		},
		{
			name: "test max interval",
			strategy: &ExponentialIntervalRetry{
				BeginInterval: time.Second,
				MaxInterval:   5 * time.Second,
				Max:           5,
			},
			wantRetryCnt: 5,
			wantInterval: func() []time.Duration {
				var res []time.Duration
				beginInterval := time.Second
				maxInterval := 5 * time.Second
				for i := 1; i <= 2; i++ {
					res = append(res, time.Duration(math.Exp2(float64(i)))*beginInterval)
				}
				res = append(res, maxInterval)
				res = append(res, maxInterval)
				res = append(res, maxInterval)
				return res
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.wantRetryCnt; i++ {
				interval, ok := tc.strategy.Next()
				assert.Equal(t, true, ok)
				assert.Equal(t, tc.wantInterval[i], interval)
			}
			_, ok := tc.strategy.Next()
			assert.Equal(t, false, ok)
		})
	}
}

func TestFixIntervalRetry_Next(t *testing.T) {
	testCases := []struct {
		name         string
		strategy     *FixIntervalRetry
		wantRetryCnt int
		wantInterval time.Duration
	}{
		{
			name: "test",
			strategy: &FixIntervalRetry{
				Interval: time.Second,
				Max:      5,
			},
			wantRetryCnt: 5,
			wantInterval: time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.wantRetryCnt; i++ {
				interval, ok := tc.strategy.Next()
				assert.Equal(t, true, ok)
				assert.Equal(t, tc.wantInterval, interval)
			}
			_, ok := tc.strategy.Next()
			assert.Equal(t, false, ok)
		})
	}
}
