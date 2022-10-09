package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPriorityArrayQueue(t *testing.T) {
	data := []int{6, 5, 4, 3, 2, 1}
	expected := []int{1, 2, 3, 4, 5, 6}
	var less Less[int]
	less = func(a, b int) bool {
		return a < b
	}
	testCases := []struct {
		name     string
		q        *PriorityArrayQueue[int]
		capacity int
	}{
		{
			name:     "默认",
			q:        NewPriorityArrayQueueFromArray(data, less),
			capacity: len(data),
		},
		{
			name:     "capacity 小于默认",
			q:        NewPriorityArrayQueueFromArray(data, less, WithNewCapacity[int](len(data)-2)),
			capacity: len(data),
		},
		{
			name:     "capacity 大于默认",
			q:        NewPriorityArrayQueueFromArray(data, less, WithNewCapacity[int](8)),
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
	var less Less[int]
	less = func(a, b int) bool {
		return a < b
	}
	testCases := []struct {
		name     string
		q        *PriorityArrayQueue[int]
		capacity int
	}{
		{
			name:     "无边界",
			q:        NewBoundlessPriorityArrayQueue(less),
			capacity: 0,
		},
		{
			name:     "有边界 ",
			q:        NewPriorityArrayQueue(len(data), less),
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

func TestPriorityArrayQueue_Dequeue(t *testing.T) {
	data := []int{6, 5, 4, 3, 2, 1}
	expected := []int{1, 2, 3, 4, 5, 6}
	var less Less[int]
	less = func(a, b int) bool {
		return a < b
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
			q:         NewBoundlessPriorityArrayQueue[int](less),
			data:      data,
			wantSlice: expected,
			wantErr:   ErrEmptyQueue,
		},
		{
			name:      "有边界",
			capacity:  0,
			q:         NewPriorityArrayQueue[int](len(data), less),
			data:      data,
			wantSlice: expected,
			wantErr:   ErrEmptyQueue,
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
	var less Less[int]
	less = func(a, b int) bool {
		return a < b
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
			q:        NewPriorityArrayQueueFromArray[int](data, less),
			data:     data,
			wantErr:  ErrOutOfCapacity,
		},
		{
			name:      "队列不满",
			capacity:  len(data),
			q:         NewPriorityArrayQueue[int](len(data), less),
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
