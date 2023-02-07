package stack

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestArrayDeque_AddLast(t *testing.T) {
	testClasses := []struct {
		name    string
		deque   *ArrayDeque[int]
		loopNum int
		wantErr error
		wantLen int
	}{
		{
			name:    "向尾部添加，并对比元素数量",
			deque:   NewArrayDeque[int](),
			loopNum: 1 << 16,
			wantLen: 1 << 16,
		},
		{
			name:    "向末尾添加，直到达到队列的最大容量",
			deque:   NewArrayDequeWithCap[int](1 << 29),
			loopNum: (1 << 29) - 1,
			wantErr: ErrOutOfCapacity,
			wantLen: (1 << 29) - 1,
		},
	}
	for _, tc := range testClasses {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.loopNum; i++ {
				err := tc.deque.AddLast(i)
				if err != nil {
					t.Logf("error: %v", err)
					assert.Equal(t, tc.wantErr, err)
					break
				}
			}
			assert.Equal(t, tc.wantLen, tc.deque.Len())
		})
	}
}

func TestArrayDeque_AddFirst(t *testing.T) {
	testClasses := []struct {
		name    string
		deque   *ArrayDeque[int]
		loopNum int
		wantErr error
		wantLen int
	}{
		{
			name:    "向头部添加，并对比元素数量",
			deque:   NewArrayDeque[int](),
			loopNum: 1 << 16,
			wantLen: 1 << 16,
		},
		{
			name:    "向头部添加，直到达到队列的最大容量",
			deque:   NewArrayDequeWithCap[int](1 << 29),
			loopNum: (1 << 29) - 1,
			wantErr: ErrOutOfCapacity,
			wantLen: (1 << 29) - 1,
		},
	}
	for _, tc := range testClasses {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.loopNum; i++ {
				err := tc.deque.AddFirst(i)
				if err != nil {
					assert.Equal(t, tc.wantErr, err)
					break
				}
			}
			assert.Equal(t, tc.wantLen, tc.deque.Len())
		})
	}
}

func TestArrayDeque_GetFirst(t *testing.T) {
	testClass := []struct {
		name    string
		deque   *ArrayDeque[int]
		newVal  int
		wantErr error
	}{
		{
			name:    "get first where deque is empty",
			deque:   NewArrayDeque[int](),
			wantErr: ErrEmpty,
		},
		{
			name:    "get first where deque not empty",
			deque:   NewArrayDeque[int](),
			newVal:  100,
			wantErr: nil,
		},
	}
	for _, tc := range testClass {
		val := tc.newVal
		if val != 0 {
			err := tc.deque.AddLast(val)
			assert.Nil(t, err)
		}
		first, err := tc.deque.GetFirst()
		if err != nil {
			assert.Equal(t, tc.wantErr, err)
		} else {
			assert.Equal(t, first, tc.newVal)
		}
	}
}

func TestArrayDeque_GetLast(t *testing.T) {
	testClass := []struct {
		name    string
		deque   *ArrayDeque[int]
		newVal  int
		wantErr error
	}{
		{
			name:    "get first where deque is empty",
			deque:   NewArrayDeque[int](),
			wantErr: ErrEmpty,
		},
		{
			name:    "get first where deque not empty",
			deque:   NewArrayDeque[int](),
			newVal:  100,
			wantErr: nil,
		},
	}
	for _, tc := range testClass {
		val := tc.newVal
		if val != 0 {
			err := tc.deque.AddFirst(val)
			assert.Nil(t, err)
		}
		first, err := tc.deque.GetLast()
		if err != nil {
			assert.Equal(t, tc.wantErr, err)
		} else {
			assert.Equal(t, first, tc.newVal)
		}
	}
}

func TestArrayDeque_RemoveFirst(t *testing.T) {
	testClasses := []struct {
		name     string
		deque    *ArrayDeque[int]
		wantErr  error
		newValue int
	}{
		{
			name:    "remove first where deque is empty",
			deque:   NewArrayDeque[int](),
			wantErr: ErrEmpty,
		},
		{
			name:     "remove a value",
			deque:    NewArrayDeque[int](),
			wantErr:  nil,
			newValue: 1,
		},
	}
	for _, tc := range testClasses {
		value := tc.newValue
		if value != 0 {
			err := tc.deque.AddLast(tc.newValue)
			assert.Nil(t, err)
		}
		first, err := tc.deque.RemoveFirst()
		if err != nil {
			assert.Equal(t, tc.wantErr, err)
			assert.True(t, tc.deque.IsEmpty())
		} else {
			assert.Equal(t, first, tc.newValue)
		}
	}
}

