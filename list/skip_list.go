package list

import (
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/internal/list"
)

func NewSkipList[T any](compare ekit.Comparator[T]) *SkipList[T] {
	pq := &SkipList[T]{}
	pq.skiplist = list.NewSkipList[T](compare)
	return pq
}

func (sl *SkipList[T]) Peek() (T, error) {
	return sl.skiplist.Peek()
}

type SkipList[T any] struct {
	skiplist *list.SkipList[T]
}

func (sl *SkipList[T]) Get(index int) (T, error) {
	return sl.skiplist.Get(index)
}

func (sl *SkipList[T]) Search(target T) bool {
	return sl.skiplist.Search(target)
}

func (sl *SkipList[T]) AsSlice() []T {
	return sl.skiplist.AsSlice()
}

func (sl *SkipList[T]) Len() int {
	return sl.skiplist.Len()
}

func (sl *SkipList[T]) Cap() int {
	return sl.Len()
}

func (sl *SkipList[T]) Insert(Val T) {
	sl.skiplist.Insert(Val)
}

func (sl *SkipList[T]) DeleteElement(target T) bool {
	return sl.skiplist.DeleteElement(target)
}
