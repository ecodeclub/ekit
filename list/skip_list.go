package list

import (
	"golang.org/x/exp/rand"
	"time"
)

// skip list

const (
	FactorP  = 0.25 // fraction of FactorP of the nodes with level i pointers also have level i + 1 pointers.
	MaxLevel = 32
)

// for FactorP  = 0.25 and  MaxLevel = 32 the list can contain 2^64 elements

// the level of a list is the maximum level currently in the list (or 1 if the list is empty).
// the header of a list has forward pointers of the header  at levels higher than the current maximum level of the list point to nil.
type skipListNode[T any] struct {
	Val     T
	Forward []*skipListNode[T]
}

type Comparator[T any] func(a, b T) int

type SkipList[T any] struct {
	header  *skipListNode[T]
	level   int
	compare Comparator[T]
	size    int
}

func newSkipListNode[T any](Val T, level int) *skipListNode[T] {
	return &skipListNode[T]{Val, make([]*skipListNode[T], level)}
}

func SkipListToSlice[T any](sl *SkipList[T]) []T {
	curr := sl.header
	slice := make([]T, 0, sl.size)
	for curr.Forward[0] != nil {
		slice = append(slice, curr.Forward[0].Val)
		curr = curr.Forward[0]
	}
	return slice
}

func NewSkipListFromSlice[T any](slice []T, compare Comparator[T]) *SkipList[T] {
	sl := NewSkipList[T](compare)
	for _, n := range slice {
		sl.Insert(n)
	}
	return sl
}

func NewSkipList[T any](compare Comparator[T]) *SkipList[T] {
	return &SkipList[T]{
		header: &skipListNode[T]{
			Forward: make([]*skipListNode[T], MaxLevel),
		},
		level:   1,
		compare: compare,
	}
}

// levels are generated without reference to the number of elements in the list
func (sl *SkipList[T]) randomLevel() int {
	level := 1
	rnd := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	for rnd.Float64() < FactorP {
		level++
	} // here the level generated for 8 elements may larger than 8, but the probabilistic is low
	if level < MaxLevel {
		return level
	}
	return MaxLevel

}

func (sl *SkipList[T]) Search(target T) bool {
	curr := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for curr.Forward[i] != nil && sl.compare(curr.Forward[i].Val, target) < 0 {
			curr = curr.Forward[i]
		}
	}
	curr = curr.Forward[0] // 第0层 包含所有元素
	return curr != nil && sl.compare(curr.Val, target) == 0
}

// update[i] contains a pointer to the left of the location of the insertion/deletion for level i
// if an insertion generates a skipListNode with a level greater than the previous maximum level of the list, we update the maximum level of the list and initialize the appropriate portions of the update vector.
// after each deletion, we check to see if we need to update the level of list.

func (sl *SkipList[T]) Insert(Val T) {
	update := make([]*skipListNode[T], MaxLevel)
	curr := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for curr.Forward[i] != nil && sl.compare(curr.Forward[i].Val, Val) < 0 {
			curr = curr.Forward[i]
		}
		update[i] = curr
	}

	level := sl.randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.header
		}
		sl.level = level
	}

	// 插入新节点
	newNode := newSkipListNode[T](Val, level)
	for i := 0; i < level; i++ {
		newNode.Forward[i] = update[i].Forward[i]
		update[i].Forward[i] = newNode
	}

	sl.size += 1

}

func (sl *SkipList[T]) Delete(target T) bool {
	update := make([]*skipListNode[T], MaxLevel)
	curr := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for curr.Forward[i] != nil && sl.compare(curr.Forward[i].Val, target) < 0 {
			curr = curr.Forward[i]
		}
		update[i] = curr
	}
	node := curr.Forward[0]
	if node == nil || sl.compare(node.Val, target) != 0 {
		return false
	}
	// 删除target结点
	for i := 0; i < sl.level && update[i].Forward[i] == node; i++ {
		update[i].Forward[i] = node.Forward[i]
	}

	// 更新层级
	for sl.level > 1 && sl.header.Forward[sl.level-1] == nil {
		sl.level--
	}
	sl.size -= 1
	return true
}
