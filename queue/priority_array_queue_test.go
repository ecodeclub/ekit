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
	"testing"

	"github.com/gotomicro/ekit"

	"github.com/stretchr/testify/assert"
)

func TestNewPriorityArrayQueue(t *testing.T) {
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
		q        *PriorityArrayQueue[int]
		capacity int
	}{
		{
			name:     "默认",
			q:        NewPriorityArrayQueueFromArray(data, compare),
			capacity: len(data),
		},
		{
			name:     "capacity 小于默认",
			q:        NewPriorityArrayQueueFromArray(data, compare, WithNewCapacity[int](len(data)-2)),
			capacity: len(data),
		},
		{
			name:     "capacity 大于默认",
			q:        NewPriorityArrayQueueFromArray(data, compare, WithNewCapacity[int](8)),
			capacity: 8,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, len(data), tc.q.Len())
			assert.Equal(t, tc.capacity, tc.q.Cap())
			res := make([]int, 0, 6)
			for tc.q.Len() > 0 {
				el, err := tc.q.Dequeue()
				assert.Nil(t, err)
				if err != nil {
					return
				}
				res = append(res, el)
			}
			assert.Equal(t, expected, res)
		})

	}

}

func TestNewEmptyPriorityArrayQueue(t *testing.T) {
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
		q        *PriorityArrayQueue[int]
		capacity int
	}{
		{
			name:     "无边界",
			q:        NewBoundlessPriorityArrayQueue(compare),
			capacity: 0,
		},
		{
			name:     "有边界 ",
			q:        NewPriorityArrayQueue(len(data), compare),
			capacity: len(data),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, 0, tc.q.Len())
			for _, d := range data {
				err := tc.q.Enqueue(d)
				assert.Nil(t, err)
				if err != nil {
					return
				}
			}
			assert.Equal(t, tc.capacity, tc.q.Cap())
			assert.Equal(t, len(data), tc.q.Len())
			res := make([]int, 0, len(data))
			for tc.q.Len() > 0 {
				el, err := tc.q.Dequeue()
				assert.Nil(t, err)
				if err != nil {
					return
				}
				res = append(res, el)
			}
			assert.Equal(t, expected, res)
		})

	}

}

func TestPriorityArrayQueue_Peek(t *testing.T) {
	data := []int{6, 5, 4, 3, 2, 1}
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
		capacity int
		q        *PriorityArrayQueue[int]
		data     []int
		wantErr  error
	}{
		{
			name:     "无边界",
			capacity: 0,
			q:        NewBoundlessPriorityArrayQueue[int](compare),
			data:     data,
			wantErr:  errEmptyQueue,
		},
		{
			name:     "有边界",
			capacity: 0,
			q:        NewPriorityArrayQueue[int](len(data), compare),
			data:     data,
			wantErr:  errEmptyQueue,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, el := range tc.data {
				err := tc.q.Enqueue(el)
				assert.Nil(t, err)
				if err != nil {
					return
				}
			}
			for tc.q.Len() > 0 {
				peek, err := tc.q.Peek()
				assert.Nil(t, err)
				el, _ := tc.q.Dequeue()
				assert.Equal(t, el, peek)
			}
			_, err := tc.q.Peek()
			assert.Equal(t, tc.wantErr, err)
		})

	}
}

func TestPriorityArrayQueue_Dequeue(t *testing.T) {
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
		name      string
		capacity  int
		q         *PriorityArrayQueue[int]
		data      []int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "无边界",
			capacity:  0,
			q:         NewBoundlessPriorityArrayQueue[int](compare),
			data:      data,
			wantSlice: expected,
			wantErr:   errEmptyQueue,
		},
		{
			name:      "有边界",
			capacity:  0,
			q:         NewPriorityArrayQueue[int](len(data), compare),
			data:      data,
			wantSlice: expected,
			wantErr:   errEmptyQueue,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, el := range tc.data {
				err := tc.q.Enqueue(el)
				assert.Nil(t, err)
				if err != nil {
					return
				}
			}
			res := make([]int, 0, len(tc.data))
			for tc.q.Len() > 0 {
				el, err := tc.q.Dequeue()
				assert.Nil(t, err)
				if err != nil {
					return
				}
				res = append(res, el)
			}
			assert.Equal(t, tc.wantSlice, res)
			_, err := tc.q.Dequeue()
			assert.Equal(t, tc.wantErr, err)
		})

	}
}

func TestPriorityArrayQueue_Enqueue(t *testing.T) {
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
		name      string
		capacity  int
		q         *PriorityArrayQueue[int]
		data      []int
		wantSlice []int
		wantErr   error
	}{
		{
			name:     "队列满",
			capacity: len(data),
			q:        NewPriorityArrayQueueFromArray[int](data, compare),
			data:     data,
			wantErr:  errOutOfCapacity,
		},
		{
			name:      "队列不满",
			capacity:  len(data),
			q:         NewPriorityArrayQueue[int](len(data), compare),
			data:      data,
			wantSlice: expected,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.q.Len() == tc.q.capacity {
				err := tc.q.Enqueue(1)
				assert.Equal(t, tc.wantErr, err)
				if err != nil {
					return
				}
			}
			for _, el := range tc.data {
				err := tc.q.Enqueue(el)
				assert.Nil(t, err)
				if err != nil {
					return
				}
			}
		})

	}
}

func TestPriorityArrayQueue_Shrink(t *testing.T) {
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
		originCap   int
		enqueueLoop int
		dequeueLoop int
		expectCap   int
	}{
		{
			name:        "小于64",
			originCap:   32,
			enqueueLoop: 6,
			dequeueLoop: 5,
			expectCap:   32,
		},
		{
			name:        "小于2048, 不足1/4",
			originCap:   1000,
			enqueueLoop: 20,
			dequeueLoop: 5,
			expectCap:   61,
		},
		{
			name:        "小于2048, 超过1/4",
			originCap:   1000,
			enqueueLoop: 400,
			dequeueLoop: 5,
			expectCap:   1000,
		},
		{
			name:        "大于2048，不足一半",
			originCap:   3000,
			enqueueLoop: 400,
			dequeueLoop: 40,
			expectCap:   936,
		},
		{
			name:        "大于2048，不足一半",
			originCap:   3000,
			enqueueLoop: 60,
			dequeueLoop: 40,
			expectCap:   57,
		},
		{
			name:        "大于2048，大于一半",
			originCap:   3000,
			enqueueLoop: 2000,
			dequeueLoop: 5,
			expectCap:   3000,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewPriorityArrayQueue[int](tc.originCap, compare)
			for i := 0; i < tc.enqueueLoop; i++ {
				err := q.Enqueue(i)
				if err != nil {
					return
				}
			}
			for i := 0; i < tc.dequeueLoop; i++ {
				_, err := q.Dequeue()
				if err != nil {
					return
				}
			}
			assert.Equal(t, tc.expectCap, q.Cap())
		})
	}
}