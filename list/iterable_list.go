package list

import (
	"strings"

	"github.com/ecodeclub/ekit/internal/errs"
)

type IterableListImpl[T any] struct {
	List[T]
	hasIter bool
}

func NewIterableList[T any](list List[T]) *IterableListImpl[T] {
	return &IterableListImpl[T]{
		List: list,
	}
}

func (l *IterableListImpl[T]) GetIter() (Iter[T], bool) {
	if l.hasIter {
		return nil, false
	}
	l.hasIter = true
	return &IterImpl[T]{
		iterableList:   l,
		index:          -1,
		deletedIndices: make([]int, 0, 4),
	}, true
}

func (l *IterableListImpl[T]) releaseIter(f func(list List[T])) {
	l.hasIter = false
	f(l.List)
}

func (l *IterableListImpl[T]) Add(index int, t T) error {
	if l.hasIter {
		return errs.ErrNotEditableDuringIterating
	}
	return l.List.Add(index, t)
}

func (l *IterableListImpl[T]) Delete(index int) (T, error) {
	if l.hasIter {
		var empty T
		return empty, errs.ErrNotEditableDuringIterating
	}
	return l.List.Delete(index)
}

type IterImpl[T any] struct {
	iterableList   IterableList[T]
	index          int
	deletedIndices []int
}

func (i *IterImpl[T]) Next() (T, bool) {
	if i.index >= i.iterableList.Len() {
		var empty T
		return empty, false
	}
	i.index++
	ele, err := i.iterableList.Get(i.index)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "ekit: 下标超出范围") {
			panic(err)
		}
		// 释放Iter
		i.iterableList.releaseIter(i.doDelete)
		var empty T
		return empty, false
	}
	return ele, true
}

func (i *IterImpl[T]) Delete() {
	if i.index < 0 || i.index >= i.iterableList.Len() {
		return
	}
	i.deletedIndices = append(i.deletedIndices, i.index)
}

func (i *IterImpl[T]) Release() {
	i.iterableList.releaseIter(i.doDelete)
	i.index = i.iterableList.Len()
}

func (i *IterImpl[T]) doDelete(l List[T]) {
	prev := -1
	for j := len(i.deletedIndices) - 1; j >= 0; j-- {
		idx := i.deletedIndices[j]
		if idx != prev {
			if _, err := l.Delete(idx); err != nil {
				panic(err)
			}
			prev = idx
		}
	}
	i.deletedIndices = nil
}