func TestArrayDeque_RemoveLast(t *testing.T) {
	testClasses := []struct {
		name     string
		deque    *ArrayDeque[int]
		wantErr  error
		newValue int
	}{
		{
			name:    "remove last where deque is empty",
			deque:   NewArrayDeque[int](),
			wantErr: ErrEmpty,
		},
		{
			name:     "remove a value",
			deque:    NewArrayDeque[int](),
			wantErr:  nil,
			newValue: 1,
		},
	}
	for _, tc := range testClasses {
		value := tc.newValue
		if value != 0 {
			err := tc.deque.AddFirst(tc.newValue)
			assert.Nil(t, err)
		}
		first, err := tc.deque.RemoveLast()
		if err != nil {
			assert.Equal(t, tc.wantErr, err)
			assert.True(t, tc.deque.IsEmpty())
		} else {
			assert.Equal(t, first, tc.newValue)
		}
	}
}

func TestArrayDeque_AsStack(t *testing.T) {
	testClasses := []struct {
		name         string
		deque        *ArrayDeque[int]
		compareStack *list.List
		pushNum      int
		wantErr      error
	}{
		{
			name:         "push and pop",
			deque:        NewArrayDeque[int](),
			compareStack: list.New(),
			pushNum:      10000,
			wantErr:      nil,
		},
	}
	for _, tc := range testClasses {
		deque := tc.deque
		stack := tc.compareStack
		for i := 0; i < tc.pushNum; i++ {
			err := deque.Push(i)
			assert.Nil(t, err)
			stack.PushBack(i)
		}
		for tc.deque.Len() > 0 {
			pop, err := deque.Pop()
			assert.Nil(t, err)
			compare := stack.Remove(stack.Back()).(int)
			assert.Equal(t, compare, pop)
		}
		assert.True(t, deque.IsEmpty())
		assert.True(t, stack.Len() == 0)
	}
}

func TestArrayDeque_AsQueue(t *testing.T) {
	testClasses := []struct {
		name         string
		deque        *ArrayDeque[int]
		compareStack *list.List
		pushNum      int
		wantErr      error
	}{
		{
			name:         "push and pop",
			deque:        NewArrayDeque[int](),
			compareStack: list.New(),
			pushNum:      10000,
			wantErr:      nil,
		},
	}
	for _, tc := range testClasses {
		deque := tc.deque
		stack := tc.compareStack
		for i := 0; i < tc.pushNum; i++ {
			err := deque.Enqueue(i)
			assert.Nil(t, err)
			stack.PushBack(i)
		}
		for tc.deque.Len() > 0 {
			front, err := deque.Dequeue()
			assert.Nil(t, err)
			compare := stack.Remove(stack.Front()).(int)
			assert.Equal(t, compare, front)
		}
		t.Logf("deque is empty: %v", deque.IsEmpty())
		assert.True(t, deque.IsEmpty())
		assert.True(t, stack.Len() == 0)
	}
}

func TestArrayDeque_Cap(t *testing.T) {
	testClasses := []struct {
		name        string
		deque       *ArrayDeque[int]
		newVals     []int
		expectedCap int
	}{
		{
			name:        "test init cap",
			deque:       NewArrayDeque[int](),
			newVals:     make([]int, 0, 0),
			expectedCap: 16,
		},
		{
			name:        "test init with expected cap",
			deque:       NewArrayDequeWithCap[int](20),
			newVals:     make([]int, 0, 0),
			expectedCap: 32,
		},
		{
			name:        "test capacity grow up",
			deque:       NewArrayDeque[int](),
			newVals:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expectedCap: 32,
		},
		{
			name:        "test init with expected cap and not grow",
			deque:       NewArrayDequeWithCap[int](17),
			newVals:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			expectedCap: 32,
		},
	}
	for _, tc := range testClasses {
		deque := tc.deque
		for _, val := range tc.newVals {
			err := deque.Push(val)
			assert.Nil(t, err)
		}
		assert.Equal(t, tc.expectedCap, deque.Cap())
	}
}

