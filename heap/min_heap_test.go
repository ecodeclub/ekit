package heap

import (
	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	initSize = 8 //测试用的初始堆大小
)

// 用int类型的测试
func compareInt() ekit.Comparator[int] {
	return ekit.ComparatorRealNumber[int]
}

// 比较两个小根堆是否相同
func compareTwoMinHeap(src *MinHeap[int], dst *MinHeap[int]) bool {
	if src.Size() != dst.Size() {
		return false
	}
	for i := 0; i < src.Size(); i++ {
		if src.data[i] != dst.data[i] {
			return false
		}
	}
	return true
}

func TestNewMinHeap(t *testing.T) {
	testCases := []struct {
		name    string
		compare ekit.Comparator[int]
		wantErr error
	}{
		{
			name:    "不传compare，报错",
			compare: nil,
			wantErr: ErrMinHeapComparatorIsNull,
		},
		{
			name:    "正常new",
			compare: compareInt(),
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewMinHeap[int](tc.compare, initSize)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestAdd(t *testing.T) {
	testCases := []struct {
		name      string
		startHeap func() *MinHeap[int]
		add       int
		wantHeap  func() *MinHeap[int]
	}{
		{
			name: "堆内元素数量0，新增1",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				return mh
			},
			add: 10,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10}
				mh.size = 1
				return mh
			},
		},
		{
			name: "堆内元素数量1，根结点左边新增1，不交换",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10}
				mh.size = 1
				return mh
			},
			add: 20,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20}
				mh.size = 2
				return mh
			},
		},
		{
			name: "堆内元素数量1，根结点左边新增1，交换",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{20}
				mh.size = 1
				return mh
			},
			add: 10,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20}
				mh.size = 2
				return mh
			},
		},
		{
			name: "堆内元素数量2，根结点右边新增1，不交换",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20}
				mh.size = 2
				return mh
			},
			add: 30,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30}
				mh.size = 3
				return mh
			},
		},
		{
			name: "堆内元素数量2，根结点右边新增1，交换",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{20, 30}
				mh.size = 2
				return mh
			},
			add: 10,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 30, 20}
				mh.size = 3
				return mh
			},
		},
		{
			name: "堆内元素数量3，新增1，不交换",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30}
				mh.size = 3
				return mh
			},
			add: 40,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30, 40}
				mh.size = 4
				return mh
			},
		},
		{
			name: "堆内元素数量3，新增1，交换1层",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30}
				mh.size = 3
				return mh
			},
			add: 15,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 15, 30, 20}
				mh.size = 4
				return mh
			},
		},
		{
			name: "堆内元素数量3，新增1，交换2层",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30}
				mh.size = 3
				return mh
			},
			add: 5,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{5, 10, 30, 20}
				mh.size = 4
				return mh
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startHeap := tc.startHeap()
			startHeap.Add(tc.add)
			wantHeap := tc.wantHeap()
			assert.Equal(t, compareTwoMinHeap(startHeap, wantHeap), true)
		})
	}
}

func TestGetTop(t *testing.T) {
	testCases := []struct {
		name      string
		startHeap func() *MinHeap[int]
		wantTop   int
		wantErr   error
	}{
		{
			name: "堆内元素数量0，报错",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				return mh
			},
			wantTop: 0,
			wantErr: ErrMinHeapIsEmpty,
		},
		{
			name: "堆内元素数量1，获取后1",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10}
				mh.size = 1
				return mh
			},
			wantTop: 10,
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startHeap := tc.startHeap()
			top, err := startHeap.GetTop()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantTop, top)
		})
	}
}

func TestExtractTop(t *testing.T) {
	testCases := []struct {
		name      string
		startHeap func() *MinHeap[int]
		wantTop   int
		wantErr   error
		wantHeap  func() *MinHeap[int]
	}{
		{
			name: "堆内元素数量0，报错",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				return mh
			},
			wantTop: 0,
			wantErr: ErrMinHeapIsEmpty,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				return mh
			},
		},
		{
			name: "堆内元素数量1，提取后0",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10}
				mh.size = 1
				return mh
			},
			wantTop: 10,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				return mh
			},
		},
		{
			name: "堆内元素数量2，提取后1，交换1层",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20}
				mh.size = 2
				return mh
			},
			wantTop: 10,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{20}
				mh.size = 1
				return mh
			},
		},
		{
			name: "堆内元素数量3，提取后2，交换1层",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30}
				mh.size = 3
				return mh
			},
			wantTop: 10,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{20, 30}
				mh.size = 2
				return mh
			},
		},
		{
			name: "堆内元素数量4，提取后3，交换1层",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30, 40}
				mh.size = 4
				return mh
			},
			wantTop: 10,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{20, 40, 30}
				mh.size = 3
				return mh
			},
		},
		{
			name: "堆内元素数量5，提取后4，交换2层",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30, 40, 50}
				mh.size = 5
				return mh
			},
			wantTop: 10,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{20, 40, 30, 50}
				mh.size = 4
				return mh
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startHeap := tc.startHeap()
			top, err := startHeap.ExtractTop()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantTop, top)
			wantHeap := tc.wantHeap()
			assert.Equal(t, compareTwoMinHeap(startHeap, wantHeap), true)
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name      string
		startHeap func() *MinHeap[int]
		delete    int
		wantErr   error
		wantHeap  func() *MinHeap[int]
	}{
		{
			name: "堆内元素数量0，报错",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				return mh
			},
			delete:  0,
			wantErr: ErrValueNotInMinHeap,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				return mh
			},
		},
		{
			name: "堆内元素数量1，找不到，报错",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10}
				mh.size = 1
				return mh
			},
			delete:  20,
			wantErr: ErrValueNotInMinHeap,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10}
				mh.size = 1
				return mh
			},
		},
		{
			name: "堆内元素数量1，删除1",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10}
				mh.size = 1
				return mh
			},
			delete:  10,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				return mh
			},
		},
		{
			name: "堆内元素数量2，删除1，删除堆顶",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20}
				mh.size = 2
				return mh
			},
			delete:  10,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{20}
				mh.size = 1
				return mh
			},
		},
		{
			name: "堆内元素数量4，删除1，删除堆顶，往左边换",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30, 40}
				mh.size = 4
				return mh
			},
			delete:  10,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{20, 40, 30}
				mh.size = 3
				return mh
			},
		},
		{
			name: "堆内元素数量4，删除1，删除堆顶，往右边换",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 30, 20, 40}
				mh.size = 4
				return mh
			},
			delete:  10,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{20, 30, 40}
				mh.size = 3
				return mh
			},
		},
		{
			name: "堆内元素数量4，删除1，删除堆顶，换两层",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30, 40, 50}
				mh.size = 5
				return mh
			},
			delete:  10,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{20, 40, 30, 50}
				mh.size = 4
				return mh
			},
		},
		{
			name: "堆内元素数量4，删除1，删除子堆顶",
			startHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 20, 30, 40}
				mh.size = 4
				return mh
			},
			delete:  20,
			wantErr: nil,
			wantHeap: func() *MinHeap[int] {
				mh, _ := NewMinHeap[int](compareInt(), initSize)
				mh.data = []int{10, 40, 30}
				mh.size = 3
				return mh
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startHeap := tc.startHeap()
			err := startHeap.Delete(tc.delete)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			wantHeap := tc.wantHeap()
			assert.Equal(t, compareTwoMinHeap(startHeap, wantHeap), true)
		})
	}
}
