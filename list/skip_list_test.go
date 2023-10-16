package list

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func compareInt(a, b int) int {
	if a < b {
		return -1
	} else if a == b {
		return 0
	} else {
		return 1
	}
}

func TestNewSkipList(t *testing.T) {
	testCases := []struct {
		name       string
		compare    Comparator[int]
		level      int
		wantHeader *skipListNode[int]
		wantLevel  int
		wantSlice  []int
		wantErr    error
		wantSize   int
	}{
		{
			name:       "new skip list",
			compare:    compareInt,
			level:      1,
			wantLevel:  1,
			wantHeader: newSkipListNode[int](0, MaxLevel),
			wantSlice:  []int{},
			wantSize:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList(tc.compare)
			assert.Equal(t, tc.wantLevel, sl.level)
			assert.Equal(t, tc.wantHeader, sl.header)
			assert.Equal(t, tc.wantSlice, SkipListToSlice[int](sl))
			assert.Equal(t, tc.wantSize, sl.size)

		})
	}
}

func TestNewSkipListFromSlice(t *testing.T) {
	testCases := []struct {
		name    string
		compare Comparator[int]
		level   int
		slice   []int

		wantSlice []int
		wantErr   error
		wantSize  int
	}{
		{
			name:    "new skip list",
			compare: compareInt,
			level:   1,
			slice:   []int{1, 2, 3},

			wantSlice: []int{1, 2, 3},
			wantSize:  3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipListFromSlice[int](tc.slice, tc.compare)
			assert.Equal(t, tc.wantSlice, SkipListToSlice[int](sl))
			assert.Equal(t, tc.wantSize, sl.size)

		})
	}
}

//func TestSkipListToSlice(t *testing.T) {
//
//}

func TestSkipList_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		skiplist  *SkipList[int]
		compare   Comparator[int]
		value     int
		wantSlice []int
		wantSize  int
		wantRes   bool
	}{
		{
			name:      "delete 2 from [1,3]",
			compare:   compareInt,
			skiplist:  NewSkipListFromSlice[int]([]int{1, 3}, compareInt),
			value:     2,
			wantSlice: []int{1, 3},
			wantSize:  2,
			wantRes:   false,
		},
		{
			name:      "delete 1 from [1,3]",
			compare:   compareInt,
			skiplist:  NewSkipListFromSlice[int]([]int{1, 3}, compareInt),
			value:     1,
			wantSlice: []int{3},
			wantSize:  1,
			wantRes:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok := tc.skiplist.Delete(tc.value)
			assert.Equal(t, tc.wantSize, tc.skiplist.size)
			assert.Equal(t, tc.wantSlice, SkipListToSlice[int](tc.skiplist))
			assert.Equal(t, tc.wantRes, ok)
		})
	}
}

func TestSkipList_Insert(t *testing.T) {
	testCases := []struct {
		name      string
		skiplist  *SkipList[int]
		compare   Comparator[int]
		value     int
		wantSlice []int
		wantSize  int
	}{
		{
			name:      "insert 2 into [1,3]",
			compare:   compareInt,
			skiplist:  NewSkipListFromSlice[int]([]int{1, 3}, compareInt),
			value:     2,
			wantSlice: []int{1, 2, 3},
			wantSize:  3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.skiplist.Insert(tc.value)
			assert.Equal(t, tc.wantSize, tc.skiplist.size)
			assert.Equal(t, tc.wantSlice, SkipListToSlice[int](tc.skiplist))
		})
	}
}

func TestSkipList_Search(t *testing.T) {
	testCases := []struct {
		name      string
		skiplist  *SkipList[int]
		compare   Comparator[int]
		value     int
		wantSlice []int
		wantSize  int
		wantRes   bool
	}{
		{
			name:      "search 2 into [1,3]",
			compare:   compareInt,
			skiplist:  NewSkipListFromSlice[int]([]int{1, 3}, compareInt),
			value:     2,
			wantSlice: []int{1, 3},
			wantSize:  2,
			wantRes:   false,
		},
		{
			name:      "search 1 into [1,3]",
			compare:   compareInt,
			skiplist:  NewSkipListFromSlice[int]([]int{1, 3}, compareInt),
			value:     1,
			wantSlice: []int{1, 3},
			wantSize:  2,
			wantRes:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok := tc.skiplist.Search(tc.value)
			assert.Equal(t, tc.wantSize, tc.skiplist.size)
			assert.Equal(t, tc.wantSlice, SkipListToSlice[int](tc.skiplist))
			assert.Equal(t, tc.wantRes, ok)
		})
	}
}

func TestSkipList_randomLevel(t *testing.T) {
	sl := NewSkipListFromSlice[int]([]int{1, 2, 3}, compareInt)
	fmt.Println(sl.randomLevel())
}