func TestArrayDeque_RandomOps(t *testing.T) {
	testClasses := []struct {
		name string
		ops  func(arrayDeque *ArrayDeque[int], list *list.List)
	}{
		{
			name: "AddLast",
			ops: func(arrayDeque *ArrayDeque[int], list *list.List) {
				val := rand.Int()
				err := arrayDeque.AddLast(val)
				if err != nil {
					assert.Equal(t, ErrOutOfCapacity, err)
				} else {
					list.PushBack(val)
				}
			},
		},
		{
			name: "GetLast",
			ops: func(arrayDeque *ArrayDeque[int], list *list.List) {
				last, err := arrayDeque.GetLast()
				if err != nil {
					assert.Equal(t, ErrEmpty, err)
					assert.True(t, arrayDeque.IsEmpty())
				} else {
					val := list.Back().Value.(int)
					assert.Equal(t, val, last)
				}
			},
		},
		{
			name: "RemoveLast",
			ops: func(arrayDeque *ArrayDeque[int], list *list.List) {
				last, err := arrayDeque.RemoveLast()
				if err != nil {
					assert.Equal(t, ErrEmpty, err)
					assert.True(t, arrayDeque.IsEmpty())
				} else {
					val := list.Remove(list.Back()).(int)
					assert.Equal(t, val, last)
				}
			},
		},
		{
			name: "AddFirst",
			ops: func(arrayDeque *ArrayDeque[int], list *list.List) {
				val := rand.Int()
				err := arrayDeque.AddFirst(val)
				if err != nil {
					assert.Equal(t, ErrOutOfCapacity, err)
				} else {
					list.PushFront(val)
				}
			},
		},
		{
			name: "GetFirst",
			ops: func(arrayDeque *ArrayDeque[int], list *list.List) {
				first, err := arrayDeque.GetFirst()
				if err != nil {
					assert.Equal(t, ErrEmpty, err)
					assert.True(t, arrayDeque.IsEmpty())
				} else {
					val := list.Front().Value.(int)
					assert.Equal(t, val, first)
				}
			},
		},
		{
			name: "RemoveFirst",
			ops: func(arrayDeque *ArrayDeque[int], list *list.List) {
				first, err := arrayDeque.RemoveFirst()
				if err != nil {
					assert.Equal(t, ErrEmpty, err)
					assert.True(t, arrayDeque.IsEmpty())
				} else {
					val := list.Remove(list.Front()).(int)
					assert.Equal(t, val, first)
				}
			},
		},
		{
			name: "IsEmpty",
			ops: func(arrayDeque *ArrayDeque[int], list *list.List) {
				empty := arrayDeque.IsEmpty()
				listEmpty := list.Len() == 0
				assert.Equal(t, empty, listEmpty)
			},
		},
	}
	NewArrayDeque[int]()
	opsNum := 10000000
	opsCount := len(testClasses)
	deque := NewArrayDeque[int]()
	stack := list.New()
	for i := 0; i < opsNum; i++ {
		intn := rand.Int()
		ops := intn % opsCount
		testClass := testClasses[ops]
		testClass.ops(deque, stack)
	}
}

func BenchmarkArrayDeque_RandomOps(b *testing.B) {
	deque := NewArrayDeque[int]()
	stack := list.New()
	b.Run("deque", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ops := rand.Int()
			arrayDequeRandomOps(deque, ops)
		}
	})
	b.Run("stack", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ops := rand.Int()
			listRandomOps(stack, ops)
		}
	})
}

func arrayDequeRandomOps(deque *ArrayDeque[int], ops int) {
	oprations := []func(d *ArrayDeque[int]){
		func(d1 *ArrayDeque[int]) {
			val := rand.Int()
			d1.AddLast(val)
		},
		func(d2 *ArrayDeque[int]) {
			d2.GetLast()
		},
		func(d3 *ArrayDeque[int]) {
			d3.RemoveLast()
		},
		func(d4 *ArrayDeque[int]) {
			val := rand.Int()
			d4.AddFirst(val)
		},
		func(d5 *ArrayDeque[int]) {
			d5.GetFirst()
		},
		func(d6 *ArrayDeque[int]) {
			d6.RemoveFirst()
		},
	}
	idx := ops % len(oprations)
	oprations[idx](deque)
}

func listRandomOps(stack *list.List, ops int) {
	oprations := []func(d *list.List){
		func(d1 *list.List) {
			val := rand.Int()
			d1.PushFront(val)
		},
		func(d2 *list.List) {
			d2.Back()
		},
		func(d3 *list.List) {
			if d3.Len() > 0 {
				d3.Remove(d3.Back())
			}
		},
		func(d4 *list.List) {
			val := rand.Int()
			d4.PushBack(val)
		},
		func(d5 *list.List) {
			d5.Front()
		},
		func(d6 *list.List) {
			if d6.Len() > 0 {
				d6.Remove(d6.Front())
			}
		},
	}
	idx := ops % len(oprations)
	oprations[idx](stack)
}
