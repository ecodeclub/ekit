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

package queue

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/internal/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errOutOfCapacity = queue.ErrOutOfCapacity
	errEmptyQueue    = queue.ErrEmptyQueue
)

func TestNewConcurrentPriorityQueue(t *testing.T) {
	testCases := []struct {
		name     string
		q        *ConcurrentPriorityQueue[int]
		capacity int
		data     []int
		expect   []int
	}{
		{
			name:     "无边界",
			q:        NewConcurrentPriorityQueue(0, ekit.ComparatorRealNumber[int]),
			capacity: 0,
			data:     []int{6, 5, 4, 3, 2, 1},
			expect:   []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "有边界 ",
			q:        NewConcurrentPriorityQueue(6, ekit.ComparatorRealNumber[int]),
			capacity: 6,
			data:     []int{6, 5, 4, 3, 2, 1},
			expect:   []int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, 0, tc.q.Len())
			for _, d := range tc.data {
				require.NoError(t, tc.q.Enqueue(d))
			}
			assert.Equal(t, tc.capacity, tc.q.Cap())
			assert.Equal(t, len(tc.data), tc.q.Len())
			res := make([]int, 0, len(tc.data))
			for tc.q.Len() > 0 {
				head, err := tc.q.Peek()
				require.NoError(t, err)
				el, err := tc.q.Dequeue()
				require.NoError(t, err)
				assert.Equal(t, head, el)
				res = append(res, el)
			}
			assert.Equal(t, tc.expect, res)
		})

	}

}

// 多个go routine 执行入队操作，完成后，主携程把元素逐一出队，只要有序，可以认为并发入队没问题
func TestConcurrentPriorityQueue_Enqueue(t *testing.T) {
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
			q := NewConcurrentPriorityQueue[int](tc.capacity, ekit.ComparatorRealNumber[int])
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
			assert.Equal(t, tc.errCount, len(errChan))
			prev := -1
			for q.Len() > 0 {
				el, _ := q.Dequeue()
				assert.Less(t, prev, el)

				// 入队元素总数小于capacity时，应该所有元素都入队了，出队顺序应该依次加1
				if prev > -1 && len(errChan) == 0 {
					assert.Equal(t, prev+1, el)
				}
				prev = el
			}
		})

	}
}

// 预先入队一组数据，通过测试多个协程并发出队时，每个协程内出队元素有序，间接确认并发安全
func TestConcurrentPriorityQueue_Dequeue(t *testing.T) {
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
			q := NewConcurrentPriorityQueue[int](tc.total, ekit.ComparatorRealNumber[int])
			for i := tc.total; i > 0; i-- {
				require.NoError(t, q.Enqueue(i))
			}

			resultChan := make(chan int, tc.concurrency*tc.perRoutine)
			disOrderChan := make(chan bool, tc.concurrency*tc.perRoutine)
			errChan := make(chan error, tc.errCount)
			wg := sync.WaitGroup{}
			wg.Add(tc.concurrency)

			for i := 0; i < tc.concurrency; i++ {
				go func() {
					prev := -1
					for i := 0; i < tc.perRoutine; i++ {
						el, err := q.Dequeue()
						if err != nil {
							// 如果出队报错，把错误放到error通道，以便后续检查错误的内容和数量是否符合预期
							errChan <- err
						} else {
							// 如果出队不报错，则检查出队结果是否符合预期
							resultChan <- el
							if prev >= el {
								disOrderChan <- false
							}
							prev = el
						}

					}
					wg.Done()
				}()
			}
			wg.Wait()
			close(resultChan)
			close(errChan)
			close(disOrderChan)

			// 检查并发出队的元素数量，是否符合预期
			assert.Equal(t, tc.remain, q.Len())

			// 检查所有协程中的执行错误，是否符合预期
			assert.Equal(t, tc.errCount, len(errChan))
			for err := range errChan {
				assert.Equal(t, tc.wantErr, err)
			}

			// 每个协程内部，出队元素应该有序，检查是否发现无序的情况
			assert.Equal(t, 0, len(disOrderChan))

			// 每个协程的每次出队操作，出队元素都应该不同，检查是否符合预期
			resultSet := make(map[int]bool)
			for el := range resultChan {
				_, ok := resultSet[el]
				assert.Equal(t, false, ok)
				resultSet[el] = true
			}

		})

	}
}

// 测试同时并发出入队。只要并发安全，并发出入队后的剩余元素数量+报错数量应该符合预期
// TODO 有待设计更好的并发出入队测试方案
func TestConcurrentPriorityQueue_EnqueueDequeue(t *testing.T) {
	testCases := []struct {
		name    string
		enqueue int
		dequeue int
		remain  int
	}{
		{
			name:    "出队等于入队",
			enqueue: 50,
			dequeue: 50,
			remain:  0,
		},
		{
			name:    "出队小于入队",
			enqueue: 50,
			dequeue: 40,
			remain:  10,
		},
		{
			name:    "出队大于入队",
			enqueue: 50,
			dequeue: 60,
			remain:  -10,
		},
	}
	for _, tt := range testCases {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			q := NewConcurrentPriorityQueue[int](0, ekit.ComparatorRealNumber[int])
			errChan := make(chan error, tc.dequeue)
			wg := sync.WaitGroup{}
			wg.Add(tc.enqueue + tc.dequeue)
			go func() {
				for i := 0; i < tc.enqueue; i++ {
					go func(i int) {
						require.NoError(t, q.Enqueue(i))
						wg.Done()
					}(i)
				}
			}()
			go func() {
				for i := 0; i < tc.dequeue; i++ {
					_, err := q.Dequeue()
					if err != nil {
						errChan <- err
					}
					wg.Done()
				}
			}()

			wg.Wait()
			close(errChan)
			assert.Equal(t, tc.remain, q.Len()-len(errChan))
		})
	}
}

func ExampleNewConcurrentPriorityQueue() {
	q := NewConcurrentPriorityQueue[int](10, ekit.ComparatorRealNumber[int])
	_ = q.Enqueue(3)
	_ = q.Enqueue(2)
	_ = q.Enqueue(1)
	var vals []int
	val, _ := q.Dequeue()
	vals = append(vals, val)
	val, _ = q.Dequeue()
	vals = append(vals, val)
	val, _ = q.Dequeue()
	vals = append(vals, val)
	fmt.Println(vals)
	// Output:
	// [1 2 3]
}
