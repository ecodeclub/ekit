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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ecodeclub/ekit"

	"github.com/stretchr/testify/assert"
)

func TestNewPriorityQueue(t *testing.T) {
	data := []int{6, 5, 4, 3, 2, 1}
	testCases := []struct {
		name     string
		q        *PriorityQueue[int]
		capacity int
		data     []int
		expected []int
	}{
		{
			name:     "无边界",
			q:        NewPriorityQueue(0, compare()),
			capacity: 0,
			data:     data,
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "有边界 ",
			q:        NewPriorityQueue(len(data), compare()),
			capacity: len(data),
			data:     data,
			expected: []int{1, 2, 3, 4, 5, 6},
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
			assert.Equal(t, tc.expected, res)
		})

	}

}

func TestPriorityQueue_Peek(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
		data     []int
		wantErr  error
	}{
		{
			name:     "有数据",
			capacity: 0,
			data:     []int{6, 5, 4, 3, 2, 1},
			wantErr:  ErrEmptyQueue,
		},
		{
			name:     "无数据",
			capacity: 0,
			data:     []int{},
			wantErr:  ErrEmptyQueue,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewPriorityQueue[int](tc.capacity, compare())
			for _, el := range tc.data {
				err := q.Enqueue(el)
				require.NoError(t, err)
			}
			for q.Len() > 0 {
				peek, err := q.Peek()
				assert.NoError(t, err)
				el, _ := q.Dequeue()
				assert.Equal(t, el, peek)
			}
			_, err := q.Peek()
			assert.Equal(t, tc.wantErr, err)
		})

	}
}

func TestPriorityQueue_Enqueue(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
		data     []int
		element  int
		wantErr  error
	}{
		{
			name:     "有界空队列",
			capacity: 10,
			data:     []int{},
			element:  10,
		},
		{
			name:     "有界满队列",
			capacity: 6,
			data:     []int{6, 5, 4, 3, 2, 1},
			element:  10,
			wantErr:  ErrOutOfCapacity,
		},
		{
			name:     "有界非空不满队列",
			capacity: 12,
			data:     []int{6, 5, 4, 3, 2, 1},
			element:  10,
		},
		{
			name:     "无界空队列",
			capacity: 0,
			data:     []int{},
			element:  10,
		},
		{
			name:     "无界非空队列",
			capacity: 0,
			data:     []int{6, 5, 4, 3, 2, 1},
			element:  10,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := priorityQueueOf(tc.capacity, tc.data, compare())
			require.NotNil(t, q)
			err := q.Enqueue(tc.element)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.capacity, q.Cap())
		})

	}
}

func TestPriorityQueue_EnqueueElement(t *testing.T) {
	testCases := []struct {
		name      string
		data      []int
		element   int
		wantSlice []int
	}{
		{
			name:      "新加入的元素是最大的",
			data:      []int{10, 8, 7, 6, 2},
			element:   20,
			wantSlice: []int{0, 2, 6, 8, 10, 7, 20},
		},
		{
			name:      "新加入的元素是最小的",
			data:      []int{10, 8, 7, 6, 2},
			element:   1,
			wantSlice: []int{0, 1, 6, 2, 10, 7, 8},
		},
		{
			name:      "新加入的元素子区间中",
			data:      []int{10, 8, 7, 6, 2},
			element:   5,
			wantSlice: []int{0, 2, 6, 5, 10, 7, 8},
		},
		{
			name:      "新加入的元素与已有元素相同",
			data:      []int{10, 8, 7, 6, 2},
			element:   6,
			wantSlice: []int{0, 2, 6, 6, 10, 7, 8},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := priorityQueueOf(0, tc.data, compare())
			require.NotNil(t, q)
			err := q.Enqueue(tc.element)
			require.NoError(t, err)
			assert.Equal(t, tc.wantSlice, q.data)
		})

	}
}

func TestPriorityQueue_EnqueueHeapStruct(t *testing.T) {
	data := []int{6, 5, 4, 3, 2, 1}
	testCases := []struct {
		name      string
		capacity  int
		data      []int
		wantSlice []int
		pivot     int
		pivotData []int
	}{
		{
			name:      "队列满",
			capacity:  len(data),
			data:      data,
			pivot:     2,
			pivotData: []int{0, 4, 6, 5},
			wantSlice: []int{0, 1, 3, 2, 6, 4, 5},
		},
		{
			name:      "队列不满",
			capacity:  len(data) * 2,
			data:      data,
			pivot:     3,
			pivotData: []int{0, 3, 4, 5, 6},
			wantSlice: []int{0, 1, 3, 2, 6, 4, 5},
		},
		{
			name:      "无界队列",
			capacity:  0,
			data:      data,
			pivot:     3,
			pivotData: []int{0, 3, 4, 5, 6},
			wantSlice: []int{0, 1, 3, 2, 6, 4, 5},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewPriorityQueue[int](tc.capacity, compare())
			for i, el := range tc.data {
				require.NoError(t, q.Enqueue(el))
				// 检查中途堆结构堆调整，是否符合预期
				if i == tc.pivot {
					assert.Equal(t, tc.pivotData, q.data)
				}
			}
			// 检查最终堆结构，是否符合预期
			assert.Equal(t, tc.wantSlice, q.data)
		})

	}
}

