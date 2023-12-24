package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterableList_IterNext(t *testing.T) {
	testCase := struct {
		name      string
		list      IterableList[int]
		wantSlice []int
	}{
		name:      "range the list",
		list:      NewIterableList[int](NewArrayListOf([]int{1, 2, 3})),
		wantSlice: []int{1, 2, 3},
	}
	t.Run(testCase.name, func(t *testing.T) {
		iter, ok := testCase.list.GetIter()
		assert.Equal(t, true, ok)
		slice := make([]int, 0, testCase.list.Len())
		for {
			ele, ok := iter.Next()
			if !ok {
				break
			}
			slice = append(slice, ele)
		}
		assert.Equal(t, testCase.wantSlice, slice)
	})
}

func TestIterableList_DuplicatedIter(t *testing.T) {
	testCase := struct {
		name string
		list IterableList[int]
	}{
		name: "range the list",
		list: NewIterableList[int](NewArrayListOf([]int{1, 2, 3})),
	}
	t.Run(testCase.name, func(t *testing.T) {
		iter, ok := testCase.list.GetIter()
		assert.Equal(t, true, ok)
		_, ok = testCase.list.GetIter()
		assert.Equal(t, false, ok)
		for {
			_, ok := iter.Next()
			if !ok {
				break
			}
		}
		iter, ok = testCase.list.GetIter()
		assert.Equal(t, true, ok)
	})
}

func TestIterableList_ModifyDuringIterating(t *testing.T) {
	testCase := struct {
		name       string
		list       IterableList[int]
		wantSlice1 []int
		wantSlice2 []int
	}{
		name:       "range the list",
		list:       NewIterableList[int](NewArrayListOf([]int{1, 2, 3})),
		wantSlice1: []int{1, 2, 3},
		wantSlice2: []int{4, 1, 2, 3},
	}
	t.Run(testCase.name, func(t *testing.T) {
		iter, ok := testCase.list.GetIter()
		assert.Equal(t, true, ok)
		err := testCase.list.Add(0, 4)
		assert.Equal(t, true, err != nil)
		slice := make([]int, 0, testCase.list.Len())
		for {
			ele, ok := iter.Next()
			if !ok {
				break
			}
			slice = append(slice, ele)
		}
		assert.Equal(t, testCase.wantSlice1, slice)
		err = testCase.list.Add(0, 4)
		assert.Equal(t, nil, err)
		assert.Equal(t, testCase.wantSlice2, testCase.list.AsSlice())
	})
}

func TestIterableList_IterRelease(t *testing.T) {
	testCase := struct {
		name       string
		list       IterableList[int]
		wantSlice1 []int
		wantSlice2 []int
		wantSlice3 []int
	}{
		name:       "range the list",
		list:       NewIterableList[int](NewArrayListOf([]int{1, 2, 3})),
		wantSlice1: []int{1, 2, 3},
		wantSlice2: []int{4, 1, 2, 3},
		wantSlice3: []int{4, 2, 3},
	}
	t.Run(testCase.name, func(t *testing.T) {
		iter, ok := testCase.list.GetIter()
		assert.Equal(t, true, ok)
		err := testCase.list.Add(0, 4)
		assert.Equal(t, true, err != nil)
		_, err = testCase.list.Delete(0)
		assert.Equal(t, true, err != nil)
		iter.Release()
		assert.Equal(t, testCase.wantSlice1, testCase.list.AsSlice())
		err = testCase.list.Add(0, 4)
		assert.Equal(t, nil, err)
		assert.Equal(t, testCase.wantSlice2, testCase.list.AsSlice())
		_, err2 := testCase.list.Delete(1)
		assert.Equal(t, nil, err2)
		assert.Equal(t, testCase.wantSlice3, testCase.list.AsSlice())
		iter.Release() // Release again
		assert.Equal(t, testCase.wantSlice3, testCase.list.AsSlice())
	})
}

func TestIterableList_IterDelete(t *testing.T) {
	testCase := struct {
		name      string
		list      IterableList[int]
		wantSlice []int
	}{
		name:      "range the list",
		list:      NewIterableList[int](NewArrayListOf([]int{1, 2, 3, 4})),
		wantSlice: []int{2, 4},
	}
	t.Run(testCase.name, func(t *testing.T) {
		iter, ok := testCase.list.GetIter()
		assert.Equal(t, true, ok)
		for {
			ele, ok := iter.Next()
			if !ok {
				break
			}
			if ele&1 == 1 {
				iter.Delete()
				iter.Delete()
			}
		}
		assert.Equal(t, testCase.wantSlice, testCase.list.AsSlice())
	})
}

func TestIterableList_IterDeleteAndRelease(t *testing.T) {
	testCase := struct {
		name      string
		list      IterableList[int]
		wantSlice []int
	}{
		name:      "range the list",
		list:      NewIterableList[int](NewArrayListOf([]int{1, 2, 3, 4})),
		wantSlice: []int{2, 3, 4},
	}
	t.Run(testCase.name, func(t *testing.T) {
		iter, ok := testCase.list.GetIter()
		assert.Equal(t, true, ok)
		for i := 0; i < testCase.list.Len()>>1; i++ {
			ele, ok := iter.Next()
			if !ok {
				break
			}
			if ele&1 == 1 {
				iter.Delete()
			}
		}
		iter.Release()
		assert.Equal(t, testCase.wantSlice, testCase.list.AsSlice())
	})
}
