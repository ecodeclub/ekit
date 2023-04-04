package retry

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBackoffIntervalRetryStrategy_Next(t *testing.T) {
	type fields struct {
		initialInterval time.Duration
		maxInterval     time.Duration
		maxRetries      int32
		multiplier      float64

		interval time.Duration
		retries  int32
	}

	tests := []struct {
		name        string
		fields      fields
		curInterval time.Duration
		isContinue  bool
	}{
		{
			name: "init case, retries 0 and interval is 0",
			fields: fields{
				initialInterval: time.Second,
				maxInterval:     30 * time.Second,
				maxRetries:      3,
				multiplier:      2,
			},
			curInterval: 2 * time.Second,
			isContinue:  true,
		},
		{
			name: "init case, retries 0",
			fields: fields{
				initialInterval: time.Second,
				maxInterval:     30 * time.Second,
				maxRetries:      3,
				multiplier:      2,
				interval:        time.Second,
			},
			curInterval: 2 * time.Second,
			isContinue:  true,
		},
		{
			name: "init case, interval is 0",
			fields: fields{
				initialInterval: time.Second,
				maxInterval:     30 * time.Second,
				maxRetries:      3,
				multiplier:      2,
				retries:         1,
			},
			curInterval: 2 * time.Second,
			isContinue:  true,
		},
		{
			name: "interval equal to maxInterval after the increase",
			fields: fields{
				initialInterval: time.Second,
				maxInterval:     32 * time.Second,
				maxRetries:      5,
				multiplier:      2,
				interval:        16 * time.Second,
			},
			curInterval: 32 * time.Second,
			isContinue:  true,
		},
		{
			name: "interval over maxInterval after the increase",
			fields: fields{
				initialInterval: time.Second,
				maxInterval:     30 * time.Second,
				maxRetries:      5,
				multiplier:      2,
				interval:        16 * time.Second,
			},
			curInterval: 0 * time.Second,
			isContinue:  false,
		},
		{
			name: "retries equal maxRetries after the increase",
			fields: fields{
				initialInterval: time.Second,
				maxInterval:     32 * time.Second,
				maxRetries:      5,
				multiplier:      2,
				retries:         4,
				interval:        16 * time.Second,
			},
			curInterval: 32 * time.Second,
			isContinue:  true,
		},
		{
			name: "retries over maxRetries after the increase",
			fields: fields{
				initialInterval: time.Second,
				maxInterval:     30 * time.Second,
				maxRetries:      3,
				multiplier:      2,
				retries:         3,
			},
			curInterval: 0,
			isContinue:  false,
		},
		{
			name: "maxRetries equals to 0",
			fields: fields{
				initialInterval: time.Second,
				maxInterval:     30 * time.Second,
				maxRetries:      0,
				multiplier:      2,
				interval:        2 * time.Second,
			},
			curInterval: 4 * time.Second,
			isContinue:  true,
		},
		{
			name: "maxRetries equals to 0, interval is 0",
			fields: fields{
				initialInterval: time.Second,
				maxInterval:     30 * time.Second,
				maxRetries:      0,
				multiplier:      2,
			},
			curInterval: 2 * time.Second,
			isContinue:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.interval == 0 {
				tt.fields.interval = tt.fields.initialInterval
			}
			bo := &BackoffIntervalRetryStrategy{
				initialInterval: tt.fields.initialInterval,
				maxInterval:     tt.fields.maxInterval,
				maxRetries:      tt.fields.maxRetries,
				multiplier:      tt.fields.multiplier,
				retries:         tt.fields.retries,
				interval:        tt.fields.interval,
			}
			interval, isContinue := bo.Next()
			assert.Equal(t, tt.curInterval, interval)
			assert.Equal(t, tt.isContinue, isContinue)
		})
	}
}