func TestPriorityQueue_Dequeue(t *testing.T) {
	testCases := []struct {
		name      string
		data      []int
		wantErr   error
		wantVal   int
		wantSlice []int
	}{
		{
			name:    "空队列",
			data:    []int{},
			wantErr: ErrEmptyQueue,
		},
		{
			name:      "只有一个元素",
			data:      []int{10},
			wantVal:   10,
			wantSlice: []int{0},
		},
		{
			name:      "many",
			data:      []int{6, 5, 4, 3, 2, 1},
			wantVal:   1,
			wantSlice: []int{0, 2, 3, 5, 6, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := priorityQueueOf(0, tc.data, compare())
			require.NotNil(t, q)
			val, err := q.Dequeue()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSlice, q.data)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestPriorityQueue_DequeueComplexCheck(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
		data     []int
		pivot    int
		want     []int
	}{
		{
			name:     "无边界",
			capacity: 0,
			data:     []int{6, 5, 4, 3, 2, 1},
			pivot:    2,
			want:     []int{0, 4, 6, 5},
		},
		{
			name:     "有边界",
			capacity: 6,
			data:     []int{6, 5, 4, 3, 2, 1},
			pivot:    3,
			want:     []int{0, 5, 6},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := priorityQueueOf(tc.capacity, tc.data, compare())
			require.NotNil(t, q)
			i := 0
			prev := -1
			for q.Len() > 0 {
				el, err := q.Dequeue()
				require.NoError(t, err)
				// 检查中途出队后，堆结构堆调整是否符合预期
				if i == tc.pivot {
					assert.Equal(t, tc.want, q.data)
				}
				// 检查出队是否有序
				assert.LessOrEqual(t, prev, el)
				prev = el
				i++
			}
		})

	}
}

func TestPriorityQueue_Shrink(t *testing.T) {
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
		sliceCap    int
	}{
		{
			name:        "有界，小于64",
			originCap:   32,
			enqueueLoop: 6,
			dequeueLoop: 5,
			expectCap:   32,
			sliceCap:    33,
		},
		{
			name:        "有界，小于2048, 不足1/4",
			originCap:   1000,
			enqueueLoop: 20,
			dequeueLoop: 5,
			expectCap:   1000,
			sliceCap:    1001,
		},
		{
			name:        "有界，小于2048, 超过1/4",
			originCap:   1000,
			enqueueLoop: 400,
			dequeueLoop: 5,
			expectCap:   1000,
			sliceCap:    1001,
		},
		{
			name:        "有界，大于2048，不足一半",
			originCap:   3000,
			enqueueLoop: 400,
			dequeueLoop: 40,
			expectCap:   3000,
			sliceCap:    3001,
		},
		{
			name:        "有界，大于2048，大于一半",
			originCap:   3000,
			enqueueLoop: 2000,
			dequeueLoop: 5,
			expectCap:   3000,
			sliceCap:    3001,
		},
		{
			name:        "无界，小于64",
			originCap:   0,
			enqueueLoop: 30,
			dequeueLoop: 5,
			expectCap:   0,
			sliceCap:    64,
		},
		{
			name:        "无界，小于2048, 不足1/4",
			originCap:   0,
			enqueueLoop: 2000,
			dequeueLoop: 1990,
			expectCap:   0,
			sliceCap:    50,
		},
		{
			name:        "无界，小于2048, 超过1/4",
			originCap:   0,
			enqueueLoop: 2000,
			dequeueLoop: 600,
			expectCap:   0,
			sliceCap:    2560,
		},
		{
			name:        "无界，大于2048，不足一半",
			originCap:   0,
			enqueueLoop: 3000,
			dequeueLoop: 2000,
			expectCap:   0,
			sliceCap:    1331,
		},
		{
			name:        "无界，大于2048，大于一半",
			originCap:   0,
			enqueueLoop: 3000,
			dequeueLoop: 5,
			expectCap:   0,
			sliceCap:    3408,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewPriorityQueue[int](tc.originCap, compare)
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
			assert.Equal(t, tc.sliceCap, cap(q.data))
		})
	}
}

func priorityQueueOf(capacity int, data []int, compare ekit.Comparator[int]) *PriorityQueue[int] {
	q := NewPriorityQueue[int](capacity, compare)
	for _, el := range data {
		err := q.Enqueue(el)
		if err != nil {
			return nil
		}
	}
	return q
}

func compare() ekit.Comparator[int] {
	return func(a, b int) int {
		if a < b {
			return -1
		}
		if a == b {
			return 0
		}
		return 1
	}
}
