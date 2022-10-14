// Copyright 2021 gotomicro
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

package queue

import (
	"sync"
	"testing"
	"time"

	"github.com/gotomicro/ekit"
	"github.com/gotomicro/ekit/internal/queue"
	"github.com/stretchr/testify/assert"
)

var (
	errOutOfCapacity = queue.ErrOutOfCapacity
	errEmptyQueue    = queue.ErrEmptyQueue
)

func TestNewConcurrentPriorityQueue(t *testing.T) {
	data := []int{6, 5, 4, 3, 2, 1}
	expected := []int{1, 2, 3, 4, 5, 6}
	var compare ekit.Comparator[int] = func(a, b int) int {
		if a < b {
			return -1
		}
		if a == b {
			return 0
		}
		return 1
	}
	testCases := []struct {
		name     string
		q        *ConcurrentPriorityQueue[int]
		capacity int
	}{
		{
			name:     "无边界",
			q:        NewConcurrentPriorityQueue(0, compare),
			capacity: 0,
		},
		{
			name:     "有边界 ",
			q:        NewConcurrentPriorityQueue(len(data), compare),
			capacity: len(data),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, 0, tc.q.Len())
			for _, d := range data {
				err := tc.q.Enqueue(d)
				assert.NoError(t, err)
				if err != nil {
					return
				}
			}
			assert.Equal(t, tc.capacity, tc.q.Cap())
			assert.Equal(t, len(data), tc.q.Len())
			res := make([]int, 0, len(data))
			for tc.q.Len() > 0 {
				el, err := tc.q.Dequeue()
				assert.NoError(t, err)
				if err != nil {
					return
				}
				res = append(res, el)
			}
			assert.Equal(t, expected, res)
		})

	}

}

func TestConcurrentPriorityQueue_Enqueue(t *testing.T) {
	var compare ekit.Comparator[int] = func(a, b int) int {
		if a < b {
			return -1
		}
		if a == b {
			return 0
		}
		return 1
	}

	testCases := []struct {
		name        string
		capacity    int
		concurrency int
		perRoutine  int
		wantSlice   []int
		remain      int
		wantErr     error
		errCount    int
	}{
		{
			name:        "不超过capacity",
			capacity:    1100,
			concurrency: 100,
			perRoutine:  10,
		},
		{
			name:        "超过capacity",
			capacity:    1000,
			concurrency: 101,
			perRoutine:  10,
			wantErr:     errOutOfCapacity,
			errCount:    10,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewConcurrentPriorityQueue[int](tc.capacity, compare)
			wg := sync.WaitGroup{}
			wg.Add(tc.concurrency)
			errChan := make(chan error, tc.capacity)
			for i := tc.concurrency; i > 0; i-- {
				go func(i int) {
					start := i * 10
					for j := 0; j < tc.perRoutine; j++ {
						err := q.Enqueue(start + j)
						if err != nil {
							errChan <- err
						}
					}
					wg.Done()
				}(i)
			}
			wg.Wait()
			prev := -1
			for q.Len() > 0 {
				el, _ := q.Dequeue()
				assert.Less(t, prev, el)
				prev = el
			}
			assert.Equal(t, tc.errCount, len(errChan))
		})

	}
}

func TestConcurrentPriorityQueue_Dequeue(t *testing.T) {
	var compare ekit.Comparator[int] = func(a, b int) int {
		if a < b {
			return -1
		}
		if a == b {
			return 0
		}
		return 1
	}

	testCases := []struct {
		name        string
		total       int
		concurrency int
		perRoutine  int
		wantSlice   []int
		remain      int
		wantErr     error
		errCount    int
	}{
		{
			name:        "入队大于出队",
			total:       910,
			concurrency: 100,
			perRoutine:  9,
			remain:      10,
			wantErr:     errEmptyQueue,
		},
		{
			name:        "入队小于出队",
			total:       900,
			concurrency: 101,
			perRoutine:  9,
			remain:      0,
			wantErr:     errEmptyQueue,
			errCount:    9,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewConcurrentPriorityQueue[int](tc.total, compare)
			resultChan := make(chan int, tc.concurrency*tc.perRoutine)
			errChan := make(chan error, tc.errCount)
			for i := tc.total; i > 0; i-- {
				err := q.Enqueue(i)
				assert.NoError(t, err)
				if err != nil {
					return
				}
			}
			wg := sync.WaitGroup{}
			wg.Add(tc.concurrency * tc.perRoutine)
			for i := 0; i < tc.concurrency; i++ {
				go func() {
					for i := 0; i < tc.perRoutine; i++ {
						el, err := q.Dequeue()
						if err != nil {
							assert.Equal(t, tc.wantErr, err)
							errChan <- err
						}
						go func() {
							// TODO: 更合理的计算时间
							d := el
							time.Sleep(time.Millisecond * time.Duration(d))
							resultChan <- el
							wg.Done()
						}()
					}
				}()
			}
			wg.Wait()
			assert.Equal(t, tc.remain, q.Len())
			assert.Equal(t, tc.errCount, len(errChan))
			close(resultChan)
			prev := -1
			for {
				el, ok := <-resultChan
				if !ok {
					break
				}
				// assert.Less(t, prev, el)
				if prev > el {
					return
				}
				prev = el
			}

		})

	}
}
